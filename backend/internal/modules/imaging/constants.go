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

const (
	RejectionCodeDuplicateInstance         = "DICOM_DUPLICATE_INSTANCE"
	RejectionCodeInvalidDICOM              = "INVALID_DICOM"
	RejectionCodeUnsupportedTransferSyntax = "DICOM_UNSUPPORTED_TRANSFER_SYNTAX"
	RejectionCodeMissingIdentifiers        = "DICOM_MISSING_IDENTIFIERS"
)

const (
	MaxZIPBytes                 int64 = 2 << 30
	MaxZIPEntries                     = 2000
	MaxZIPUncompressedBytes     int64 = 4 << 30
	MaxZIPCompressionRatio      int64 = 200
	ZIPCompressionRatioMinBytes int64 = 1 << 20
)

const warningMessageKeyPrefix = "imaging.upload.warning."
