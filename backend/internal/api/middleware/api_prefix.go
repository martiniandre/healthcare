package middleware

import (
	"net/http"
	"strings"
)

func APIPrefixRewrite(next http.Handler) http.Handler {
	return http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		if strings.HasPrefix(httpRequest.URL.Path, "/api/") && !strings.HasPrefix(httpRequest.URL.Path, "/api/v1/") {
			httpRequest.URL.Path = "/api/v1/" + strings.TrimPrefix(httpRequest.URL.Path, "/api/")
		}
		next.ServeHTTP(httpResponseWriter, httpRequest)
	})
}
