package executor

//
//import (
//	"bytes"
//	"context"
//	"fmt"
//	"github.com/shirou/gopsutil/process"
//	"github.com/skmonir/mango-ui/backend/judge-framework/cache"
//	"github.com/skmonir/mango-ui/backend/judge-framework/dto"
//	"github.com/skmonir/mango-ui/backend/judge-framework/models"
//	"github.com/skmonir/mango-ui/backend/judge-framework/utils"
//	"github.com/skmonir/mango-ui/backend/socket"
//	"os/exec"
//	"strings"
//	"sync"
//	"time"
//)
//
//type executionRequest struct {
//	index    int
//	testcase models.Testcase
//}
//
//var executionCompleteChan = make(chan executionRequest)
//
//func getVerdict(testcase *models.Testcase) {
//	errorMsg := ""
//	if testcase.ExecResult.Error != nil {
//		errorMsg = testcase.ExecResult.Error.Error()
//	}
//
//	if strings.Contains(errorMsg, "segmentation fault") {
//		testcase.ExecResult.Verdict = "RE"
//	} else if (testcase.ExecResult.Runtime > (testcase.TimeLimit * 1000)) || strings.Contains(errorMsg, "killed") {
//		testcase.ExecResult.Verdict = "TLE"
//		testcase.ExecResult.Runtime = testcase.TimeLimit * 1000
//	} else if utils.ConvertMemoryInMb(testcase.ExecResult.Memory) > testcase.MemoryLimit {
//		testcase.ExecResult.Verdict = "MLE"
//	} else if testcase.ExecResult.Status == "error" {
//		testcase.ExecResult.Verdict = "RE"
//	} else if testcase.Output == testcase.ExecResult.Output {
//		testcase.ExecResult.Verdict = "OK"
//	} else {
//		testcase.ExecResult.Verdict = "WA"
//	}
//}
//
//func executeSourceBinary(request executionRequest, command string) {
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(request.testcase.TimeLimit)*time.Second)
//	defer cancel()
//
//	inputBuffer := bytes.NewBuffer([]byte(request.testcase.Input))
//	var outputBuffer bytes.Buffer
//
//	cmd := exec.CommandContext(ctx, command)
//	// cmd.Stderr = os.Stderr
//	cmd.Stdin = inputBuffer
//	cmd.Stdout = &outputBuffer
//
//	maxMemory := uint64(0)
//
//	completeExecution := func(err error) {
//		if err != nil {
//			request.testcase.ExecResult.Status = "error"
//		} else {
//			request.testcase.ExecResult.Status = "success"
//		}
//		request.testcase.ExecResult.Message = "done"
//		request.testcase.ExecResult.Error = err
//		request.testcase.ExecResult.Memory = maxMemory / 1024 // converting to KB
//		request.testcase.ExecResult.Runtime = cmd.ProcessState.UserTime().Milliseconds()
//	}
//
//	// if err := cmd.Run(); err != nil {
//	// 	response.Status = false
//	// 	response.Error = err
//	// } else {
//	// 	response.Status = true
//	// }
//
//	// response.ExitCode = cmd.ProcessState.ExitCode()
//	// response.Memory = uint64(cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss)
//	// response.Runtime = cmd.ProcessState.UserTime().Milliseconds()
//
//	if err := cmd.Start(); err != nil {
//		completeExecution(err)
//		executionCompleteChan <- request
//		return
//	}
//
//	pid := int32(cmd.Process.Pid)
//	ch := make(chan error)
//	go func() {
//		ch <- cmd.Wait()
//	}()
//	running := true
//	for running {
//		select {
//		case err := <-ch:
//			completeExecution(err)
//			if err != nil {
//				executionCompleteChan <- request
//				return
//			}
//			running = false
//		default:
//			p, err := process.NewProcess(pid)
//			if err == nil {
//				m, err := p.MemoryInfo()
//				if err == nil && m.RSS > maxMemory {
//					maxMemory = m.RSS
//				}
//			}
//		}
//	}
//
//	if request.testcase.ExecResult.Status == "error" {
//		executionCompleteChan <- request
//		return
//	}
//
//	request.testcase.ExecResult.Status = "success"
//	request.testcase.ExecResult.Message = "done"
//	request.testcase.ExecResult.Output = utils.TrimIO(outputBuffer.String())
//
//	executionCompleteChan <- request
//}
//
//func Execute(prob models.Problem) []dto.ProblemExecutionResult {
//	judgeConfig := cache.GetJudgeConfigFromCache()
//
//	folderPath := fmt.Sprintf("%v/%v/%v/source", strings.TrimRight(judgeConfig.WorkspaceDirectory, "/"), prob.Platform, prob.ContestId)
//	filePathWithoutExt := folderPath + "/" + prob.Label
//
//	if !utils.IsFileExist(filePathWithoutExt) {
//		prob.TestStatus = "Binary fileService not found!"
//		return prob
//	}
//
//	prob.TestStatus = "Running testcases"
//	socket.PublishTestMessage(prob)
//
//	wg := sync.WaitGroup{}
//	wg.Add(1)
//
//	go func() {
//		for {
//			select {
//			case execResp := <-executionCompleteChan:
//				mu := sync.Mutex{}
//				mu.Lock()
//				prob.Testcases[execResp.index] = execResp.testcase
//				getVerdict(&prob.Testcases[execResp.index])
//				socket.PublishTestMessage(prob)
//				untestedProblemExists := false
//				for i := 0; i < len(prob.Testcases); i++ {
//					untestedProblemExists = untestedProblemExists || prob.Testcases[i].ExecResult.Message == "running"
//				}
//				if !untestedProblemExists {
//					wg.Done()
//					fmt.Println("All cases are tested")
//					return
//				}
//				mu.Unlock()
//			}
//		}
//	}()
//	for i := 0; i < len(prob.Testcases); i++ {
//		prob.Testcases[i].ExecResult.Message = "running"
//	}
//
//	socket.PublishTestMessage(prob)
//
//	for i, testcase := range prob.Testcases {
//		go executeSourceBinary(
//			executionRequest{
//				index:    i,
//				testcase: testcase,
//			},
//			filePathWithoutExt,
//		)
//	}
//	wg.Wait()
//	return prob
//}
