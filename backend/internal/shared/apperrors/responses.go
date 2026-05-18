package apperrors

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func BadRequest(message string, fieldViolations map[string]string) error {
	if len(fieldViolations) == 0 {
		return status.Error(codes.InvalidArgument, message)
	}

	violations := make([]string, 0, len(fieldViolations))
	for field, description := range fieldViolations {
		violations = append(violations, fmt.Sprintf("%s: %s", field, description))
	}

	return status.Errorf(codes.InvalidArgument, "%s — %s", message, strings.Join(violations, ", "))
}

func NotFound(resourceType string) error {
	return status.Errorf(codes.NotFound, "%s not found", resourceType)
}

func AlreadyExists(resourceType string) error {
	return status.Errorf(codes.AlreadyExists, "%s already exists", resourceType)
}

func Unauthenticated(reason string) error {
	return status.Errorf(codes.Unauthenticated, "unauthenticated: %s", reason)
}

func PermissionDenied(reason string) error {
	return status.Errorf(codes.PermissionDenied, "permission denied: %s", reason)
}

func InternalError(reason string) error {
	return status.Errorf(codes.Internal, "internal error: %s", reason)
}

func Unavailable(reason string) error {
	return status.Errorf(codes.Unavailable, "service unavailable: %s", reason)
}

func Unimplemented(methodName string) error {
	return status.Errorf(codes.Unimplemented, "%s is not implemented", methodName)
}

func ResourceExhausted(reason string) error {
	return status.Errorf(codes.ResourceExhausted, "resource exhausted: %s", reason)
}

func DeadlineExceeded(reason string) error {
	return status.Errorf(codes.DeadlineExceeded, "deadline exceeded: %s", reason)
}
