package config

type JudgeConfig struct {
	WorkspaceDirectory string                    `json:"workspaceDirectory"`
	SourceDirectory    string                    `json:"sourceDirectory"`
	Author             string                    `json:"author"`
	ActiveLang         string                    `json:"activeLang"`
	LangConfigs        map[string]LanguageConfig `json:"langConfigs"`
}

type LanguageConfig struct {
	Lang               string `json:"lang"`
	CompilationCommand string `json:"compilationCommand"`
	CompilationFlags   string `json:"compilationFlags"`
	ExecutionCommand   string `json:"executionCommand"`
	ExecutionFlags     string `json:"executionFlags"`
	TemplatePath       string `json:"templatePath"`
	FileExtension      string `json:"fileExtension"`
}
