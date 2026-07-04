package health

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HTTPHandler struct {
	databaseConnectionPool *pgxpool.Pool
	cacheConnectionClient  *redis.Client
}

func NewHTTPHandler(databaseConnectionPool *pgxpool.Pool, cacheConnectionClient *redis.Client) *HTTPHandler {
	return &HTTPHandler{
		databaseConnectionPool: databaseConnectionPool,
		cacheConnectionClient:  cacheConnectionClient,
	}
}

func (healthHTTPHandler *HTTPHandler) RegisterRoutes(httpServeMux *http.ServeMux) {
	httpServeMux.Handle("GET /health", http.HandlerFunc(healthHTTPHandler.HealthCheck))
}

type healthCheckResult struct {
	Status  string `json:"status"`
	Service string `json:"service,omitempty"`
	Error   string `json:"error,omitempty"`
}

type healthCheckResponse struct {
	Status string             `json:"status"`
	Checks []healthCheckResult `json:"checks"`
}

func (healthHTTPHandler *HTTPHandler) HealthCheck(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	healthResponse := healthCheckResponse{
		Status: "ok",
		Checks: make([]healthCheckResult, 0),
	}

	overallUnhealthy := false

	if healthHTTPHandler.databaseConnectionPool != nil {
		databasePingError := healthHTTPHandler.databaseConnectionPool.Ping(httpRequest.Context())
		if databasePingError != nil {
			overallUnhealthy = true
			healthResponse.Checks = append(healthResponse.Checks, healthCheckResult{
				Status:  "unhealthy",
				Service: "database",
				Error:   databasePingError.Error(),
			})
		} else {
			healthResponse.Checks = append(healthResponse.Checks, healthCheckResult{
				Status:  "ok",
				Service: "database",
			})
		}
	}

	if healthHTTPHandler.cacheConnectionClient != nil {
		cachePingError := healthHTTPHandler.cacheConnectionClient.Ping(context.Background()).Err()
		if cachePingError != nil {
			overallUnhealthy = true
			healthResponse.Checks = append(healthResponse.Checks, healthCheckResult{
				Status:  "unhealthy",
				Service: "cache",
				Error:   cachePingError.Error(),
			})
		} else {
			healthResponse.Checks = append(healthResponse.Checks, healthCheckResult{
				Status:  "ok",
				Service: "cache",
			})
		}
	}

	if overallUnhealthy {
		healthResponse.Status = "unhealthy"
	}

	httpResponseWriter.Header().Set("Content-Type", "application/json")

	if overallUnhealthy {
		httpResponseWriter.WriteHeader(http.StatusServiceUnavailable)
	} else {
		httpResponseWriter.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(httpResponseWriter).Encode(healthResponse)
}
