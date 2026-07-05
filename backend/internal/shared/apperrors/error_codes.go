package apperrors

import (
	"errors"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError struct {
	GRPCCode codes.Code
	HTTPCode int
	Message  string
}

func (appError AppError) Error() string {
	return appError.Message
}

func (appError AppError) ToGRPC() error {
	return status.Error(appError.GRPCCode, appError.Message)
}

func (appError AppError) WithFields(fields map[string]string) error {
	if len(fields) == 0 {
		return appError.ToGRPC()
	}

	var details []string
	for key, value := range fields {
		details = append(details, key+" "+value)
	}

	return status.Error(appError.GRPCCode, appError.Message+" — "+strings.Join(details, ", "))
}

var (
	ErrBadRequest = AppError{
		GRPCCode: codes.InvalidArgument,
		HTTPCode: http.StatusBadRequest,
		Message:  "invalid request parameters",
	}

	ErrInvalidCredentials = AppError{
		GRPCCode: codes.Unauthenticated,
		HTTPCode: http.StatusUnauthorized,
		Message:  "invalid credentials",
	}

	ErrMissingToken = AppError{
		GRPCCode: codes.Unauthenticated,
		HTTPCode: http.StatusUnauthorized,
		Message:  "missing authentication token",
	}

	ErrTokenExpired = AppError{
		GRPCCode: codes.Unauthenticated,
		HTTPCode: http.StatusUnauthorized,
		Message:  "token expired or invalid",
	}

	ErrCSRFTokenMismatch = AppError{
		GRPCCode: codes.Unauthenticated,
		HTTPCode: http.StatusUnauthorized,
		Message:  "csrf token mismatch",
	}

	ErrPermissionDenied = AppError{
		GRPCCode: codes.PermissionDenied,
		HTTPCode: http.StatusForbidden,
		Message:  "permission denied",
	}

	ErrMethodNotInPermissionMatrix = AppError{
		GRPCCode: codes.PermissionDenied,
		HTTPCode: http.StatusForbidden,
		Message:  "method not registered in permission matrix",
	}

	ErrUserNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "user not found",
	}

	ErrUserAlreadyExists = AppError{
		GRPCCode: codes.AlreadyExists,
		HTTPCode: http.StatusConflict,
		Message:  "user already exists",
	}

	ErrEmployeeNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "employee not found",
	}

	ErrPatientNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "patient not found",
	}

	ErrPatientAlreadyExists = AppError{
		GRPCCode: codes.AlreadyExists,
		HTTPCode: http.StatusConflict,
		Message:  "patient already exists",
	}

	ErrEncounterNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "encounter not found",
	}

	ErrObservationNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "observation not found",
	}

	ErrConditionNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "condition not found",
	}

	ErrAllergyIntoleranceNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "allergy intolerance not found",
	}

	ErrMedicationRequestNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "medication request not found",
	}

	ErrDiagnosticReportNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "diagnostic report not found",
	}

	ErrImagingStudyNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "imaging study not found",
	}

	ErrInvalidDICOM = AppError{
		GRPCCode: codes.InvalidArgument,
		HTTPCode: http.StatusBadRequest,
		Message:  "invalid dicom file structure or preamble signature",
	}

	ErrRateLimitExceeded = AppError{
		GRPCCode: codes.ResourceExhausted,
		HTTPCode: http.StatusTooManyRequests,
		Message:  "rate limit exceeded, try again later",
	}

	ErrRequestTimeout = AppError{
		GRPCCode: codes.DeadlineExceeded,
		HTTPCode: http.StatusGatewayTimeout,
		Message:  "request timeout exceeded",
	}

	ErrInternalServer = AppError{
		GRPCCode: codes.Internal,
		HTTPCode: http.StatusInternalServerError,
		Message:  "internal server error",
	}

	ErrServiceUnavailable = AppError{
		GRPCCode: codes.Unavailable,
		HTTPCode: http.StatusServiceUnavailable,
		Message:  "service unavailable",
	}

	ErrNotImplemented = AppError{
		GRPCCode: codes.Unimplemented,
		HTTPCode: http.StatusNotImplemented,
		Message:  "method not implemented",
	}

	ErrInvalidPasscode = AppError{
		GRPCCode: codes.PermissionDenied,
		HTTPCode: http.StatusForbidden,
		Message:  "invalid passcode for telemetry room",
	}

	ErrRoomNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "telemetry room not found",
	}

	ErrBedNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "telemetry bed not found",
	}

	ErrNotificationNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "notification not found",
	}

	ErrNotFound = AppError{
		GRPCCode: codes.NotFound,
		HTTPCode: http.StatusNotFound,
		Message:  "resource not found",
	}
)

func ToGRPCStatus(err error) error {
	if err == nil {
		return nil
	}

	var appError AppError
	if errors.As(err, &appError) {
		return appError.ToGRPC()
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	return status.Error(codes.Internal, err.Error())
}
