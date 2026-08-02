package render

import (
	"errors"
	"net/http"

	"github.com/healthcare/backend/internal/shared/apperrors"
)

func ErrorFromAppError(httpResponseWriter http.ResponseWriter, err error) {
	var appError apperrors.AppError
	if errors.As(err, &appError) {
		Error(httpResponseWriter, appError.HTTPCode, appError.Message)
		return
	}
	Error(httpResponseWriter, http.StatusInternalServerError, apperrors.ErrInternalServer.Message)
}
