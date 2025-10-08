package admin

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/auth"
	"github.com/grvbrk/async0_server/internal/store/admin"
	"github.com/grvbrk/async0_server/internal/utils"
)

type AdminSolutionHandler struct {
	AdminSolutionStore admin.AdminSolutionStore
	Logger             *log.Logger
	Oauth              *auth.AdminGoogleOauth
}

func NewAdminSolutionHandler(adminSolutionStore admin.AdminSolutionStore, logger *log.Logger, oauth *auth.AdminGoogleOauth) *AdminSolutionHandler {
	return &AdminSolutionHandler{
		AdminSolutionStore: adminSolutionStore,
		Logger:             logger,
		Oauth:              oauth,
	}
}

func (as *AdminSolutionHandler) HandlerGetSolutionsByProblemID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	problemID, err := uuid.Parse(id)
	if err != nil {
		as.Logger.Println("Error parsing problem id", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	solutions, err := as.AdminSolutionStore.GetSolutionsByProblemID(problemID)
	if err != nil {
		as.Logger.Println("Error getting solutions by problem id", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": solutions})
}
