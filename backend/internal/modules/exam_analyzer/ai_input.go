package exam_analyzer

import (
	"errors"
	"fmt"
	"io"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/filecontent"
)

type aiInput struct {
	FileName string
	Bytes    []byte
	Report   *filecontent.AnonymizationReport
}

func prepareAIInput(analysis *ExamAnalysis, fileBytes []byte) (aiInput, error) {
	if !analysis.Anonymized {
		return aiInput{FileName: analysis.FileName, Bytes: fileBytes}, nil
	}

	kind := filecontent.Detect(headOf(fileBytes))
	cleanBytes, report, stripErr := filecontent.StripMetadata(kind, fileBytes)
	if stripErr != nil {
		return aiInput{}, stripErr
	}
	report.AddStep(filecontent.StepFileNameNeutralized)
	report.Finalize()

	return aiInput{
		FileName: filecontent.NeutralFileName(kind),
		Bytes:    cleanBytes,
		Report:   &report,
	}, nil
}

func headOf(content []byte) []byte {
	if len(content) > filecontent.HeadBytes {
		return content[:filecontent.HeadBytes]
	}
	return content
}

func readBoundedExamFile(fileContent io.Reader) ([]byte, error) {
	fileBytes, readErr := io.ReadAll(io.LimitReader(fileContent, filecontent.MaxFileBytes+1))
	if readErr != nil {
		return nil, fmt.Errorf("failed to read uploaded exam file: %w", readErr)
	}
	if len(fileBytes) > filecontent.MaxFileBytes {
		return nil, filecontent.ErrFileTooLarge
	}
	return fileBytes, nil
}

func mapFileContentError(err error) error {
	switch {
	case errors.Is(err, filecontent.ErrEmptyFile):
		return apperrors.ErrFileEmpty
	case errors.Is(err, filecontent.ErrContentExtensionMismatch):
		return apperrors.ErrFileContentMismatch
	case errors.Is(err, filecontent.ErrCorruptedFile):
		return apperrors.ErrFileCorrupted
	case errors.Is(err, filecontent.ErrFileTooLarge):
		return apperrors.ErrPayloadTooLarge
	default:
		return err
	}
}
