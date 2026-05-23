package imaging

import "errors"

var (
	ErrImagingStudyNotFound = errors.New("imaging study not found")
	ErrInvalidDICOM         = errors.New("invalid dicom file")
	ErrDICOMTooLarge        = errors.New("dicom file exceeds maximum allowed size")
)

const (
	ImagingStudyStatusPending    = "PENDING"
	ImagingStudyStatusProcessing = "PROCESSING"
	ImagingStudyStatusProcessed  = "PROCESSED"
	ImagingStudyStatusFailed     = "FAILED"
)
