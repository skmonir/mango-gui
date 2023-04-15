package parser

import (
	"github.com/anaskhan96/soup"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type AtcoderParser struct {
	Parser
}

const acHost = "https://atcoder.jp"

func (parser *AtcoderParser) ParseProblemListOnContestPage() []models.Problem {
	defer utils.PanicRecovery()

	var problemList []models.Problem

	html, err := utils.GetBody(parser.httpClient, parser.getContestUrl())
	if err != nil {
		log.Println("Error occurred while trying to fetch the contest page")
		return problemList
	}

	body := soup.HTMLParse(string(html))

	problemElements := body.
		FindStrict("div", "id", "contest-nav-tabs").
		FindNextElementSibling().Find("table", "class", "table").
		Children()[3].
		FindAll("a")

	for index, row := range problemElements {
		if index%2 == 1 {
			label := strings.ToUpper(strings.TrimSpace(problemElements[index-1].Text()))
			name := strings.TrimSpace(row.Text())
			url := strings.ToLower(acHost + row.Attrs()["href"])
			log.Println(label, name, url)
			problemList = append(problemList, models.Problem{
				Platform:  "atcoder",
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

func (parser *AtcoderParser) ParseProblemConstraints(doc soup.Root) (int64, uint64) {
	defer utils.PanicRecovery()

	constraintElement := doc.
		Find("div", "id", "task-statement").FindPrevElementSibling()
	constraintText := strings.TrimSpace(constraintElement.Text())
	values := regexp.MustCompile("[0-9]+").FindAllString(constraintText, -1)

	var timeLimit int64 = 2       // default time-limit
	var memoryLimit uint64 = 1024 // default memory-limit

	if tl, err := strconv.Atoi(values[0]); err == nil {
		timeLimit = int64(tl)
	}
	if ml, err := strconv.Atoi(values[1]); err == nil {
		memoryLimit = uint64(ml)
	}

	return timeLimit, memoryLimit
}

func (parser *AtcoderParser) ParseProblemSamples(doc soup.Root) []models.Testcase {
	defer utils.PanicRecovery()

	var testcases []models.Testcase

	preElements := doc.FindAll("pre")
	ios := []string{}
	for _, elem := range preElements {
		if len(strings.TrimSpace(elem.Text())) > 0 {
			ios = append(ios, strings.TrimSpace(elem.Text()))
		}
	}
	ios = ios[0 : len(ios)/2] // io is parsed twice

	for index, io := range ios {
		if index%2 == 1 {
			testcases = append(testcases, models.Testcase{
				Input:  utils.TrimIO(ios[index-1]),
				Output: utils.TrimIO(io),
			})
		}
	}
	return testcases
}

func (parser *AtcoderParser) getContestUrl() string {
	return acHost + "/contests" + "/" + parser.contestId + "/" + "tasks"
}
