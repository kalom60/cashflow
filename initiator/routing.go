package initiator

import (
	"github.com/kalom60/cashflow/internal/glue/payment"
	"github.com/kalom60/cashflow/platform/logger"
	"github.com/labstack/echo/v4"
)

func initRoute(eg *echo.Group, handler *Handler, logger logger.Logger) {
	payment.RegisterPaymentRoutes(eg, handler.Payment, logger)
}
