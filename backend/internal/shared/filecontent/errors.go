package filecontent

import "errors"

var (
	ErrEmptyFile                = errors.New("file is empty")
	ErrContentExtensionMismatch = errors.New("file content does not match its extension")
	ErrCorruptedFile            = errors.New("file content is corrupted or incomplete")
	ErrFileTooLarge             = errors.New("file exceeds the maximum processable size")
)
