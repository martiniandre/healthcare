# --- Estágio de Compilação ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Instalar dependências necessárias para compilação (se houver CGO)
RUN apk add --no-cache git gcc musl-dev

# Copiar arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código fonte do backend
COPY . .

# Compilar o binário (estático)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/healthcare-api ./cmd/api/main.go

# --- Estágio de Execução ---
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copiar o binário compilado do estágio anterior
COPY --from=builder /app/healthcare-api /app/healthcare-api

# Copiar o arquivo .env se ele existir (ou carregar via docker-compose)
COPY --from=builder /app/.env* /app/

EXPOSE 50051 8080

CMD ["/app/healthcare-api"]
