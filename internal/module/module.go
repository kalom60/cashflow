package module

import (
	"context"

	"github.com/google/uuid"
	"github.com/kalom60/cashflow/internal/constant/dto"
)

type Payment interface {
	CreatePayment(ctx context.Context, req dto.Payment) (dto.Payment, error)
	GetPaymentByID(ctx context.Context, id uuid.UUID) (dto.Payment, error)
}
