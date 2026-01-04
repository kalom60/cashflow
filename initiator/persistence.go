package initiator

import (
	"github.com/kalom60/cashflow/internal/constant/model/persistencedb"
	"github.com/kalom60/cashflow/internal/storage"
	outboxevent "github.com/kalom60/cashflow/internal/storage/outbox_event"
	"github.com/kalom60/cashflow/internal/storage/payment"
	"github.com/kalom60/cashflow/platform/logger"
	"github.com/spf13/viper"
)

type Persistance struct {
	Payement    storage.Payment
	OutboxEvent storage.OutboxEvent
}

func initPersistence(persistencedb *persistencedb.PersistenceDB, log logger.Logger) *Persistance {
	limit := viper.GetInt("app.limit")
	paymentStorage := payment.Init(log, persistencedb)
	outboxEventStorage := outboxevent.Init(log, persistencedb, limit)

	return &Persistance{
		Payement:    paymentStorage,
		OutboxEvent: outboxEventStorage,
	}
}
