package admin

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/auth"
	"github.com/grvbrk/async0_server/internal/models"
	"github.com/grvbrk/async0_server/internal/store/admin"
	"github.com/grvbrk/async0_server/internal/utils"
)

type AdminProblemHandler struct {
	AdminProblemStore admin.AdminProblemStore
	Logger            *log.Logger
	Oauth             *auth.AdminGoogleOauth
}

func NewAdminProblemHandler(adminProblemStore admin.AdminProblemStore, logger *log.Logger, oauth *auth.AdminGoogleOauth) *AdminProblemHandler {
	return &AdminProblemHandler{
		AdminProblemStore: adminProblemStore,
		Logger:            logger,
		Oauth:             oauth,
	}
}

func (ap *AdminProblemHandler) HandlerGetAllProblems(w http.ResponseWriter, r *http.Request) {

	problems, err := ap.AdminProblemStore.GetAllProblems()
	if err != nil {
		ap.Logger.Println("Error getting all problems from store", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": problems})
}

type TestCaseBody struct {
	UI       string `json:"ui"`
	Input    string `json:"input"`
	Output   string `json:"output"`
	Position int    `json:"position"`
}

type SolutionBody struct {
	Title           string `json:"title"`
	Hint            string `json:"hint"`
	Description     string `json:"description"`
	Code            string `json:"code"`
	CodeExplanation string `json:"code_explanation"`
	Notes           string `json:"notes"`
	TimeComplexity  string `json:"time_complexity"`
	SpaceComplexity string `json:"space_complexity"`
	DifficultyLevel string `json:"difficulty_level"`
	DisplayOrder    int    `json:"display_order"`
	Author          string `json:"author"`
	IsActive        bool   `json:"is_active"`
}

type ProblemBody struct {
	Name          string         `json:"name"`
	ProblemNumber *int           `json:"problem_number"`
	Slug          string         `json:"slug"`
	Description   string         `json:"description"`
	Link          string         `json:"link"`
	Difficulty    string         `json:"difficulty"`
	StarterCode   any            `json:"starter_code"`
	SolutionCode  any            `json:"solution_code"`
	TimeLimit     int            `json:"time_limit"`
	MemoryLimit   int            `json:"memory_limit"`
	IsActive      bool           `json:"is_active"`
	Topics        []string       `json:"topics"`
	Lists         []string       `json:"lists"`
	TestCases     []TestCaseBody `json:"testcases"`
	Solutions     []SolutionBody `json:"solutions"`
}

func (ap *AdminProblemHandler) HandlerCreateProblem(w http.ResponseWriter, r *http.Request) {
	var problemBody ProblemBody
	err := json.NewDecoder(r.Body).Decode(&problemBody)
	if err != nil {
		ap.Logger.Println("Error decoding problem body", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	problem := models.Problem{
		Name:          problemBody.Name,
		ProblemNumber: problemBody.ProblemNumber,
		Slug:          problemBody.Slug,
		Description:   problemBody.Description,
		Link:          problemBody.Link,
		Difficulty:    problemBody.Difficulty,
		StarterCode:   problemBody.StarterCode,
		SolutionCode:  problemBody.SolutionCode,
		TimeLimit:     problemBody.TimeLimit,
		MemoryLimit:   problemBody.MemoryLimit,
		IsActive:      problemBody.IsActive,
	}

	var topicIDs []uuid.UUID
	for _, topicSlug := range problemBody.Topics {
		topicID, err := uuid.Parse(topicSlug)
		if err != nil {
			ap.Logger.Println("Error parsing topic slug", err)
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
			return
		}

		topicIDs = append(topicIDs, topicID)
	}

	var listIDs []uuid.UUID
	for _, listSlug := range problemBody.Lists {
		listID, err := uuid.Parse(listSlug)
		if err != nil {
			ap.Logger.Println("Error parsing list slug", err)
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
			return
		}

		listIDs = append(listIDs, listID)
	}

	var testcases []models.Testcase
	for _, tc := range problemBody.TestCases {
		testcases = append(testcases, models.Testcase{
			UI:       tc.UI,
			Input:    tc.Input,
			Output:   tc.Output,
			Position: tc.Position,
		})
	}

	var solutions []models.Solution
	for _, solution := range problemBody.Solutions {
		solutions = append(solutions, models.Solution{
			Title:           solution.Title,
			Hint:            solution.Hint,
			Description:     solution.Description,
			Code:            solution.Code,
			CodeExplanation: solution.CodeExplanation,
			Notes:           solution.Notes,
			TimeComplexity:  solution.TimeComplexity,
			SpaceComplexity: solution.SpaceComplexity,
			DifficultyLevel: solution.DifficultyLevel,
			DisplayOrder:    solution.DisplayOrder,
			Author:          solution.Author,
			IsActive:        solution.IsActive,
		})
	}

	err = ap.AdminProblemStore.CreateProblem(problem, listIDs, topicIDs, testcases, solutions)
	if err != nil {
		ap.Logger.Println("Error creating problem", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Successfully created problem"})
}

func (ap *AdminProblemHandler) HandlerGetProblemByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	problemID, err := uuid.Parse(id)
	if err != nil {
		ap.Logger.Println("Error parsing problem id", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	problem, err := ap.AdminProblemStore.GetProblemByID(problemID)
	if err != nil {
		ap.Logger.Println("Error getting problem by id", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": problem})
}

func (ap *AdminProblemHandler) HandlerUpdateProblem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	problemID, err := uuid.Parse(id)
	if err != nil {
		ap.Logger.Println("Error parsing problem id", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	var problemBody ProblemBody
	err = json.NewDecoder(r.Body).Decode(&problemBody)
	if err != nil {
		ap.Logger.Println("Error decoding problem body", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	problem := models.Problem{
		Name:          problemBody.Name,
		ProblemNumber: problemBody.ProblemNumber,
		Slug:          problemBody.Slug,
		Description:   problemBody.Description,
		Link:          problemBody.Link,
		Difficulty:    problemBody.Difficulty,
		StarterCode:   problemBody.StarterCode,
		SolutionCode:  problemBody.SolutionCode,
		TimeLimit:     problemBody.TimeLimit,
		MemoryLimit:   problemBody.MemoryLimit,
		IsActive:      problemBody.IsActive,
	}

	var topicIDs []uuid.UUID
	for _, topicSlug := range problemBody.Topics {
		topicID, err := uuid.Parse(topicSlug)
		if err != nil {
			ap.Logger.Println("Error parsing topic slug", err)
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
			return
		}

		topicIDs = append(topicIDs, topicID)
	}

	var listIDs []uuid.UUID
	for _, listSlug := range problemBody.Lists {
		listID, err := uuid.Parse(listSlug)
		if err != nil {
			ap.Logger.Println("Error parsing list slug", err)
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
			return
		}

		listIDs = append(listIDs, listID)
	}

	var testcases []models.Testcase
	for _, tc := range problemBody.TestCases {
		testcases = append(testcases, models.Testcase{
			UI:       tc.UI,
			Input:    tc.Input,
			Output:   tc.Output,
			Position: tc.Position,
		})
	}

	var solutions []models.Solution
	for _, solution := range problemBody.Solutions {
		solutions = append(solutions, models.Solution{
			Title:           solution.Title,
			Hint:            solution.Hint,
			Description:     solution.Description,
			Code:            solution.Code,
			CodeExplanation: solution.CodeExplanation,
			Notes:           solution.Notes,
			TimeComplexity:  solution.TimeComplexity,
			SpaceComplexity: solution.SpaceComplexity,
			DifficultyLevel: solution.DifficultyLevel,
			DisplayOrder:    solution.DisplayOrder,
			Author:          solution.Author,
			IsActive:        solution.IsActive,
		})
	}

	err = ap.AdminProblemStore.UpdateProblem(problemID, problem, listIDs, topicIDs, testcases, solutions)
	if err != nil {
		ap.Logger.Println("Error updating problem", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Successfully updated problem"})
}
