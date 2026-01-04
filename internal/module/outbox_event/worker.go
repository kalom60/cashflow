package outboxevent

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kalom60/cashflow/internal/storage"
	"github.com/kalom60/cashflow/platform/logger"
	"go.uber.org/zap"
)

type OutboxEventWorker struct {
	logger             logger.Logger
	outboxEventStorage storage.OutboxEvent
	interval           time.Duration
}

func Init(logger logger.Logger, outboxEventStorage storage.OutboxEvent, interval time.Duration) *OutboxEventWorker {
	return &OutboxEventWorker{
		logger:             logger,
		outboxEventStorage: outboxEventStorage,
		interval:           interval,
	}
}

func (w *OutboxEventWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.logger.Info(ctx, "Starting Global Outbox Worker...")

	for {
		select {
		case <-ctx.Done():
			w.logger.Info(ctx, "Stopping Global Outbox Worker...")
			return
		case <-ticker.C:
			w.processEvents(ctx)
		}
	}
}

func (w *OutboxEventWorker) processEvents(ctx context.Context) {
	tx, err := w.outboxEventStorage.BeginTx(ctx)
	if err != nil {
		return
	}
	defer tx.Rollback(ctx)

	events, err := w.outboxEventStorage.GetPendingOutboxEventsForUpdate(ctx, tx)
	if err != nil {
		return
	}

	for _, event := range events {
		var payload map[string]any
		if err := json.Unmarshal(event.Payload.Bytes, &payload); err != nil {
			w.logger.Named("OutboxEventWorker-ProcessEvents").Error(ctx, "failed to Unmarshal event payload", zap.Any("event_id", event.ID), zap.Error(err))
			w.logger.Named("OutboxEventWorker-ProcessEvents").Info(ctx, "deleting corrupted outbox event", zap.Any("event_id", event.ID), zap.Error(err))
			w.outboxEventStorage.DeleteOutboxEvent(ctx, tx, event.ID)
			continue
		}

		// TODO: handle messaging

		if err := w.outboxEventStorage.DeleteOutboxEvent(ctx, tx, event.ID); err != nil {
			w.logger.Named("OutboxEventWorker-ProcessEvents-DeleteOutboxEvent").Error(ctx, "failed to delete outbox event", zap.Any("event_id", event.ID), zap.Error(err))
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		w.logger.Named("OutboxEventWorker-ProcessEvents-Commit").Error(ctx, "failed to commit transaction", zap.Error(err))
	}
}
