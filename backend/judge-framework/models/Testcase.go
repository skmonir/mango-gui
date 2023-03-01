package models

type Testcase struct {
	Input            string `json:"input"`
	Output           string `json:"output"`
	TimeLimit        int64  `json:"timeLimit"`
	MemoryLimit      uint64 `json:"memoryLimit"`
	SourceBinaryPath string `json:"sourceBinaryPath"`
	InputFilePath    string `json:"inputFilePath"`
	OutputFilePath   string `json:"outputFilePath"`
}
