package render

import (
	"errors"
	"net/http"

	"github.com/healthcare/backend/internal/shared/apperrors"
)

func ErrorFromAppError(httpResponseWriter http.ResponseWriter, err error) {
	var appError apperrors.AppError
	if errors.As(err, &appError) {
		if appError.Code != "" {
			JSON(httpResponseWriter, appError.HTTPCode, map[string]string{"error": appError.Message, "code": appError.Code})
			return
		}
		Error(httpResponseWriter, appError.HTTPCode, appError.Message)
		return
	}
	Error(httpResponseWriter, http.StatusInternalServerError, apperrors.ErrInternalServer.Message)
}
