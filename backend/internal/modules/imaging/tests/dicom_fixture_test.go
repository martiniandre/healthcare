package tests

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"time"

	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/modules/imaging/mocks"
	"github.com/healthcare/backend/internal/shared/dicom"
	"github.com/healthcare/backend/internal/shared/storage"
)

const (
	fixtureStudyInstanceUID = "1.2.826.0.1.3680043.8.498.2001"
	fixtureSeriesUID        = "1.2.826.0.1.3680043.8.498.3001"
	fixtureSOPInstanceUID   = "1.2.826.0.1.3680043.8.498.4001"
	fixtureSOPClassUID      = "1.2.840.10008.5.1.4.1.1.4"
	fixturePatientFhirID    = "patient-123"
)

type dicomFixture struct {
	transferSyntax    string
	sopInstanceUID    string
	studyInstanceUID  string
	seriesInstanceUID string
	modality          string
	patientID         string
	patientName       string
	studyDescription  string
	studyDate         string
	studyTime         string
	accessionNumber   string
	pixelPaddingBytes int
}

func newDICOMFixture() dicomFixture {
	return dicomFixture{
		transferSyntax:    dicom.TransferSyntaxExplicitVRLittleEndian,
		sopInstanceUID:    fixtureSOPInstanceUID,
		studyInstanceUID:  fixtureStudyInstanceUID,
		seriesInstanceUID: fixtureSeriesUID,
		modality:          "MR",
		patientID:         "529.982.247-25",
		patientName:       "SILVA^MARIA",
		studyDescription:  "BRAIN WITHOUT CONTRAST",
		studyDate:         "20260115",
		studyTime:         "101500",
		accessionNumber:   "ACC-0042",
	}
}

func writeFixtureElement(targetBuffer *bytes.Buffer, groupNumber uint16, elementNumber uint16, representationCode string, value []byte, explicit bool) {
	tagBytes := make([]byte, 4)
	binary.LittleEndian.PutUint16(tagBytes[0:2], groupNumber)
	binary.LittleEndian.PutUint16(tagBytes[2:4], elementNumber)
	targetBuffer.Write(tagBytes)

	paddedValue := value
	if len(paddedValue)%2 == 1 {
		paddedValue = append(append([]byte{}, value...), 0)
	}

	extendedWidth := representationCode == "OB" || representationCode == "OW" || representationCode == "UN" || representationCode == "SQ"
	if explicit {
		targetBuffer.WriteString(representationCode)
		if extendedWidth {
			targetBuffer.Write([]byte{0, 0})
			lengthField := make([]byte, 4)
			binary.LittleEndian.PutUint32(lengthField, uint32(len(paddedValue)))
			targetBuffer.Write(lengthField)
		} else {
			lengthField := make([]byte, 2)
			binary.LittleEndian.PutUint16(lengthField, uint16(len(paddedValue)))
			targetBuffer.Write(lengthField)
		}
	} else {
		lengthField := make([]byte, 4)
		binary.LittleEndian.PutUint32(lengthField, uint32(len(paddedValue)))
		targetBuffer.Write(lengthField)
	}
	targetBuffer.Write(paddedValue)
}

func (fixture dicomFixture) bytes() []byte {
	var fileBuffer bytes.Buffer
	fileBuffer.Write(make([]byte, 128))
	fileBuffer.WriteString("DICM")

	writeFixtureElement(&fileBuffer, 0x0002, 0x0002, "UI", []byte(fixtureSOPClassUID), true)
	if fixture.sopInstanceUID != "" {
		writeFixtureElement(&fileBuffer, 0x0002, 0x0003, "UI", []byte(fixture.sopInstanceUID), true)
	}
	writeFixtureElement(&fileBuffer, 0x0002, 0x0010, "UI", []byte(fixture.transferSyntax), true)

	explicit := fixture.transferSyntax != dicom.TransferSyntaxImplicitVRLittleEndian
	optionalElements := []struct {
		group   uint16
		element uint16
		vr      string
		value   string
	}{
		{0x0008, 0x0020, "DA", fixture.studyDate},
		{0x0008, 0x0030, "TM", fixture.studyTime},
		{0x0008, 0x0050, "SH", fixture.accessionNumber},
		{0x0008, 0x0060, "CS", fixture.modality},
		{0x0008, 0x1030, "LO", fixture.studyDescription},
		{0x0010, 0x0010, "PN", fixture.patientName},
		{0x0010, 0x0020, "LO", fixture.patientID},
		{0x0020, 0x000D, "UI", fixture.studyInstanceUID},
		{0x0020, 0x000E, "UI", fixture.seriesInstanceUID},
	}
	for _, element := range optionalElements {
		if element.value != "" {
			writeFixtureElement(&fileBuffer, element.group, element.element, element.vr, []byte(element.value), explicit)
		}
	}

	writeFixtureElement(&fileBuffer, 0x7FE0, 0x0010, "OW", append([]byte{1, 2, 3, 4}, make([]byte, fixture.pixelPaddingBytes)...), explicit)
	return fileBuffer.Bytes()
}

func (fixture dicomFixture) reader() io.Reader {
	return bytes.NewReader(fixture.bytes())
}

type recordingStorage struct {
	storage.StorageClient
	uploadedObjects map[string][]byte
	deletedKeys     []string
	uploadError     error
}

func newRecordingStorage() *recordingStorage {
	return &recordingStorage{
		StorageClient:   storage.NewStorageClient("test-bucket", "test"),
		uploadedObjects: make(map[string][]byte),
	}
}

func (recording *recordingStorage) Upload(ctx context.Context, key string, reader io.Reader, contentType string) error {
	uploadedBytes, readError := io.ReadAll(reader)
	if readError != nil {
		return readError
	}
	if recording.uploadError != nil {
		return recording.uploadError
	}
	recording.uploadedObjects[key] = uploadedBytes
	return nil
}

func (recording *recordingStorage) Delete(ctx context.Context, key string) error {
	recording.deletedKeys = append(recording.deletedKeys, key)
	delete(recording.uploadedObjects, key)
	return nil
}

type uploadHarness struct {
	repository *mocks.MockRepository
	storage    *recordingStorage
	identity   *mocks.MockPatientIdentityReader
	service    imaging.Service
}

func newUploadHarness() *uploadHarness {
	repository := mocks.NewMockRepository()
	recording := newRecordingStorage()
	identityReader := &mocks.MockPatientIdentityReader{Identity: imaging.PatientIdentity{
		FullNames:   []string{"Maria Silva"},
		DocumentIDs: []string{"52998224725"},
	}}
	return &uploadHarness{
		repository: repository,
		storage:    recording,
		identity:   identityReader,
		service:    imaging.NewService(repository, recording, nil, nil, identityReader),
	}
}

func (harness *uploadHarness) upload(fixture dicomFixture, title string, modality string) (*imaging.UploadDICOMResult, error) {
	return harness.service.UploadDICOMStream(context.Background(), imaging.UploadDICOMInput{
		PatientFhirID: fixturePatientFhirID,
		Title:         title,
		Modality:      modality,
		FileName:      "scan.dcm",
	}, fixture.reader())
}

func warningCodes(warnings []imaging.UploadWarning) []string {
	codes := make([]string, 0, len(warnings))
	for _, warning := range warnings {
		codes = append(codes, warning.Code)
	}
	return codes
}

func expectedStudyTime(layout string, value string) time.Time {
	parsedTime, _ := time.ParseInLocation(layout, value, time.UTC)
	return parsedTime
}
