package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/auth"
	"github.com/grvbrk/async0_server/internal/middlewares"
	"github.com/grvbrk/async0_server/internal/models"
	"github.com/grvbrk/async0_server/internal/store"
	"github.com/grvbrk/async0_server/internal/utils"
)

type SubmissionBody struct {
	Code string `json:"code"`
}

type Judge0Submission struct {
	LanguageID     int    `json:"language_id"`
	SourceCode     string `json:"source_code"`
	Stdin          string `json:"stdin,omitempty"`
	ExpectedOutput string `json:"expected_output,omitempty"`
}

type Judge0BatchRequest struct {
	Submissions []Judge0Submission `json:"submissions"`
}

type Status struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Judge0Result struct {
	Stdout        *string `json:"stdout"`
	Time          *string `json:"time"`
	Memory        *int    `json:"memory"`
	Stderr        *string `json:"stderr"`
	Token         string  `json:"token"`
	CompileOutput *string `json:"compile_output"`
	Message       *string `json:"message"`
	Status        Status  `json:"status"`
}

type RunSubmissionResponse struct {
	StatusID     string `json:"status_id"`
	StatusDesc   string `json:"status_description"`
	Result       string `json:"result,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type SubmissionHandler struct {
	SubmissionStore store.SubmissionStore
	TestcaseStore   store.TestcaseStore
	Logger          *log.Logger
	Oauth           *auth.GoogleOauth
}

func NewSubmissionHandler(submissionStore store.SubmissionStore, testcaseStore store.TestcaseStore, logger *log.Logger, oauth *auth.GoogleOauth) *SubmissionHandler {
	return &SubmissionHandler{
		SubmissionStore: submissionStore,
		TestcaseStore:   testcaseStore,
		Logger:          logger,
		Oauth:           oauth,
	}
}

func (ph *SubmissionHandler) HandlerGetSubmissionsByProblemID(w http.ResponseWriter, r *http.Request) {

	user, ok := middlewares.GetUserFromContext(r)
	if !ok {
		ph.Logger.Println("No user found in context. Sending an empty array of submissions")
		utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": []models.Submission{}})
		return
	}

	problemID, err := uuid.Parse(chi.URLParam(r, "problemID"))
	if err != nil {
		ph.Logger.Println("Error parsing problem id", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	submissions, err := ph.SubmissionStore.GetSubmissionsByProblemID(user.ID, problemID)
	if err != nil {
		ph.Logger.Println("Error getting submissions by problem id", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": submissions})

}

func (ph *SubmissionHandler) HandlerSubmitSubmission(w http.ResponseWriter, r *http.Request) {

	user, ok := middlewares.GetUserFromContext(r)
	if !ok {
		ph.Logger.Println("No user found in context")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Not Authorized"})
		return
	}

	id := chi.URLParam(r, "id")
	problemID, err := uuid.Parse(id)
	if err != nil {
		ph.Logger.Println("Error parsing problem id", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	var body SubmissionBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ph.Logger.Println("Error decoding submission body", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	testcases, err := ph.TestcaseStore.GetTestcasesByProblemID(problemID)
	if err != nil {
		ph.Logger.Println("Error getting testcases", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	var submissions []Judge0Submission

	for _, testcase := range testcases {
		inputTemplate := `
			%s
			try {
				const result = %s;
				console.log(JSON.stringify(result));
			} catch (error) {
				console.error('Runtime Error:', error.message);
				process.exit(1);
			}
		`

		sourceCode := fmt.Sprintf(inputTemplate, body.Code, testcase.Input)
		expectedOutput := strings.ReplaceAll(strings.TrimSpace(testcase.Output), " ", "")

		submission := Judge0Submission{
			LanguageID:     63,
			SourceCode:     sourceCode,
			ExpectedOutput: expectedOutput,
		}

		submissions = append(submissions, submission)
	}

	batchRequest := Judge0BatchRequest{
		Submissions: submissions,
	}

	jsonBody, err := json.Marshal(batchRequest)
	if err != nil {
		ph.Logger.Println("Error marshalling batch request", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/submissions/batch?base64_encoded=false&wait=false", os.Getenv("JUDGE0_URL")), "application/json", bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		ph.Logger.Println("Error submitting batch request", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	defer resp.Body.Close()

	var batchResponse []map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&batchResponse); err != nil {
		ph.Logger.Println("Error decoding batch response", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	tokens := make([]string, len(batchResponse))
	for i, submission := range batchResponse {
		tokens[i] = submission["token"]
	}

	results, err := pollJudge0BatchResults(tokens)
	if err != nil {
		ph.Logger.Println("Error polling judge0 results", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	result := formatMultipleJudge0Results(results, testcases)

	err = ph.SubmissionStore.CreateSubmission(user.ID, problemID, body.Code, result)
	if err != nil {
		ph.Logger.Println("Error creating submission", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": result})

}

func pollJudge0BatchResults(tokens []string) ([]Judge0Result, error) {
	tokensParam := strings.Join(tokens, ",")
	maxRetries := 30
	baseDelay := 500 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := http.Get(fmt.Sprintf("%s/submissions/batch?tokens=%s&base64_encoded=false", os.Getenv("JUDGE0_URL"), tokensParam))
		if err != nil {
			return nil, fmt.Errorf("error sending get request: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		var judge0Response struct {
			Submissions []Judge0Result `json:"submissions"`
		}

		err = json.Unmarshal(body, &judge0Response)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling response body: %w, raw response: %s", err, string(body))
		}

		results := judge0Response.Submissions

		if len(results) != len(tokens) {
			return nil, fmt.Errorf("expected %d results, got %d", len(tokens), len(results))
		}

		allComplete := true
		for _, result := range results {
			if result.Status.ID == 1 || result.Status.ID == 2 {
				allComplete = false
				break
			}
		}

		if allComplete {
			return results, nil
		}

		delay := time.Duration(float64(baseDelay) * (1.5 * float64(attempt)))
		if delay > 5*time.Second {
			delay = 5 * time.Second
		}
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("polling timeout exceeded")
}

func formatMultipleJudge0Results(results []Judge0Result, testCases []models.Testcase) models.SubmitSubmissionResponse {
	passedTests := 0
	formattedResults := make([]models.TestcaseResult, len(results))

	statusDescriptions := map[int]string{
		1:  "In Queue",
		2:  "Processing",
		3:  "Accepted",
		4:  "Wrong Answer",
		5:  "Time Limit Exceeded",
		6:  "Compilation Error",
		7:  "Runtime Error (SIGSEGV)",
		8:  "Runtime Error (SIGXFSZ)",
		9:  "Runtime Error (SIGFPE)",
		10: "Runtime Error (SIGABRT)",
		11: "Runtime Error (NZEC)",
		12: "Runtime Error (Other)",
		13: "Internal Error",
		14: "Exec Format Error",
	}

	for i, result := range results {
		statusDesc := statusDescriptions[result.Status.ID]
		if statusDesc == "" {
			statusDesc = fmt.Sprintf("Unknown Status (%d)", result.Status.ID)
		}

		actualOutput := ""
		if result.Stdout != nil {
			actualOutput = strings.TrimSpace(*result.Stdout)
		}

		normalize := func(s string) string {
			s = strings.TrimSpace(s)
			s = strings.ReplaceAll(s, " ", "")
			return s
		}

		actualOutputNorm := normalize(actualOutput)
		expectedOutputNorm := normalize(testCases[i].Output)

		passed := result.Status.ID == 3 && actualOutputNorm == expectedOutputNorm
		if passed {
			passedTests++
		}

		tcTime := ""
		if result.Time != nil {
			tcTime = *result.Time
		}

		tcMemory := 0
		if result.Memory != nil {
			tcMemory = *result.Memory
		}

		formattedResults[i] = models.TestcaseResult{
			TCPass:           passed,
			TCStatusID:       result.Status.ID,
			TCStatus:         statusDesc,
			TCTime:           tcTime,
			TCMemory:         tcMemory,
			TCOutput:         actualOutput,
			TCExpectedOutput: expectedOutputNorm,
		}
	}

	overallStatusID := 4
	overallStatus := models.StatusRE
	if passedTests == len(results) {
		overallStatusID = 3
		overallStatus = models.StatusAC
	}

	return models.SubmitSubmissionResponse{
		OverallStatusID:  overallStatusID,
		OverallStatus:    overallStatus,
		PassedTestcases:  passedTests,
		TotalTestcases:   len(results),
		TestcasesResults: formattedResults,
	}
}

func (ph *SubmissionHandler) HandlerRunSubmission(w http.ResponseWriter, r *http.Request) {

	var body SubmissionBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ph.Logger.Println("Error decoding submission body", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	var submission = Judge0Submission{
		LanguageID: 63,
		SourceCode: body.Code,
	}

	jsonBody, err := json.Marshal(submission)
	if err != nil {
		ph.Logger.Println("Error marshalling submission", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/submissions?base64_encoded=false&wait=false", os.Getenv("JUDGE0_URL")), "application/json", bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		ph.Logger.Println("Error submitting request", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	defer resp.Body.Close()

	type Response struct {
		Token string `json:"token"`
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		ph.Logger.Println("Error decoding response", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	pollResult, err := pollJudge0SingleResult(response.Token)
	if err != nil {
		ph.Logger.Println("Error polling judge0 results", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": formatJudge0Result(pollResult)})

}

func pollJudge0SingleResult(token string) (Judge0Result, error) {
	maxRetries := 30
	baseDelay := 500 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := http.Get(fmt.Sprintf("%s/submissions/%s?base64_encoded=false", os.Getenv("JUDGE0_URL"), token))
		if err != nil {
			return Judge0Result{}, fmt.Errorf("error sending get request: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return Judge0Result{}, fmt.Errorf("error reading response body: %w", err)
		}

		var result Judge0Result
		err = json.Unmarshal(body, &result)
		if err != nil {
			return Judge0Result{}, fmt.Errorf("error unmarshalling response body: %w, raw response: %s", err, string(body))
		}

		if result.Status.ID != 1 && result.Status.ID != 2 {
			return result, nil
		}

		delay := time.Duration(float64(baseDelay) * (1.5 * float64(attempt)))
		if delay > 5*time.Second {
			delay = 5 * time.Second
		}
		time.Sleep(delay)
	}

	return Judge0Result{}, fmt.Errorf("polling timeout exceeded")
}

func formatJudge0Result(result Judge0Result) RunSubmissionResponse {

	var statusMap = map[int]string{
		1: "In Queue",
		2: "Processing",
		3: "Accepted",
		4: "Wrong Answer",
		5: "Time Limit Exceeded",
		6: "Compilation Error",
	}

	desc, ok := statusMap[result.Status.ID]
	if !ok {
		desc = "Execution Error"
	}

	resp := RunSubmissionResponse{
		StatusID:   fmt.Sprintf("%d", result.Status.ID),
		StatusDesc: desc,
	}

	if result.Stdout != nil && *result.Stdout != "" {
		resp.Result = *result.Stdout
	}

	if result.Stderr != nil && *result.Stderr != "" {
		resp.ErrorMessage = *result.Stderr
	} else if result.CompileOutput != nil && *result.CompileOutput != "" {
		resp.ErrorMessage = *result.CompileOutput
	} else if result.Message != nil && *result.Message != "" {
		resp.ErrorMessage = *result.Message
	}

	return resp
}
