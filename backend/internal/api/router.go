package api

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/healthcare/backend/internal/api/middleware"
	_ "github.com/healthcare/backend/cmd/api/docs"
)

type RouteRegisterer interface {
	RegisterRoutes(mux *http.ServeMux)
}

func NewRouter(secureCookies bool, registerers ...RouteRegisterer) http.Handler {
	httpServeMux := http.NewServeMux()

	httpServeMux.Handle("GET /swagger/", httpSwagger.Handler())

	for _, registerer := range registerers {
		registerer.RegisterRoutes(httpServeMux)
	}

	handlerPipeline := middleware.CORS(secureCookies)(httpServeMux)
	handlerPipeline = middleware.APIPrefixRewrite(handlerPipeline)
	handlerPipeline = middleware.Recovery(handlerPipeline)
	handlerPipeline = middleware.RequestID(handlerPipeline)
	handlerPipeline = middleware.Logger(handlerPipeline)

	return handlerPipeline
}
