package api

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/healthcare/backend/cmd/api/docs"
	"github.com/healthcare/backend/internal/api/middleware"
)

type RouteRegisterer interface {
	RegisterRoutes(mux *http.ServeMux)
}

func NewRouter(secureCookies bool, swaggerEnabled bool, registerers ...RouteRegisterer) http.Handler {
	httpServeMux := http.NewServeMux()

	if swaggerEnabled {
		httpServeMux.Handle("GET /swagger/", httpSwagger.Handler())
	}

	for _, registerer := range registerers {
		registerer.RegisterRoutes(httpServeMux)
	}

	handlerPipeline := middleware.CORS(secureCookies)(httpServeMux)
	handlerPipeline = middleware.RateLimit(handlerPipeline)
	handlerPipeline = middleware.APIPrefixRewrite(handlerPipeline)
	handlerPipeline = middleware.Recovery(handlerPipeline)
	handlerPipeline = middleware.RequestID(handlerPipeline)
	handlerPipeline = middleware.AuditTrail(handlerPipeline)
	handlerPipeline = middleware.Logger(handlerPipeline)

	return handlerPipeline
}
