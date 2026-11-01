package imaging

type UploadWarning struct {
	Code       string
	MessageKey string
}

func newUploadWarning(code string) UploadWarning {
	return UploadWarning{Code: code, MessageKey: warningMessageKeyPrefix + code}
}

type UploadDICOMResult struct {
	Study             *ImagingStudy
	SeriesInstanceUID string
	SOPInstanceUID    string
	Warnings          []UploadWarning
}
