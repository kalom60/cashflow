package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type OutboxStatus string

const (
	OutboxStatusPending OutboxStatus = "PENDING"
	OutboxStatusSent    OutboxStatus = "SENT"
	OutboxStatusFailed  OutboxStatus = "FAILED"
)

type OutboxEvent struct {
	ID        uuid.UUID    `json:"id"`
	Payload   pgtype.JSONB `json:"payload"`
	Status    OutboxStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
