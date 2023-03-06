package cacheServices

import (
	"encoding/json"
	"fmt"
	"github.com/skmonir/mango-ui/backend/judge-framework/models"

	"github.com/skmonir/mango-ui/backend/judge-framework/cache"
	"github.com/skmonir/mango-ui/backend/judge-framework/dto"
)

func GetExecutionResult(platform string, cid string, label string) (bool, dto.ProblemExecutionResult) {
	fmt.Println("Fetching execution result from cache...")

	key := fmt.Sprintf("ProblemExecutionResult:%v.%v.%v", platform, cid, label)
	execResultObjStr := cache.GetGlobalCache().Get(key)

	var execResult dto.ProblemExecutionResult
	if execResultObjStr != "" {
		if err := json.Unmarshal([]byte(execResultObjStr), &execResult); err != nil {
			fmt.Println(err)
			return false, execResult
		}
		fmt.Println("Fetched execution result from cache")
	} else {
		return false, execResult
	}
	return true, execResult
}

func SaveExecutionResult(platform string, cid string, label string, execResult dto.ProblemExecutionResult) {
	fmt.Println("Saving execution result into cache...")

	data, err := json.MarshalIndent(execResult, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Saved execution result into cache")

	key := fmt.Sprintf("ProblemExecutionResult:%v.%v.%v", platform, cid, label)
	cache.GetGlobalCache().Set(key, string(data))
}

func UpdateTestcaseIntoExecutionResult(platform string, cid string, label string, testcase models.Testcase) {
	fmt.Println("Saving execution result into cache...")

	ok, execResult := GetExecutionResult(platform, cid, label)
	if !ok {
		return
	}
	isFound := false
	for i := 0; i < len(execResult.TestcaseExecutionDetailsList); i++ {
		if execResult.TestcaseExecutionDetailsList[i].Testcase.InputFilePath == testcase.OutputFilePath {
			isFound = true
			execResult.TestcaseExecutionDetailsList[i].Testcase = testcase
		}
	}
	if !isFound {
		execResult.TestcaseExecutionDetailsList = append(execResult.TestcaseExecutionDetailsList, dto.TestcaseExecutionDetails{
			Status:   "none",
			Testcase: testcase,
		})
	}
	SaveExecutionResult(platform, cid, label, execResult)
}
