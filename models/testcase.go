package models

type Testcase struct {
	Input       string
	Output      string
	TimeLimit   int64
	MemoryLimit uint64
}
