package payment_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kalom60/cashflow/internal/constant/dto"
	"github.com/kalom60/cashflow/internal/storage"
	"github.com/kalom60/cashflow/internal/storage/payment"
	"github.com/kalom60/cashflow/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var (
	ctx   context.Context
	store storage.Payment

	paymentIDETB uuid.UUID
	paymentIDUSD uuid.UUID
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	testDB := testutils.SetupTestDB()
	store = payment.Init(testutils.NewTestLogger(), &testDB)

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

	resp, err := store.CreatePayment(ctx, req)
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

	resp, err := store.CreatePayment(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Reference, resp.Reference)
	assert.Equal(t, dto.USD, resp.Currency)
	assert.Equal(t, dto.PENDING, resp.Status)
	assert.Equal(t, req.Amount, resp.Amount)

	paymentIDUSD = resp.ID
}

func TestGetPaymentByIDETB(t *testing.T) {
	resp, err := store.GetPaymentByID(ctx, paymentIDETB)
	assert.NoError(t, err)
	assert.Equal(t, dto.ETB, resp.Currency)
	assert.Equal(t, dto.PENDING, resp.Status)
}

func TestGetPaymentByIDUSD(t *testing.T) {
	resp, err := store.GetPaymentByID(ctx, paymentIDUSD)
	assert.NoError(t, err)
	assert.Equal(t, dto.USD, resp.Currency)
	assert.Equal(t, dto.PENDING, resp.Status)
}

func TestGetPaymentByNewID(t *testing.T) {
	_, err := store.GetPaymentByID(ctx, uuid.New())
	assert.Error(t, err)
}

func TestUpdatePaymentETBToSuccess(t *testing.T) {
	err := store.UpdatePaymentStatus(ctx, paymentIDETB, dto.SUCCESS)
	assert.NoError(t, err)
}

func TestGetPaymentByIDETBAfterSuccessUpdate(t *testing.T) {
	resp, err := store.GetPaymentByID(ctx, paymentIDETB)
	assert.NoError(t, err)
	assert.Equal(t, dto.ETB, resp.Currency)
	assert.Equal(t, dto.SUCCESS, resp.Status)
}

func TestUpdatePaymentUSDToSuccess(t *testing.T) {
	err := store.UpdatePaymentStatus(ctx, paymentIDUSD, dto.SUCCESS)
	assert.NoError(t, err)
}

func TestGetPaymentByIDUSDAfterSuccessUpdate(t *testing.T) {
	resp, err := store.GetPaymentByID(ctx, paymentIDUSD)
	assert.NoError(t, err)
	assert.Equal(t, dto.USD, resp.Currency)
	assert.Equal(t, dto.SUCCESS, resp.Status)
}

func TestUpdatePaymentETBToFailed(t *testing.T) {
	err := store.UpdatePaymentStatus(ctx, paymentIDETB, dto.FAILED)
	assert.NoError(t, err)
}

func TestGetPaymentByIDETBAfterFailedUpdate(t *testing.T) {
	resp, err := store.GetPaymentByID(ctx, paymentIDETB)
	assert.NoError(t, err)
	assert.Equal(t, dto.ETB, resp.Currency)
	assert.Equal(t, dto.FAILED, resp.Status)
}

func TestUpdatePaymentUSDToFailed(t *testing.T) {
	err := store.UpdatePaymentStatus(ctx, paymentIDUSD, dto.FAILED)
	assert.NoError(t, err)
}

func TestGetPaymentByIDUSDAfterFailedUpdate(t *testing.T) {
	resp, err := store.GetPaymentByID(ctx, paymentIDUSD)
	assert.NoError(t, err)
	assert.Equal(t, dto.USD, resp.Currency)
	assert.Equal(t, dto.FAILED, resp.Status)
}
