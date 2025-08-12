package admin

import (
	"log"
	"net/http"

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

func (al *AdminListHandler) HandlerGetAllLists(w http.ResponseWriter, r *http.Request) {

	lists, err := al.AdminListStore.GetAllLists()
	if err != nil {
		al.Logger.Println("Error getting all lists from store", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": lists})

}
