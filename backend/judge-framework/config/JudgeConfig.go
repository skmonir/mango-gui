package config

type JudgeConfig struct {
	WorkspaceDirectory string           `json:"workspaceDirectory"`
	SourceDirectory    string           `json:"sourceDirectory"`
	Author             string           `json:"author"`
	ActiveLanguage     LanguageConfig   `json:"activeLanguage"`
	LanguageConfigs    []LanguageConfig `json:"languageConfigs"`
}

type LanguageConfig struct {
	Lang               string `json:"lang"`
	CompilationCommand string `json:"compilationCommand"`
	CompilationArgs    string `json:"compilationArgs"`
	TemplatePath       string `json:"templatePath"`
	FileExtension      string `json:"fileExtension"`
}
