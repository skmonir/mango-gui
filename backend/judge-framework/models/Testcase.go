package models

type Testcase struct {
	Input              string   `json:"input"`
	Output             string   `json:"output"`
	TimeLimit          int64    `json:"timeLimit"`
	MemoryLimit        uint64   `json:"memoryLimit"`
	InputFilePath      string   `json:"inputFilePath"`
	OutputFilePath     string   `json:"outputFilePath"`
	ExecOutputFilePath string   `json:"execOutputFilePath"`
	SourceBinaryPath   string   `json:"sourceBinaryPath"`
	ExecutionCommand   []string `json:"executionCommand"`
}
