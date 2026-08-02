package healthcare

import (
	"errors"
	"fmt"
)

type NotFoundError struct {
	ResourceType string
	ResourceID   string
}

func (notFoundError *NotFoundError) Error() string {
	if notFoundError.ResourceID == "" {
		return fmt.Sprintf("resource %s not found", notFoundError.ResourceType)
	}
	return fmt.Sprintf("resource %s/%s not found", notFoundError.ResourceType, notFoundError.ResourceID)
}

func IsNotFound(err error) bool {
	var notFoundError *NotFoundError
	return errors.As(err, &notFoundError)
}
