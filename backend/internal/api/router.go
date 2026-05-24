package api

import (
	"net/http"

	"github.com/healthcare/backend/internal/api/middleware"
)

type RouteRegisterer interface {
	RegisterRoutes(mux *http.ServeMux)
}

func NewRouter(secureCookies bool, registerers ...RouteRegisterer) http.Handler {
	httpServeMux := http.NewServeMux()

	for _, registerer := range registerers {
		registerer.RegisterRoutes(httpServeMux)
	}

	handlerPipeline := middleware.CORS(secureCookies)(httpServeMux)
	handlerPipeline = middleware.Recovery(handlerPipeline)
	handlerPipeline = middleware.RequestID(handlerPipeline)
	handlerPipeline = middleware.Logger(handlerPipeline)

	return handlerPipeline
}
