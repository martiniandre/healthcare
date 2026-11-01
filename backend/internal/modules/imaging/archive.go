package imaging

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/dicom"
)

const (
	zipLocalFileHeaderSignature = "PK\x03\x04"
	zipArchiveExtension         = ".zip"
	macOSMetadataDirectory      = "__MACOSX"
	macOSFinderMetadataFile     = ".DS_Store"
	maximumDisplayNameRunes     = 255
	archiveSpoolPattern         = "dicom-archive-*.zip"
	uploadedDICOMContentType    = "application/dicom"
	duplicateWithinArchiveText  = "SOPInstanceUID repeats an instance already accepted from this upload"
	duplicateStoredText         = "SOPInstanceUID is already stored for this study"
)

type ArchiveLimits struct {
	MaxArchiveBytes          int64
	MaxEntries               int
	MaxUncompressedBytes     int64
	MaxCompressionRatio      int64
	CompressionRatioMinBytes int64
}

func DefaultArchiveLimits() ArchiveLimits {
	return ArchiveLimits{
		MaxArchiveBytes:          MaxZIPBytes,
		MaxEntries:               MaxZIPEntries,
		MaxUncompressedBytes:     MaxZIPUncompressedBytes,
		MaxCompressionRatio:      MaxZIPCompressionRatio,
		CompressionRatioMinBytes: ZIPCompressionRatioMinBytes,
	}
}

type archiveCandidate struct {
	fileName  string
	entry     *zip.File
	validated *validatedDICOMMetadata
	warnings  []UploadWarning
}

func isArchiveUpload(leadingBytes []byte, fileName string) bool {
	if bytes.HasPrefix(leadingBytes, []byte(zipLocalFileHeaderSignature)) {
		return true
	}
	return strings.EqualFold(filepath.Ext(fileName), zipArchiveExtension)
}

func (serviceInstance *service) uploadArchive(ctx context.Context, input UploadDICOMInput, source io.Reader) (*UploadDICOMResult, error) {
	spoolFile, archiveSize, spoolError := serviceInstance.spoolArchive(source)
	if spoolError != nil {
		return nil, spoolError
	}
	defer discardSpoolFile(spoolFile)

	archiveReader, openError := zip.NewReader(spoolFile, archiveSize)
	if openError != nil {
		return nil, apperrors.WithDetail(apperrors.ErrZIPInvalid, "the file is not a readable ZIP archive")
	}

	candidates, rejections, scanError := serviceInstance.scanArchive(input, archiveReader)
	if scanError != nil {
		return nil, scanError
	}
	if len(candidates) == 0 {
		return nil, apperrors.ErrDICOMNoValidInstances
	}

	studyInstanceUID, studyError := archiveStudyInstanceUID(candidates)
	if studyError != nil {
		return nil, studyError
	}

	existingStudy, resolveError := serviceInstance.resolveExistingStudy(ctx, input.PatientFhirID, studyInstanceUID)
	if resolveError != nil {
		return nil, resolveError
	}

	identity, identityLoadError := serviceInstance.loadPatientIdentity(ctx, input.PatientFhirID)
	if identityLoadError != nil {
		return nil, identityLoadError
	}
	var warnings []UploadWarning
	for _, candidate := range candidates {
		identityWarnings, identityError := evaluateInstancePatient(input.PatientFhirID, identity, candidate.validated.source)
		if identityError != nil {
			return nil, identityError
		}
		warnings = appendUniqueWarnings(warnings, candidate.warnings...)
		warnings = appendUniqueWarnings(warnings, identityWarnings...)
	}

	accepted, duplicateRejections, filterError := serviceInstance.filterDuplicateInstances(ctx, candidates, existingStudy)
	if filterError != nil {
		return nil, filterError
	}
	rejections = append(rejections, duplicateRejections...)
	if len(accepted) == 0 {
		return nil, apperrors.ErrDICOMNoValidInstances
	}

	studyID := uuid.New()
	if existingStudy != nil {
		studyID = existingStudy.ID
	}
	instances := make([]*InstanceIndex, 0, len(accepted))
	uploadedKeys := make([]string, 0, len(accepted))
	for _, candidate := range accepted {
		instance := candidate.validated.instanceIndex
		instance.ID = uuid.New()
		instance.StudyID = studyID
		instance.FileName = candidate.fileName
		instance.ObjectKey = instanceObjectKey(input.PatientFhirID, studyID, instance.ID)
		uploadedKeys = append(uploadedKeys, instance.ObjectKey)
		if uploadError := serviceInstance.uploadArchiveEntry(ctx, candidate.entry, instance.ObjectKey); uploadError != nil {
			serviceInstance.discardStoredObjects(ctx, uploadedKeys)
			return nil, uploadError
		}
		instances = append(instances, instance)
	}

	firstAccepted := accepted[0].validated
	study, persistError := serviceInstance.persistInstances(ctx, input, existingStudy, studyID, firstAccepted, accepted[0].fileName, instances)
	if persistError != nil {
		serviceInstance.discardStoredObjects(ctx, uploadedKeys)
		return nil, persistError
	}

	if enqueueError := serviceInstance.enqueueStudyProcessing(ctx, study.ID); enqueueError != nil {
		return nil, enqueueError
	}

	return &UploadDICOMResult{
		Study:             study,
		StudyCreated:      existingStudy == nil,
		SeriesInstanceUID: firstAccepted.seriesInstanceUID,
		SOPInstanceUID:    firstAccepted.sopInstanceUID,
		Warnings:          warnings,
		InstancesAccepted: len(accepted),
		InstancesRejected: rejections,
	}, nil
}

func (serviceInstance *service) spoolArchive(source io.Reader) (*os.File, int64, error) {
	spoolFile, createError := os.CreateTemp(serviceInstance.settings.SpoolDirectory, archiveSpoolPattern)
	if createError != nil {
		return nil, 0, fmt.Errorf("failed to create archive spool file: %w", createError)
	}

	maximumBytes := serviceInstance.settings.ArchiveLimits.MaxArchiveBytes
	spooledBytes, copyError := io.Copy(spoolFile, io.LimitReader(source, maximumBytes+1))
	if copyError != nil {
		discardSpoolFile(spoolFile)
		return nil, 0, fmt.Errorf("failed to spool archive: %w", copyError)
	}
	if spooledBytes > maximumBytes {
		discardSpoolFile(spoolFile)
		return nil, 0, apperrors.WithDetail(apperrors.ErrZIPLimitsExceeded, "the archive is larger than the allowed size")
	}
	return spoolFile, spooledBytes, nil
}

func discardSpoolFile(spoolFile *os.File) {
	spoolName := spoolFile.Name()
	_ = spoolFile.Close()
	if removeError := os.Remove(spoolName); removeError != nil && !errors.Is(removeError, fs.ErrNotExist) {
		slog.Warn("failed to remove archive spool file", "error", removeError)
	}
}

func (serviceInstance *service) scanArchive(input UploadDICOMInput, archiveReader *zip.Reader) ([]archiveCandidate, []RejectedInstance, error) {
	limits := serviceInstance.settings.ArchiveLimits
	if len(archiveReader.File) > limits.MaxEntries {
		return nil, nil, apperrors.WithDetail(apperrors.ErrZIPLimitsExceeded, "the archive contains too many entries")
	}

	var dataEntries []*zip.File
	var declaredTotal int64
	for _, entry := range archiveReader.File {
		if isUnsafeArchiveEntry(entry) {
			return nil, nil, apperrors.ErrZIPUnsafeEntry
		}
		if isIgnoredArchiveEntry(entry.Name) || entry.FileInfo().IsDir() {
			continue
		}

		if entry.UncompressedSize64 > uint64(MaxDICOMUploadBytes) {
			return nil, nil, apperrors.WithDetail(apperrors.ErrZIPLimitsExceeded, "an archive entry is larger than a single DICOM upload")
		}
		entrySize := int64(entry.UncompressedSize64)
		declaredTotal += entrySize
		if declaredTotal > limits.MaxUncompressedBytes {
			return nil, nil, apperrors.WithDetail(apperrors.ErrZIPLimitsExceeded, "the archive expands beyond the allowed total size")
		}
		if exceedsCompressionRatio(entry, limits) {
			return nil, nil, apperrors.WithDetail(apperrors.ErrZIPLimitsExceeded, "an archive entry exceeds the allowed compression ratio")
		}
		dataEntries = append(dataEntries, entry)
	}

	candidates := []archiveCandidate{}
	rejections := []RejectedInstance{}
	for _, entry := range dataEntries {
		fileName := archiveEntryDisplayName(entry.Name)
		metadata, inspectionError := readArchiveEntryMetadata(entry)
		if inspectionError != nil {
			rejections = append(rejections, rejectedInstanceFromError(fileName, inspectionError))
			continue
		}
		prepared, prepareError := prepareInstance(metadata, input.Modality)
		if prepareError != nil {
			rejections = append(rejections, rejectedInstanceFromError(fileName, prepareError))
			continue
		}
		candidates = append(candidates, archiveCandidate{
			fileName:  fileName,
			entry:     entry,
			validated: prepared.validated,
			warnings:  prepared.warnings,
		})
	}
	return candidates, rejections, nil
}

func isUnsafeArchiveEntry(entry *zip.File) bool {
	if entry.Mode()&fs.ModeSymlink != 0 {
		return true
	}
	return hasUnsafeArchiveEntryName(entry.Name)
}

func hasUnsafeArchiveEntryName(entryName string) bool {
	if entryName == "" || strings.HasPrefix(entryName, "/") || strings.ContainsAny(entryName, "\\:\x00") {
		return true
	}
	for _, segment := range strings.Split(entryName, "/") {
		if segment == ".." {
			return true
		}
	}
	return false
}

func isIgnoredArchiveEntry(entryName string) bool {
	for _, segment := range strings.Split(entryName, "/") {
		if segment == macOSMetadataDirectory {
			return true
		}
	}
	return path.Base(entryName) == macOSFinderMetadataFile
}

func exceedsCompressionRatio(entry *zip.File, limits ArchiveLimits) bool {
	uncompressedBytes := int64(entry.UncompressedSize64)
	if uncompressedBytes < limits.CompressionRatioMinBytes {
		return false
	}
	compressedBytes := int64(entry.CompressedSize64)
	if compressedBytes <= 0 {
		return true
	}
	return uncompressedBytes/compressedBytes > limits.MaxCompressionRatio
}

func archiveEntryDisplayName(entryName string) string {
	sanitizedName := sanitizeLeafName(entryName)
	if sanitizedName == "" {
		sanitizedName = "instance" + defaultDICOMExtension
	}
	return truncateRunes(sanitizedName, maximumDisplayNameRunes)
}

func readArchiveEntryMetadata(entry *zip.File) (*dicom.InstanceMetadata, error) {
	entryReader, openError := entry.Open()
	if openError != nil {
		return nil, openError
	}
	defer entryReader.Close()

	inspection, inspectionError := inspectDICOMHeader(nil, entryReader)
	if inspectionError != nil {
		return nil, inspectionError
	}
	return inspection.metadata, nil
}

func archiveStudyInstanceUID(candidates []archiveCandidate) (string, error) {
	studyInstanceUID := candidates[0].validated.studyInstanceUID
	for _, candidate := range candidates[1:] {
		if candidate.validated.studyInstanceUID != studyInstanceUID {
			return "", apperrors.WithDetail(apperrors.ErrZIPInvalid, "the archive must contain instances of a single study")
		}
	}
	return studyInstanceUID, nil
}

func (serviceInstance *service) filterDuplicateInstances(ctx context.Context, candidates []archiveCandidate, existingStudy *ImagingStudy) ([]archiveCandidate, []RejectedInstance, error) {
	var accepted []archiveCandidate
	var rejections []RejectedInstance
	seenSOPInstanceUIDs := make(map[string]bool)
	for _, candidate := range candidates {
		sopInstanceUID := candidate.validated.sopInstanceUID
		if seenSOPInstanceUIDs[sopInstanceUID] {
			rejections = append(rejections, duplicateRejection(candidate.fileName, duplicateWithinArchiveText))
			continue
		}
		if existingStudy != nil {
			alreadyStored, lookupError := serviceInstance.dbRepository.HasInstance(ctx, existingStudy.ID, sopInstanceUID)
			if lookupError != nil {
				return nil, nil, fmt.Errorf("failed to check duplicate imaging instance: %w", lookupError)
			}
			if alreadyStored {
				rejections = append(rejections, duplicateRejection(candidate.fileName, duplicateStoredText))
				continue
			}
		}
		seenSOPInstanceUIDs[sopInstanceUID] = true
		accepted = append(accepted, candidate)
	}
	return accepted, rejections, nil
}

func duplicateRejection(fileName string, reason string) RejectedInstance {
	return RejectedInstance{FileName: fileName, Code: RejectionCodeDuplicateInstance, Reason: reason}
}

func (serviceInstance *service) uploadArchiveEntry(ctx context.Context, entry *zip.File, objectKey string) error {
	entryReader, openError := entry.Open()
	if openError != nil {
		return apperrors.WithDetail(apperrors.ErrZIPInvalid, "an archive entry could not be read")
	}
	defer entryReader.Close()

	guardedReader := &archiveEntryReader{
		reader:        entryReader,
		declaredBytes: int64(entry.UncompressedSize64),
	}
	uploadError := serviceInstance.storageClient.Upload(ctx, objectKey, guardedReader, uploadedDICOMContentType)
	if uploadError == nil {
		return nil
	}
	if errors.Is(uploadError, zip.ErrChecksum) || errors.Is(uploadError, zip.ErrFormat) || errors.Is(uploadError, zip.ErrAlgorithm) {
		return apperrors.WithDetail(apperrors.ErrZIPInvalid, "an archive entry is corrupted")
	}
	return fmt.Errorf("failed to upload dicom to storage: %w", uploadError)
}

type archiveEntryReader struct {
	reader        io.Reader
	declaredBytes int64
	consumedBytes int64
}

func (entryReader *archiveEntryReader) Read(buffer []byte) (int, error) {
	remainingAllowance := entryReader.declaredBytes - entryReader.consumedBytes + 1
	if remainingAllowance <= 0 {
		return 0, apperrors.ErrZIPLimitsExceeded
	}
	if int64(len(buffer)) > remainingAllowance {
		buffer = buffer[:remainingAllowance]
	}
	bytesRead, readError := entryReader.reader.Read(buffer)
	entryReader.consumedBytes += int64(bytesRead)
	if entryReader.consumedBytes > entryReader.declaredBytes {
		return bytesRead, apperrors.ErrZIPLimitsExceeded
	}
	return bytesRead, readError
}
