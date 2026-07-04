package interceptor

import (
	"context"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const RequestIDMetadataKey = "x-request-id"

func UnaryRequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		requestID := extractRequestIDFromMetadata(ctx)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, ctxkeys.RequestIDKey, requestID)
		ctx = context.WithValue(ctx, ctxkeys.CorrelationIDKey, requestID)
		grpc.SetHeader(ctx, metadata.Pairs(RequestIDMetadataKey, requestID))

		return handler(ctx, req)
	}
}

func extractRequestIDFromMetadata(ctx context.Context) string {
	incomingMetadata, metadataOk := metadata.FromIncomingContext(ctx)
	if !metadataOk {
		return ""
	}
	values := incomingMetadata.Get(RequestIDMetadataKey)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}
