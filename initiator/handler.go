package initiator

import (
	"github.com/kalom60/cashflow/internal/handler"
	"github.com/kalom60/cashflow/internal/handler/payment"
	"github.com/kalom60/cashflow/platform/logger"
)

type Handler struct {
	Payment handler.Payment
}

func initHandler(module *Module, log logger.Logger) *Handler {
	return &Handler{
		Payment: payment.Init(log, module.Payment),
	}
}
