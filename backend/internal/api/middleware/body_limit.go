package middleware

import (
	"net/http"
	"strings"
)

const DefaultMaxJSONBodyBytes int64 = 1 << 20

func isStreamingUpload(httpRequest *http.Request) bool {
	contentType := strings.ToLower(httpRequest.Header.Get("Content-Type"))
	return strings.HasPrefix(contentType, "multipart/")
}

func BodyLimit(maxBodyBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
			if httpRequest.Body != nil && httpRequest.Body != http.NoBody && !isStreamingUpload(httpRequest) {
				if httpRequest.ContentLength > maxBodyBytes {
					http.Error(httpResponseWriter, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
					return
				}
				httpRequest.Body = http.MaxBytesReader(httpResponseWriter, httpRequest.Body, maxBodyBytes)
			}
			next.ServeHTTP(httpResponseWriter, httpRequest)
		})
	}
}
