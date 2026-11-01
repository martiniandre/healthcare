package apperrors

type DetailedError struct {
	Base   AppError
	Detail string
}

func WithDetail(base AppError, detail string) error {
	return DetailedError{Base: base, Detail: detail}
}

func (detailedError DetailedError) Error() string {
	if detailedError.Detail == "" {
		return detailedError.Base.Message
	}
	return detailedError.Base.Message + ": " + detailedError.Detail
}

func (detailedError DetailedError) Unwrap() error {
	return detailedError.Base
}

func (detailedError DetailedError) As(target any) bool {
	appErrorTarget, isAppErrorTarget := target.(*AppError)
	if !isAppErrorTarget {
		return false
	}
	enrichedError := detailedError.Base
	enrichedError.Message = detailedError.Error()
	*appErrorTarget = enrichedError
	return true
}
