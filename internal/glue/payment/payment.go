package payment

import (
	"net/http"

	"github.com/kalom60/cashflow/internal/glue/routing"
	"github.com/kalom60/cashflow/internal/handler"
	"github.com/kalom60/cashflow/platform/logger"
	"github.com/labstack/echo/v4"
)

func RegisterPaymentRoutes(
	group *echo.Group,
	paymentHandler handler.Payment,
	log logger.Logger,
) {

	payments := []routing.Route{
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/payments",
			Handler: paymentHandler.CreatePayment,
		}, {
			Method:  http.MethodGet,
			Path:    "/api/v1/payments/:id",
			Handler: paymentHandler.GetPaymentDetails,
		},
	}

	routing.RegisterRoute(group, payments, log)
}
