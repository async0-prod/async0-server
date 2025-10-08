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

type AdminTestcaseHandler struct {
	AdminTestcaseStore admin.AdminTestcaseStore
	Logger             *log.Logger
	Oauth              *auth.AdminGoogleOauth
}

func NewAdminTestcaseHandler(adminTestcaseStore admin.AdminTestcaseStore, logger *log.Logger, oauth *auth.AdminGoogleOauth) *AdminTestcaseHandler {
	return &AdminTestcaseHandler{
		AdminTestcaseStore: adminTestcaseStore,
		Logger:             logger,
		Oauth:              oauth,
	}
}

func (at *AdminTestcaseHandler) HandlerGetTestcasesByProblemID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	problemID, err := uuid.Parse(id)
	if err != nil {
		at.Logger.Println("Error parsing problem id", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	testcases, err := at.AdminTestcaseStore.GetTestcasesByProblemID(problemID)
	if err != nil {
		at.Logger.Println("Error getting testcases by problem id", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": testcases})
}
