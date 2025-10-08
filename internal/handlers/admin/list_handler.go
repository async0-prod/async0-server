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

type AdminListHandler struct {
	AdminListStore admin.AdminListStore
	Logger         *log.Logger
	Oauth          *auth.AdminGoogleOauth
}

func NewAdminListHandler(adminListStore admin.AdminListStore, logger *log.Logger, oauth *auth.AdminGoogleOauth) *AdminListHandler {
	return &AdminListHandler{
		AdminListStore: adminListStore,
		Logger:         logger,
		Oauth:          oauth,
	}
}

func (ah *AdminListHandler) HandlerGetAllLists(w http.ResponseWriter, r *http.Request) {
	lists, err := ah.AdminListStore.GetAllLists()
	if err != nil {
		ah.Logger.Println("Error getting all lists", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": lists})
}

func (ah *AdminListHandler) HandlerGetListsByProblemID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	problemID, err := uuid.Parse(id)
	if err != nil {
		ah.Logger.Println("Error parsing problem id", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	lists, err := ah.AdminListStore.GetListsByProblemID(problemID)
	if err != nil {
		ah.Logger.Println("Error getting lists by problem id", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": lists})
}
