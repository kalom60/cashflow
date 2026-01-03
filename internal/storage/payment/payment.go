package payment

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/kalom60/cashflow/internal/constant/dto"
	customErrors "github.com/kalom60/cashflow/internal/constant/errors"
	"github.com/kalom60/cashflow/internal/constant/model/db"
	"github.com/kalom60/cashflow/internal/constant/model/persistencedb"
	"github.com/kalom60/cashflow/internal/storage"
	"github.com/kalom60/cashflow/platform/logger"
	"go.uber.org/zap"
)

type paymentStore struct {
	logger        logger.Logger
	persistencedb *persistencedb.PersistenceDB
}

func Init(logger logger.Logger, persistencedb *persistencedb.PersistenceDB) storage.Payment {
	return &paymentStore{
		logger:        logger,
		persistencedb: persistencedb,
	}
}

func (ps *paymentStore) CreatePayment(ctx context.Context, payment dto.Payment) (dto.Payment, error) {
	row, err := ps.persistencedb.Queries.CreatePayment(ctx, db.CreatePaymentParams{
		Reference: payment.Reference,
		Amount:    payment.Amount,
		Currency:  db.PaymentCurrency(payment.Currency),
		Status:    db.PaymentStatus(payment.Status),
		CreatedAt: payment.CreatedAt,
	})
	if err != nil {
		ps.logger.Named("Create-Payment-Store").Error(ctx, "failed to create payment", zap.Error(err))
		err = customErrors.ErrUnableToCreate.New("failed to create payment")
		return dto.Payment{}, err
	}

	payment.ID = row.ID
	payment.CreatedAt = row.CreatedAt
	payment.UpdatedAt = row.UpdatedAt

	return payment, nil
}

func (ps *paymentStore) GetPaymentByID(ctx context.Context, id uuid.UUID) (dto.Payment, error) {
	row, err := ps.persistencedb.Queries.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			ps.logger.Named("Get-Payment-By-ID-Store").Error(ctx, "no row found", zap.Any("id", id), zap.Error(err))
			err = customErrors.ErrResourceNotFound.New("no row found")
			return dto.Payment{}, err
		}
		ps.logger.Named("Get-Payment-By-ID-Store").Error(ctx, "failed to get payment by id", zap.Any("id", id), zap.Error(err))
		err = customErrors.ErrUnableToGet.New("failed to get payment by id")
		return dto.Payment{}, err
	}

	return dto.Payment{
		ID:        row.ID,
		Reference: row.Reference,
		Amount:    row.Amount,
		Currency:  dto.PaymentCurrency(row.Currency),
		Status:    dto.PaymentStatus(row.Status),
		CreatedAt: row.CreatedAt,
	}, nil
}

func (ps *paymentStore) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status dto.PaymentStatus) error {
	_, err := ps.persistencedb.Queries.UpdatePaymentStatus(ctx, db.UpdatePaymentStatusParams{
		ID:     id,
		Status: db.PaymentStatus(status),
	})
	if err != nil {
		ps.logger.Named("Update-Payment-Status-Store").Error(ctx, "failed to update payment status", zap.Any("id", id), zap.Error(err))
		err = customErrors.ErrUnableToUpdate.New("failed to update payment")
		return err
	}

	return nil
}
