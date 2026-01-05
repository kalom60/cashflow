package outboxevent

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kalom60/cashflow/internal/storage"
	"github.com/kalom60/cashflow/platform/logger"
	"github.com/kalom60/cashflow/platform/messaging"
	"go.uber.org/zap"
)

type OutboxEventWorker struct {
	logger             logger.Logger
	outboxEventStorage storage.OutboxEvent
	msgClient          messaging.MessagingClient
	interval           time.Duration
}

func Init(logger logger.Logger, outboxEventStorage storage.OutboxEvent, msgClient messaging.MessagingClient, interval time.Duration) *OutboxEventWorker {
	return &OutboxEventWorker{
		logger:             logger,
		outboxEventStorage: outboxEventStorage,
		msgClient:          msgClient,
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
		// Handle messaging
		paymentID, ok := payload["id"].(string)
		if !ok {
			w.logger.Named("OutboxEventWorker-ProcessEvents").Error(ctx, "payload missing payment id", zap.Any("event_id", event.ID))
			w.outboxEventStorage.DeleteOutboxEvent(ctx, tx, event.ID)
			continue
		}

		if err := w.msgClient.PublishPayment(ctx, paymentID); err != nil {
			w.logger.Named("OutboxEventWorker-ProcessEvents-Publish").Error(ctx, "failed to publish message", zap.Any("payment_id", paymentID), zap.Error(err))
			return
		}

		if err := w.outboxEventStorage.DeleteOutboxEvent(ctx, tx, event.ID); err != nil {
			w.logger.Named("OutboxEventWorker-ProcessEvents-DeleteOutboxEvent").Error(ctx, "failed to delete outbox event", zap.Any("event_id", event.ID), zap.Error(err))
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		w.logger.Named("OutboxEventWorker-ProcessEvents-Commit").Error(ctx, "failed to commit transaction", zap.Error(err))
	}
}
