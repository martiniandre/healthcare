package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (recorder *statusRecorder) WriteHeader(code int) {
	recorder.statusCode = code
	recorder.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestStartTime := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: httpResponseWriter,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, httpRequest)

		slog.Info("http request",
			"method", httpRequest.Method,
			"path", httpRequest.URL.Path,
			"status", recorder.statusCode,
			"duration_ms", time.Since(requestStartTime).Milliseconds(),
			"request_id", GetRequestID(httpRequest.Context()),
		)
	})
}
