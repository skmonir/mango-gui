package services

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/skmonir/mango-ui/backend/judge-framework/cache"
//	"github.com/skmonir/mango-ui/backend/judge-framework/dto"
//	"github.com/skmonir/mango-ui/backend/judge-framework/executor"
//	"github.com/skmonir/mango-ui/backend/judge-framework/utils"
//	"github.com/skmonir/mango-ui/backend/socket"
//	"strings"
//	"time"
//)
//
//func RunTest(platform string, cid string, label string) dto.ProblemExecutionResult {
//	//prob := fileService.LoadProblem(platform, cid, label)
//
//	socket.PublishStatusMessage("test_status", "Compiling source code", "info")
//	err := executor.Compile(platform, cid, label)
//	if err != "" {
//		socket.PublishStatusMessage("test_status", "Compilation error!", "error")
//		execResponse.CompilationError = err
//		return execResponse
//	}
//	socket.PublishStatusMessage("test_status", "Compilation successful!", "success")
//
//	time.Sleep(500 * time.Millisecond)
//
//	judgeConfig := cache.GetJudgeConfigFromCache()
//	folderPath := fmt.Sprintf("%v/%v/%v/source", strings.TrimRight(judgeConfig.WorkspaceDirectory, "/"), platform, cid)
//	filePathWithoutExt := folderPath + "/" + label
//
//	if !utils.IsFileExist(filePathWithoutExt) {
//		socket.PublishStatusMessage("test_status", "Binary fileService not found!", "error")
//		return prob
//	}
//
//	prob = executor.Execute(prob)
//
//	totalPassed, totalTests := 0, len(prob.Testcases)
//	for _, test := range prob.Testcases {
//		if test.ExecResult.Verdict == "OK" {
//			totalPassed++
//		}
//	}
//	prob.TestStatus = fmt.Sprintf("%v/%v Tests Passed", totalPassed, totalTests)
//	prob.IsPassed = totalPassed == totalTests
//	socket.PublishTestMessage(prob)
//
//	return prob
//}
