package filecontent

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var pngKeptChunkTypes = map[string]bool{
	"IHDR": true,
	"PLTE": true,
	"IDAT": true,
	"IEND": true,
	"tRNS": true,
	"gAMA": true,
	"cHRM": true,
	"sRGB": true,
	"iCCP": true,
	"sBIT": true,
	"bKGD": true,
	"hIST": true,
	"pHYs": true,
	"sPLT": true,
	"cICP": true,
	"acTL": true,
	"fcTL": true,
	"fdAT": true,
}

func stripPNG(content []byte, report AnonymizationReport) ([]byte, AnonymizationReport, error) {
	if !bytes.HasPrefix(content, pngSignature) {
		return nil, AnonymizationReport{}, fmt.Errorf("%w: png signature missing", ErrCorruptedFile)
	}

	cleaned := make([]byte, 0, len(content))
	cleaned = append(cleaned, pngSignature...)
	offset := len(pngSignature)
	removedChunks := 0
	trailingRemoved := false
	foundEnd := false

	for offset+pngChunkHeaderBytes+pngChunkTrailerBytes <= len(content) && !foundEnd {
		chunkLength := int(binary.BigEndian.Uint32(content[offset:]))
		chunkType := string(content[offset+4 : offset+8])
		chunkEnd := offset + pngChunkHeaderBytes + chunkLength + pngChunkTrailerBytes
		if chunkLength < 0 || chunkEnd > len(content) {
			return nil, AnonymizationReport{}, fmt.Errorf("%w: png chunk %q out of bounds", ErrCorruptedFile, chunkType)
		}

		if pngKeptChunkTypes[chunkType] {
			cleaned = append(cleaned, content[offset:chunkEnd]...)
		} else {
			removedChunks++
		}
		if chunkType == "IEND" {
			foundEnd = true
			trailingRemoved = chunkEnd < len(content)
		}
		offset = chunkEnd
	}

	if !foundEnd {
		return nil, AnonymizationReport{}, fmt.Errorf("%w: png IEND chunk missing", ErrCorruptedFile)
	}

	if removedChunks > 0 {
		report.AddStep(StepMetadataRemoved)
	} else {
		report.AddStep(StepNoMetadataFound)
	}
	if trailingRemoved {
		report.AddStep(StepTrailingDataRemoved)
	}
	return cleaned, report, nil
}
