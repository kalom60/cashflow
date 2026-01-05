package payment

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kalom60/cashflow/internal/constant/dto"
	"github.com/kalom60/cashflow/internal/constant/response"
	"github.com/kalom60/cashflow/internal/handler"
	"github.com/kalom60/cashflow/internal/module"
	"github.com/kalom60/cashflow/platform/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type paymentHandler struct {
	logger        logger.Logger
	paymentModule module.Payment
}

func Init(logger logger.Logger, paymentModule module.Payment) handler.Payment {
	return &paymentHandler{
		logger:        logger,
		paymentModule: paymentModule,
	}
}

// CreatePayment godoc
//
//	@Summary		Create a new payment
//	@Description	Creates a new payment record and initiates processing via RabbitMQ
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			payment	body		dto.CreatePaymentRequest	true	"Payment creation request"
//	@Success		201		{object}	dto.CreatePaymentResponse
//	@Failure		400		{object}	map[string]string	"Invalid input"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/api/v1/payments [post]
func (ph *paymentHandler) CreatePayment(c echo.Context) error {
	var req dto.CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		ph.logger.Named("PaymentHandler-CreatePayment-Bind").Error(c.Request().Context(), "failed to bind request", zap.Any("error", err.Error()))
		return response.SendErrorResponse(c, 400, "invalid request payload")
	}

	if err := req.Validate(); err != nil {
		ph.logger.Named("PaymentHandler-CreatePayment-Validate").Error(c.Request().Context(), "validation failed", zap.Any("error", err.Error()))
		return response.SendErrorResponse(c, 400, err.Error())
	}

	payment, err := ph.paymentModule.CreatePayment(c.Request().Context(), req.ToPayment())

	if err != nil {
		ph.logger.Named("PaymentHandler-CreatePayment-Module").Error(c.Request().Context(), "failed to create payment", zap.Any("error", err.Error()))
		return response.SendErrorResponseFormated(c, err)
	}

	return response.SendSuccessResponse(c, http.StatusCreated, dto.CreatePaymentResponse{
		ID:     payment.ID,
		Status: payment.Status,
	})
}

// GetPaymentDetails godoc
//
//	@Summary		Get payment details
//	@Description	Retrieves details of a payment by its unique ID
//	@Tags			Payments
//	@Produce		json
//	@Param			id	path		string	true	"Payment ID"
//	@Success		200	{object}	dto.GetPaymentDetailsResponse
//	@Failure		400	{object}	map[string]string	"Invalid ID format"
//	@Failure		404	{object}	map[string]string	"Payment not found"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/api/v1/payments/{id} [get]
func (ph *paymentHandler) GetPaymentDetails(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.SendErrorResponse(c, 400, "invalid payment id format")
	}

	payment, err := ph.paymentModule.GetPaymentByID(c.Request().Context(), id)
	if err != nil {
		ph.logger.Named("PaymentHandler-GetPaymentDetails-Module").Error(c.Request().Context(), "failed to get payment", zap.Any("id", id), zap.Any("error", err.Error()))
		return response.SendErrorResponseFormated(c, err)
	}

	return response.SendSuccessResponse(c, http.StatusOK, dto.GetPaymentDetailsResponse{
		ID:        payment.ID,
		Amount:    payment.Amount,
		Currency:  payment.Currency,
		Reference: payment.Reference,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt,
	})
}
