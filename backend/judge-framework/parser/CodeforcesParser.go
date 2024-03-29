package parser

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type CodeforcesParser struct {
	Parser
}

const cfHost = "https://codeforces.com"

func (parser *CodeforcesParser) ParseProblemListOnContestPage() []models.Problem {
	defer utils.PanicRecovery()

	var problemList []models.Problem

	html, err := utils.GetBody(parser.httpClient, parser.getContestUrl())
	if err != nil {
		log.Println("Error occurred while trying to fetch the contest page")
		return problemList
	}

	body := soup.HTMLParse(string(html))

	problemElements := body.Find("select", "name", "submittedProblemIndex").Children()

	for _, link := range problemElements {
		fullName := link.Text()
		dividerIndex := strings.Index(fullName, "-")
		if dividerIndex != -1 {
			label := strings.ToUpper(strings.TrimSpace(fullName[0:dividerIndex]))
			name := strings.TrimSpace(fullName[dividerIndex+1:])
			url := parser.getProblemUrl(label)
			fmt.Println(label, name, url)
			problemList = append(problemList, models.Problem{
				Platform:  "codeforces",
				ContestId: parser.contestId,
				Label:     label,
				Name:      name,
				Url:       url,
				Status:    "none",
			})
		}
	}
	return problemList
}

func (parser *CodeforcesParser) ParseProblemConstraints(doc soup.Root) (int64, uint64) {
	defer utils.PanicRecovery()

	constraints := []int{2, 512}
	tags := []string{"time-limit", "memory-limit"}

	for index, tag := range tags {
		constElement := doc.Find("div", "class", tag)
		constText := strings.TrimSpace(constElement.Text())
		values := regexp.MustCompile("[0-9]+").FindAllString(constText, -1)
		if val, err := strconv.Atoi(values[0]); err == nil {
			constraints[index] = val
		}
	}

	return int64(constraints[0]), uint64(constraints[1])
}

func (parser *CodeforcesParser) ParseProblemSamples(doc soup.Root) []models.Testcase {
	defer utils.PanicRecovery()

	var testcases []models.Testcase

	sampleElement := doc.Find("div", "class", "sample-test")
	inputs := parser.extractSamples(sampleElement, "input")
	outputs := parser.extractSamples(sampleElement, "output")

	for i := 0; i < len(inputs); i++ {
		testcases = append(testcases, models.Testcase{
			Input:  utils.TrimIO(inputs[i]),
			Output: utils.TrimIO(outputs[i]),
		})
	}
	return testcases
}

func (parser *CodeforcesParser) extractSamples(sampleElement soup.Root, ioType string) []string {
	defer utils.PanicRecovery()

	ioElements := sampleElement.FindAll("div", "class", ioType)

	var ios []string
	for _, element := range ioElements {
		preElement := element.Find("pre")
		innerText := strings.TrimSpace(preElement.Text())
		if innerText != "" {
			innerText = strings.TrimSpace(preElement.HTML())
			innerText = strings.Replace(innerText, "<pre>", "", -1)
			innerText = strings.Replace(innerText, "</pre>", "", -1)
			innerText = strings.Replace(innerText, "<br>", "\n", -1)
			innerText = strings.Replace(innerText, "<br/>", "\n", -1)
			ios = append(ios, strings.TrimSpace(innerText))
		} else {
			io := ""
			for _, line := range preElement.Children() {
				if len(io) > 0 {
					io += "\n"
				}
				innerText = strings.TrimSpace(line.FullText())
				innerText = strings.Replace(innerText, "<br>", "\n", -1)
				innerText = strings.Replace(innerText, "<br/>", "\n", -1)
				if len(strings.TrimSpace(innerText)) > 0 {
					io += strings.TrimSpace(innerText)
				}
			}
			if len(io) > 0 {
				ios = append(ios, io)
			}
		}
	}
	return ios
}

func (parser *CodeforcesParser) getContestUrl() string {
	cid, _ := strconv.Atoi(parser.contestId)
	if cid > 100000 {
		return cfHost + "/gym" + "/" + parser.contestId
	}
	return strings.ToLower(cfHost + "/contest" + "/" + parser.contestId)
}

func (parser *CodeforcesParser) getProblemUrl(label string) string {
	return strings.ToLower(parser.getContestUrl() + "/" + "problem" + "/" + label)
}
