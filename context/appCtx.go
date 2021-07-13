package context

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type HeaderUiCtx struct {
	WorkSpaceDirChooser *widget.Button
	CurrentContestField *widget.Label
	CurrentOnlineJudge  *widget.Label
}

type ParserUiCtx struct {
	OnlineJudgeOptions         *widget.Select
	ContestIdInputField        *widget.Entry
	ParsedProblemStatus        *[]string
	ParsedProblemListContainer *fyne.Container
}

type TesterUiCtx struct {
}

type AppCtx struct {
	Config      *AppConfig
	HeaderUi    *HeaderUiCtx
	ParserUi    *ParserUiCtx
	TesterUi    *TesterUiCtx
	ProgressBar *widget.ProgressBar
}
