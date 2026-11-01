package tests

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/modules/imaging/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/dicom"
	"github.com/healthcare/backend/internal/shared/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_UploadDICOMStream_Success(testingInstance *testing.T) {
	harness := newUploadHarness()

	result, uploadError := harness.upload(newDICOMFixture(), "Brain MRI", "MR")

	require.NoError(testingInstance, uploadError)
	assert.Equal(testingInstance, fixturePatientFhirID, result.Study.PatientFhirID)
	assert.Equal(testingInstance, "MR", result.Study.Modality)
	assert.Equal(testingInstance, "Brain MRI", result.Study.Title)
	assert.Equal(testingInstance, imaging.ImagingStudyStatusReceiving, result.Study.Status)
	assert.Equal(testingInstance, fixtureStudyInstanceUID, result.Study.StudyInstanceUID)
	assert.Equal(testingInstance, "ACC-0042", result.Study.AccessionNumber)
	assert.Equal(testingInstance, "BRAIN WITHOUT CONTRAST", result.Study.StudyDescription)
	require.NotNil(testingInstance, result.Study.StudyDateTime)
	assert.Equal(testingInstance, expectedStudyTime("20060102150405", "20260115101500"), *result.Study.StudyDateTime)
	assert.Equal(testingInstance, fixtureSeriesUID, result.SeriesInstanceUID)
	assert.Equal(testingInstance, fixtureSOPInstanceUID, result.SOPInstanceUID)
	assert.Empty(testingInstance, result.Warnings)

	instanceIndex, indexError := harness.repository.GetInstanceIndex(context.Background(), result.Study.ID)
	require.NoError(testingInstance, indexError)
	assert.Equal(testingInstance, fixtureSOPInstanceUID, instanceIndex.SOPInstanceUID)
	assert.Equal(testingInstance, fixtureSOPClassUID, instanceIndex.SOPClassUID)
	assert.Equal(testingInstance, dicom.TransferSyntaxExplicitVRLittleEndian, instanceIndex.TransferSyntaxUID)

	retrievedStudy, queryError := harness.service.GetImagingStudy(context.Background(), result.Study.ID.String())
	require.NoError(testingInstance, queryError)
	assert.Equal(testingInstance, result.Study.ID, retrievedStudy.ID)
}

func TestService_UploadDICOMStream_ImplicitVRIsAccepted(testingInstance *testing.T) {
	harness := newUploadHarness()
	fixture := newDICOMFixture()
	fixture.transferSyntax = dicom.TransferSyntaxImplicitVRLittleEndian

	result, uploadError := harness.upload(fixture, "", "")

	require.NoError(testingInstance, uploadError)
	assert.Equal(testingInstance, fixtureStudyInstanceUID, result.Study.StudyInstanceUID)
	assert.Equal(testingInstance, "MR", result.Study.Modality)
}

func TestService_UploadDICOMStream_FileModalityWinsAndWarnsWhenUserDiffers(testingInstance *testing.T) {
	harness := newUploadHarness()

	result, uploadError := harness.upload(newDICOMFixture(), "Brain", "CT")

	require.NoError(testingInstance, uploadError)
	assert.Equal(testingInstance, "MR", result.Study.Modality)
	assert.Equal(testingInstance, []string{imaging.WarningCodeModalityOverridden}, warningCodes(result.Warnings))
}

func TestService_UploadDICOMStream_UserModalityIsFallbackWhenFileHasNone(testingInstance *testing.T) {
	harness := newUploadHarness()
	fixture := newDICOMFixture()
	fixture.modality = ""

	result, uploadError := harness.upload(fixture, "Brain", "ct")

	require.NoError(testingInstance, uploadError)
	assert.Equal(testingInstance, "CT", result.Study.Modality)
	assert.Empty(testingInstance, result.Warnings)
}

func TestService_UploadDICOMStream_RejectsWhenNoModalityAnywhere(testingInstance *testing.T) {
	harness := newUploadHarness()
	fixture := newDICOMFixture()
	fixture.modality = ""

	_, uploadError := harness.upload(fixture, "Brain", "")

	assert.ErrorIs(testingInstance, uploadError, apperrors.ErrInvalidDICOM)
	assert.Contains(testingInstance, uploadError.Error(), "Modality is missing")
}

func TestService_UploadDICOMStream_TitleDerivation(testingInstance *testing.T) {
	testCases := []struct {
		name          string
		description   string
		studyDate     string
		typedTitle    string
		expectedTitle string
	}{
		{name: "modality and description", description: "KNEE", studyDate: "20260115", expectedTitle: "MR - KNEE"},
		{name: "modality and date when no description", description: "", studyDate: "20260115", expectedTitle: "MR - 2026-01-15"},
		{name: "modality only", description: "", studyDate: "", expectedTitle: "MR"},
		{name: "typed title wins", description: "KNEE", studyDate: "20260115", typedTitle: "  My title  ", expectedTitle: "My title"},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(testingInstance *testing.T) {
			harness := newUploadHarness()
			fixture := newDICOMFixture()
			fixture.studyDescription = testCase.description
			fixture.studyDate = testCase.studyDate

			result, uploadError := harness.upload(fixture, testCase.typedTitle, "")

			require.NoError(testingInstance, uploadError)
			assert.Equal(testingInstance, testCase.expectedTitle, result.Study.Title)
		})
	}
}

func TestService_UploadDICOMStream_StudyDateWithoutTimeUsesMidnight(testingInstance *testing.T) {
	harness := newUploadHarness()
	fixture := newDICOMFixture()
	fixture.studyTime = ""

	result, uploadError := harness.upload(fixture, "", "")

	require.NoError(testingInstance, uploadError)
	require.NotNil(testingInstance, result.Study.StudyDateTime)
	assert.Equal(testingInstance, expectedStudyTime("20060102150405", "20260115000000"), *result.Study.StudyDateTime)
}

func TestService_UploadDICOMStream_InvalidDICOMCases(testingInstance *testing.T) {
	testCases := []struct {
		name           string
		mutate         func(fixture *dicomFixture)
		expectedReason string
	}{
		{name: "missing study instance uid", mutate: func(fixture *dicomFixture) { fixture.studyInstanceUID = "" }, expectedReason: "StudyInstanceUID is missing"},
		{name: "malformed study instance uid", mutate: func(fixture *dicomFixture) { fixture.studyInstanceUID = "not-a-uid" }, expectedReason: "StudyInstanceUID is not a valid UID"},
		{name: "missing sop instance uid", mutate: func(fixture *dicomFixture) { fixture.sopInstanceUID = "" }, expectedReason: "SOPInstanceUID is missing"},
		{name: "unsupported transfer syntax", mutate: func(fixture *dicomFixture) { fixture.transferSyntax = dicom.TransferSyntaxJPEGBaseline }, expectedReason: "JPEG Baseline"},
		{name: "big endian transfer syntax", mutate: func(fixture *dicomFixture) { fixture.transferSyntax = dicom.TransferSyntaxExplicitVRBigEndian }, expectedReason: "Explicit VR Big Endian"},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(testingInstance *testing.T) {
			harness := newUploadHarness()
			fixture := newDICOMFixture()
			testCase.mutate(&fixture)

			_, uploadError := harness.upload(fixture, "", "")

			assert.ErrorIs(testingInstance, uploadError, apperrors.ErrInvalidDICOM)
			var appError apperrors.AppError
			require.True(testingInstance, errors.As(uploadError, &appError))
			assert.Equal(testingInstance, "INVALID_DICOM", appError.Code)
			assert.Contains(testingInstance, appError.Message, testCase.expectedReason)
			assert.Empty(testingInstance, harness.storage.uploadedObjects)
			assert.Empty(testingInstance, harness.repository.Studies)
		})
	}
}

func TestService_UploadDICOMStream_InvalidMagicBytes(testingInstance *testing.T) {
	harness := newUploadHarness()
	invalidDICOMBytes := make([]byte, 200)
	copy(invalidDICOMBytes[128:132], []byte("JPEG"))

	_, uploadError := harness.service.UploadDICOMStream(context.Background(), imaging.UploadDICOMInput{
		PatientFhirID: fixturePatientFhirID,
		Title:         "Brain MRI",
		Modality:      "MR",
	}, bytes.NewReader(invalidDICOMBytes))

	assert.ErrorIs(testingInstance, uploadError, apperrors.ErrInvalidDICOM)
	assert.Contains(testingInstance, uploadError.Error(), "magic bytes DICM signature missing")
}

func TestService_UploadDICOMStream_MagicOnlyFileHasNoMetaGroup(testingInstance *testing.T) {
	harness := newUploadHarness()
	magicOnlyBytes := make([]byte, 200)
	copy(magicOnlyBytes[128:132], []byte("DICM"))

	_, uploadError := harness.service.UploadDICOMStream(context.Background(), imaging.UploadDICOMInput{
		PatientFhirID: fixturePatientFhirID,
		Title:         "Brain MRI",
		Modality:      "MR",
	}, bytes.NewReader(magicOnlyBytes))

	assert.ErrorIs(testingInstance, uploadError, apperrors.ErrInvalidDICOM)
}

func TestService_UploadDICOMStream_TooSmall(testingInstance *testing.T) {
	harness := newUploadHarness()

	_, uploadError := harness.service.UploadDICOMStream(context.Background(), imaging.UploadDICOMInput{
		PatientFhirID: fixturePatientFhirID,
		Title:         "Brain MRI",
		Modality:      "MR",
	}, bytes.NewReader([]byte("too-small")))

	assert.ErrorIs(testingInstance, uploadError, apperrors.ErrInvalidDICOM)
	assert.Contains(testingInstance, uploadError.Error(), "preamble is too small")
}

func TestService_UploadDICOMStream_InvalidInput(testingInstance *testing.T) {
	harness := newUploadHarness()

	_, uploadError := harness.service.UploadDICOMStream(context.Background(), imaging.UploadDICOMInput{
		PatientFhirID: "",
		Modality:      "NOT_A_MODALITY",
	}, bytes.NewReader([]byte("anything")))

	var appError apperrors.AppError
	require.ErrorAs(testingInstance, uploadError, &appError)
	assert.Equal(testingInstance, 400, appError.HTTPCode)
	assert.Contains(testingInstance, appError.Message, "patient_fhir_id")
	assert.Contains(testingInstance, appError.Message, "modality")
	assert.NotContains(testingInstance, appError.Message, "title")
}

func TestService_UploadDICOMStream_PatientIdentityRules(testingInstance *testing.T) {
	testCases := []struct {
		name             string
		mutateFixture    func(fixture *dicomFixture)
		patientIdentity  imaging.PatientIdentity
		expectedWarnings []string
		expectMismatch   bool
	}{
		{
			name:             "formatted cpf matches",
			mutateFixture:    func(fixture *dicomFixture) {},
			patientIdentity:  imaging.PatientIdentity{FullNames: []string{"Maria Silva"}, DocumentIDs: []string{"52998224725"}},
			expectedWarnings: []string{},
		},
		{
			name:             "cpf match wins over a differently written name",
			mutateFixture:    func(fixture *dicomFixture) { fixture.patientName = "MARIA S" },
			patientIdentity:  imaging.PatientIdentity{FullNames: []string{"Maria Silva"}, DocumentIDs: []string{"52998224725"}},
			expectedWarnings: []string{},
		},
		{
			name:            "different cpf is rejected",
			mutateFixture:   func(fixture *dicomFixture) { fixture.patientID = "111.444.777-35" },
			patientIdentity: imaging.PatientIdentity{FullNames: []string{"Maria Silva"}, DocumentIDs: []string{"52998224725"}},
			expectMismatch:  true,
		},
		{
			name:            "different cpf is rejected even when the name matches",
			mutateFixture:   func(fixture *dicomFixture) { fixture.patientID = "11144477735" },
			patientIdentity: imaging.PatientIdentity{FullNames: []string{"Maria Silva"}, DocumentIDs: []string{"52998224725"}},
			expectMismatch:  true,
		},
		{
			name:             "name only matches accent and case insensitive",
			mutateFixture:    func(fixture *dicomFixture) { fixture.patientID = "MRN-77"; fixture.patientName = "JOSÉ^CONCEIÇÃO" },
			patientIdentity:  imaging.PatientIdentity{FullNames: []string{"Jose Conceicao"}, DocumentIDs: []string{"52998224725"}},
			expectedWarnings: []string{},
		},
		{
			name:             "name matches with middle name and particle differences",
			mutateFixture:    func(fixture *dicomFixture) { fixture.patientID = ""; fixture.patientName = "SILVA^MARIA DE FATIMA" },
			patientIdentity:  imaging.PatientIdentity{FullNames: []string{"Maria Silva"}},
			expectedWarnings: []string{},
		},
		{
			name:            "clearly different name is rejected",
			mutateFixture:   func(fixture *dicomFixture) { fixture.patientID = "MRN-77"; fixture.patientName = "SANTOS^PEDRO" },
			patientIdentity: imaging.PatientIdentity{FullNames: []string{"Maria Silva"}, DocumentIDs: []string{"52998224725"}},
			expectMismatch:  true,
		},
		{
			name:            "same first name but different surname is rejected",
			mutateFixture:   func(fixture *dicomFixture) { fixture.patientID = ""; fixture.patientName = "SANTOS^MARIA" },
			patientIdentity: imaging.PatientIdentity{FullNames: []string{"Maria Silva"}},
			expectMismatch:  true,
		},
		{
			name:             "patient without identifiers warns",
			mutateFixture:    func(fixture *dicomFixture) {},
			patientIdentity:  imaging.PatientIdentity{},
			expectedWarnings: []string{imaging.WarningCodePatientIdentityMissing},
		},
		{
			name:             "non document patient id without names warns",
			mutateFixture:    func(fixture *dicomFixture) { fixture.patientID = "MRN-77"; fixture.patientName = "" },
			patientIdentity:  imaging.PatientIdentity{FullNames: []string{"Maria Silva"}, DocumentIDs: []string{"52998224725"}},
			expectedWarnings: []string{imaging.WarningCodePatientIdentityMissing},
		},
		{
			name:             "anonymized by ANON patient id",
			mutateFixture:    func(fixture *dicomFixture) { fixture.patientID = "ANON-4821"; fixture.patientName = "ANONYMOUS" },
			patientIdentity:  imaging.PatientIdentity{FullNames: []string{"Maria Silva"}, DocumentIDs: []string{"52998224725"}},
			expectedWarnings: []string{imaging.WarningCodeAnonymizedStudy},
		},
		{
			name:             "anonymized by blank patient id and name",
			mutateFixture:    func(fixture *dicomFixture) { fixture.patientID = ""; fixture.patientName = "" },
			patientIdentity:  imaging.PatientIdentity{FullNames: []string{"Maria Silva"}, DocumentIDs: []string{"52998224725"}},
			expectedWarnings: []string{imaging.WarningCodeAnonymizedStudy},
		},
		{
			name:             "dicom patient id equal to the fhir id matches",
			mutateFixture:    func(fixture *dicomFixture) { fixture.patientID = "PATIENT-123"; fixture.patientName = "OTHER^NAME" },
			patientIdentity:  imaging.PatientIdentity{FullNames: []string{"Maria Silva"}},
			expectedWarnings: []string{},
		},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(testingInstance *testing.T) {
			harness := newUploadHarness()
			harness.identity.Identity = testCase.patientIdentity
			fixture := newDICOMFixture()
			testCase.mutateFixture(&fixture)

			result, uploadError := harness.upload(fixture, "", "")

			if testCase.expectMismatch {
				assert.ErrorIs(testingInstance, uploadError, apperrors.ErrDICOMPatientMismatch)
				assert.Empty(testingInstance, harness.storage.uploadedObjects)
				assert.Empty(testingInstance, harness.repository.Studies)
				return
			}
			require.NoError(testingInstance, uploadError)
			assert.Equal(testingInstance, testCase.expectedWarnings, warningCodes(result.Warnings))
		})
	}
}

func TestService_UploadDICOMStream_MismatchErrorDoesNotEchoPatientData(testingInstance *testing.T) {
	harness := newUploadHarness()
	fixture := newDICOMFixture()
	fixture.patientID = "111.444.777-35"
	fixture.patientName = "CONFIDENCIAL^PACIENTE"

	_, uploadError := harness.upload(fixture, "", "")

	var appError apperrors.AppError
	require.True(testingInstance, errors.As(uploadError, &appError))
	assert.Equal(testingInstance, "DICOM_PATIENT_MISMATCH", appError.Code)
	assert.Equal(testingInstance, 422, appError.HTTPCode)
	for _, forbiddenFragment := range []string{"111", "CONFIDENCIAL", "PACIENTE", "52998224725", "Maria"} {
		assert.NotContains(testingInstance, uploadError.Error(), forbiddenFragment)
	}
}

func TestService_UploadDICOMStream_WithoutIdentityReaderWarnsIdentityMissing(testingInstance *testing.T) {
	repository := mocks.NewMockRepository()
	imagingService := imaging.NewService(repository, newRecordingStorage(), nil, nil, nil)

	result, uploadError := imagingService.UploadDICOMStream(context.Background(), imaging.UploadDICOMInput{
		PatientFhirID: fixturePatientFhirID,
	}, newDICOMFixture().reader())

	require.NoError(testingInstance, uploadError)
	assert.Equal(testingInstance, []string{imaging.WarningCodePatientIdentityMissing}, warningCodes(result.Warnings))
}

func TestService_UploadDICOMStream_PatientLookupFailures(testingInstance *testing.T) {
	harness := newUploadHarness()
	harness.identity.Err = apperrors.ErrPatientNotFound

	_, notFoundError := harness.upload(newDICOMFixture(), "", "")
	assert.ErrorIs(testingInstance, notFoundError, apperrors.ErrPatientNotFound)

	harness.identity.Err = errors.New("fhir unavailable")
	_, unavailableError := harness.upload(newDICOMFixture(), "", "")
	require.Error(testingInstance, unavailableError)
	assert.Empty(testingInstance, harness.storage.uploadedObjects)
	assert.Empty(testingInstance, harness.repository.Studies)
}

func TestService_UploadDICOMStream_DuplicateStudyIsRejectedBeforeStorage(testingInstance *testing.T) {
	harness := newUploadHarness()
	_, firstError := harness.upload(newDICOMFixture(), "", "")
	require.NoError(testingInstance, firstError)
	objectsAfterFirstUpload := len(harness.storage.uploadedObjects)

	_, secondError := harness.upload(newDICOMFixture(), "", "")

	assert.ErrorIs(testingInstance, secondError, apperrors.ErrDICOMDuplicateInstance)
	var appError apperrors.AppError
	require.True(testingInstance, errors.As(secondError, &appError))
	assert.Equal(testingInstance, "DICOM_DUPLICATE_INSTANCE", appError.Code)
	assert.Equal(testingInstance, 409, appError.HTTPCode)
	assert.Len(testingInstance, harness.storage.uploadedObjects, objectsAfterFirstUpload)
	assert.Len(testingInstance, harness.repository.Studies, 1)
}

type racingDuplicateRepository struct {
	*mocks.MockRepository
}

func (racing *racingDuplicateRepository) CreateImagingStudy(ctx context.Context, study *imaging.ImagingStudy, instanceIndex *imaging.InstanceIndex) error {
	return apperrors.ErrDICOMDuplicateInstance
}

func TestService_UploadDICOMStream_RaceOnUniqueIndexDiscardsStoredObject(testingInstance *testing.T) {
	recording := newRecordingStorage()
	imagingService := imaging.NewService(&racingDuplicateRepository{MockRepository: mocks.NewMockRepository()}, recording, nil, nil, nil)

	_, uploadError := imagingService.UploadDICOMStream(context.Background(), imaging.UploadDICOMInput{
		PatientFhirID: fixturePatientFhirID,
	}, newDICOMFixture().reader())

	assert.ErrorIs(testingInstance, uploadError, apperrors.ErrDICOMDuplicateInstance)
	assert.Len(testingInstance, recording.deletedKeys, 1)
	assert.Empty(testingInstance, recording.uploadedObjects)
}

func TestService_UploadDICOMStream_RepositoryFailureDiscardsStoredObject(testingInstance *testing.T) {
	harness := newUploadHarness()
	harness.repository.Err = errors.New("database down")

	_, uploadError := harness.upload(newDICOMFixture(), "", "")

	require.Error(testingInstance, uploadError)
	assert.Empty(testingInstance, harness.storage.uploadedObjects)
}

func TestService_UploadDICOMStream_StorageFailureIsPropagated(testingInstance *testing.T) {
	harness := newUploadHarness()
	harness.storage.uploadError = errors.New("bucket unavailable")

	_, uploadError := harness.upload(newDICOMFixture(), "", "")

	require.Error(testingInstance, uploadError)
	assert.Contains(testingInstance, uploadError.Error(), "failed to upload dicom to storage")
	assert.Empty(testingInstance, harness.repository.Studies)
}

func TestService_UploadDICOMStream_StoresTheOriginalBytesUnchanged(testingInstance *testing.T) {
	harness := newUploadHarness()
	fixture := newDICOMFixture()
	fixture.pixelPaddingBytes = 6 << 20
	originalBytes := fixture.bytes()

	result, uploadError := harness.upload(fixture, "", "")

	require.NoError(testingInstance, uploadError)
	require.Len(testingInstance, harness.storage.uploadedObjects, 1)
	for objectKey, storedBytes := range harness.storage.uploadedObjects {
		assert.Contains(testingInstance, objectKey, result.Study.ID.String())
		assert.True(testingInstance, bytes.Equal(originalBytes, storedBytes))
	}
}

type countingReader struct {
	reader   io.Reader
	consumed int64
}

func (counting *countingReader) Read(buffer []byte) (int, error) {
	bytesRead, readError := counting.reader.Read(buffer)
	counting.consumed += int64(bytesRead)
	return bytesRead, readError
}

type consumptionProbeRepository struct {
	*mocks.MockRepository
	probe               func() int64
	consumedBeforeStore int64
}

func (probing *consumptionProbeRepository) FindImagingStudyByStudyInstanceUID(ctx context.Context, studyInstanceUID string) (*imaging.ImagingStudy, error) {
	probing.consumedBeforeStore = probing.probe()
	return probing.MockRepository.FindImagingStudyByStudyInstanceUID(ctx, studyInstanceUID)
}

func TestService_UploadDICOMStream_ReadsOnlyABoundedHeaderBeforeStreamingToStorage(testingInstance *testing.T) {
	fixture := newDICOMFixture()
	fixture.pixelPaddingBytes = 24 << 20
	source := &countingReader{reader: fixture.reader()}
	probingRepository := &consumptionProbeRepository{MockRepository: mocks.NewMockRepository(), probe: func() int64 { return source.consumed }}
	recording := newRecordingStorage()
	imagingService := imaging.NewService(probingRepository, recording, nil, nil, nil)

	_, uploadError := imagingService.UploadDICOMStream(context.Background(), imaging.UploadDICOMInput{PatientFhirID: fixturePatientFhirID}, source)

	require.NoError(testingInstance, uploadError)
	assert.Less(testingInstance, probingRepository.consumedBeforeStore, int64(1<<20))
	assert.Greater(testingInstance, source.consumed, int64(24<<20))
}

func TestService_GetImagingStudy_NotFound(testingInstance *testing.T) {
	harness := newUploadHarness()

	study, queryError := harness.service.GetImagingStudy(context.Background(), "no-such-uuid")
	var appError apperrors.AppError
	assert.ErrorAs(testingInstance, queryError, &appError)
	assert.Equal(testingInstance, 400, appError.HTTPCode)
	assert.Nil(testingInstance, study)

	_, notFoundError := harness.service.GetImagingStudy(context.Background(), uuid.Nil.String())
	assert.ErrorIs(testingInstance, notFoundError, apperrors.ErrImagingStudyNotFound)
}

func TestService_GetDownloadURL(testingInstance *testing.T) {
	harness := newUploadHarness()
	harness.service = imaging.NewService(harness.repository, storage.NewStorageClient("test-bucket", "test"), nil, nil, harness.identity)

	result, uploadError := harness.upload(newDICOMFixture(), "Brain MRI", "MR")
	require.NoError(testingInstance, uploadError)

	downloadURL, expiresAt, queryError := harness.service.GetDownloadURL(context.Background(), result.Study.ID.String())
	assert.NoError(testingInstance, queryError)
	assert.NotEmpty(testingInstance, downloadURL)
	assert.True(testingInstance, expiresAt.After(result.Study.CreatedAt))
}
