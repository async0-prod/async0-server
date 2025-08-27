package models

type CardAnalytics struct {
	Name              string `json:"name"`
	TotalQuestions    int    `json:"total_questions"`
	TotalSolved       int    `json:"total_solved"`
	TotalSolutions    int    `json:"total_solutions"`
	TotalUserAttempts int    `json:"total_user_attempts"`

	TotalEasyQ   int `json:"total_easy_q"`
	TotalMediumQ int `json:"total_medium_q"`
	TotalHardQ   int `json:"total_hard_q"`

	TotalEasySolved   int `json:"total_easy_solved"`
	TotalMediumSolved int `json:"total_medium_solved"`
	TotalHardSolved   int `json:"total_hard_solved"`
}
