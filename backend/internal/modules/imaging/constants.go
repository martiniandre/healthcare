package imaging

const (
	ImagingStudyStatusReceiving  = "RECEIVING"
	ImagingStudyStatusProcessing = "PROCESSING"
	ImagingStudyStatusReady      = "READY"
	ImagingStudyStatusPartial    = "PARTIAL"
	ImagingStudyStatusFailed     = "FAILED"
)

const (
	WarningCodePatientIdentityMissing = "patient_identity_missing"
	WarningCodeAnonymizedStudy        = "anonymized_study"
	WarningCodeModalityOverridden     = "modality_overridden"
)

const warningMessageKeyPrefix = "imaging.upload.warning."
