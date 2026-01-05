migrate-down:
	- migrate -database postgres://postgres:secret@localhost:5432/cashflow?sslmode=disable -path internal/constant/query/schemas -verbose down $(N)
migrate-up:
	- migrate -database postgres://postgres:secret@localhost:5432/cashflow?sslmode=disable -path internal/constant/query/schemas -verbose up
migrate-down-test:
	- migrate -database postgres://postgres:secret@localhost:5433/cashflow?sslmode=disable -path internal/constant/query/schemas -verbose down
migrate-up-test:
	- migrate -database postgres://postgres:secret@localhost:5433/cashflow?sslmode=disable -path internal/constant/query/schemas -verbose up
migrate-create:
	- migrate create -ext sql -dir internal/constant/query/schemas -tz "UTC" $(name)

swagger:
	-swag fmt && swag init -g cmd/main.go

run:
	go run cmd/main.go

sqlc:
	cd ./config && sqlc generate

air:
	@echo "Running air..."
	air -c .air.toml

up:
	@echo "Starting Docker images..."
	docker-compose -f docker-compose.yaml up --build -d
	@echo "Docker images started!"
down:
	@echo "Stopping docker compose..."
	docker-compose -f docker-compose.yaml down
	@echo "Done!"

test-env:
	docker-compose --profile test up -d test_db

test:
	go test -v $(path) | grep -v '"level"' | grep -v 'Error #'
