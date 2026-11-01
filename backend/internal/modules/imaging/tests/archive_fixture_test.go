package tests

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/modules/imaging/mocks"
	"github.com/stretchr/testify/require"
)

type zipEntry struct {
	name    string
	content []byte
	symlink bool
}

func buildZIPArchive(testingInstance *testing.T, entries []zipEntry) []byte {
	testingInstance.Helper()

	var archiveBuffer bytes.Buffer
	archiveWriter := zip.NewWriter(&archiveBuffer)
	for _, entry := range entries {
		header := &zip.FileHeader{Name: entry.name, Method: zip.Deflate}
		if entry.symlink {
			header.SetMode(fs.ModeSymlink | 0o777)
		}
		entryWriter, createError := archiveWriter.CreateHeader(header)
		require.NoError(testingInstance, createError)
		_, writeError := entryWriter.Write(entry.content)
		require.NoError(testingInstance, writeError)
	}
	require.NoError(testingInstance, archiveWriter.Close())
	return archiveBuffer.Bytes()
}

func instanceFixtureWithSOP(sopSuffix string) dicomFixture {
	fixture := newDICOMFixture()
	fixture.sopInstanceUID = fixtureSOPInstanceUID + "." + sopSuffix
	return fixture
}

func instanceObjectKeyFor(patientFhirID string, studyID uuid.UUID, instanceID uuid.UUID) string {
	return "pacs/" + patientFhirID + "/" + studyID.String() + "/" + instanceID.String() + ".dcm"
}

func newUploadHarnessWithSettings(settings imaging.ServiceSettings) *uploadHarness {
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
		service:    imaging.NewServiceWithSettings(repository, recording, nil, nil, identityReader, settings),
	}
}

func uploadArchiveThrough(harness *uploadHarness, archiveBytes []byte, fileName string) (*imaging.UploadDICOMResult, error) {
	return harness.service.UploadDICOMStream(context.Background(), imaging.UploadDICOMInput{
		PatientFhirID: fixturePatientFhirID,
		FileName:      fileName,
	}, bytes.NewReader(archiveBytes))
}

type spoolObservingStorage struct {
	*recordingStorage
	spoolDirectory      string
	entriesDuringUpload []int
}

func (observing *spoolObservingStorage) Upload(ctx context.Context, key string, reader io.Reader, contentType string) error {
	spoolEntries, readError := os.ReadDir(observing.spoolDirectory)
	if readError != nil {
		return readError
	}
	observing.entriesDuringUpload = append(observing.entriesDuringUpload, len(spoolEntries))
	return observing.recordingStorage.Upload(ctx, key, reader, contentType)
}

func spoolDirectoryEntryCount(testingInstance *testing.T, directory string) int {
	testingInstance.Helper()
	spoolEntries, readError := os.ReadDir(directory)
	require.NoError(testingInstance, readError)
	return len(spoolEntries)
}

func findRejection(rejections []imaging.RejectedInstance, fileName string) (imaging.RejectedInstance, bool) {
	for _, rejection := range rejections {
		if rejection.FileName == fileName {
			return rejection, true
		}
	}
	return imaging.RejectedInstance{}, false
}
