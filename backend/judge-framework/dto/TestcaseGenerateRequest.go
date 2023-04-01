package dto

type TestcaseGenerateRequest struct {
	IsForParsedProblem  bool   `json:"isForParsedProblem"`
	ParsedProblemUrl    string `json:"parsedProblemUrl"`
	FileNum             int    `json:"fileNum"`
	FileMode            string `json:"fileMode"`
	FileName            string `json:"fileName"`
	TestPerFile         int    `json:"testPerFile"`
	SerialFrom          int    `json:"serialFrom"`
	InputDirectoryPath  string `json:"inputDirectoryPath"`
	OutputDirectoryPath string `json:"outputDirectoryPath"`
	GenerationProcess   string `json:"generationProcess"`
	GeneratorScriptPath string `json:"generatorScriptPath"`
	TgenScriptContent   string `json:"tgenScriptContent"`
}
