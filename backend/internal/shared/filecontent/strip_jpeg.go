package filecontent

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	jpegMarkerPrefix      = 0xff
	jpegMarkerSOI         = 0xd8
	jpegMarkerEOI         = 0xd9
	jpegMarkerSOS         = 0xda
	jpegMarkerCOM         = 0xfe
	jpegMarkerAPP0        = 0xe0
	jpegMarkerAPP2        = 0xe2
	jpegMarkerAPP14       = 0xee
	jpegMarkerAPP15       = 0xef
	jpegMarkerTEM         = 0x01
	jpegMarkerRSTFirst    = 0xd0
	jpegMarkerRSTLast     = 0xd7
	jpegSegmentLengthSize = 2
)

var (
	jpegICCIdentifier   = []byte("ICC_PROFILE\x00")
	jpegJFIFIdentifier  = []byte("JFIF\x00")
	jpegAdobeIdentifier = []byte("Adobe")
)

func jpegSegmentIsKept(marker byte, payload []byte) bool {
	if marker == jpegMarkerCOM {
		return false
	}
	if marker < jpegMarkerAPP0 || marker > jpegMarkerAPP15 {
		return true
	}
	switch marker {
	case jpegMarkerAPP0:
		return bytes.HasPrefix(payload, jpegJFIFIdentifier)
	case jpegMarkerAPP2:
		return bytes.HasPrefix(payload, jpegICCIdentifier)
	case jpegMarkerAPP14:
		return bytes.HasPrefix(payload, jpegAdobeIdentifier)
	default:
		return false
	}
}

func stripJPEG(content []byte, report AnonymizationReport) ([]byte, AnonymizationReport, error) {
	if len(content) < 4 || content[0] != jpegMarkerPrefix || content[1] != jpegMarkerSOI {
		return nil, AnonymizationReport{}, fmt.Errorf("%w: jpeg SOI missing", ErrCorruptedFile)
	}

	cleaned := make([]byte, 0, len(content))
	cleaned = append(cleaned, content[0], content[1])
	offset := 2
	removedSegments := 0
	trailingRemoved := false
	finished := false

	for offset < len(content) && !finished {
		if content[offset] != jpegMarkerPrefix {
			return nil, AnonymizationReport{}, fmt.Errorf("%w: jpeg marker expected at offset %d", ErrCorruptedFile, offset)
		}
		for offset < len(content) && content[offset] == jpegMarkerPrefix {
			offset++
		}
		if offset >= len(content) {
			return nil, AnonymizationReport{}, fmt.Errorf("%w: jpeg truncated after marker prefix", ErrCorruptedFile)
		}
		marker := content[offset]
		offset++

		if marker == jpegMarkerEOI {
			cleaned = append(cleaned, jpegMarkerPrefix, jpegMarkerEOI)
			trailingRemoved = offset < len(content)
			finished = true
			continue
		}
		if marker == jpegMarkerTEM || (marker >= jpegMarkerRSTFirst && marker <= jpegMarkerRSTLast) {
			cleaned = append(cleaned, jpegMarkerPrefix, marker)
			continue
		}

		if offset+jpegSegmentLengthSize > len(content) {
			return nil, AnonymizationReport{}, fmt.Errorf("%w: jpeg segment length truncated", ErrCorruptedFile)
		}
		segmentLength := int(binary.BigEndian.Uint16(content[offset:]))
		segmentEnd := offset + segmentLength
		if segmentLength < jpegSegmentLengthSize || segmentEnd > len(content) {
			return nil, AnonymizationReport{}, fmt.Errorf("%w: jpeg segment length out of bounds", ErrCorruptedFile)
		}

		if marker == jpegMarkerSOS {
			cleaned = append(cleaned, jpegMarkerPrefix, marker)
			cleaned = append(cleaned, content[offset:segmentEnd]...)
			scanData, foundEnd := splitJPEGScanData(content[segmentEnd:])
			cleaned = append(cleaned, scanData...)
			trailingRemoved = foundEnd && segmentEnd+len(scanData) < len(content)
			finished = true
			continue
		}

		payload := content[offset+jpegSegmentLengthSize : segmentEnd]
		if jpegSegmentIsKept(marker, payload) {
			cleaned = append(cleaned, jpegMarkerPrefix, marker)
			cleaned = append(cleaned, content[offset:segmentEnd]...)
		} else {
			removedSegments++
		}
		offset = segmentEnd
	}

	if !finished {
		return nil, AnonymizationReport{}, fmt.Errorf("%w: jpeg ended without SOS or EOI", ErrCorruptedFile)
	}

	if removedSegments > 0 {
		report.AddStep(StepMetadataRemoved)
	} else {
		report.AddStep(StepNoMetadataFound)
	}
	if trailingRemoved {
		report.AddStep(StepTrailingDataRemoved)
	}
	return cleaned, report, nil
}

func splitJPEGScanData(remaining []byte) ([]byte, bool) {
	endIndex := bytes.Index(remaining, []byte{jpegMarkerPrefix, jpegMarkerEOI})
	if endIndex < 0 {
		return remaining, false
	}
	return remaining[:endIndex+2], true
}
