package parser

import (
	"github.com/anaskhan96/soup"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type CodemarshalParser struct {
	Parser
}

const cmHost = "https://algo.codemarshal.org/"

func (parser *CodemarshalParser) ParseProblemListOnContestPage() []models.Problem {
	defer utils.PanicRecovery()

	var problemList []models.Problem

	html, err := utils.GetHtmlBody(parser.getContestUrl())
	if err != nil {
		logger.Error("Error occurred while trying to fetch the contest page")
		return problemList
	}

	body := soup.HTMLParse(html)

	problemElements := body.Find("div", "class", "panel-problems").Children()[1].FindAll("a")

	for _, elem := range problemElements {
		childs := elem.Children()
		label := strings.TrimSpace(childs[1].Text())
		name := strings.TrimSpace(childs[2].Text())
		url := strings.ToLower(cmHost + elem.Attrs()["href"])
		log.Println(label, name, url)
		problemList = append(problemList, models.Problem{
			Platform:  "codemarshal",
			ContestId: parser.contestId,
			Label:     label,
			Name:      name,
			Url:       url,
			Status:    "none",
		})
	}
	return problemList
}

func (parser *CodemarshalParser) ParseProblemConstraints(doc soup.Root) (int64, uint64) {
	defer utils.PanicRecovery()

	constraints := []int{2, 2048}
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

func (parser *CodemarshalParser) ParseProblemSamples(doc soup.Root) []models.Testcase {
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

func (parser *CodemarshalParser) extractSamples(sampleElement soup.Root, ioType string) []string {
	defer utils.PanicRecovery()

	ioElements := sampleElement.FindAll("div", "class", ioType)

	var ios []string
	for _, element := range ioElements {
		preElement := element.Find("pre")
		if strings.TrimSpace(preElement.Text()) != "" {
			ios = append(ios, strings.TrimSpace(preElement.Text()))
		} else {
			io := ""
			for _, line := range preElement.Children() {
				if len(io) > 0 {
					io += "\n"
				}
				if len(strings.TrimSpace(line.Text())) > 0 {
					io += strings.TrimSpace(line.Text())
				}
			}
			if len(io) > 0 {
				ios = append(ios, io)
			}
		}
	}
	return ios
}

func (parser *CodemarshalParser) getContestUrl() string {
	return cmHost + "/contests/" + parser.contestId
}

func (parser *CodemarshalParser) getProblemUrl(label string) string {
	return parser.getContestUrl() + "/problems/" + label
}
