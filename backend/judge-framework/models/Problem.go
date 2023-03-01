package models

type Problem struct {
	Platform    string `json:"platform"`
	ContestId   string `json:"contestId"`
	Label       string `json:"label"`
	Name        string `json:"name"`
	TimeLimit   int64  `json:"timeLimit"`
	MemoryLimit uint64 `json:"memoryLimit"`
	Url         string `json:"url"`
	Status      string `json:"status"`
}
