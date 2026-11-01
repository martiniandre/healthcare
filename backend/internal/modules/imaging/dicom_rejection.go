package imaging

import (
	"errors"
	"fmt"

	"github.com/healthcare/backend/internal/shared/apperrors"
)

type instanceRejection struct {
	code   string
	reason string
}

func newInstanceRejection(code string, reason string) *instanceRejection {
	return &instanceRejection{code: code, reason: reason}
}

func (rejection *instanceRejection) Error() string {
	return rejection.reason
}

func unsupportedTransferSyntaxRejection(description string) error {
	return newInstanceRejection(RejectionCodeUnsupportedTransferSyntax, fmt.Sprintf(
		"transfer syntax %s is not supported; send an uncompressed Implicit or Explicit VR Little Endian file",
		description,
	))
}

func singleFileUploadError(failure error) error {
	var rejection *instanceRejection
	if errors.As(failure, &rejection) {
		return apperrors.WithDetail(apperrors.ErrInvalidDICOM, rejection.reason)
	}
	return failure
}

func rejectedInstanceFromError(fileName string, failure error) RejectedInstance {
	var rejection *instanceRejection
	if errors.As(failure, &rejection) {
		return RejectedInstance{FileName: fileName, Code: rejection.code, Reason: rejection.reason}
	}
	return RejectedInstance{FileName: fileName, Code: RejectionCodeInvalidDICOM, Reason: failure.Error()}
}
