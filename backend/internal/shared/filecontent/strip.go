package filecontent

import "fmt"

func StripMetadata(kind Kind, content []byte) ([]byte, AnonymizationReport, error) {
	if len(content) == 0 {
		return nil, AnonymizationReport{}, ErrEmptyFile
	}
	if len(content) > MaxFileBytes {
		return nil, AnonymizationReport{}, ErrFileTooLarge
	}

	report := newReport()

	switch kind {
	case KindJPEG:
		return stripJPEG(content, report)
	case KindPNG:
		return stripPNG(content, report)
	case KindPDF:
		return stripPDF(content, report)
	default:
		return nil, AnonymizationReport{}, fmt.Errorf("%w: unsupported kind", ErrCorruptedFile)
	}
}
