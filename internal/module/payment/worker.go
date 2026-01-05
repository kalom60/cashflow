package payment

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"

	"github.com/google/uuid"
	"github.com/kalom60/cashflow/internal/constant/dto"
	"github.com/kalom60/cashflow/internal/storage"
	"github.com/kalom60/cashflow/platform/logger"
	"github.com/kalom60/cashflow/platform/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type PaymentWorker struct {
	logger         logger.Logger
	paymentStorage storage.Payment
	msgClient      messaging.MessagingClient
}

func NewPaymentWorker(logger logger.Logger, paymentStorage storage.Payment, msgClient messaging.MessagingClient) *PaymentWorker {
	return &PaymentWorker{
		logger:         logger,
		paymentStorage: paymentStorage,
		msgClient:      msgClient,
	}
}

func (pw *PaymentWorker) Start(ctx context.Context) {
	pw.logger.Info(ctx, "Starting Payment Status Consumer...")

	msgs, err := pw.msgClient.ConsumePayments(ctx)
	if err != nil {
		pw.logger.Named("PaymentWorker-Start").Fatal(ctx, "failed to start consuming payments", zap.Error(err))
	}

	go func() {
		for msg := range msgs {
			pw.processMessage(ctx, msg)
		}
	}()
}

func (pw *PaymentWorker) processMessage(ctx context.Context, msg amqp.Delivery) {
	var body map[string]string
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		pw.logger.Named("PaymentWorker-ProcessMessage").Error(ctx, "failed to unmarshal message body", zap.Error(err))
		_ = msg.Nack(false, false)
		return
	}

	paymentIDStr, ok := body["payment_id"]
	if !ok {
		pw.logger.Named("PaymentWorker-ProcessMessage").Error(ctx, "message body missing payment_id")
		_ = msg.Nack(false, false)
		return
	}

	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		pw.logger.Named("PaymentWorker-ProcessMessage").Error(ctx, "failed to parse payment_id", zap.String("payment_id", paymentIDStr), zap.Error(err))
		_ = msg.Nack(false, false)
		return
	}

	pw.logger.Info(ctx, "Processing payment message", zap.String("payment_id", paymentID.String()))

	// Start Transaction for row-level locking and status check
	tx, err := pw.paymentStorage.BeginTx(ctx)
	if err != nil {
		pw.logger.Named("PaymentWorker-ProcessMessage-BeginTx").Error(ctx, "failed to begin transaction", zap.Error(err))
		_ = msg.Nack(false, true)
		return
	}
	defer tx.Rollback(ctx)

	// Lock the row and get current status
	payment, err := pw.paymentStorage.GetPaymentByIDForUpdate(ctx, tx, paymentID)
	if err != nil {
		pw.logger.Named("PaymentWorker-ProcessMessage-Lock").Error(ctx, "failed to lock payment for update", zap.String("payment_id", paymentID.String()), zap.Error(err))
		_ = msg.Nack(false, true) // Retry
		return
	}

	// Idempotency check: only process if PENDING
	if payment.Status != dto.PENDING {
		pw.logger.Info(ctx, "Payment already processed, skipping", zap.String("payment_id", paymentID.String()), zap.String("current_status", string(payment.Status)))
		_ = msg.Ack(false)
		return
	}

	// Simulate randomized processing result
	var status dto.PaymentStatus
	n, _ := rand.Int(rand.Reader, big.NewInt(100))
	if n.Int64() < 50 {
		status = dto.SUCCESS
	} else {
		status = dto.FAILED
	}

	pw.logger.Info(ctx, "Simulated processing result", zap.String("payment_id", paymentID.String()), zap.String("status", string(status)))

	if err := pw.paymentStorage.UpdatePaymentStatusWithTx(ctx, tx, paymentID, status); err != nil {
		pw.logger.Named("PaymentWorker-ProcessMessage-UpdateStatus").Error(ctx, "failed to update payment status", zap.String("payment_id", paymentID.String()), zap.Error(err))
		_ = msg.Nack(false, true) // Retry
		return
	}

	if err := tx.Commit(ctx); err != nil {
		pw.logger.Named("PaymentWorker-ProcessMessage-Commit").Error(ctx, "failed to commit transaction", zap.Error(err))
		_ = msg.Nack(false, true)
		return
	}

	if err := msg.Ack(false); err != nil {
		pw.logger.Named("PaymentWorker-ProcessMessage-Ack").Error(ctx, "failed to acknowledge message", zap.String("payment_id", paymentID.String()), zap.Error(err))
	}
}
