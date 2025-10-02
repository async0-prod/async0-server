package models

import (
	"time"

	"github.com/google/uuid"
)

type SubmissionStatus string

const (
	StatusAC      SubmissionStatus = "AC"
	StatusWA      SubmissionStatus = "WA"
	StatusTLE     SubmissionStatus = "TLE"
	StatusMLE     SubmissionStatus = "MLE"
	StatusRE      SubmissionStatus = "RE"
	StatusCE      SubmissionStatus = "CE"
	StatusPE      SubmissionStatus = "PE"
	StatusPending SubmissionStatus = "PENDING"
	StatusRunning SubmissionStatus = "RUNNING"
)

type Submission struct {
	ID              uuid.UUID        `json:"id"`
	UserID          uuid.UUID        `json:"user_id"`
	ProblemID       uuid.UUID        `json:"problem_id"`
	Code            string           `json:"code"`
	Status          SubmissionStatus `json:"status"`
	Runtime         *int             `json:"runtime"`
	MemoryUsed      *int             `json:"memory_used"`
	TotalTestcases  *int             `json:"total_testcases"`
	PassedTestcases *int             `json:"passed_testcases"`
	FailedTestcases *int             `json:"failed_testcases"`
	CreatedAt       time.Time        `json:"created_at"`
}

type TestcaseResult struct {
	TCPass           bool   `json:"tc_pass"`
	TCStatusID       int    `json:"tc_status_id"`
	TCStatus         string `json:"tc_status"`
	TCTime           string `json:"tc_time"`
	TCMemory         int    `json:"tc_memory"`
	TCOutput         string `json:"tc_output"`
	TCExpectedOutput string `json:"tc_expected_output"`
}

type SubmitSubmissionResponse struct {
	OverallStatusID  int              `json:"overall_status_id"`
	OverallStatus    SubmissionStatus `json:"overall_status"`
	PassedTestcases  int              `json:"passed_testcases"`
	TotalTestcases   int              `json:"total_testcases"`
	TestcasesResults []TestcaseResult `json:"testcases_results"`
}
