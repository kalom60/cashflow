# Build stage
FROM --platform=$BUILDPLATFORM golang:1.25.0-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY . .
RUN go mod download

# Build the binary for the target architecture
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o cashflow ./cmd/main.go

# Development stage
FROM golang:1.25.0-alpine AS dev
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY . .
RUN go mod download
CMD ["air", "-c", ".air.toml"]

# Runtime stage
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/cashflow .
COPY config/ /app/config/
COPY internal/constant/query/schemas/ /app/internal/constant/query/schemas/
RUN chmod +x ./cashflow
CMD ["./cashflow"]
