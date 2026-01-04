package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/kalom60/cashflow/internal/constant/dto"
)

type Payment interface {
	CreatePayment(ctx context.Context, payment dto.Payment) (dto.Payment, error)
	GetPaymentByID(ctx context.Context, id uuid.UUID) (dto.Payment, error)
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status dto.PaymentStatus) error
}

type OutboxEvent interface {
	GetPendingOutboxEventsForUpdate(ctx context.Context) ([]dto.OutboxEvent, error)
	UpdateOutboxStatus(ctx context.Context, id uuid.UUID, status dto.OutboxStatus) error
	DeleteOutboxEvent(ctx context.Context, id uuid.UUID) error
}
