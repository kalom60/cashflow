package outboxevent

import (
	"context"

	"github.com/google/uuid"
	"github.com/kalom60/cashflow/internal/constant/dto"
	customErrors "github.com/kalom60/cashflow/internal/constant/errors"
	"github.com/kalom60/cashflow/internal/constant/model/db"
	"github.com/kalom60/cashflow/internal/constant/model/persistencedb"
	"github.com/kalom60/cashflow/internal/storage"
	"github.com/kalom60/cashflow/platform/logger"
	"go.uber.org/zap"
)

type outboxEventStore struct {
	logger        logger.Logger
	persistencedb *persistencedb.PersistenceDB
	limit         int
}

func Init(logger logger.Logger, persistencedb *persistencedb.PersistenceDB, limit int) storage.OutboxEvent {
	return &outboxEventStore{
		logger:        logger,
		persistencedb: persistencedb,
		limit:         limit,
	}
}

func (oes *outboxEventStore) GetPendingOutboxEventsForUpdate(ctx context.Context) ([]dto.OutboxEvent, error) {
	rows, err := oes.persistencedb.Queries.GetPendingOutboxEventsForUpdate(ctx, int32(oes.limit))
	if err != nil {
		oes.logger.Named("OutboxEventStore-GetPendingOutboxEventsForUpdate").Error(ctx, "failed to get outbox events", zap.Error(err))
		return nil, customErrors.ErrUnableToGet.New("failed to get outbox events")
	}

	events := make([]dto.OutboxEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, dto.OutboxEvent{
			ID:        row.ID,
			Payload:   row.Payload,
			Status:    dto.OutboxStatus(row.Status),
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		})
	}

	return events, nil
}

func (oes *outboxEventStore) UpdateOutboxStatus(ctx context.Context, id uuid.UUID, status dto.OutboxStatus) error {
	rowsAffected, err := oes.persistencedb.Queries.UpdateOutboxStatus(ctx, db.UpdateOutboxStatusParams{
		ID:     id,
		Status: db.OutboxStatus(status),
	})
	if err != nil {
		oes.logger.Named("OutboxEventStore-UpdateOutboxStatus").Error(ctx, "failed to update outbox event", zap.Any("id", id), zap.Error(err))
		return customErrors.ErrUnableToUpdate.New("failed to update outbox event")
	}

	if rowsAffected == 0 {
		oes.logger.Named("OutboxEventStore-UpdateOutboxStatus").Error(ctx, "no row found", zap.Any("id", id), zap.Error(err))
		return customErrors.ErrResourceNotFound.New("no row found")
	}

	return nil
}

func (oes *outboxEventStore) DeleteOutboxEvent(ctx context.Context, id uuid.UUID) error {
	rowsAffected, err := oes.persistencedb.Queries.DeleteOutboxEvent(ctx, id)
	if err != nil {
		oes.logger.Named("OutboxEventStore-DeleteOutboxEvent").Error(ctx, "failed to delete outbox event", zap.Any("id", id), zap.Error(err))
		return customErrors.ErrUnableToUpdate.New("failed to delete outbox event")
	}

	if rowsAffected == 0 {
		oes.logger.Named("OutboxEventStore-DeleteOutboxEvent").Error(ctx, "no row found", zap.Any("id", id), zap.Error(err))
		return customErrors.ErrResourceNotFound.New("no row found")
	}

	return nil
}
