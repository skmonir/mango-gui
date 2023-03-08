package dto

import "github.com/skmonir/mango-gui/backend/judge-framework/models"

type TestcaseExecutionDetails struct {
	Status                  string                         `json:"status"`
	Testcase                models.Testcase                `json:"testcase"`
	TestcaseExecutionResult models.TestcaseExecutionResult `json:"testcaseExecutionResult"`
}

type ProblemExecutionResult struct {
	CompilationError             string                     `json:"compilationError"`
	TestcaseExecutionDetailsList []TestcaseExecutionDetails `json:"testcaseExecutionDetailsList"`
}
