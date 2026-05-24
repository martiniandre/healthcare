package render

import (
	"encoding/json"
	"net/http"
)

func JSON(httpResponseWriter http.ResponseWriter, statusCode int, payload interface{}) {
	httpResponseWriter.Header().Set("Content-Type", "application/json")
	httpResponseWriter.WriteHeader(statusCode)
	if payload != nil {
		json.NewEncoder(httpResponseWriter).Encode(payload)
	}
}

func Error(httpResponseWriter http.ResponseWriter, statusCode int, message string) {
	JSON(httpResponseWriter, statusCode, map[string]string{"error": message})
}
