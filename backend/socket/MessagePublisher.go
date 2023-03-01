package socket

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("publishing test result.....")

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
