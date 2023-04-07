package socket

import (
	"encoding/json"
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/config"
	"github.com/skmonir/mango-gui/backend/judge-framework/dto"

	"github.com/skmonir/mango-gui/backend/judge-framework/models"
)

type Message struct {
	Key     string `json:"key"`
	Content string `json:"content"`
}

type ParseProblemPublishDetails struct {
	Status   string         `json:"status"`
	Metadata models.Problem `json:"metadata"`
}

func PublishParseMessage(parsedProblemList []models.Problem) {
	fmt.Println("publishing parse result.....")

	parseResponseListJson, err := json.Marshal(parsedProblemList)
	if err != nil {
		fmt.Println(err)
		return
	}

	broadcastMessage(Message{
		Key:     "parse_problems_event",
		Content: string(parseResponseListJson),
	})
}

func PublishStatusMessage(topic string, message string, messageType string) {
	fmt.Println("publishing test status message.....")

	messageContent := struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	}{
		Message: message,
		Type:    messageType,
	}

	messageContentJson, err := json.Marshal(messageContent)
	if err != nil {
		return
	}

	broadcastMessage(Message{
		Key:     topic,
		Content: string(messageContentJson),
	})
}

func PublishExecutionResult(execResult dto.ProblemExecutionResult, socketEvent string) {
	if len(socketEvent) == 0 {
		return
	}

	if socketEvent == "test_exec_result_event" {
		PublishPreviousRunStatus(execResult)
	}

	fmt.Println("publishing execution result.....")

	execResultJson, err := json.Marshal(execResult)
	if err != nil {
		fmt.Println(err)
		return
	}

	broadcastMessage(Message{
		Key:     socketEvent,
		Content: string(execResultJson),
	})
}

func PublishPreviousRunStatus(execResult dto.ProblemExecutionResult) {
	if execResult.CompilationError != "" {
		PublishStatusMessage("test_status", "Compilation error!", "error")
		return
	}
	testcaseExecutionDetailsList := execResult.TestcaseExecutionDetailsList
	totalPassed, totalFailed, totalTests := 0, 0, len(testcaseExecutionDetailsList)
	for _, execDetails := range testcaseExecutionDetailsList {
		if execDetails.Status == "success" {
			if execDetails.TestcaseExecutionResult.Verdict == "AC" {
				totalPassed++
			} else {
				totalFailed++
			}
		}
	}

	testStatus := fmt.Sprintf("{\"total\": %v,\"passed\": %v,\"failed\": %v}", totalTests, totalPassed, totalFailed)
	PublishStatusMessage("test_status", testStatus, "test_stat")
}

func PublishAppConfig(conf config.JudgeConfig) {
	fmt.Println("publishing app config.....")

	confJson, err := json.Marshal(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	broadcastMessage(Message{
		Key:     "app_conf_get_event",
		Content: string(confJson),
	})
}
