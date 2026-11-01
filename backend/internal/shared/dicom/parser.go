package dicom

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	preambleLength          = 128
	magicBytesSignature     = "DICM"
	preambleAndMagicLength  = preambleLength + len(magicBytesSignature)
	undefinedElementLength  = 0xFFFFFFFF
	defaultScanByteBudget   = 1 << 20
	maximumReasonableFrames = 100000
	headerFixedBytes        = 8
)

var (
	ErrNotDICOM             = errors.New("not a DICOM Part 10 file")
	ErrMetaGroupMalformed   = errors.New("DICOM file meta group is malformed")
	ErrTransferSyntaxAbsent = errors.New("DICOM file meta group has no TransferSyntaxUID")
)

type valueRepresentation struct {
	code              string
	usesExtendedWidth bool
}

type elementHeader struct {
	tagValue            uint32
	valueLength         uint32
	lengthFieldWidth    uint32
	valueRepresentation valueRepresentation
}

var extendedWidthValueRepresentations = map[string]bool{
	"OB": true,
	"OD": true,
	"OF": true,
	"OL": true,
	"OV": true,
	"OW": true,
	"SQ": true,
	"UC": true,
	"UN": true,
	"UR": true,
	"UT": true,
}

type FileMeta struct {
	TransferSyntaxUID          string
	MediaStorageSOPClassUID    string
	MediaStorageSOPInstanceUID string
	ImplementationClassUID     string
	ImplementationVersionName  string
}

type InstanceMetadata struct {
	FileMeta
	StudyInstanceUID          string
	SeriesInstanceUID         string
	SOPInstanceUID            string
	Modality                  string
	StudyDescription          string
	SeriesDescription         string
	SeriesNumber              int
	InstanceNumber            int
	Rows                      int
	Columns                   int
	BitsAllocated             int
	BitsStored                int
	HighBit                   int
	PixelRepresentation       int
	NumberOfFrames            int
	PixelSpacingRow           string
	PixelSpacingColumn        string
	ImagePositionPatient      []string
	ImageOrientationPatient   []string
	StudyDate                 string
	StudyTime                 string
	AccessionNumber           string
	ScanTruncatedBeforePixels bool
}

func (metadata *InstanceMetadata) PixelSpacing() (float64, float64, bool) {
	rowSpacing, rowError := strconv.ParseFloat(strings.TrimSpace(metadata.PixelSpacingRow), 64)
	columnSpacing, columnError := strconv.ParseFloat(strings.TrimSpace(metadata.PixelSpacingColumn), 64)
	if rowError != nil || columnError != nil {
		return 0, 0, false
	}
	if rowSpacing <= 0 || columnSpacing <= 0 {
		return 0, 0, false
	}
	return rowSpacing, columnSpacing, true
}

func (metadata *InstanceMetadata) SOPClassUID() string {
	if trimmedClassUID := strings.TrimSpace(metadata.MediaStorageSOPClassUID); trimmedClassUID != "" {
		return trimmedClassUID
	}
	return SOPClassSecondaryCapture
}

func (metadata *InstanceMetadata) SOPUID() string {
	if trimmedInstanceUID := strings.TrimSpace(metadata.SOPInstanceUID); trimmedInstanceUID != "" {
		return trimmedInstanceUID
	}
	return strings.TrimSpace(metadata.MediaStorageSOPInstanceUID)
}

func (metadata *InstanceMetadata) FrameCount() int {
	if metadata.NumberOfFrames < 1 {
		return 1
	}
	if metadata.NumberOfFrames > maximumReasonableFrames {
		return maximumReasonableFrames
	}
	return metadata.NumberOfFrames
}

func ParseInstance(sourceReader io.Reader) (*InstanceMetadata, error) {
	bufferedReader := bufio.NewReaderSize(sourceReader, 64*1024)

	preamble := make([]byte, preambleAndMagicLength)
	readBytes, readError := io.ReadFull(bufferedReader, preamble)
	if readError != nil && !errors.Is(readError, io.EOF) && !errors.Is(readError, io.ErrUnexpectedEOF) {
		return nil, readError
	}
	if readBytes < preambleAndMagicLength || string(preamble[preambleLength:]) != magicBytesSignature {
		return nil, ErrNotDICOM
	}

	fileMeta, pendingDatasetTag, metaError := parseFileMetaGroup(bufferedReader)
	if metaError != nil {
		return nil, metaError
	}
	if strings.TrimSpace(fileMeta.TransferSyntaxUID) == "" {
		return nil, ErrTransferSyntaxAbsent
	}

	metadata := &InstanceMetadata{FileMeta: fileMeta}
	metadata.SOPInstanceUID = fileMeta.MediaStorageSOPInstanceUID
	parseDataset(bufferedReader, metadata, pendingDatasetTag)
	return metadata, nil
}

func parseFileMetaGroup(sourceReader *bufio.Reader) (FileMeta, uint32, error) {
	var fileMeta FileMeta

	for {
		tagValue, tagError := readTag(sourceReader)
		if tagError != nil {
			if isStreamEnd(tagError) {
				return fileMeta, 0, nil
			}
			return fileMeta, 0, fmt.Errorf("%w: %v", ErrMetaGroupMalformed, tagError)
		}

		if tagValue>>16 != 0x0002 {
			return fileMeta, tagValue, nil
		}

		header, headerError := finishElementHeader(sourceReader, tagValue, true)
		if headerError != nil {
			if isStreamEnd(headerError) {
				return fileMeta, 0, nil
			}
			return fileMeta, 0, fmt.Errorf("%w: %v", ErrMetaGroupMalformed, headerError)
		}

		valueBytes, valueError := readPaddedValue(sourceReader, header)
		if valueError != nil {
			return fileMeta, 0, fmt.Errorf("%w: %v", ErrMetaGroupMalformed, valueError)
		}

		switch tagValue {
		case 0x00020001:
			fileMeta.ImplementationVersionName = decodeDicomString(valueBytes)
		case 0x00020002:
			fileMeta.MediaStorageSOPClassUID = decodeDicomString(valueBytes)
		case 0x00020003:
			fileMeta.MediaStorageSOPInstanceUID = decodeDicomString(valueBytes)
		case 0x00020010:
			fileMeta.TransferSyntaxUID = decodeDicomString(valueBytes)
		case 0x00020012:
			fileMeta.ImplementationClassUID = decodeDicomString(valueBytes)
		}
	}
}

func parseDataset(sourceReader *bufio.Reader, metadata *InstanceMetadata, pendingTagValue uint32) {
	explicitValueRepresentation := strings.TrimSpace(metadata.TransferSyntaxUID) == TransferSyntaxExplicitVRLittleEndian
	remainingBudget := uint32(defaultScanByteBudget)
	tagValue := pendingTagValue

	for remainingBudget > 0 {
		currentTagValue := tagValue
		tagValue = 0

		if currentTagValue == 0 {
			readTagValue, tagError := readTag(sourceReader)
			if tagError != nil {
				metadata.ScanTruncatedBeforePixels = true
				return
			}
			currentTagValue = readTagValue
		}

		if currentTagValue == 0x7FE00010 {
			metadata.ScanTruncatedBeforePixels = false
			return
		}

		header, headerError := finishElementHeader(sourceReader, currentTagValue, explicitValueRepresentation)
		if headerError != nil {
			metadata.ScanTruncatedBeforePixels = true
			return
		}

		if header.valueLength == undefinedElementLength {
			metadata.ScanTruncatedBeforePixels = true
			return
		}

		headerWidth := headerFixedBytes + header.lengthFieldWidth
		if header.valueLength+headerWidth > remainingBudget {
			metadata.ScanTruncatedBeforePixels = true
			return
		}
		remainingBudget -= headerWidth

		if isCollectedDatasetTag(currentTagValue) {
			valueBytes, valueError := readPaddedValue(sourceReader, header)
			if valueError != nil {
				metadata.ScanTruncatedBeforePixels = true
				return
			}
			applyDatasetValue(metadata, currentTagValue, valueBytes)
			continue
		}

		if skipError := skipValue(sourceReader, header.valueLength); skipError != nil {
			metadata.ScanTruncatedBeforePixels = true
			return
		}
	}

	metadata.ScanTruncatedBeforePixels = true
}

func readTag(sourceReader *bufio.Reader) (uint32, error) {
	tagBytes := make([]byte, 4)
	if _, readError := io.ReadFull(sourceReader, tagBytes); readError != nil {
		return 0, readError
	}
	groupNumber := binary.LittleEndian.Uint16(tagBytes[0:2])
	elementNumber := binary.LittleEndian.Uint16(tagBytes[2:4])
	return uint32(groupNumber)<<16 | uint32(elementNumber), nil
}

func finishElementHeader(sourceReader *bufio.Reader, tagValue uint32, explicitValueRepresentation bool) (elementHeader, error) {
	header := elementHeader{tagValue: tagValue}

	if !explicitValueRepresentation {
		lengthBytes := make([]byte, 4)
		if _, readError := io.ReadFull(sourceReader, lengthBytes); readError != nil {
			return elementHeader{}, readError
		}
		header.valueLength = binary.LittleEndian.Uint32(lengthBytes)
		header.lengthFieldWidth = 4
		return header, nil
	}

	representationBytes := make([]byte, 2)
	if _, readError := io.ReadFull(sourceReader, representationBytes); readError != nil {
		return elementHeader{}, readError
	}
	representation := valueRepresentation{code: string(representationBytes)}
	header.valueRepresentation = representation

	if extendedWidthValueRepresentations[representation.code] {
		extendedLengthBytes := make([]byte, 6)
		if _, readError := io.ReadFull(sourceReader, extendedLengthBytes); readError != nil {
			return elementHeader{}, readError
		}
		header.valueLength = binary.LittleEndian.Uint32(extendedLengthBytes[2:6])
		header.lengthFieldWidth = 4
		return header, nil
	}

	lengthBytes := make([]byte, 2)
	if _, readError := io.ReadFull(sourceReader, lengthBytes); readError != nil {
		return elementHeader{}, readError
	}
	header.valueLength = uint32(binary.LittleEndian.Uint16(lengthBytes))
	header.lengthFieldWidth = 2
	return header, nil
}

func readPaddedValue(sourceReader *bufio.Reader, header elementHeader) ([]byte, error) {
	valueBytes := make([]byte, header.valueLength)
	if _, readError := io.ReadFull(sourceReader, valueBytes); readError != nil {
		return nil, readError
	}
	if header.valueLength%2 == 1 {
		padByte := make([]byte, 1)
		if _, readError := io.ReadFull(sourceReader, padByte); readError != nil {
			return nil, readError
		}
	}
	return valueBytes, nil
}

func skipValue(sourceReader *bufio.Reader, valueLength uint32) error {
	if valueLength > 0 {
		if _, skipError := io.CopyN(io.Discard, sourceReader, int64(valueLength)); skipError != nil {
			return skipError
		}
	}
	if valueLength%2 == 1 {
		padByte := make([]byte, 1)
		if _, readError := io.ReadFull(sourceReader, padByte); readError != nil {
			return readError
		}
	}
	return nil
}

func isStreamEnd(readError error) bool {
	return errors.Is(readError, io.EOF) || errors.Is(readError, io.ErrUnexpectedEOF)
}

func isCollectedDatasetTag(tagValue uint32) bool {
	switch tagValue {
	case 0x00080020, 0x00080030, 0x00080050, 0x00080060, 0x00081030, 0x0008103E,
		0x0020000D, 0x0020000E, 0x00200011, 0x00200013,
		0x00200032, 0x00200037, 0x00280008, 0x00280010, 0x00280011,
		0x00280030, 0x00280100, 0x00280101, 0x00280102, 0x00280103:
		return true
	}
	return false
}

func applyDatasetValue(metadata *InstanceMetadata, tagValue uint32, valueBytes []byte) {
	decodedValue := decodeDicomString(valueBytes)

	switch tagValue {
	case 0x00080020:
		metadata.StudyDate = decodedValue
	case 0x00080030:
		metadata.StudyTime = decodedValue
	case 0x00080050:
		metadata.AccessionNumber = decodedValue
	case 0x00080060:
		metadata.Modality = decodedValue
	case 0x00081030:
		metadata.StudyDescription = decodedValue
	case 0x0008103E:
		metadata.SeriesDescription = decodedValue
	case 0x0020000D:
		metadata.StudyInstanceUID = decodedValue
	case 0x0020000E:
		metadata.SeriesInstanceUID = decodedValue
	case 0x00200011:
		metadata.SeriesNumber = parseDicomInteger(decodedValue)
	case 0x00200013:
		metadata.InstanceNumber = parseDicomInteger(decodedValue)
	case 0x00200032:
		metadata.ImagePositionPatient = splitDicomDecimalString(decodedValue)
	case 0x00200037:
		metadata.ImageOrientationPatient = splitDicomDecimalString(decodedValue)
	case 0x00280008:
		metadata.NumberOfFrames = parseDicomInteger(decodedValue)
	case 0x00280010:
		metadata.Rows = decodeUnsignedShortValue(valueBytes, 0)
	case 0x00280011:
		metadata.Columns = decodeUnsignedShortValue(valueBytes, 0)
	case 0x00280030:
		spacingValues := splitDicomDecimalString(decodedValue)
		if len(spacingValues) == 2 {
			metadata.PixelSpacingRow = spacingValues[0]
			metadata.PixelSpacingColumn = spacingValues[1]
		}
	case 0x00280100:
		metadata.BitsAllocated = decodeUnsignedShortValue(valueBytes, 0)
	case 0x00280101:
		metadata.BitsStored = decodeUnsignedShortValue(valueBytes, 0)
	case 0x00280102:
		metadata.HighBit = decodeUnsignedShortValue(valueBytes, 0)
	case 0x00280103:
		metadata.PixelRepresentation = decodeUnsignedShortValue(valueBytes, 0)
	}
}

func decodeUnsignedShortValue(valueBytes []byte, valueIndex int) int {
	valueOffset := valueIndex * 2
	if valueOffset+2 > len(valueBytes) {
		return 0
	}
	return int(binary.LittleEndian.Uint16(valueBytes[valueOffset : valueOffset+2]))
}

func decodeDicomString(valueBytes []byte) string {
	return strings.TrimSpace(strings.TrimRight(string(valueBytes), "\x00"))
}

func parseDicomInteger(rawValue string) int {
	trimmedValue := decodeDicomString([]byte(rawValue))
	if trimmedValue == "" {
		return 0
	}
	if parsedValue, parseError := strconv.Atoi(trimmedValue); parseError == nil {
		return parsedValue
	}
	if parsedValue, parseError := strconv.ParseFloat(trimmedValue, 64); parseError == nil {
		return int(parsedValue)
	}
	return 0
}

func splitDicomDecimalString(rawValue string) []string {
	trimmedValue := decodeDicomString([]byte(rawValue))
	if trimmedValue == "" {
		return nil
	}
	parts := strings.Split(trimmedValue, `\`)
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		values = append(values, strings.TrimSpace(part))
	}
	return values
}
