package executor

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/shirou/gopsutil/process"
	"github.com/skmonir/mango-ui/backend/judge-framework/dto"
	"github.com/skmonir/mango-ui/backend/judge-framework/models"
	"github.com/skmonir/mango-ui/backend/judge-framework/utils"
)

type testExecutionResult struct {
	index          int
	testExecResult models.TestcaseExecutionResult
}

var executionCompleteChan = make(chan testExecutionResult)

func getVerdict(testcase *models.Testcase) {
	// errorMsg := ""
	// if testcase.ExecResult.Error != nil {
	// 	errorMsg = testcase.ExecResult.Error.Error()
	// }

	// if strings.Contains(errorMsg, "segmentation fault") {
	// 	testcase.ExecResult.Verdict = "RE"
	// } else if (testcase.ExecResult.Runtime > (testcase.TimeLimit * 1000)) || strings.Contains(errorMsg, "killed") {
	// 	testcase.ExecResult.Verdict = "TLE"
	// 	testcase.ExecResult.Runtime = testcase.TimeLimit * 1000
	// } else if utils.ConvertMemoryInMb(testcase.ExecResult.Memory) > testcase.MemoryLimit {
	// 	testcase.ExecResult.Verdict = "MLE"
	// } else if testcase.ExecResult.Status == "error" {
	// 	testcase.ExecResult.Verdict = "RE"
	// } else if testcase.Output == testcase.ExecResult.Output {
	// 	testcase.ExecResult.Verdict = "OK"
	// } else {
	// 	testcase.ExecResult.Verdict = "WA"
	// }
}

func executeSourceBinary(index int, testcase models.Testcase) {
	utils.PanicRecovery()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(testcase.TimeLimit)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, testcase.SourceBinaryPath)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	inputFile, _ := os.Open(testcase.InputFilePath)
	io.Copy(stdin, inputFile)
	stdin.Close()

	outputFile, _ := os.OpenFile(testcase.UserOutputFilePath, os.O_WRONLY, os.ModeAppend)

	testExecRes := testExecutionResult{
		index:          index,
		testExecResult: models.TestcaseExecutionResult{},
	}

	maxMemory := uint64(0)
	completeExecution := func(err error) {
		if err != nil {
			fmt.Println(err.Error())
			testExecRes.testExecResult.ExecutionError = err.Error()
		}
		testExecRes.testExecResult.ConsumedMemory = maxMemory / 1024 // converting to KB
		testExecRes.testExecResult.ConsumedTime = cmd.ProcessState.UserTime().Milliseconds()
	}

	if err := cmd.Start(); err != nil {
		completeExecution(err)
		executionCompleteChan <- testExecRes
		return
	}

	io.Copy(io.MultiWriter(outputFile, os.Stdout), stdout)

	pid := int32(cmd.Process.Pid)
	ch := make(chan error)
	go func() {
		ch <- cmd.Wait()
	}()
	running := true
	for running {
		select {
		case err := <-ch:
			completeExecution(err)
			if err != nil {
				executionCompleteChan <- testExecRes
				return
			}
			running = false
		default:
			p, err := process.NewProcess(pid)
			if err == nil {
				m, err := p.MemoryInfo()
				if err == nil && m.RSS > maxMemory {
					maxMemory = m.RSS
				}
			}
		}
	}

	executionCompleteChan <- testExecRes
}

func Execute(execResult dto.ProblemExecutionResult) dto.ProblemExecutionResult {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			select {
			case testExecutionResult := <-executionCompleteChan:
				mu := sync.Mutex{}
				mu.Lock()
				fmt.Println("exec done", len(execResult.TestcaseExecutionDetailsList), testExecutionResult.index)
				execResult.TestcaseExecutionDetailsList[testExecutionResult.index].TestcaseExecutionResult = testExecutionResult.testExecResult
				execResult.TestcaseExecutionDetailsList[testExecutionResult.index].Status = "success"
				// getVerdict(&execResult.TestcaseExecutionDetailsList[testExecutionResult.index])
				// socket.PublishTestMessage(prob)
				untestedProblemExists := false
				for i := 0; i < len(execResult.TestcaseExecutionDetailsList); i++ {
					untestedProblemExists = untestedProblemExists || execResult.TestcaseExecutionDetailsList[i].Status == "running"
				}
				if !untestedProblemExists {
					wg.Done()
					fmt.Println("All cases are tested")
					return
				}
				mu.Unlock()
			}
		}
	}()

	for i := 0; i < len(execResult.TestcaseExecutionDetailsList); i++ {
		go executeSourceBinary(i, execResult.TestcaseExecutionDetailsList[i].Testcase)
	}
	wg.Wait()
	return execResult
}
