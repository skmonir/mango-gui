package config

import "github.com/skmonir/mango-gui/backend/judge-framework/models"

type JudgeConfig struct {
	AppVersion         string                             `json:"appVersion"`
	WorkspaceDirectory string                             `json:"workspaceDirectory"`
	SourceDirectory    string                             `json:"sourceDirectory"`
	Author             string                             `json:"author"`
	JudgeAccInfo       map[string]models.JudgeAccountInfo `json:"judgeAccInfo"`
	ActiveTestingLang  string                             `json:"activeTestingLang"`
	TestingLangConfigs map[string]LanguageConfig          `json:"testingLangConfigs"`
	EditorPreference   EditorPreferences                  `json:"editorPreference"`
}

type LanguageConfig struct {
	Lang                string `json:"lang"`
	CompilationCommand  string `json:"compilationCommand"`
	CompilationFlags    string `json:"compilationFlags"`
	ExecutionCommand    string `json:"executionCommand"`
	ExecutionFlags      string `json:"executionFlags"`
	DefaultTemplatePath string `json:"defaultTemplatePath"`
	UserTemplatePath    string `json:"userTemplatePath"`
	FileExtension       string `json:"fileExtension"`
}

type EditorPreferences struct {
	Theme    string `json:"theme"`
	FontSize string `json:"fontSize"`
	TabSize  string `json:"tabSize"`
}
