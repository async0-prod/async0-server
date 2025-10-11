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

type SolutionHandler struct {
	SolutionStore store.SolutionStore
	Logger        *log.Logger
	Oauth         *auth.GoogleOauth
}

func NewSolutionHandler(solutionStore store.SolutionStore, logger *log.Logger, oauth *auth.GoogleOauth) *SolutionHandler {
	return &SolutionHandler{
		SolutionStore: solutionStore,
		Logger:        logger,
		Oauth:         oauth,
	}
}

func (sh *SolutionHandler) HandlerGetSolutionsByProblemID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	problemID, err := uuid.Parse(id)
	if err != nil {
		sh.Logger.Println("Error parsing problem id", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	solutions, err := sh.SolutionStore.GetSolutionsByProblemID(problemID)
	if err != nil {
		sh.Logger.Println("Error getting solutions by problem id", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": solutions})

}
