package config

type JudgeConfig struct {
	AppVersion         string                    `json:"appVersion"`
	WorkspaceDirectory string                    `json:"workspaceDirectory"`
	SourceDirectory    string                    `json:"sourceDirectory"`
	Author             string                    `json:"author"`
	ActiveLang         string                    `json:"activeLang"`
	LangConfigs        map[string]LanguageConfig `json:"langConfigs"`
	EditorPreference   EditorPreferences         `json:"editorPreference"`
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
	FontSize int    `json:"fontSize"`
	TabSize  int    `json:"tabSize"`
}
