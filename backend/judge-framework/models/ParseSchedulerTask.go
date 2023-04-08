package models

import "time"

type ParseSchedulerTask struct {
	Id        string    `json:"id"`
	Url       string    `json:"url"`
	StartTime time.Time `json:"startTime"`
	Stage     string    `json:"stage"`
}
