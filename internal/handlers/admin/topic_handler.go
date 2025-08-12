package admin

import (
	"log"
	"net/http"

	"github.com/grvbrk/async0_server/internal/auth"
	"github.com/grvbrk/async0_server/internal/store/admin"
	"github.com/grvbrk/async0_server/internal/utils"
)

type AdminTopicHandler struct {
	AdminTopicStore admin.AdminTopicStore
	Logger          *log.Logger
	Oauth           *auth.AdminGoogleOauth
}

func NewAdminTopicHandler(adminTopicStore admin.AdminTopicStore, logger *log.Logger, oauth *auth.AdminGoogleOauth) *AdminTopicHandler {
	return &AdminTopicHandler{
		AdminTopicStore: adminTopicStore,
		Logger:          logger,
		Oauth:           oauth,
	}
}

func (at *AdminTopicHandler) HandlerGetAllTopics(w http.ResponseWriter, r *http.Request) {

	topics, err := at.AdminTopicStore.GetAllTopics()
	if err != nil {
		at.Logger.Println("Error getting all topics from store", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": topics})

}
