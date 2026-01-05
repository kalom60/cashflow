# Cashflow - Payment Gateway Module

A robust Payment Gateway Module built with Go, PostgreSQL, and RabbitMQ. Features include asynchronous payment processing, idempotency, row-level locking, and a scalable worker pool.

## Features

- **Asynchronous Processing**: Payments are processed in the background via RabbitMQ.
- **Idempotency**: Prevents duplicate processing of the same payment.
- **Concurrency Safety**: Uses PostgreSQL `SELECT ... FOR UPDATE` for row-level locking.
- **Scalable Worker Pool**: Configurable worker goroutines for high throughput.
- **Swagger Documentation**: Interactive API documentation.
- **Structured Error Handling**: Multi-level error responses with clear, concise messages.

## Prerequisites

- [Go 1.25+](https://golang.org/dl/)
- [Docker](https://www.docker.com/products/docker-desktop) & [Docker Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.gnu.org/software/make/)

## Getting Started

### 1. Start Infrastructure

The easiest way to run the entire stack (API, DB, RabbitMQ) is using Docker Compose:

```bash
make up
```

This will:

- Build the Go application in development mode (with hot reloading via `Air`).
- Start PostgreSQL on port `5432`.
- Start RabbitMQ with Management UI on port `15672` (Guest/Guest).
- Start the API server on port `8282`.

### 2. View API Documentation

Once started, you can access the Swagger UI at:
[http://localhost:8282/swagger/index.html](http://localhost:8282/swagger/index.html)

## Development

### Useful Commands

- **Stop all services**: `make down`
- **Run migrations**: `make migrate-up`
- **Rebuild Swagger**: `make swagger`
- **Generate SQLC code**: `make sqlc`
- **Run local dev server (no Docker)**: `make air` (Requires local PG/RabbitMQ)

### Configuration

Configuration is managed in `config/config.yaml`. Key settings:

- `app.worker_count`: Number of concurrent workers (if using internal pool).
- `app.interval`: Duration (15s) defining how often the worker checks for pending events.
- `workerpool.max_workers`: Max workers for the centralized pool.
- `rabbitmq.url`: RabbitMQ connection string.

## Testing

### Automated Tests

To run the test suite:

1. **Start test environment**:
   ```bash
   make test-env
   ```
2. **Run tests**:
   ```bash
   make test path=./internal/storage/...
   ```

### Manual Verification Examples

#### Create a Payment

```bash
curl -X POST http://localhost:8282/payment \
  -H "Content-Type: application/json" \
  -d '{
    "reference": "unique-UUID-ref",
    "amount": 100.50,
    "currency": "USD"
  }'
```

#### Get Payment Status

```bash
curl http://localhost:8282/payment/{PAYMENT_ID}
```

#### Duplicate Reference Error

Try sending the same `reference` twice. You should receive a `400 Bad Request`:

```json
{
  "message": "reference should be unique"
}
```

## Architecture

- **initiator/**: App entry point and dependency injection.
- **internal/handler/**: HTTP controllers (Echo framework).
- **internal/module/**: Business logic and workers.
- **internal/storage/**: Data persistence (PGX + SQLC).
- **platform/**: Shared utilities (Logger, Messaging, WorkerPool).
