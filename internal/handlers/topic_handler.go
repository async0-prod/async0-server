package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/auth"
	"github.com/grvbrk/async0_server/internal/models"
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
	// Extract and validate list_id
	listIDStr := r.URL.Query().Get("list_id")
	if listIDStr == "" {
		th.Logger.Printf("Missing required parameter: list_id")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "list_id parameter is required",
		})
		return
	}

	listID, err := uuid.Parse(listIDStr)
	if err != nil {
		th.Logger.Printf("Invalid list_id format: %s, error: %v", listIDStr, err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "invalid list_id format",
		})
		return
	}

	showProblems := true
	problemsParam := r.URL.Query().Get("problems")
	if problemsParam != "" {
		parsed, err := strconv.ParseBool(problemsParam)
		if err != nil {
			th.Logger.Printf("Invalid problems parameter: %s, error: %v", problemsParam, err)
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
				"error": "problems parameter must be true or false",
			})
			return
		}
		showProblems = parsed
	}

	if !showProblems {
		topics, err := th.TopicStore.GetAllTopicsByListID(listID)
		if err != nil {
			th.Logger.Printf("Error getting topics for list %s: %v", listID, err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
				"error": "Internal Server Error",
			})
			return
		}

		if topics == nil {
			topics = []models.Topic{}
		}

		utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": topics})
		return
	}

	topics, err := th.TopicStore.GetAllTopicsAndProblemsByListID(listID)
	if err != nil {
		th.Logger.Printf("Error getting topics with problems for list %s: %v", listID, err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Internal Server Error",
		})
		return
	}

	if topics == nil {
		topics = []store.TopicsWithProblems{}
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": topics})
}
