package initiator

import (
	"context"

	"github.com/kalom60/cashflow/internal/module"
	outboxevent "github.com/kalom60/cashflow/internal/module/outbox_event"
	"github.com/kalom60/cashflow/internal/module/payment"
	"github.com/kalom60/cashflow/platform/logger"
	"github.com/kalom60/cashflow/platform/messaging"
	"github.com/spf13/viper"
)

type Module struct {
	Payment     module.Payment
	OutboxEvent outboxevent.OutboxEventWorker
}

func initModule(
	persistence *Persistance,
	msgClient messaging.MessagingClient,
	log logger.Logger,
) *Module {
	paymentStorage := persistence.Payement
	outboxEventStorage := persistence.OutboxEvent

	paymentModule := payment.Init(log, paymentStorage)

	interval := viper.GetDuration("app.interval")
	outboxEventModule := outboxevent.Init(log, outboxEventStorage, msgClient, interval)

	// Start Global Outbox Worker
	go outboxEventModule.Start(context.Background())

	// Start Payment Status Consumer
	paymentWorker := payment.NewPaymentWorker(log, paymentStorage, msgClient)
	paymentWorker.Start(context.Background())

	return &Module{
		Payment:     paymentModule,
		OutboxEvent: *outboxEventModule,
	}
}
