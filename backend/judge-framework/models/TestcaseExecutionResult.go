package models

type TestcaseExecutionResult struct {
	ExecutionError error  `json:"executionError"`
	Output         string `json:"output"`
	Verdict        string `json:"verdict"`
	ConsumedTime   int64  `json:"consumedTime"`
	ConsumedMemory uint64 `json:"consumedMemory"`
}
