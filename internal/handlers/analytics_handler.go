package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/auth"
	"github.com/grvbrk/async0_server/internal/middlewares"
	"github.com/grvbrk/async0_server/internal/store"
	"github.com/grvbrk/async0_server/internal/utils"
)

type AnalyticsHandler struct {
	AnalyticsStore store.AnalyticsStore
	Logger         *log.Logger
	Oauth          *auth.GoogleOauth
}

func NewAnalyticsHandler(logger *log.Logger, oauth *auth.GoogleOauth, analyticsStore store.AnalyticsStore) *AnalyticsHandler {
	return &AnalyticsHandler{
		AnalyticsStore: analyticsStore,
		Logger:         logger,
		Oauth:          oauth,
	}
}

func (h *AnalyticsHandler) HandlerGetCardAnalyticsByListID(w http.ResponseWriter, r *http.Request) {
	listIDstr := chi.URLParam(r, "listID")

	listID, err := uuid.Parse(listIDstr)
	if err != nil {
		h.Logger.Println("Error parsing list id")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	user, ok := middlewares.GetUserFromContext(r)
	if !ok {
		analytics, err := h.AnalyticsStore.GetCardAnalyticsByListIDNoUser(listID)
		if err != nil {
			h.Logger.Println("Error getting card analytics")
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": analytics})
		return
	}

	fmt.Println(user.ID, "reached ere")

	analytics, err := h.AnalyticsStore.GetCardAnalyticsByListID(user.ID, listID)
	if err != nil {
		h.Logger.Println("Error getting card analytics")
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": analytics})
}
