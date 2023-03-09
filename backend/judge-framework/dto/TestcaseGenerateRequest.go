package dto

type TestcaseGenerateRequest struct {
	ProblemUrl          string `json:"problemUrl"`
	FileNum             int    `json:"fileNum"`
	FileMode            string `json:"fileMode"`
	FileName            string `json:"fileName"`
	TestPerFile         int    `json:"testPerFile"`
	SerialFrom          int    `json:"serialFrom"`
	InputDirectoryPath  string `json:"inputDirectoryPath"`
	GenerationProcess   string `json:"generationProcess"`
	GeneratorScriptPath string `json:"generatorScriptPath"`
	TgenScriptContent   string `json:"tgenScriptContent"`
}
