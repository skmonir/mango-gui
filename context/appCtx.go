package context

import (
	"fyne.io/fyne/v2/widget"
)

type HeaderUiCtx struct {
	WorkSpaceDirChooser *widget.Button
	ParsedProblemStatus *[]string
	CurrentContestField *widget.Label
}

type ParserUiCtx struct {
	OnlineJudgeOptions  *widget.Select
	ContestIdInputField *widget.Entry
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
