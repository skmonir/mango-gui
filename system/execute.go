package system

import (
	"bytes"
	"context"
	"os/exec"
	"time"

	"github.com/shirou/gopsutil/process"
	appContext "github.com/skmonir/mango-gui/context"
	"github.com/skmonir/mango-gui/models"
)

func getExecutionCommand(config appContext.AppConfig, problemId string) string {
	command := config.GetSourceFilePathWithoutExt(problemId)
	return command
}

func ExecuteSourceBinary(config appContext.AppConfig, testcase models.Testcase, problemId string) models.ExecutionResult {
	command := getExecutionCommand(config, problemId)

	var response models.ExecutionResult
	response.Test = testcase

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(testcase.TimeLimit)*time.Second)
	defer cancel()

	input_buffer := bytes.NewBuffer([]byte(testcase.Input))
	var output_buffer bytes.Buffer

	cmd := exec.CommandContext(ctx, command)
	// cmd.Stderr = os.Stderr
	cmd.Stdin = input_buffer
	cmd.Stdout = &output_buffer

	maxMemory := uint64(0)

	completeExecution := func(err error) {
		response.Status = (err == nil)
		response.Error = err
		response.Memory = maxMemory
		response.Runtime = cmd.ProcessState.UserTime().Milliseconds()
	}

	if err := cmd.Start(); err != nil {
		completeExecution(err)
		return response
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
			completeExecution(err)
			if err != nil {
				return response
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

	if !response.Status {
		return response
	}

	response.Output = TrimIO(output_buffer.String())

	return response
}
