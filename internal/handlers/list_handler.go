package handlers

import (
	"log"
	"net/http"

	"github.com/grvbrk/async0_server/internal/auth"
	"github.com/grvbrk/async0_server/internal/store"
	"github.com/grvbrk/async0_server/internal/utils"
)

type ListHandler struct {
	ListStore store.ListStore
	Logger    *log.Logger
	Oauth     *auth.GoogleOauth
}

func NewListHandler(listStore store.ListStore, logger *log.Logger, oauth *auth.GoogleOauth) *ListHandler {
	return &ListHandler{
		ListStore: listStore,
		Logger:    logger,
		Oauth:     oauth,
	}
}

func (lh *ListHandler) HandlerGetAllLists(w http.ResponseWriter, r *http.Request) {

	lists, err := lh.ListStore.GetAllLists()
	if err != nil {
		lh.Logger.Println("Error getting all lists from store", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": lists})

}
