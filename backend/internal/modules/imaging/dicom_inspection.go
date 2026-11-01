package imaging

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/dicom"
)

const (
	dicomHeaderScanLimitBytes = 2 << 20
	maximumUIDLength          = 64
	maximumAccessionLength    = 64
	maximumDescriptionLength  = 255
	maximumTitleLength        = 255
)

var (
	dicomUIDPattern      = regexp.MustCompile(`^[0-9]+(\.[0-9]+)*$`)
	dicomModalityPattern = regexp.MustCompile(`^[A-Z]{2,16}$`)
	dicomTimeDigits      = regexp.MustCompile(`^[0-9]+`)
)

type dicomHeaderInspection struct {
	metadata      *dicom.InstanceMetadata
	consumedBytes []byte
}

type errorRecordingReader struct {
	reader      io.Reader
	recordedErr error
}

func (recording *errorRecordingReader) Read(buffer []byte) (int, error) {
	bytesRead, readError := recording.reader.Read(buffer)
	if readError != nil && !errors.Is(readError, io.EOF) {
		recording.recordedErr = readError
	}
	return bytesRead, readError
}

func inspectDICOMHeader(headerBytes []byte, remainingStream io.Reader) (*dicomHeaderInspection, error) {
	recordingStream := &errorRecordingReader{reader: io.LimitReader(remainingStream, dicomHeaderScanLimitBytes)}
	consumedBuffer := &bytes.Buffer{}
	scanSource := io.TeeReader(io.MultiReader(bytes.NewReader(headerBytes), recordingStream), consumedBuffer)

	metadata, parseError := dicom.ParseInstance(scanSource)
	if recordingStream.recordedErr != nil {
		return nil, fmt.Errorf("failed to read dicom header: %w", recordingStream.recordedErr)
	}
	if parseError != nil {
		return nil, describeParseFailure(parseError)
	}

	return &dicomHeaderInspection{metadata: metadata, consumedBytes: consumedBuffer.Bytes()}, nil
}

func describeParseFailure(parseError error) error {
	switch {
	case errors.Is(parseError, dicom.ErrNotDICOM):
		return apperrors.WithDetail(apperrors.ErrInvalidDICOM, "file is not a DICOM Part 10 file")
	case errors.Is(parseError, dicom.ErrTransferSyntaxAbsent):
		return apperrors.WithDetail(apperrors.ErrInvalidDICOM, "file meta group has no transfer syntax")
	case errors.Is(parseError, dicom.ErrMetaGroupMalformed):
		return apperrors.WithDetail(apperrors.ErrInvalidDICOM, "file meta group is malformed")
	}
	return apperrors.WithDetail(apperrors.ErrInvalidDICOM, "file could not be parsed")
}

type validatedDICOMMetadata struct {
	studyInstanceUID  string
	seriesInstanceUID string
	sopInstanceUID    string
	modality          string
	studyDescription  string
	accessionNumber   string
	studyDateTime     *time.Time
	instanceIndex     *InstanceIndex
	source            *dicom.InstanceMetadata
}

func validateDICOMMetadata(metadata *dicom.InstanceMetadata) (*validatedDICOMMetadata, error) {
	if !dicom.IsSupportedTransferSyntax(metadata.TransferSyntaxUID) {
		return nil, apperrors.WithDetail(apperrors.ErrInvalidDICOM, fmt.Sprintf(
			"transfer syntax %s is not supported; send an uncompressed Implicit or Explicit VR Little Endian file",
			dicom.DescribeTransferSyntax(metadata.TransferSyntaxUID),
		))
	}

	studyInstanceUID := strings.TrimSpace(metadata.StudyInstanceUID)
	if studyInstanceUID == "" {
		return nil, apperrors.WithDetail(apperrors.ErrInvalidDICOM, "StudyInstanceUID is missing")
	}
	if !isValidDICOMUID(studyInstanceUID) {
		return nil, apperrors.WithDetail(apperrors.ErrInvalidDICOM, "StudyInstanceUID is not a valid UID")
	}

	sopInstanceUID := metadata.SOPUID()
	if sopInstanceUID == "" {
		return nil, apperrors.WithDetail(apperrors.ErrInvalidDICOM, "SOPInstanceUID is missing")
	}
	if !isValidDICOMUID(sopInstanceUID) {
		return nil, apperrors.WithDetail(apperrors.ErrInvalidDICOM, "SOPInstanceUID is not a valid UID")
	}

	seriesInstanceUID := strings.TrimSpace(metadata.SeriesInstanceUID)
	if seriesInstanceUID == "" {
		seriesInstanceUID = studyInstanceUID + ".1"
	} else if !isValidDICOMUID(seriesInstanceUID) {
		return nil, apperrors.WithDetail(apperrors.ErrInvalidDICOM, "SeriesInstanceUID is not a valid UID")
	}

	validated := &validatedDICOMMetadata{
		studyInstanceUID:  studyInstanceUID,
		seriesInstanceUID: seriesInstanceUID,
		sopInstanceUID:    sopInstanceUID,
		modality:          normalizeFileModality(metadata.Modality),
		studyDescription:  truncateRunes(strings.TrimSpace(metadata.StudyDescription), maximumDescriptionLength),
		accessionNumber:   truncateRunes(strings.TrimSpace(metadata.AccessionNumber), maximumAccessionLength),
		studyDateTime:     parseStudyDateTime(metadata.StudyDate, metadata.StudyTime),
		source:            metadata,
	}
	validated.instanceIndex = buildInstanceIndex(validated, metadata)
	return validated, nil
}

func isValidDICOMUID(uid string) bool {
	return len(uid) <= maximumUIDLength && dicomUIDPattern.MatchString(uid)
}

func normalizeFileModality(rawModality string) string {
	candidate := strings.ToUpper(strings.TrimSpace(rawModality))
	if dicomModalityPattern.MatchString(candidate) {
		return candidate
	}
	return ""
}

func buildInstanceIndex(validated *validatedDICOMMetadata, metadata *dicom.InstanceMetadata) *InstanceIndex {
	rowSpacing, columnSpacing, _ := metadata.PixelSpacing()
	return &InstanceIndex{
		SeriesInstanceUID:  validated.seriesInstanceUID,
		SeriesNumber:       metadata.SeriesNumber,
		SeriesDescription:  truncateRunes(strings.TrimSpace(metadata.SeriesDescription), maximumDescriptionLength),
		Modality:           validated.modality,
		SOPInstanceUID:     validated.sopInstanceUID,
		SOPClassUID:        metadata.SOPClassUID(),
		TransferSyntaxUID:  strings.TrimSpace(metadata.TransferSyntaxUID),
		InstanceNumber:     metadata.InstanceNumber,
		FrameCount:         metadata.FrameCount(),
		Rows:               metadata.Rows,
		Columns:            metadata.Columns,
		PixelSpacingRow:    rowSpacing,
		PixelSpacingColumn: columnSpacing,
	}
}

func parseStudyDateTime(rawDate string, rawTime string) *time.Time {
	trimmedDate := strings.TrimSpace(rawDate)
	if len(trimmedDate) != 8 {
		return nil
	}

	timeDigits := dicomTimeDigits.FindString(strings.TrimSpace(rawTime))
	for len(timeDigits) < 6 {
		timeDigits += "0"
	}
	timeDigits = timeDigits[:6]

	parsedDateTime, parseError := time.ParseInLocation("20060102150405", trimmedDate+timeDigits, time.UTC)
	if parseError != nil {
		return nil
	}
	return &parsedDateTime
}

func buildStudyTitle(modality string, studyDescription string, studyDateTime *time.Time) string {
	switch {
	case studyDescription != "" && modality != "":
		return truncateRunes(modality+" - "+studyDescription, maximumTitleLength)
	case studyDescription != "":
		return truncateRunes(studyDescription, maximumTitleLength)
	case studyDateTime != nil && modality != "":
		return modality + " - " + studyDateTime.Format("2006-01-02")
	}
	return modality
}

func truncateRunes(value string, maximumRunes int) string {
	if utf8.RuneCountInString(value) <= maximumRunes {
		return value
	}
	return string([]rune(value)[:maximumRunes])
}
