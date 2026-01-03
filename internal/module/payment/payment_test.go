package payment_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kalom60/cashflow/internal/constant/dto"
	"github.com/kalom60/cashflow/internal/module"
	paymentModule "github.com/kalom60/cashflow/internal/module/payment"
	"github.com/kalom60/cashflow/internal/storage"
	paymentStorage "github.com/kalom60/cashflow/internal/storage/payment"
	"github.com/kalom60/cashflow/platform/logger"
	"github.com/kalom60/cashflow/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var (
	ctx     context.Context
	store   storage.Payment
	log     logger.Logger
	pModule module.Payment

	paymentIDETB uuid.UUID
	paymentIDUSD uuid.UUID
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	testDB := testutils.SetupTestDB()
	log = testutils.NewTestLogger()
	store = paymentStorage.Init(log, &testDB)
	pModule = paymentModule.Init(log, store)

	code := m.Run()
	os.Exit(code)
}

func TestCreatePaymentETB(t *testing.T) {
	req := dto.Payment{
		Reference: uuid.New(),
		Amount:    decimal.NewFromFloat(100000),
		Currency:  dto.ETB,
		Status:    dto.PENDING,
		CreatedAt: time.Now(),
	}

	resp, err := pModule.CreatePayment(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Reference, resp.Reference)
	assert.Equal(t, dto.ETB, resp.Currency)
	assert.Equal(t, dto.PENDING, resp.Status)
	assert.Equal(t, req.Amount, resp.Amount)

	paymentIDETB = resp.ID
}

func TestCreatePaymentUSD(t *testing.T) {
	req := dto.Payment{
		Reference: uuid.New(),
		Amount:    decimal.NewFromFloat(1000000),
		Currency:  dto.USD,
		Status:    dto.PENDING,
		CreatedAt: time.Now(),
	}

	resp, err := pModule.CreatePayment(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Reference, resp.Reference)
	assert.Equal(t, dto.USD, resp.Currency)
	assert.Equal(t, dto.PENDING, resp.Status)
	assert.Equal(t, req.Amount, resp.Amount)

	paymentIDUSD = resp.ID
}

func TestGetPaymentByIDETB(t *testing.T) {
	resp, err := pModule.GetPaymentByID(ctx, paymentIDETB)
	assert.NoError(t, err)
	assert.Equal(t, dto.ETB, resp.Currency)
	assert.Equal(t, dto.PENDING, resp.Status)
}

func TestGetPaymentByIDUSD(t *testing.T) {
	resp, err := pModule.GetPaymentByID(ctx, paymentIDUSD)
	assert.NoError(t, err)
	assert.Equal(t, dto.USD, resp.Currency)
	assert.Equal(t, dto.PENDING, resp.Status)
}

func TestGetPaymentByNewID(t *testing.T) {
	_, err := pModule.GetPaymentByID(ctx, uuid.New())
	assert.Error(t, err)
}
