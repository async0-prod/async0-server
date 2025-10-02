package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/auth"
	"github.com/grvbrk/async0_server/internal/store"
	"github.com/grvbrk/async0_server/internal/utils"
)

type TestcaseHandler struct {
	TestcaseStore store.TestcaseStore
	Logger        *log.Logger
	Oauth         *auth.GoogleOauth
}

func NewTestcaseHandler(testcaseStore store.TestcaseStore, logger *log.Logger, oauth *auth.GoogleOauth) *TestcaseHandler {
	return &TestcaseHandler{
		TestcaseStore: testcaseStore,
		Logger:        logger,
		Oauth:         oauth,
	}
}

func (ph *TestcaseHandler) HandlerGetTestcaseByProblemID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	problemID, err := uuid.Parse(id)
	if err != nil {
		ph.Logger.Println("Error parsing problem id", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	testcase, err := ph.TestcaseStore.GetTestcasesByProblemID(problemID)
	if err != nil {
		ph.Logger.Println("Error getting testcase by slug", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": testcase})
}
