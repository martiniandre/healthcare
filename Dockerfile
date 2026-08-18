FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/api /app/api
COPY --from=builder /app/migrations /app/migrations

RUN printf '#!/bin/sh\nif [ -n "$GCP_KEY_B64" ]; then\n  mkdir -p /tmp/gcp\n  echo "$GCP_KEY_B64" | base64 -d > /tmp/gcp/key.json\n  export GOOGLE_APPLICATION_CREDENTIALS=/tmp/gcp/key.json\nfi\nexec /app/api\n' > /app/entrypoint.sh && chmod +x /app/entrypoint.sh

USER appuser

EXPOSE 50051 8080

ENTRYPOINT ["/app/entrypoint.sh"]
