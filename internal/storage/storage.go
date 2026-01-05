package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/kalom60/cashflow/internal/constant/dto"
)

type Payment interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	CreatePayment(ctx context.Context, payment dto.Payment) (dto.Payment, error)
	GetPaymentByID(ctx context.Context, id uuid.UUID) (dto.Payment, error)
	GetPaymentByIDForUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID) (dto.Payment, error)
	UpdatePaymentStatusWithTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, status dto.PaymentStatus) error
}

type OutboxEvent interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)

	GetPendingOutboxEventsForUpdate(ctx context.Context, tx pgx.Tx) ([]dto.OutboxEvent, error)
	UpdateOutboxStatus(ctx context.Context, tx pgx.Tx, id uuid.UUID, status dto.OutboxStatus) error
	DeleteOutboxEvent(ctx context.Context, tx pgx.Tx, id uuid.UUID) error
}
