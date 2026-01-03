package payment

import (
	"context"

	"github.com/google/uuid"
	"github.com/kalom60/cashflow/internal/constant/dto"
	"github.com/kalom60/cashflow/internal/module"
	"github.com/kalom60/cashflow/internal/storage"
	"github.com/kalom60/cashflow/platform/logger"
)

type paymentModule struct {
	logger         logger.Logger
	paymentStorage storage.Payment
}

func Init(logger logger.Logger, paymentStorage storage.Payment) module.Payment {
	return &paymentModule{
		logger:         logger,
		paymentStorage: paymentStorage,
	}
}

func (pm *paymentModule) CreatePayment(ctx context.Context, req dto.Payment) (dto.Payment, error) {
	payment, err := pm.paymentStorage.CreatePayment(ctx, req)
	if err != nil {
		return dto.Payment{}, err
	}

	return payment, nil
}

func (pm *paymentModule) GetPaymentByID(ctx context.Context, id uuid.UUID) (dto.Payment, error) {
	payment, err := pm.paymentStorage.GetPaymentByID(ctx, id)
	if err != nil {
		return dto.Payment{}, err
	}

	return payment, nil
}
