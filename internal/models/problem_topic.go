package models

import (
	"time"

	"github.com/google/uuid"
)

type ProblemTopic struct {
	ProblemID uuid.UUID `json:"problem_id"`
	TopicID   uuid.UUID `json:"topic_id"`
	CreatedAt time.Time `json:"created_at"`
}
