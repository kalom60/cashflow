package dto

import (
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
