package tests

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/modules/imaging/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/dicom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func requireAppError(testingInstance *testing.T, uploadError error, expectedCode string, expectedHTTPCode int) {
	testingInstance.Helper()
	var appError apperrors.AppError
	require.Error(testingInstance, uploadError)
	require.True(testingInstance, errors.As(uploadError, &appError), "expected an AppError, got %v", uploadError)
	assert.Equal(testingInstance, expectedCode, appError.Code)
	assert.Equal(testingInstance, expectedHTTPCode, appError.HTTPCode)
}

func TestService_UploadDICOMStream_SecondInstanceAppendsToExistingStudy(testingInstance *testing.T) {
	harness := newUploadHarness()
	first, firstError := harness.upload(instanceFixtureWithSOP("1"), "Brain MRI", "MR")
	require.NoError(testingInstance, firstError)

	second, secondError := harness.upload(instanceFixtureWithSOP("2"), "Ignored title", "MR")

	require.NoError(testingInstance, secondError)
	assert.False(testingInstance, second.StudyCreated)
	assert.Equal(testingInstance, first.Study.ID, second.Study.ID)
	assert.Equal(testingInstance, imaging.ImagingStudyStatusReceiving, second.Study.Status)
	assert.Equal(testingInstance, "Brain MRI", second.Study.Title)
	assert.Equal(testingInstance, 1, second.InstancesAccepted)
	assert.Len(testingInstance, harness.repository.Instances[first.Study.ID], 2)
	assert.Len(testingInstance, harness.storage.uploadedObjects, 2)
}

func TestService_UploadDICOMStream_SameStudyOfAnotherPatientIsRejected(testingInstance *testing.T) {
	harness := newUploadHarness()
	_, firstError := harness.upload(instanceFixtureWithSOP("1"), "", "")
	require.NoError(testingInstance, firstError)
	objectsBefore := len(harness.storage.uploadedObjects)
	writesBefore := harness.repository.WriteCount

	_, secondError := harness.service.UploadDICOMStream(context.Background(), imaging.UploadDICOMInput{
		PatientFhirID: "patient-999",
		FileName:      "scan.dcm",
	}, instanceFixtureWithSOP("2").reader())

	assert.ErrorIs(testingInstance, secondError, apperrors.ErrDICOMPatientMismatch)
	requireAppError(testingInstance, secondError, "DICOM_PATIENT_MISMATCH", 422)
	assert.Len(testingInstance, harness.storage.uploadedObjects, objectsBefore)
	assert.Equal(testingInstance, writesBefore, harness.repository.WriteCount)
}

func TestService_Archive_StoresEveryInstanceUnderItsOwnKey(testingInstance *testing.T) {
	harness := newUploadHarness()
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "series-a/first.dcm", content: instanceFixtureWithSOP("1").bytes()},
		{name: "series-a/second.dcm", content: instanceFixtureWithSOP("2").bytes()},
	})

	result, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	require.NoError(testingInstance, uploadError)
	assert.True(testingInstance, result.StudyCreated)
	assert.Equal(testingInstance, 2, result.InstancesAccepted)
	assert.Empty(testingInstance, result.InstancesRejected)
	require.Len(testingInstance, harness.storage.uploadedObjects, 2)

	storedInstances, listError := harness.repository.ListInstances(context.Background(), result.Study.ID)
	require.NoError(testingInstance, listError)
	require.Len(testingInstance, storedInstances, 2)
	fixtureBytesBySOP := map[string][]byte{
		fixtureSOPInstanceUID + ".1": instanceFixtureWithSOP("1").bytes(),
		fixtureSOPInstanceUID + ".2": instanceFixtureWithSOP("2").bytes(),
	}
	for _, storedInstance := range storedInstances {
		assert.Equal(testingInstance, instanceObjectKeyFor(fixturePatientFhirID, result.Study.ID, storedInstance.ID), storedInstance.ObjectKey)
		assert.True(testingInstance, bytes.Equal(fixtureBytesBySOP[storedInstance.SOPInstanceUID], harness.storage.uploadedObjects[storedInstance.ObjectKey]))
	}
	assert.Equal(testingInstance, "first.dcm", storedInstances[0].FileName)
	assert.Equal(testingInstance, "second.dcm", storedInstances[1].FileName)
}

func TestService_Archive_AppendsToExistingStudyOfSamePatient(testingInstance *testing.T) {
	harness := newUploadHarness()
	first, firstError := harness.upload(instanceFixtureWithSOP("1"), "", "")
	require.NoError(testingInstance, firstError)
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "second.dcm", content: instanceFixtureWithSOP("2").bytes()},
		{name: "third.dcm", content: instanceFixtureWithSOP("3").bytes()},
	})

	result, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	require.NoError(testingInstance, uploadError)
	assert.False(testingInstance, result.StudyCreated)
	assert.Equal(testingInstance, first.Study.ID, result.Study.ID)
	assert.Equal(testingInstance, imaging.ImagingStudyStatusReceiving, result.Study.Status)
	assert.Equal(testingInstance, 2, result.InstancesAccepted)
	assert.Len(testingInstance, harness.repository.Instances[first.Study.ID], 3)
}

func TestService_Archive_RepeatedInstancesAreRejectedPerInstance(testingInstance *testing.T) {
	harness := newUploadHarness()
	first, firstError := harness.upload(instanceFixtureWithSOP("1"), "", "")
	require.NoError(testingInstance, firstError)
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "already-stored.dcm", content: instanceFixtureWithSOP("1").bytes()},
		{name: "new.dcm", content: instanceFixtureWithSOP("2").bytes()},
	})

	result, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	require.NoError(testingInstance, uploadError)
	assert.Equal(testingInstance, first.Study.ID, result.Study.ID)
	assert.Equal(testingInstance, 1, result.InstancesAccepted)
	require.Len(testingInstance, result.InstancesRejected, 1)
	assert.Equal(testingInstance, "already-stored.dcm", result.InstancesRejected[0].FileName)
	assert.Equal(testingInstance, imaging.RejectionCodeDuplicateInstance, result.InstancesRejected[0].Code)
}

func TestService_Archive_AllDuplicatesReturnNoValidInstances(testingInstance *testing.T) {
	harness := newUploadHarness()
	_, firstError := harness.upload(instanceFixtureWithSOP("1"), "", "")
	require.NoError(testingInstance, firstError)
	writesBefore := harness.repository.WriteCount
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "already-stored.dcm", content: instanceFixtureWithSOP("1").bytes()},
	})

	_, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	requireAppError(testingInstance, uploadError, "DICOM_NO_VALID_INSTANCES", 422)
	assert.Len(testingInstance, harness.storage.uploadedObjects, 1)
	assert.Equal(testingInstance, writesBefore, harness.repository.WriteCount)
}

func TestService_Archive_PatientMismatchRejectsWholeArchive(testingInstance *testing.T) {
	harness := newUploadHarness()
	mismatched := instanceFixtureWithSOP("2")
	mismatched.patientID = "111.444.777-35"
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "belongs.dcm", content: instanceFixtureWithSOP("1").bytes()},
		{name: "does-not-belong.dcm", content: mismatched.bytes()},
	})

	_, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	assert.ErrorIs(testingInstance, uploadError, apperrors.ErrDICOMPatientMismatch)
	requireAppError(testingInstance, uploadError, "DICOM_PATIENT_MISMATCH", 422)
	assert.Empty(testingInstance, harness.storage.uploadedObjects)
	assert.Empty(testingInstance, harness.storage.deletedKeys)
	assert.Empty(testingInstance, harness.repository.Studies)
	assert.Equal(testingInstance, 0, harness.repository.WriteCount)
}

func TestService_Archive_PerInstanceProblemsKeepTheRest(testingInstance *testing.T) {
	harness := newUploadHarness()
	withoutStudy := instanceFixtureWithSOP("5")
	withoutStudy.studyInstanceUID = ""
	jpegCompressed := instanceFixtureWithSOP("6")
	jpegCompressed.transferSyntax = dicom.TransferSyntaxJPEGBaseline
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "ok-1.dcm", content: instanceFixtureWithSOP("1").bytes()},
		{name: "notes.txt", content: []byte("plain text, not dicom")},
		{name: "no-study.dcm", content: withoutStudy.bytes()},
		{name: "jpeg.dcm", content: jpegCompressed.bytes()},
		{name: "repeat.dcm", content: instanceFixtureWithSOP("1").bytes()},
		{name: "ok-2.dcm", content: instanceFixtureWithSOP("2").bytes()},
	})

	result, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	require.NoError(testingInstance, uploadError)
	assert.Equal(testingInstance, 2, result.InstancesAccepted)
	assert.Len(testingInstance, harness.storage.uploadedObjects, 2)
	require.Len(testingInstance, result.InstancesRejected, 4)
	expectedCodes := map[string]string{
		"notes.txt":    imaging.RejectionCodeInvalidDICOM,
		"no-study.dcm": imaging.RejectionCodeMissingIdentifiers,
		"jpeg.dcm":     imaging.RejectionCodeUnsupportedTransferSyntax,
		"repeat.dcm":   imaging.RejectionCodeDuplicateInstance,
	}
	for fileName, expectedCode := range expectedCodes {
		rejection, found := findRejection(result.InstancesRejected, fileName)
		require.True(testingInstance, found, "missing rejection for %s", fileName)
		assert.Equal(testingInstance, expectedCode, rejection.Code)
		assert.NotEmpty(testingInstance, rejection.Reason)
	}
}

func TestService_Archive_ZeroValidInstancesIsRejected(testingInstance *testing.T) {
	harness := newUploadHarness()
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "notes.txt", content: []byte("plain text")},
		{name: "readme.pdf", content: []byte("%PDF-1.4")},
	})

	_, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	requireAppError(testingInstance, uploadError, "DICOM_NO_VALID_INSTANCES", 422)
	assert.Empty(testingInstance, harness.storage.uploadedObjects)
	assert.Empty(testingInstance, harness.repository.Studies)
}

func TestService_Archive_UnsafeEntryNamesAreRejected(testingInstance *testing.T) {
	unsafeNames := []string{"../evil.dcm", "series/../../evil.dcm", "/etc/evil.dcm", "C:evil.dcm", "series\\evil.dcm"}
	for _, unsafeName := range unsafeNames {
		testingInstance.Run(unsafeName, func(testingInstance *testing.T) {
			harness := newUploadHarness()
			archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
				{name: "ok.dcm", content: instanceFixtureWithSOP("1").bytes()},
				{name: unsafeName, content: instanceFixtureWithSOP("2").bytes()},
			})

			_, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

			requireAppError(testingInstance, uploadError, "ZIP_UNSAFE_ENTRY", 422)
			assert.Empty(testingInstance, harness.storage.uploadedObjects)
			assert.Empty(testingInstance, harness.repository.Studies)
		})
	}
}

func TestService_Archive_SymlinkEntryIsRejected(testingInstance *testing.T) {
	harness := newUploadHarness()
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "link.dcm", content: []byte("/etc/passwd"), symlink: true},
	})

	_, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	requireAppError(testingInstance, uploadError, "ZIP_UNSAFE_ENTRY", 422)
}

func TestService_Archive_IgnoresDirectoriesAndMacOSMetadata(testingInstance *testing.T) {
	harness := newUploadHarness()
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "series/"},
		{name: "__MACOSX/series/._ok.dcm", content: []byte("resource fork")},
		{name: "series/.DS_Store", content: []byte("finder")},
		{name: "series/ok.dcm", content: instanceFixtureWithSOP("1").bytes()},
	})

	result, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	require.NoError(testingInstance, uploadError)
	assert.Equal(testingInstance, 1, result.InstancesAccepted)
	assert.Empty(testingInstance, result.InstancesRejected)
}

func TestService_Archive_TooManyEntriesIsRejected(testingInstance *testing.T) {
	settings := imaging.DefaultServiceSettings()
	settings.ArchiveLimits.MaxEntries = 2
	harness := newUploadHarnessWithSettings(settings)
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "a.dcm", content: instanceFixtureWithSOP("1").bytes()},
		{name: "b.dcm", content: instanceFixtureWithSOP("2").bytes()},
		{name: "c.dcm", content: instanceFixtureWithSOP("3").bytes()},
	})

	_, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	requireAppError(testingInstance, uploadError, "ZIP_LIMITS_EXCEEDED", 413)
	assert.Empty(testingInstance, harness.storage.uploadedObjects)
}

func TestService_Archive_TooLargeArchiveIsRejected(testingInstance *testing.T) {
	settings := imaging.DefaultServiceSettings()
	settings.ArchiveLimits.MaxArchiveBytes = 64
	harness := newUploadHarnessWithSettings(settings)
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{{name: "big.dcm", content: instanceFixtureWithSOP("1").bytes()}})

	_, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	requireAppError(testingInstance, uploadError, "ZIP_LIMITS_EXCEEDED", 413)
	assert.Empty(testingInstance, harness.storage.uploadedObjects)
}

func TestService_Archive_UncompressedTotalIsBounded(testingInstance *testing.T) {
	settings := imaging.DefaultServiceSettings()
	settings.ArchiveLimits.MaxUncompressedBytes = 2048
	harness := newUploadHarnessWithSettings(settings)
	fixture := instanceFixtureWithSOP("1")
	fixture.pixelPaddingBytes = 4096
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{{name: "big.dcm", content: fixture.bytes()}})

	_, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	requireAppError(testingInstance, uploadError, "ZIP_LIMITS_EXCEEDED", 413)
	assert.Empty(testingInstance, harness.storage.uploadedObjects)
}

func TestService_Archive_CompressionRatioBombIsRejected(testingInstance *testing.T) {
	harness := newUploadHarness()
	bombPayload := make([]byte, 4<<20)
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{{name: "bomb.dcm", content: bombPayload}})

	_, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	requireAppError(testingInstance, uploadError, "ZIP_LIMITS_EXCEEDED", 413)
	assert.Empty(testingInstance, harness.storage.uploadedObjects)
	assert.Empty(testingInstance, harness.repository.Studies)
}

func TestService_Archive_NonZIPPayloadNamedZipIsInvalid(testingInstance *testing.T) {
	harness := newUploadHarness()

	_, uploadError := uploadArchiveThrough(harness, instanceFixtureWithSOP("1").bytes(), "study.zip")

	requireAppError(testingInstance, uploadError, "ZIP_INVALID", 422)
	assert.Empty(testingInstance, harness.storage.uploadedObjects)
}

func TestService_Archive_InstancesOfTwoStudiesAreRejected(testingInstance *testing.T) {
	harness := newUploadHarness()
	otherStudy := instanceFixtureWithSOP("2")
	otherStudy.studyInstanceUID = fixtureStudyInstanceUID + ".9"
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "a.dcm", content: instanceFixtureWithSOP("1").bytes()},
		{name: "b.dcm", content: otherStudy.bytes()},
	})

	_, uploadError := uploadArchiveThrough(harness, archiveBytes, "study.zip")

	requireAppError(testingInstance, uploadError, "ZIP_INVALID", 422)
	assert.Empty(testingInstance, harness.storage.uploadedObjects)
}

func TestService_SingleFileIsStreamedWithoutSpooling(testingInstance *testing.T) {
	spoolDirectory := testingInstance.TempDir()
	settings := imaging.DefaultServiceSettings()
	settings.SpoolDirectory = spoolDirectory
	observing := &spoolObservingStorage{recordingStorage: newRecordingStorage(), spoolDirectory: spoolDirectory}
	repository := mocks.NewMockRepository()
	identityReader := &mocks.MockPatientIdentityReader{Identity: imaging.PatientIdentity{FullNames: []string{"Maria Silva"}, DocumentIDs: []string{"52998224725"}}}
	service := imaging.NewServiceWithSettings(repository, observing, nil, nil, identityReader, settings)

	_, uploadError := service.UploadDICOMStream(context.Background(), imaging.UploadDICOMInput{
		PatientFhirID: fixturePatientFhirID,
		FileName:      "scan.dcm",
	}, instanceFixtureWithSOP("1").reader())

	require.NoError(testingInstance, uploadError)
	assert.Equal(testingInstance, []int{0}, observing.entriesDuringUpload)
	assert.Equal(testingInstance, 0, spoolDirectoryEntryCount(testingInstance, spoolDirectory))
}

func TestService_Archive_IsSpooledAndSpoolFileIsRemoved(testingInstance *testing.T) {
	spoolDirectory := testingInstance.TempDir()
	settings := imaging.DefaultServiceSettings()
	settings.SpoolDirectory = spoolDirectory
	observing := &spoolObservingStorage{recordingStorage: newRecordingStorage(), spoolDirectory: spoolDirectory}
	repository := mocks.NewMockRepository()
	identityReader := &mocks.MockPatientIdentityReader{Identity: imaging.PatientIdentity{FullNames: []string{"Maria Silva"}, DocumentIDs: []string{"52998224725"}}}
	service := imaging.NewServiceWithSettings(repository, observing, nil, nil, identityReader, settings)
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "a.dcm", content: instanceFixtureWithSOP("1").bytes()},
		{name: "b.dcm", content: instanceFixtureWithSOP("2").bytes()},
	})

	_, uploadError := service.UploadDICOMStream(context.Background(), imaging.UploadDICOMInput{
		PatientFhirID: fixturePatientFhirID,
		FileName:      "study.zip",
	}, bytes.NewReader(archiveBytes))

	require.NoError(testingInstance, uploadError)
	assert.Equal(testingInstance, []int{1, 1}, observing.entriesDuringUpload)
	assert.Equal(testingInstance, 0, spoolDirectoryEntryCount(testingInstance, spoolDirectory))
}

func TestService_Archive_SpoolFileIsRemovedWhenRejected(testingInstance *testing.T) {
	spoolDirectory := testingInstance.TempDir()
	settings := imaging.DefaultServiceSettings()
	settings.SpoolDirectory = spoolDirectory
	observing := &spoolObservingStorage{recordingStorage: newRecordingStorage(), spoolDirectory: spoolDirectory}
	service := imaging.NewServiceWithSettings(mocks.NewMockRepository(), observing, nil, nil, nil, settings)

	_, uploadError := service.UploadDICOMStream(context.Background(), imaging.UploadDICOMInput{
		PatientFhirID: fixturePatientFhirID,
		FileName:      "study.zip",
	}, bytes.NewReader([]byte("not a zip")))

	requireAppError(testingInstance, uploadError, "ZIP_INVALID", 422)
	assert.Equal(testingInstance, 0, spoolDirectoryEntryCount(testingInstance, spoolDirectory))
}

func TestService_GetDownloadURL_SingleInstanceUsesInstanceKey(testingInstance *testing.T) {
	harness := newUploadHarness()
	result, uploadError := harness.upload(instanceFixtureWithSOP("1"), "", "")
	require.NoError(testingInstance, uploadError)
	storedInstances, listError := harness.repository.ListInstances(context.Background(), result.Study.ID)
	require.NoError(testingInstance, listError)

	downloadURL, expiresAt, queryError := harness.service.GetDownloadURL(context.Background(), result.Study.ID.String())

	require.NoError(testingInstance, queryError)
	assert.Contains(testingInstance, downloadURL, storedInstances[0].ObjectKey)
	assert.True(testingInstance, expiresAt.After(result.Study.CreatedAt))
}

func TestService_GetDownloadURL_MultipleInstancesIsConflict(testingInstance *testing.T) {
	harness := newUploadHarness()
	result, uploadError := harness.upload(instanceFixtureWithSOP("1"), "", "")
	require.NoError(testingInstance, uploadError)
	_, secondError := harness.upload(instanceFixtureWithSOP("2"), "", "")
	require.NoError(testingInstance, secondError)

	_, _, queryError := harness.service.GetDownloadURL(context.Background(), result.Study.ID.String())

	requireAppError(testingInstance, queryError, "MULTIPLE_INSTANCES", 409)
}

func TestService_ListStudyInstanceDownloads_ReturnsEveryInstance(testingInstance *testing.T) {
	harness := newUploadHarness()
	result, uploadError := uploadArchiveThrough(harness, buildZIPArchive(testingInstance, []zipEntry{
		{name: "a.dcm", content: instanceFixtureWithSOP("1").bytes()},
		{name: "b.dcm", content: instanceFixtureWithSOP("2").bytes()},
	}), "study.zip")
	require.NoError(testingInstance, uploadError)

	downloads, listError := harness.service.ListStudyInstanceDownloads(context.Background(), result.Study.ID.String())

	require.NoError(testingInstance, listError)
	require.Len(testingInstance, downloads, 2)
	for _, download := range downloads {
		assert.Contains(testingInstance, download.DownloadURL, download.Instance.ObjectKey)
		assert.NotEmpty(testingInstance, download.Instance.FileName)
	}
}

func TestService_GetDownloadURL_LegacyStudyWithoutInstanceRowsUsesStudyKey(testingInstance *testing.T) {
	harness := newUploadHarness()
	legacyStudyID := uuid.New()
	harness.repository.Studies[legacyStudyID] = &imaging.ImagingStudy{
		ID:            legacyStudyID,
		PatientFhirID: fixturePatientFhirID,
		FileName:      "legacy-scan.dcm",
		Status:        imaging.ImagingStudyStatusReady,
	}

	downloadURL, _, queryError := harness.service.GetDownloadURL(context.Background(), legacyStudyID.String())

	require.NoError(testingInstance, queryError)
	assert.Contains(testingInstance, downloadURL, "pacs/"+fixturePatientFhirID+"/"+legacyStudyID.String()+".dcm")
}
