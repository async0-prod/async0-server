package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/auth"
	"github.com/grvbrk/async0_server/internal/store"
	"github.com/grvbrk/async0_server/internal/utils"
)

type TopicHandler struct {
	TopicStore store.TopicStore
	Logger     *log.Logger
	Oauth      *auth.GoogleOauth
}

func NewTopicHandler(topicStore store.TopicStore, logger *log.Logger, oauth *auth.GoogleOauth) *TopicHandler {
	return &TopicHandler{
		TopicStore: topicStore,
		Logger:     logger,
		Oauth:      oauth,
	}
}

func (th *TopicHandler) HandlerGetTopics(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("list_id")

	if id == "" {
		th.Logger.Println("Error getting topics: list id is required")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	listID, err := uuid.Parse(id)
	if err != nil {
		th.Logger.Println("Error parsing list id", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	showProblems := r.URL.Query().Get("problems")

	if showProblems == "" {
		th.Logger.Println("Error getting topics: show problems param is required")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	showProblemsBool, err := strconv.ParseBool(showProblems)
	if err != nil {
		th.Logger.Println("Error parsing show problems param", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"message": "Bad Request"})
		return
	}

	if !showProblemsBool {
		topics, err := th.TopicStore.GetAllTopicsByListID(listID)
		if err != nil {
			th.Logger.Println("Error getting topics by list id", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": topics})
		return
	}

	topics, err := th.TopicStore.GetAllTopicsAndProblemsByListID(listID)
	if err != nil {
		th.Logger.Println("Error getting topics with problems by list id", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": topics})

}
