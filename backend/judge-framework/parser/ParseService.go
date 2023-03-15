package parser

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/skmonir/mango-gui/backend/judge-framework/fileService"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/services"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"github.com/skmonir/mango-gui/backend/socket"
	"log"
	"strings"
	"sync"
)

var parsingCompleteChan = make(chan models.Problem)

func publishToParsingChannel(parsedProblem models.Problem) {
	parsingCompleteChan <- parsedProblem
}

func getProblemsToParse(parser IParser) []models.Problem {
	defer utils.PanicRecovery()

	problemList := services.GetProblemList(parser.GetPlatformAndContestId()) // Read from cache or json fileService
	if len(problemList) == 0 {
		problemList = parser.ParseProblemListOnContestPage()
		services.SaveProblemList(problemList) // Write into cache and json fileService
	}

	var problemsToParse = parser.FilterProblemsToParse(problemList)

	return problemsToParse
}

func parseProblem(parser IParser, problem models.Problem) {
	defer utils.PanicRecovery()

	problem.Status = "failed"

	html, err := utils.GetHtmlBody(problem.Url)
	if err != nil {
		log.Println("Failed to parse html from " + problem.Url)
		publishToParsingChannel(problem)
		return
	}

	body := soup.HTMLParse(html)

	timeLimit, memoryLimit := parser.ParseProblemConstraints(body)
	problem.TimeLimit = timeLimit
	problem.MemoryLimit = memoryLimit

	testcases := parser.ParseProblemSamples(body)
	for i, _ := range testcases {
		testcases[i].TimeLimit = timeLimit
		testcases[i].MemoryLimit = memoryLimit
	}
	fileService.SaveTestcasesIntoFiles(problem.Platform, problem.ContestId, problem.Label, testcases)

	problem.Status = "success"

	defer publishToParsingChannel(problem)
}

func parseProblems(parser IParser) []models.Problem {
	parsedProblemList := getProblemsToParse(parser)

	if parsedProblemList == nil || len(parsedProblemList) == 0 {
		log.Println("Problem list is empty")
		return []models.Problem{}
	} else {
		for i := 0; i < len(parsedProblemList); i++ {
			parsedProblemList[i].Status = "running"
		}
	}

	socket.PublishParseMessage(parsedProblemList)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			select {
			case parsedProblem := <-parsingCompleteChan:
				mu := sync.Mutex{}
				mu.Lock()
				for i := 0; i < len(parsedProblemList); i++ {
					if parsedProblemList[i].ContestId == parsedProblem.ContestId &&
						parsedProblemList[i].Label == parsedProblem.Label &&
						parsedProblemList[i].Status == "running" {
						parsedProblemList[i] = parsedProblem
						break
					}
				}
				socket.PublishParseMessage(parsedProblemList)
				unparsedProblemExists := false
				for i := 0; i < len(parsedProblemList); i++ {
					unparsedProblemExists = unparsedProblemExists || parsedProblemList[i].Status == "running"
				}
				if !unparsedProblemExists {
					wg.Done()
					fmt.Println("All problems are parsed")
					return
				}
				mu.Unlock()
			}
		}
	}()

	for i := 0; i < len(parsedProblemList); i++ {
		go parseProblem(parser, parsedProblemList[i])
	}

	wg.Wait()

	return parsedProblemList
}

func Parse(url string) []models.Problem {
	var parsedProblemList []models.Problem

	if url == "" {
		return parsedProblemList
	}

	var parser IParser
	if strings.Contains(url, "codeforces.com") {
		parser = &CodeforcesParser{
			Parser{},
		}
	} else if strings.Contains(url, "atcoder.jp") {
		parser = &AtcoderParser{
			Parser{},
		}
	} else {
		log.Println("Unknown platform")
		return parsedProblemList
	}

	if err := parser.ExtractUrlAndSetVars(url); err != nil {
		return parsedProblemList
	}

	parsedProblemList = parseProblems(parser)

	services.UpdateProblemList(parsedProblemList)
	fileService.CreateSourceFiles(parsedProblemList)
	services.UpdateProblemExecutionResultInCacheByUrl(url)
	return parsedProblemList
}
