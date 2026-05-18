package interceptor

import (
	"context"
	"time"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"google.golang.org/grpc"
)

const defaultRequestTimeout = 30 * time.Second

func UnaryTimeoutInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		timeoutCtx, cancelFunc := context.WithTimeout(ctx, defaultRequestTimeout)
		defer cancelFunc()

		resultChannel := make(chan struct {
			response interface{}
			err      error
		}, 1)

		go func() {
			response, err := handler(timeoutCtx, req)
			resultChannel <- struct {
				response interface{}
				err      error
			}{response, err}
		}()

		select {
		case result := <-resultChannel:
			return result.response, result.err
		case <-timeoutCtx.Done():
			return nil, apperrors.ErrRequestTimeout.ToGRPC()
		}
	}
}
