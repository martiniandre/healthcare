package apperrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithDetailKeepsIdentityAndExposesDetailThroughAs(testingInstance *testing.T) {
	wrappedError := fmt.Errorf("upload failed: %w", WithDetail(ErrInvalidDICOM, "transfer syntax is not supported"))

	assert.ErrorIs(testingInstance, wrappedError, ErrInvalidDICOM)

	var appError AppError
	assert.True(testingInstance, errors.As(wrappedError, &appError))
	assert.Equal(testingInstance, "INVALID_DICOM", appError.Code)
	assert.Equal(testingInstance, ErrInvalidDICOM.HTTPCode, appError.HTTPCode)
	assert.Contains(testingInstance, appError.Message, "transfer syntax is not supported")
	assert.Contains(testingInstance, appError.Message, ErrInvalidDICOM.Message)
}

func TestWithDetailEmptyDetailReturnsBaseMessage(testingInstance *testing.T) {
	assert.Equal(testingInstance, ErrInvalidDICOM.Message, WithDetail(ErrInvalidDICOM, "").Error())
}

func TestWithDetailMapsToGRPCWithDetailMessage(testingInstance *testing.T) {
	grpcError := ToGRPCStatus(WithDetail(ErrInvalidDICOM, "missing study uid"))

	assert.Contains(testingInstance, grpcError.Error(), "missing study uid")
}
