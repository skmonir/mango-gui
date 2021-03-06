package system

import (
	"bufio"
	"encoding/json"
	"errors"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/skmonir/mango-gui/context"
	"github.com/skmonir/mango-gui/models"
)

type AtcoderParser struct{}
type CodeforcesParser struct{}

func GetHtmlBody(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func FilterHtml(src []byte) []byte {
	newline := regexp.MustCompile(`<[\s/br]+?>`)
	src = newline.ReplaceAll(src, []byte("\n"))
	s := html.UnescapeString(string(src))
	return []byte(s)
}

func TrimIO(io string) string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(io))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	io = strings.Join(lines, "\n")
	return io
}

// ParseProblem parse problem to path
func ParseProblem(ctx *context.AppCtx, parser Parser, problemId string) (string, error) {
	URL := parser.GetProblemUrl(ctx.Config.CurrentContestId, problemId)
	body, err := GetHtmlBody(URL)
	if err != nil {
		return problemId, err
	}

	probName := parser.GetProblemName(body)
	timeLimit, memoryLimit := parser.GetProblemConstraints(body)

	if probName == "" {
		probName = problemId
	}

	input, output, err := parser.GetProblemSamples(body)
	if err != nil {
		return probName, err
	}

	testCaseList := make([]models.Testcase, len(input))
	for i := 0; i < len(input); i++ {
		testCaseList[i] = models.Testcase{
			Input:       string(input[i]),
			Output:      string(output[i]),
			TimeLimit:   timeLimit,
			MemoryLimit: memoryLimit,
		}
		testCaseList[i].Input = TrimIO(testCaseList[i].Input)
		testCaseList[i].Output = TrimIO(testCaseList[i].Output)
	}

	problem := models.Problem{
		Name:        probName,
		TimeLimit:   timeLimit,
		MemoryLimit: memoryLimit,
		Dataset:     testCaseList,
	}

	data, err := json.MarshalIndent(problem, "", " ")
	if err != nil {
		return probName, err
	}

	err = ctx.Config.ResolveTescasePath(problemId)
	if err != nil {
		return probName, err
	}

	testCasePath := ctx.Config.GetTestcaseFilePath(problemId)

	err = ioutil.WriteFile(testCasePath, data, 0644)
	if err != nil {
		return probName, err
	}

	return probName, err
}

func ParseContest(ctx *context.AppCtx, parser Parser) error {
	URL := parser.GetContestUrl(ctx.Config.CurrentContestId)
	body, err := GetHtmlBody(URL)
	if err != nil {
		return err
	}

	problemIdList := parser.GetProblemList(body)
	if len(problemIdList) == 0 {
		return errors.New("no problem found")
	}
	ctx.ProgressBar.Max = float64(len(problemIdList))

	for i := 0; i < len(problemIdList); i++ {
		status := "[PARSED] "
		problemName, err := ParseProblem(ctx, parser, problemIdList[i])
		if err != nil {
			status = "[FAILED] "
		}

		// *ctx.ParserUi.ParsedProblemStatus = append([]string{status + problemName}, (*ctx.ParserUi.ParsedProblemStatus)...)
		*ctx.ParserUi.ParsedProblemStatus = append(*ctx.ParserUi.ParsedProblemStatus, status+problemName)
		ctx.ParserUi.ParsedProblemListContainer.Refresh()
		ctx.ProgressBar.SetValue(float64(i + 1))
	}

	return nil
}

func Parse(ctx *context.AppCtx) error {
	var parser Parser

	if ctx.ParserUi.OnlineJudgeOptionSelect.Selected == "CodeForces" {
		parser = CodeforcesParser{}
	}

	contestId, problemId, err := parser.ParseContestAndProblemId(ctx.ParserUi.ContestIdInputField.Text)
	if err != nil {
		return err
	}

	if contestId == "" {
		return errors.New("please enter valid contest & problem id combination like 1512G")
	}

	ctx.ParserUi.CurrentContestId = contestId
	ctx.Config.CurrentContestId = contestId
	ctx.Config.OJ = ctx.ParserUi.OnlineJudgeOptionSelect.Selected
	if err := ctx.Config.SaveConfig(); err != nil {
		return err
	}
	ctx.HeaderUi.CurrentContestField.SetText(contestId)
	ctx.HeaderUi.CurrentOnlineJudge.SetText(ctx.ParserUi.OnlineJudgeOptionSelect.Selected)

	*ctx.ParserUi.ParsedProblemStatus = []string{}

	if problemId == "" {
		if err := ParseContest(ctx, parser); err != nil {
			return err
		}
	} else {
		ctx.ProgressBar.Max = 1

		problemName, err := ParseProblem(ctx, parser, problemId)
		status := "[PARSED] "
		if err != nil {
			status = "[FAILED] "
		}

		// *ctx.ParserUi.ParsedProblemStatus = append([]string{status + problemName}, (*ctx.ParserUi.ParsedProblemStatus)...)
		*ctx.ParserUi.ParsedProblemStatus = append(*ctx.ParserUi.ParsedProblemStatus, status+problemName)
		ctx.ParserUi.ParsedProblemListContainer.Refresh()
		ctx.ProgressBar.SetValue(1)
	}

	return nil
}
