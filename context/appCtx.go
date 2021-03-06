package context

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/skmonir/mango-gui/models"
)

type HeaderUiCtx struct {
	WorkSpaceDirChooser *widget.Button
	CurrentContestField *widget.Label
	CurrentOnlineJudge  *widget.Label
}

type ParserUiCtx struct {
	OnlineJudgeOptionSelect    *widget.Select
	ContestIdInputField        *widget.Entry
	ParsedProblemStatus        *[]string
	ParsedProblemListContainer *fyne.Container
	CurrentContestId           string
}

type TesterUiCtx struct {
	OnlineJudgeOptionSelect *widget.Select
	ContestIdInputField     *widget.Entry
	RunTestsButton          *widget.Button
	AddTestButton           *widget.Button
	RemoveTestButton        *widget.Button
	OpenSourceButton        *widget.Button
	ProblemNameListSelect   *widget.Select
	ProblemNameList         *[]string
	ProblemList             *[]models.Problem
	TestcaseList            *map[string][]models.ExecutionResult
	CurrentContestId        string
}

type AppCtx struct {
	Config             *AppConfig
	HeaderUi           *HeaderUiCtx
	ParserUi           *ParserUiCtx
	TesterUi           *TesterUiCtx
	ProgressBar        *widget.ProgressBar
	OnlineJudgeOptions *[]string
}
