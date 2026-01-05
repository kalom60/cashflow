package payment

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
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
	tx, err := ps.persistencedb.Pool.Begin(ctx)
	if err != nil {
		ps.logger.Named("PaymentStore-CreatePayment-BeginTx").Error(ctx, "failed to start transaction", zap.Error(err))
		return dto.Payment{}, customErrors.ErrUnableToCreate.New("database transaction failed")
	}
	defer tx.Rollback(ctx)

	qtx := ps.persistencedb.Queries.WithTx(tx)

	row, err := qtx.CreatePayment(ctx, db.CreatePaymentParams{
		Reference: payment.Reference,
		Amount:    payment.Amount,
		Currency:  db.PaymentCurrency(payment.Currency),
		Status:    db.PaymentStatus(payment.Status),
		CreatedAt: payment.CreatedAt,
	})
	if err != nil {
		ps.logger.Named("PaymentStore-CreatePayment-InsertPayment").Error(ctx, "failed to insert payment record", zap.Error(err))
		return dto.Payment{}, customErrors.ErrUnableToCreate.New("failed to save payment to storage")
	}

	payment.ID = row.ID
	payment.CreatedAt = row.CreatedAt
	payment.UpdatedAt = row.UpdatedAt

	payloadJson, err := json.Marshal(payment)
	if err != nil {
		ps.logger.Named("PaymentStore-CreatePayment-Marshal").Error(ctx, "failed to marshal outbox payload", zap.Error(err))
		return dto.Payment{}, err
	}

	var jsonbPayload pgtype.JSONB
	if err := jsonbPayload.Set(payloadJson); err != nil {
		ps.logger.Named("PaymentStore-CreatePayment-SetJSONB").Error(ctx, "failed to set jsonb payload", zap.Error(err))
		return dto.Payment{}, err
	}

	_, err = qtx.CreateOutboxEvent(ctx, db.CreateOutboxEventParams{
		Payload:   jsonbPayload,
		Status:    db.OutboxStatus(dto.OutboxStatusPending),
		CreatedAt: time.Now(),
	})
	if err != nil {
		ps.logger.Named("PaymentStore-CreatePayment-InsertOutbox").Error(ctx, "failed to insert outbox event", zap.Error(err))
		return dto.Payment{}, customErrors.ErrUnableToCreate.New("failed to save outbox event")
	}

	if err := tx.Commit(ctx); err != nil {
		ps.logger.Named("PaymentStore-CreatePayment-Commit").Error(ctx, "failed to commit transaction", zap.Error(err))
		return dto.Payment{}, customErrors.ErrUnableToCreate.New("final database commit failed")
	}

	return payment, nil
}

func (ps *paymentStore) GetPaymentByID(ctx context.Context, id uuid.UUID) (dto.Payment, error) {
	row, err := ps.persistencedb.Queries.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			ps.logger.Named("PaymentStore-GetPaymentByID").Error(ctx, "no row found", zap.Any("id", id), zap.Error(err))
			return dto.Payment{}, customErrors.ErrResourceNotFound.New("payment not found")
		}
		ps.logger.Named("PaymentStore-GetPaymentByID").Error(ctx, "failed to get payment by id", zap.Any("id", id), zap.Error(err))
		return dto.Payment{}, customErrors.ErrUnableToGet.New("failed to get payment")

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
		ps.logger.Named("PaymentStore-UpdatePaymentStatus").Error(ctx, "failed to update payment status", zap.Any("id", id), zap.Error(err))
		return customErrors.ErrUnableToUpdate.New("failed to update payment")

	}

	return nil
}

func (ps *paymentStore) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return ps.persistencedb.Pool.Begin(ctx)
}

func (ps *paymentStore) GetPaymentByIDForUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID) (dto.Payment, error) {
	qtx := ps.persistencedb.Queries.WithTx(tx)
	row, err := qtx.GetPaymentByIDForUpdate(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return dto.Payment{}, customErrors.ErrResourceNotFound.New("payment not found")
		}
		return dto.Payment{}, customErrors.ErrUnableToGet.New("failed to get payment for update")
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

func (ps *paymentStore) UpdatePaymentStatusWithTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, status dto.PaymentStatus) error {
	qtx := ps.persistencedb.Queries.WithTx(tx)
	_, err := qtx.UpdatePaymentStatus(ctx, db.UpdatePaymentStatusParams{
		ID:     id,
		Status: db.PaymentStatus(status),
	})
	if err != nil {
		return customErrors.ErrUnableToUpdate.New("failed to update payment status in tx")
	}
	return nil
}
