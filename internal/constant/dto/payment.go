package dto

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentCurrency string

const (
	ETB PaymentCurrency = "ETB"
	USD PaymentCurrency = "USD"
)

type PaymentStatus string

const (
	PENDING PaymentStatus = "PENDING"
	SUCCESS PaymentStatus = "SUCCESS"
	FAILED  PaymentStatus = "FAILED"
)

type Payment struct {
	ID        uuid.UUID       `json:"id"`
	Reference uuid.UUID       `json:"reference"`
	Amount    decimal.Decimal `json:"amount"`
	Currency  PaymentCurrency `json:"currency"`
	Status    PaymentStatus   `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CreatePaymentRequest struct {
	Amount    decimal.Decimal `json:"amount"`
	Currency  PaymentCurrency `json:"currency"`
	Reference uuid.UUID       `json:"reference"`
}

func (r *CreatePaymentRequest) Validate() error {
	if r.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be greater than zero")
	}

	// Max 2 decimal places as per NUMERIC(18,2)
	if r.Amount.Exponent() < -2 {
		return errors.New("amount cannot have more than 2 decimal places")
	}

	if r.Currency != ETB && r.Currency != USD {
		return fmt.Errorf("invalid currency: %s", r.Currency)
	}

	if r.Reference == uuid.Nil {
		return errors.New("reference is required")
	}

	return nil
}

func (r *CreatePaymentRequest) ToPayment() Payment {
	return Payment{
		Reference: r.Reference,
		Amount:    r.Amount,
		Currency:  r.Currency,
		Status:    PENDING,
		CreatedAt: time.Now(),
	}
}

type CreatePaymentResponse struct {
	ID     uuid.UUID     `json:"id"`
	Status PaymentStatus `json:"status"`
}

type GetPaymentDetailsResponse struct {
	ID        uuid.UUID       `json:"id"`
	Amount    decimal.Decimal `json:"amount"`
	Currency  PaymentCurrency `json:"currency"`
	Reference uuid.UUID       `json:"reference"`
	Status    PaymentStatus   `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
}
