package outboxevent_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kalom60/cashflow/internal/constant/dto"
	"github.com/kalom60/cashflow/internal/module"

	outboxeventWorker "github.com/kalom60/cashflow/internal/module/outbox_event"
	paymentModule "github.com/kalom60/cashflow/internal/module/payment"
	"github.com/kalom60/cashflow/internal/storage"
	outboxeventStorage "github.com/kalom60/cashflow/internal/storage/outbox_event"
	paymentStorage "github.com/kalom60/cashflow/internal/storage/payment"
	"github.com/kalom60/cashflow/platform/logger"
	"github.com/kalom60/cashflow/tests/testutils"
	"github.com/rabbitmq/amqp091-go"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockMessagingClient struct{}

func (m *mockMessagingClient) PublishPayment(ctx context.Context, paymentID string) error {
	return nil
}

func (m *mockMessagingClient) ConsumePayments(ctx context.Context) (<-chan amqp091.Delivery, error) {
	return nil, nil
}

func (m *mockMessagingClient) Close() error {
	return nil
}

var (
	ctx      context.Context
	log      logger.Logger
	pStore   storage.Payment
	pModule  module.Payment
	oeStore  storage.OutboxEvent
	oeWorker *outboxeventWorker.OutboxEventWorker
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	testDB := testutils.SetupTestDB()
	log = testutils.NewTestLogger()

	pStore = paymentStorage.Init(log, &testDB)
	pModule = paymentModule.Init(log, pStore)

	oeStore = outboxeventStorage.Init(log, &testDB, 100)
	oeWorker = outboxeventWorker.Init(log, oeStore, &mockMessagingClient{}, 2*time.Second)

	go oeWorker.Start(ctx)

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
}

func TestGetOutboxEventsForUpdate(t *testing.T) {
	time.Sleep(5 * time.Second)
	tx, _ := oeStore.BeginTx(ctx)
	resp, err := oeStore.GetPendingOutboxEventsForUpdate(ctx, tx)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(resp))
}
