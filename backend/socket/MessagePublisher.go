package socket

import (
	"encoding/json"
	"fmt"
	"github.com/skmonir/mango-ui/backend/judge-framework/dto"

	"github.com/skmonir/mango-ui/backend/judge-framework/models"
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

func PublishExecutionResult(execResult dto.ProblemExecutionResult) {
	PublishPreviousRunStatus(execResult)

	fmt.Println("publishing execution result.....")

	execResultJson, err := json.Marshal(execResult)
	if err != nil {
		fmt.Println(err)
		return
	}

	broadcastMessage(Message{
		Key:     "test_exec_result_event",
		Content: string(execResultJson),
	})
}

func PublishPreviousRunStatus(execResult dto.ProblemExecutionResult) {
	testcaseExecutionDetailsList := execResult.TestcaseExecutionDetailsList
	totalPassed, totalTests, isExecutedOnce := 0, len(testcaseExecutionDetailsList), false
	for _, execDetails := range testcaseExecutionDetailsList {
		isExecutedOnce = isExecutedOnce || (execDetails.Status != "none")
		if execDetails.TestcaseExecutionResult.Verdict == "AC" {
			totalPassed++
		}
	}
	if isExecutedOnce {
		testStatus := fmt.Sprintf("%v/%v Tests Passed", totalPassed, totalTests)
		if totalPassed == totalTests {
			PublishStatusMessage("test_status", testStatus, "success")
		} else {
			PublishStatusMessage("test_status", testStatus, "error")
		}
	} else if execResult.CompilationError != "" {
		PublishStatusMessage("test_status", "Compilation error!", "error")
	}
}
