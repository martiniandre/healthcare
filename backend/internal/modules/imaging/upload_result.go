package imaging

type UploadWarning struct {
	Code       string
	MessageKey string
}

func newUploadWarning(code string) UploadWarning {
	return UploadWarning{Code: code, MessageKey: warningMessageKeyPrefix + code}
}

type RejectedInstance struct {
	FileName string
	Code     string
	Reason   string
}

type UploadDICOMResult struct {
	Study             *ImagingStudy
	StudyCreated      bool
	SeriesInstanceUID string
	SOPInstanceUID    string
	Warnings          []UploadWarning
	InstancesAccepted int
	InstancesRejected []RejectedInstance
}
