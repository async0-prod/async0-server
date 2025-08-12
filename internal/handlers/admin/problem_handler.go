package admin

import (
	"encoding/json"
	"log"
	"net/http"

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
	Name     string `json:"name"`
	Input    string `json:"input"`
	Output   string `json:"output"`
	Position int    `json:"position"`
}

type ProblemBody struct {
	Name         string         `json:"name"`
	Slug         string         `json:"slug"`
	Link         string         `json:"link"`
	Difficulty   string         `json:"difficulty"`
	StarterCode  any            `json:"starter_code"`
	SolutionCode any            `json:"solution_code"`
	TimeLimit    int            `json:"time_limit"`
	MemoryLimit  int            `json:"memory_limit"`
	IsActive     bool           `json:"is_active"`
	Topics       []string       `json:"topics"`
	Lists        []string       `json:"lists"`
	TestCases    []TestCaseBody `json:"test_cases"`
}

// {
//     "name": "Two Sum",
//     "slug": "two-sum",
//     "link": "asaas",
//     "difficulty": "Easy",
//     "starter_code": {},
//     "solution_code": {},
//     "time_limit": 2000,
//     "memory_limit": 256,
//     "is_active": true,
//     "topics": [
//         "2ce5b4f9-7ebb-4ce8-9787-558c11ca86ad",
//         "660b2f88-7b8f-4009-8cb5-184c45780cdf"
//     ],
//     "lists": [
//         "daa215cc-c573-466e-94ec-902ec072c9f7",
//         "7ab83b2c-66e8-4b00-9858-f340212e481a"
//     ],
//     "testCases": [
//         {
//             "name": "asas",
//             "input": "asas",
//             "output": "asas",
// 						"position": 1,

//         }
//     ]
// }

func (ap *AdminProblemHandler) HandlerCreateProblem(w http.ResponseWriter, r *http.Request) {
	var problemBody ProblemBody
	err := json.NewDecoder(r.Body).Decode(&problemBody)
	if err != nil {
		ap.Logger.Println("Error decoding problem body", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	problem := models.Problem{
		Name:         problemBody.Name,
		Slug:         problemBody.Slug,
		Link:         problemBody.Link,
		Difficulty:   problemBody.Difficulty,
		StarterCode:  problemBody.StarterCode,
		SolutionCode: problemBody.SolutionCode,
		TimeLimit:    problemBody.TimeLimit,
		MemoryLimit:  problemBody.MemoryLimit,
		IsActive:     problemBody.IsActive,
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

	var testcases []models.TestCase
	for _, testCase := range problemBody.TestCases {
		testcases = append(testcases, models.TestCase{
			UI:       testCase.Name,
			Input:    testCase.Input,
			Output:   testCase.Output,
			Position: testCase.Position,
		})
	}

	err = ap.AdminProblemStore.CreateProblem(problem, listIDs, topicIDs, testcases)
	if err != nil {
		ap.Logger.Println("Error creating problem", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Successfully created problem"})
}
