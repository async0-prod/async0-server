package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/auth"
	"github.com/grvbrk/async0_server/internal/middlewares"
	"github.com/grvbrk/async0_server/internal/store"
	"github.com/grvbrk/async0_server/internal/utils"
)

type ProblemHandler struct {
	ProblemStore store.ProblemStore
	Logger       *log.Logger
	Oauth        *auth.GoogleOauth
}

func NewProblemHandler(problemStore store.ProblemStore, logger *log.Logger, oauth *auth.GoogleOauth) *ProblemHandler {
	return &ProblemHandler{
		ProblemStore: problemStore,
		Logger:       logger,
		Oauth:        oauth,
	}
}

func (ph *ProblemHandler) HandlerGetProblemBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	problem, err := ph.ProblemStore.GetProblemBySlug(slug)

	if err != nil {
		if errors.Is(err, store.ErrProblemNotFound) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"message": "Bad Request"})
			return
		}

		ph.Logger.Println("Error getting problem by slug:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": problem})
}

func (ph *ProblemHandler) HandlerGetTanstackTableProblems(w http.ResponseWriter, r *http.Request) {
	var userID *uuid.UUID

	user, ok := middlewares.GetUserFromContext(r)
	if ok {
		userID = &user.ID
	}

	id := chi.URLParam(r, "listID")
	listID, err := uuid.Parse(id)
	if err != nil {
		ph.Logger.Println("Error parsing list id", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	tableProblems, err := ph.ProblemStore.GetTanstackTableProblemsByListID(userID, listID)
	if err != nil {
		ph.Logger.Println("Error getting tanstack table problems", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": tableProblems})
}
