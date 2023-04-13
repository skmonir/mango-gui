package models

type JudgeAccountInfo struct {
	Handle           string `json:"handle"`
	HandleOrEmail    string `json:"handleOrEmail"`
	Password         string `json:"password"`
	SubmissionLangId string `json:"submissionLangId"`
}
