package dicom

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildExplicitElement(targetBuffer *bytes.Buffer, groupNumber uint16, elementNumber uint16, representationCode string, value []byte) {
	writeTag(targetBuffer, groupNumber, elementNumber)
	targetBuffer.WriteString(representationCode)
	if extendedWidthValueRepresentations[representationCode] {
		targetBuffer.Write([]byte{0, 0})
		lengthField := make([]byte, 4)
		binary.LittleEndian.PutUint32(lengthField, uint32(len(value)))
		targetBuffer.Write(lengthField)
	} else {
		lengthField := make([]byte, 2)
		binary.LittleEndian.PutUint16(lengthField, uint16(len(value)))
		targetBuffer.Write(lengthField)
	}
	targetBuffer.Write(value)
	writePadByte(targetBuffer, len(value))
}

func buildImplicitElement(targetBuffer *bytes.Buffer, groupNumber uint16, elementNumber uint16, value []byte) {
	writeTag(targetBuffer, groupNumber, elementNumber)
	lengthField := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthField, uint32(len(value)))
	targetBuffer.Write(lengthField)
	targetBuffer.Write(value)
	writePadByte(targetBuffer, len(value))
}

func writeTag(targetBuffer *bytes.Buffer, groupNumber uint16, elementNumber uint16) {
	groupField := make([]byte, 2)
	binary.LittleEndian.PutUint16(groupField, groupNumber)
	targetBuffer.Write(groupField)

	elementField := make([]byte, 2)
	binary.LittleEndian.PutUint16(elementField, elementNumber)
	targetBuffer.Write(elementField)
}

func writePadByte(targetBuffer *bytes.Buffer, currentLength int) {
	if currentLength%2 == 1 {
		targetBuffer.WriteByte(0)
	}
}

func writePreambleAndMagic(targetBuffer *bytes.Buffer) {
	targetBuffer.Write(bytes.Repeat([]byte{0}, preambleLength))
	targetBuffer.WriteString(magicBytesSignature)
}

func buildMetaGroup(targetBuffer *bytes.Buffer, transferSyntaxUID string, sopClassUID string, sopInstanceUID string) {
	buildExplicitElement(targetBuffer, 0x0002, 0x0001, "OB", []byte{0x00, 0x01})
	buildExplicitElement(targetBuffer, 0x0002, 0x0002, "UI", []byte(sopClassUID))
	buildExplicitElement(targetBuffer, 0x0002, 0x0003, "UI", []byte(sopInstanceUID))
	buildExplicitElement(targetBuffer, 0x0002, 0x0010, "UI", []byte(transferSyntaxUID))
	buildExplicitElement(targetBuffer, 0x0002, 0x0012, "UI", []byte("1.2.826.0.1.3680043.8.498.1"))
}

func buildExplicitVRLittleEndianInstance() []byte {
	var instanceBuffer bytes.Buffer
	writePreambleAndMagic(&instanceBuffer)
	buildMetaGroup(&instanceBuffer, TransferSyntaxExplicitVRLittleEndian, "1.2.840.10008.5.1.4.1.1.2", "1.2.826.0.1.3680043.8.498.100")

	buildExplicitElement(&instanceBuffer, 0x0008, 0x0020, "DA", []byte("20260115"))
	buildExplicitElement(&instanceBuffer, 0x0008, 0x0030, "TM", []byte("101500"))
	buildExplicitElement(&instanceBuffer, 0x0008, 0x0060, "CS", []byte("CT"))
	buildExplicitElement(&instanceBuffer, 0x0008, 0x1030, "LO", []byte("CHEST WITHOUT CONTRAST"))
	buildExplicitElement(&instanceBuffer, 0x0008, 0x103E, "LO", []byte("AXIAL 1.0MM"))
	buildExplicitElement(&instanceBuffer, 0x0020, 0x000D, "UI", []byte("1.2.826.0.1.3680043.8.498.200"))
	buildExplicitElement(&instanceBuffer, 0x0020, 0x000E, "UI", []byte("1.2.826.0.1.3680043.8.498.300"))
	buildExplicitElement(&instanceBuffer, 0x0020, 0x0011, "IS", []byte("2"))
	buildExplicitElement(&instanceBuffer, 0x0020, 0x0013, "IS", []byte("7"))
	buildExplicitElement(&instanceBuffer, 0x0020, 0x0032, "DS", []byte("-125.500\\-125.500\\-42.250"))
	buildExplicitElement(&instanceBuffer, 0x0020, 0x0037, "DS", []byte("1\\0\\0\\0\\1\\0"))
	buildExplicitElement(&instanceBuffer, 0x0028, 0x0008, "IS", []byte("1"))
	buildExplicitElement(&instanceBuffer, 0x0028, 0x0010, "US", []byte{0x00, 0x02})
	buildExplicitElement(&instanceBuffer, 0x0028, 0x0011, "US", []byte{0x00, 0x01})
	buildExplicitElement(&instanceBuffer, 0x0028, 0x0030, "DS", []byte("0.703125\\0.703125"))
	buildExplicitElement(&instanceBuffer, 0x0028, 0x0100, "US", []byte{16, 0})
	buildExplicitElement(&instanceBuffer, 0x0028, 0x0101, "US", []byte{12, 0})
	buildExplicitElement(&instanceBuffer, 0x0028, 0x0102, "US", []byte{11, 0})
	buildExplicitElement(&instanceBuffer, 0x0028, 0x0103, "US", []byte{0, 0})
	buildExplicitElement(&instanceBuffer, 0x7FE0, 0x0010, "OW", bytes.Repeat([]byte{0x00}, 512*512*2))

	return instanceBuffer.Bytes()
}

func buildImplicitVRLittleEndianInstance() []byte {
	var instanceBuffer bytes.Buffer
	writePreambleAndMagic(&instanceBuffer)
	buildMetaGroup(&instanceBuffer, TransferSyntaxImplicitVRLittleEndian, "1.2.840.10008.5.1.4.1.1.4", "1.2.826.0.1.3680043.8.498.400")

	buildImplicitElement(&instanceBuffer, 0x0008, 0x0060, []byte("MR"))
	buildImplicitElement(&instanceBuffer, 0x0020, 0x000D, []byte("1.2.826.0.1.3680043.8.498.500"))
	buildImplicitElement(&instanceBuffer, 0x0020, 0x000E, []byte("1.2.826.0.1.3680043.8.498.600"))
	buildImplicitElement(&instanceBuffer, 0x0020, 0x0011, []byte("1"))
	buildImplicitElement(&instanceBuffer, 0x0020, 0x0013, []byte("3"))
	buildImplicitElement(&instanceBuffer, 0x0028, 0x0008, []byte("240"))
	buildImplicitElement(&instanceBuffer, 0x0028, 0x0010, []byte{0x00, 0x01})
	buildImplicitElement(&instanceBuffer, 0x0028, 0x0011, []byte{0x00, 0x01})
	buildImplicitElement(&instanceBuffer, 0x0028, 0x0030, []byte("0.5\\0.5"))
	buildImplicitElement(&instanceBuffer, 0x0028, 0x0100, []byte{16, 0})
	buildImplicitElement(&instanceBuffer, 0x0028, 0x0101, []byte{16, 0})
	buildImplicitElement(&instanceBuffer, 0x0028, 0x0103, []byte{1, 0})

	return instanceBuffer.Bytes()
}

func TestParseInstanceReadsExplicitVRLittleEndianMetadata(testingInstance *testing.T) {
	metadata, parseError := ParseInstance(bytes.NewReader(buildExplicitVRLittleEndianInstance()))

	require.NoError(testingInstance, parseError)
	assert.Equal(testingInstance, TransferSyntaxExplicitVRLittleEndian, metadata.TransferSyntaxUID)
	assert.Equal(testingInstance, "1.2.840.10008.5.1.4.1.1.2", metadata.SOPClassUID())
	assert.Equal(testingInstance, "1.2.826.0.1.3680043.8.498.200", metadata.StudyInstanceUID)
	assert.Equal(testingInstance, "1.2.826.0.1.3680043.8.498.300", metadata.SeriesInstanceUID)
	assert.Equal(testingInstance, "1.2.826.0.1.3680043.8.498.100", metadata.SOPUID())
	assert.Equal(testingInstance, "CT", metadata.Modality)
	assert.Equal(testingInstance, "CHEST WITHOUT CONTRAST", metadata.StudyDescription)
	assert.Equal(testingInstance, "AXIAL 1.0MM", metadata.SeriesDescription)
	assert.Equal(testingInstance, 2, metadata.SeriesNumber)
	assert.Equal(testingInstance, 7, metadata.InstanceNumber)
	assert.Equal(testingInstance, 512, metadata.Rows)
	assert.Equal(testingInstance, 256, metadata.Columns)
	assert.Equal(testingInstance, 16, metadata.BitsAllocated)
	assert.Equal(testingInstance, 12, metadata.BitsStored)
	assert.Equal(testingInstance, 11, metadata.HighBit)
	assert.Equal(testingInstance, 0, metadata.PixelRepresentation)
	assert.Equal(testingInstance, "20260115", metadata.StudyDate)
	assert.Equal(testingInstance, "101500", metadata.StudyTime)
	assert.False(testingInstance, metadata.ScanTruncatedBeforePixels)
}

func TestParseInstanceReadsImplicitVRLittleEndianMetadata(testingInstance *testing.T) {
	metadata, parseError := ParseInstance(bytes.NewReader(buildImplicitVRLittleEndianInstance()))

	require.NoError(testingInstance, parseError)
	assert.Equal(testingInstance, TransferSyntaxImplicitVRLittleEndian, metadata.TransferSyntaxUID)
	assert.Equal(testingInstance, "MR", metadata.Modality)
	assert.Equal(testingInstance, "1.2.826.0.1.3680043.8.498.500", metadata.StudyInstanceUID)
	assert.Equal(testingInstance, "1.2.826.0.1.3680043.8.498.600", metadata.SeriesInstanceUID)
	assert.Equal(testingInstance, 1, metadata.SeriesNumber)
	assert.Equal(testingInstance, 3, metadata.InstanceNumber)
	assert.Equal(testingInstance, 256, metadata.Rows)
	assert.Equal(testingInstance, 256, metadata.Columns)
	assert.Equal(testingInstance, 240, metadata.FrameCount())
	assert.Equal(testingInstance, 1, metadata.PixelRepresentation)
}

func TestParseInstanceDecodesMultiValuedDecimalStrings(testingInstance *testing.T) {
	metadata, parseError := ParseInstance(bytes.NewReader(buildExplicitVRLittleEndianInstance()))

	require.NoError(testingInstance, parseError)
	assert.Equal(testingInstance, []string{"-125.500", "-125.500", "-42.250"}, metadata.ImagePositionPatient)
	assert.Equal(testingInstance, []string{"1", "0", "0", "0", "1", "0"}, metadata.ImageOrientationPatient)
}

func TestPixelSpacingReturnsRowAndColumnSpacing(testingInstance *testing.T) {
	metadata, parseError := ParseInstance(bytes.NewReader(buildExplicitVRLittleEndianInstance()))
	require.NoError(testingInstance, parseError)

	rowSpacing, columnSpacing, spacingPresent := metadata.PixelSpacing()

	assert.True(testingInstance, spacingPresent)
	assert.InDelta(testingInstance, 0.703125, rowSpacing, 0.0000001)
	assert.InDelta(testingInstance, 0.703125, columnSpacing, 0.0000001)
}

func TestPixelSpacingAbsentYieldsNoSpacing(testingInstance *testing.T) {
	metadata := &InstanceMetadata{}

	rowSpacing, columnSpacing, spacingPresent := metadata.PixelSpacing()

	assert.False(testingInstance, spacingPresent)
	assert.Zero(testingInstance, rowSpacing)
	assert.Zero(testingInstance, columnSpacing)
}

func TestPixelSpacingNonPositiveYieldsNoSpacing(testingInstance *testing.T) {
	metadata := &InstanceMetadata{PixelSpacingRow: "0", PixelSpacingColumn: "0.5"}

	_, _, spacingPresent := metadata.PixelSpacing()

	assert.False(testingInstance, spacingPresent)
}

func TestFrameCountDefaultsToSingleFrame(testingInstance *testing.T) {
	metadata := &InstanceMetadata{NumberOfFrames: 0}

	assert.Equal(testingInstance, 1, metadata.FrameCount())
}

func TestFrameCountClampsAbsurdMultiFrameCounts(testingInstance *testing.T) {
	metadata := &InstanceMetadata{NumberOfFrames: 999999}

	assert.Equal(testingInstance, maximumReasonableFrames, metadata.FrameCount())
}

func TestSOPUIDFallsBackToMediaStorageInstanceUID(testingInstance *testing.T) {
	metadata := &InstanceMetadata{}
	metadata.MediaStorageSOPInstanceUID = "1.2.826.0.1.3680043.8.498.900"

	assert.Equal(testingInstance, "1.2.826.0.1.3680043.8.498.900", metadata.SOPUID())
}

func TestParseInstanceRejectsNonDicomPayload(testingInstance *testing.T) {
	_, parseError := ParseInstance(strings.NewReader("this is definitely not a dicom payload"))

	assert.ErrorIs(testingInstance, parseError, ErrNotDICOM)
}

func TestParseInstanceRejectsMissingMagicBytes(testingInstance *testing.T) {
	truncatedPayload := bytes.Repeat([]byte{0}, preambleLength)

	_, parseError := ParseInstance(bytes.NewReader(truncatedPayload))

	assert.ErrorIs(testingInstance, parseError, ErrNotDICOM)
}

func TestParseInstanceRejectsPayloadShorterThanPreamble(testingInstance *testing.T) {
	_, parseError := ParseInstance(bytes.NewReader([]byte("DICM")))

	assert.ErrorIs(testingInstance, parseError, ErrNotDICOM)
}

func TestParseInstanceRejectsMetaGroupWithoutTransferSyntax(testingInstance *testing.T) {
	var instanceBuffer bytes.Buffer
	writePreambleAndMagic(&instanceBuffer)
	buildExplicitElement(&instanceBuffer, 0x0002, 0x0002, "UI", []byte("1.2.840.10008.5.1.4.1.1.2"))

	_, parseError := ParseInstance(bytes.NewReader(instanceBuffer.Bytes()))

	assert.ErrorIs(testingInstance, parseError, ErrTransferSyntaxAbsent)
}

func TestParseInstanceReadsCompressedTransferSyntaxForRejection(testingInstance *testing.T) {
	var instanceBuffer bytes.Buffer
	writePreambleAndMagic(&instanceBuffer)
	buildMetaGroup(&instanceBuffer, TransferSyntaxJPEG2000, "1.2.840.10008.5.1.4.1.1.2", "1.2.826.0.1.3680043.8.498.100")
	buildImplicitElement(&instanceBuffer, 0x0020, 0x000D, []byte("1.2.826.0.1.3680043.8.498.200"))

	metadata, parseError := ParseInstance(bytes.NewReader(instanceBuffer.Bytes()))

	require.NoError(testingInstance, parseError)
	assert.Equal(testingInstance, TransferSyntaxJPEG2000, metadata.TransferSyntaxUID)
	assert.False(testingInstance, IsSupportedTransferSyntax(metadata.TransferSyntaxUID))
}

func TestParseInstanceStopsAtUndefinedLengthSequence(testingInstance *testing.T) {
	var instanceBuffer bytes.Buffer
	writePreambleAndMagic(&instanceBuffer)
	buildMetaGroup(&instanceBuffer, TransferSyntaxExplicitVRLittleEndian, "1.2.840.10008.5.1.4.1.1.2", "1.2.826.0.1.3680043.8.498.100")
	buildExplicitElement(&instanceBuffer, 0x0008, 0x0060, "CS", []byte("US"))

	writeTag(&instanceBuffer, 0x0008, 0x1110)
	instanceBuffer.WriteString("SQ")
	instanceBuffer.Write([]byte{0, 0})
	undefinedLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(undefinedLength, undefinedElementLength)
	instanceBuffer.Write(undefinedLength)
	buildExplicitElement(&instanceBuffer, 0x0020, 0x000D, "UI", []byte("1.2.826.0.1.3680043.8.498.200"))

	metadata, parseError := ParseInstance(bytes.NewReader(instanceBuffer.Bytes()))

	require.NoError(testingInstance, parseError)
	assert.Equal(testingInstance, "US", metadata.Modality)
	assert.True(testingInstance, metadata.ScanTruncatedBeforePixels)
	assert.Empty(testingInstance, metadata.StudyInstanceUID)
}

func TestParseInstanceMarksTruncationOnShortDataset(testingInstance *testing.T) {
	var instanceBuffer bytes.Buffer
	writePreambleAndMagic(&instanceBuffer)
	buildMetaGroup(&instanceBuffer, TransferSyntaxImplicitVRLittleEndian, "1.2.840.10008.5.1.4.1.1.2", "1.2.826.0.1.3680043.8.498.100")
	buildImplicitElement(&instanceBuffer, 0x0008, 0x0060, []byte("CT"))

	metadata, parseError := ParseInstance(bytes.NewReader(instanceBuffer.Bytes()))

	require.NoError(testingInstance, parseError)
	assert.Equal(testingInstance, "CT", metadata.Modality)
	assert.True(testingInstance, metadata.ScanTruncatedBeforePixels)
}

func TestIsSupportedTransferSyntaxAcceptsOnlyLittleEndianVariants(testingInstance *testing.T) {
	assert.True(testingInstance, IsSupportedTransferSyntax(TransferSyntaxImplicitVRLittleEndian))
	assert.True(testingInstance, IsSupportedTransferSyntax(TransferSyntaxExplicitVRLittleEndian))
	assert.False(testingInstance, IsSupportedTransferSyntax(TransferSyntaxJPEGBaseline))
	assert.False(testingInstance, IsSupportedTransferSyntax(TransferSyntaxJPEG2000))
	assert.False(testingInstance, IsSupportedTransferSyntax(TransferSyntaxRLELossless))
	assert.False(testingInstance, IsSupportedTransferSyntax(TransferSyntaxExplicitVRBigEndian))
	assert.False(testingInstance, IsSupportedTransferSyntax(TransferSyntaxDeflatedExplicitVRLE))
	assert.False(testingInstance, IsSupportedTransferSyntax(""))
}

func TestDescribeTransferSyntaxNamesKnownSyntaxes(testingInstance *testing.T) {
	assert.Equal(testingInstance, "Explicit VR Little Endian (1.2.840.10008.1.2.1)", DescribeTransferSyntax(TransferSyntaxExplicitVRLittleEndian))
	assert.Equal(testingInstance, "JPEG 2000 Image Compression (1.2.840.10008.1.2.4.91)", DescribeTransferSyntax(TransferSyntaxJPEG2000))
	assert.Equal(testingInstance, "unknown", DescribeTransferSyntax(""))
	assert.Equal(testingInstance, "1.2.3.4.custom", DescribeTransferSyntax("1.2.3.4.custom"))
}
