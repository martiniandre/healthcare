package filecontent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

const (
	pngChunkHeaderBytes  = 8
	pngChunkTrailerBytes = 4
	pngHeaderChunkLength = 13
	pdfTailBytes         = 1024
	pdfEndOfFileMarker   = "%%EOF"
	pdfObjectMarker      = " obj"
)

var jpegEndOfImage = []byte{0xff, 0xd9}

func Validate(kind Kind, fileName string, content []byte) error {
	if len(content) == 0 {
		return ErrEmptyFile
	}
	if len(content) > MaxFileBytes {
		return ErrFileTooLarge
	}
	if kind == KindUnknown {
		return fmt.Errorf("%w: unrecognized file signature", ErrCorruptedFile)
	}
	if KindFromFileName(fileName) != kind {
		return fmt.Errorf("%w: content is %s", ErrContentExtensionMismatch, kind)
	}

	switch kind {
	case KindPNG:
		return validatePNG(content)
	case KindJPEG:
		return validateJPEG(content)
	case KindPDF:
		return validatePDF(content)
	default:
		return fmt.Errorf("%w: unsupported kind", ErrCorruptedFile)
	}
}

func validatePNG(content []byte) error {
	if len(content) < len(pngSignature)+pngChunkHeaderBytes+pngHeaderChunkLength+pngChunkTrailerBytes {
		return fmt.Errorf("%w: png too short", ErrCorruptedFile)
	}
	firstChunkOffset := len(pngSignature)
	if binary.BigEndian.Uint32(content[firstChunkOffset:]) != pngHeaderChunkLength || string(content[firstChunkOffset+4:firstChunkOffset+8]) != "IHDR" {
		return fmt.Errorf("%w: png IHDR chunk missing", ErrCorruptedFile)
	}
	if _, _, decodeErr := image.DecodeConfig(bytes.NewReader(content)); decodeErr != nil {
		return fmt.Errorf("%w: png header cannot be decoded", ErrCorruptedFile)
	}
	if !pngHasEndChunk(content) {
		return fmt.Errorf("%w: png IEND chunk missing", ErrCorruptedFile)
	}
	return nil
}

func pngHasEndChunk(content []byte) bool {
	offset := len(pngSignature)
	for offset+pngChunkHeaderBytes+pngChunkTrailerBytes <= len(content) {
		chunkLength := int(binary.BigEndian.Uint32(content[offset:]))
		chunkType := string(content[offset+4 : offset+8])
		nextOffset := offset + pngChunkHeaderBytes + chunkLength + pngChunkTrailerBytes
		if chunkLength < 0 || nextOffset > len(content) {
			return false
		}
		if chunkType == "IEND" {
			return true
		}
		offset = nextOffset
	}
	return false
}

func validateJPEG(content []byte) error {
	if len(content) < len(jpegSignature)+len(jpegEndOfImage) {
		return fmt.Errorf("%w: jpeg too short", ErrCorruptedFile)
	}
	if !bytes.Contains(content[2:], jpegEndOfImage) {
		return fmt.Errorf("%w: jpeg EOI marker missing", ErrCorruptedFile)
	}
	if _, _, decodeErr := image.DecodeConfig(bytes.NewReader(content)); decodeErr != nil {
		return fmt.Errorf("%w: jpeg header cannot be decoded", ErrCorruptedFile)
	}
	return nil
}

func validatePDF(content []byte) error {
	if !bytes.HasPrefix(content, pdfSignature) {
		return fmt.Errorf("%w: pdf header missing", ErrCorruptedFile)
	}
	tailStart := len(content) - pdfTailBytes
	if tailStart < 0 {
		tailStart = 0
	}
	if !bytes.Contains(content[tailStart:], []byte(pdfEndOfFileMarker)) {
		return fmt.Errorf("%w: pdf EOF marker missing", ErrCorruptedFile)
	}
	if !bytes.Contains(content, []byte(pdfObjectMarker)) {
		return fmt.Errorf("%w: pdf has no objects", ErrCorruptedFile)
	}
	return nil
}
