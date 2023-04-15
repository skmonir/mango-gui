package parser

import (
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/oj/client"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"net/http"
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
	url        string
	isContest  bool
	platform   string
	contestId  string
	httpClient *http.Client
}

func (parser *Parser) ExtractUrlAndSetVars(url string) error {
	fmt.Println("Extracting " + url)

	platform, cid, pid, _ := utils.ExtractInfoFromUrl(url)

	if platform == "" || cid == "" {
		return errors.New("url is not correct")
	}

	parser.url = strings.Trim(url, "/")
	parser.platform = platform
	parser.contestId = cid
	parser.isContest = pid == ""
	_, parser.httpClient = client.GetHttpClientByPlatform(platform)

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

	var problemsToParse []models.Problem
	for _, prob := range problemList {
		if parser.isContest || strings.Contains(prob.Url, parser.url) {
			problemsToParse = append(problemsToParse, prob)
		}
	}

	return problemsToParse
}
