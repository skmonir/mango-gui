package parser

import (
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"path/filepath"
	"strings"
)

type IParser interface {
	ExtractUrlAndSetVars(string) error
	GetProblemsJsonFilePath() string
	FilterProblemsToParse([]models.Problem) []models.Problem
	ParseProblemListOnContestPage() []models.Problem
	ParseProblemConstraints(soup.Root) (int64, uint64)
	ParseProblemSamples(soup.Root) []models.Testcase
	GetPlatformAndContestId() (string, string)
}

type Parser struct {
	url              string
	isContest        bool
	platform         string
	contestId        string
	problemListCache []models.Problem
}

func (parser *Parser) ExtractUrlAndSetVars(url string) error {
	fmt.Println("Extracting " + url)

	platform, cid, pid := utils.ExtractInfoFromUrl(url)

	if platform == "" || cid == "" {
		return errors.New("url is not correct")
	}

	parser.url = strings.Trim(url, "/")
	parser.platform = platform
	parser.contestId = cid
	parser.isContest = pid == ""

	return nil
}

func (parser *Parser) GetProblemsJsonFilePath() string {
	conf := config.GetJudgeConfigFromCache()
	jsonFilePath := filepath.Join(conf.WorkspaceDirectory, parser.platform, parser.contestId, "problems.json")
	return jsonFilePath
}

func (parser *Parser) GetPlatformAndContestId() (string, string) {
	return parser.platform, parser.contestId
}

func (parser *Parser) FilterProblemsToParse(problemList []models.Problem) []models.Problem {
	defer utils.PanicRecovery()

	parser.problemListCache = problemList

	var problemsToParse []models.Problem
	for _, prob := range problemList {
		if parser.isContest || strings.Contains(prob.Url, parser.url) {
			problemsToParse = append(problemsToParse, prob)
		}
	}

	return problemsToParse
}
