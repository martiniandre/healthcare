FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/healthcare-api ./cmd/api/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/healthcare-api /app/healthcare-api

USER appuser

EXPOSE 50051 8080

CMD ["/app/healthcare-api"]
