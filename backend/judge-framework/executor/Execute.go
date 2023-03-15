package executor

import (
	"bytes"
	"context"
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/constants"
	"github.com/skmonir/mango-gui/backend/socket"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/skmonir/mango-gui/backend/judge-framework/dto"
	"github.com/skmonir/mango-gui/backend/judge-framework/models"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
)

type executionResponse struct {
	index          int
	testExecResult dto.TestcaseExecutionResult
}

var executionCompleteChan = make(chan executionResponse)

func getVerdict(execDetails dto.TestcaseExecutionDetails, judgeOutput string) dto.TestcaseExecutionDetails {
	if strings.Contains(execDetails.TestcaseExecutionResult.ExecutionError, "segmentation fault") {
		execDetails.TestcaseExecutionResult.Verdict = "RE"
	} else if (execDetails.TestcaseExecutionResult.ConsumedTime > (execDetails.Testcase.TimeLimit * 1000)) || strings.Contains(execDetails.TestcaseExecutionResult.ExecutionError, "killed") {
		execDetails.TestcaseExecutionResult.Verdict = "TLE"
		execDetails.TestcaseExecutionResult.ConsumedTime = execDetails.Testcase.TimeLimit * 1000
	} else if utils.ConvertMemoryInMb(execDetails.TestcaseExecutionResult.ConsumedMemory) > execDetails.Testcase.MemoryLimit {
		execDetails.TestcaseExecutionResult.Verdict = "MLE"
	} else if execDetails.TestcaseExecutionResult.ExecutionError != "" {
		execDetails.TestcaseExecutionResult.Verdict = "RE"
	} else if judgeOutput == execDetails.TestcaseExecutionResult.Output {
		execDetails.TestcaseExecutionResult.Verdict = "AC"
	} else {
		execDetails.TestcaseExecutionResult.Verdict = "WA"
	}
	return execDetails
}

func executeSourceBinary(index int, testcase models.Testcase) {
	defer utils.PanicRecovery()

	fmt.Println("Running input", testcase.InputFilePath)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(testcase.TimeLimit+1)*time.Second)
	defer cancel()

	inputBuffer := bytes.NewBuffer([]byte(testcase.Input))
	var outputBuffer bytes.Buffer

	cmd := exec.CommandContext(ctx, testcase.ExecutionCommand[0], testcase.ExecutionCommand[1:]...)

	if len(testcase.InputFilePath) == 0 {
		cmd.Stdin = inputBuffer
	} else {
		stdin, _ := cmd.StdinPipe()
		inputFile, _ := os.Open(testcase.InputFilePath)
		go func() {
			io.Copy(stdin, inputFile)
			defer stdin.Close()
		}()
	}

	var stdout io.ReadCloser
	var outputFile *os.File
	if len(testcase.ExecOutputFilePath) == 0 {
		cmd.Stdout = &outputBuffer
	} else {
		stdout, _ = cmd.StdoutPipe()
		outputFile, _ = os.OpenFile(testcase.ExecOutputFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	}

	execResponse := executionResponse{
		index:          index,
		testExecResult: dto.TestcaseExecutionResult{},
	}

	maxMemory := uint64(0)
	completeExecutionFunc := func(err error) {
		if err != nil {
			fmt.Println(err.Error())
			execResponse.testExecResult.ExecutionError = err.Error()
		}
		execResponse.testExecResult.ConsumedMemory = maxMemory / 1024 // converting to KB
		execResponse.testExecResult.ConsumedTime = cmd.ProcessState.UserTime().Milliseconds()
	}

	if err := cmd.Start(); err != nil {
		completeExecutionFunc(err)
		executionCompleteChan <- execResponse
		return
	}

	if len(testcase.ExecOutputFilePath) > 0 {
		io.Copy(io.MultiWriter(outputFile, os.Stdout), stdout)
	}

	pid := int32(cmd.Process.Pid)
	ch := make(chan error)
	go func() {
		ch <- cmd.Wait()
	}()
	running := true
	for running {
		select {
		case err := <-ch:
			completeExecutionFunc(err)
			if err != nil {
				executionCompleteChan <- execResponse
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

	if execResponse.testExecResult.ExecutionError == "" {
		execResponse.testExecResult.Output = utils.TrimIO(outputBuffer.String())
	}

	executionCompleteChan <- execResponse
}

func Execute(execResult dto.ProblemExecutionResult, socketEvent string) dto.ProblemExecutionResult {
	wg := sync.WaitGroup{}
	wg.Add(1)

	judgeOutputs := make([]string, len(execResult.TestcaseExecutionDetailsList))

	go func() {
		for {
			select {
			case testExecutionResult := <-executionCompleteChan:
				mu := sync.Mutex{}
				mu.Lock()
				fmt.Println("exec done", len(execResult.TestcaseExecutionDetailsList), testExecutionResult.index)
				execResult.TestcaseExecutionDetailsList[testExecutionResult.index].TestcaseExecutionResult = testExecutionResult.testExecResult
				execResult.TestcaseExecutionDetailsList[testExecutionResult.index].Status = "success"
				execResult.TestcaseExecutionDetailsList[testExecutionResult.index] = getVerdict(
					execResult.TestcaseExecutionDetailsList[testExecutionResult.index],
					judgeOutputs[testExecutionResult.index])

				if socketEvent == "test_exec_result_event" {
					execResult.TestcaseExecutionDetailsList[testExecutionResult.index].TestcaseExecutionResult.Output = utils.ResizeIOContentForUI(
						strings.NewReader(execResult.TestcaseExecutionDetailsList[testExecutionResult.index].TestcaseExecutionResult.Output), constants.IO_MAX_ROW_FOR_UI, constants.IO_MAX_COL_FOR_UI)
				}
				socket.PublishExecutionResult(execResult, socketEvent)

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
		judgeOutputs[i] = utils.TrimIO(utils.ReadFileContent(
			execResult.TestcaseExecutionDetailsList[i].Testcase.OutputFilePath, constants.IO_MAX_ROW_FOR_TEST, constants.IO_MAX_COL_FOR_TEST))
		go executeSourceBinary(i, execResult.TestcaseExecutionDetailsList[i].Testcase)
	}
	wg.Wait()
	return execResult
}
