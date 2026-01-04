package outboxevent_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kalom60/cashflow/internal/constant/dto"
	"github.com/kalom60/cashflow/internal/storage"
	outboxevent "github.com/kalom60/cashflow/internal/storage/outbox_event"
	"github.com/kalom60/cashflow/internal/storage/payment"
	"github.com/kalom60/cashflow/tests/testutils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var (
	ctx     context.Context
	pStore  storage.Payment
	oeStore storage.OutboxEvent

	eventIDETB uuid.UUID
	eventIDUSD uuid.UUID
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	testDB := testutils.SetupTestDB()
	pStore = payment.Init(testutils.NewTestLogger(), &testDB)
	oeStore = outboxevent.Init(testutils.NewTestLogger(), &testDB, 100)

	code := m.Run()
	os.Exit(code)
}

func TestCreateOutboxEventETB(t *testing.T) {
	req := dto.Payment{
		Reference: uuid.New(),
		Amount:    decimal.NewFromFloat(100000),
		Currency:  dto.ETB,
		Status:    dto.PENDING,
		CreatedAt: time.Now(),
	}

	resp, err := pStore.CreatePayment(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Reference, resp.Reference)
	assert.Equal(t, dto.ETB, resp.Currency)
	assert.Equal(t, dto.PENDING, resp.Status)
	assert.Equal(t, req.Amount, resp.Amount)
}

func TestCreateOutboxEventUSD(t *testing.T) {
	req := dto.Payment{
		Reference: uuid.New(),
		Amount:    decimal.NewFromFloat(1000000),
		Currency:  dto.USD,
		Status:    dto.PENDING,
		CreatedAt: time.Now(),
	}

	resp, err := pStore.CreatePayment(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Reference, resp.Reference)
	assert.Equal(t, dto.USD, resp.Currency)
	assert.Equal(t, dto.PENDING, resp.Status)
	assert.Equal(t, req.Amount, resp.Amount)
}

func TestGetPendingOutboxEventsForUpdate(t *testing.T) {
	tx, _ := oeStore.BeginTx(ctx)
	resp, err := oeStore.GetPendingOutboxEventsForUpdate(ctx, tx)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp))

	for _, res := range resp {
		var data map[string]any
		if err := json.Unmarshal(res.Payload.Bytes, &data); err != nil {
			continue
		}

		currency, ok := data["currency"].(string)
		if !ok {
			continue
		}

		if currency == string(dto.ETB) {
			eventIDETB = res.ID
			continue
		}
		eventIDUSD = res.ID
	}
}

func TestUpdateOutboxStatusETB(t *testing.T) {
	tx, _ := oeStore.BeginTx(ctx)
	err := oeStore.UpdateOutboxStatus(ctx, tx, eventIDETB, dto.OutboxStatusSent)
	assert.NoError(t, err)
}

func TestUpdateOutboxStatusUSD(t *testing.T) {
	tx, _ := oeStore.BeginTx(ctx)
	err := oeStore.UpdateOutboxStatus(ctx, tx, eventIDUSD, dto.OutboxStatusFailed)
	assert.NoError(t, err)
}

func TestDeleteOutboxEventETB(t *testing.T) {
	tx, _ := oeStore.BeginTx(ctx)
	err := oeStore.DeleteOutboxEvent(ctx, tx, eventIDETB)
	assert.NoError(t, err)
}

func TestDeleteOutboxEventUSD(t *testing.T) {
	tx, _ := oeStore.BeginTx(ctx)
	err := oeStore.DeleteOutboxEvent(ctx, tx, eventIDUSD)
	assert.NoError(t, err)
}

func TestUpdateOutboxStatusETBAfterUpdate(t *testing.T) {
	tx, _ := oeStore.BeginTx(ctx)
	err := oeStore.UpdateOutboxStatus(ctx, tx, eventIDETB, dto.OutboxStatusSent)
	assert.Error(t, err)
}

func TestUpdateOutboxStatusUSDAfterUpdate(t *testing.T) {
	tx, _ := oeStore.BeginTx(ctx)
	err := oeStore.UpdateOutboxStatus(ctx, tx, eventIDUSD, dto.OutboxStatusFailed)
	assert.Error(t, err)
}

func TestGetPendingOutboxEventsForUpdateAfterUpdate(t *testing.T) {
	tx, _ := oeStore.BeginTx(ctx)
	resp, err := oeStore.GetPendingOutboxEventsForUpdate(ctx, tx)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(resp))
}

func TestDeleteOutboxEventETBAfterDelete(t *testing.T) {
	tx, _ := oeStore.BeginTx(ctx)
	err := oeStore.DeleteOutboxEvent(ctx, tx, eventIDETB)
	assert.Error(t, err)
}

func TestDeleteOutboxEventUSDAfterDelete(t *testing.T) {
	tx, _ := oeStore.BeginTx(ctx)
	err := oeStore.DeleteOutboxEvent(ctx, tx, eventIDUSD)
	assert.Error(t, err)
}
