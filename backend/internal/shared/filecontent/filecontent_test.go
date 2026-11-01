package filecontent

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func TestDetect(testingInstance *testing.T) {
	testCases := []struct {
		name     string
		head     []byte
		expected Kind
	}{
		{"png", buildPlainPNG(testingInstance), KindPNG},
		{"jpeg", buildPlainJPEG(testingInstance), KindJPEG},
		{"pdf", buildPDF(false), KindPDF},
		{"text", []byte("hello world"), KindUnknown},
		{"empty", nil, KindUnknown},
		{"short png signature", []byte{0x89, 'P', 'N'}, KindUnknown},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(testingInstance *testing.T) {
			head := testCase.head
			if len(head) > HeadBytes {
				head = head[:HeadBytes]
			}
			if detected := Detect(head); detected != testCase.expected {
				testingInstance.Fatalf("expected %q, got %q", testCase.expected, detected)
			}
		})
	}
}

func TestValidate(testingInstance *testing.T) {
	plainPNG := buildPlainPNG(testingInstance)
	plainJPEG := buildPlainJPEG(testingInstance)
	plainPDF := buildPDF(false)

	pngWithoutEnd := plainPNG[:len(plainPNG)-12]
	pngBadHeader := append([]byte{}, plainPNG...)
	copy(pngBadHeader[12:16], "XXXX")
	jpegTruncated := plainJPEG[:len(plainJPEG)/2]
	pdfWithoutEOF := bytes.ReplaceAll(plainPDF, []byte("%%EOF"), []byte("%%EOX"))
	pdfWithoutObjects := append([]byte("%PDF-1.4\n"), []byte(strings.Repeat("x", 40)+"\n%%EOF\n")...)
	pdfEOFOutsideTail := append(append([]byte{}, plainPDF...), bytes.Repeat([]byte(" "), 2000)...)

	testCases := []struct {
		name     string
		fileName string
		content  []byte
		expected error
	}{
		{"valid png", "scan.png", plainPNG, nil},
		{"valid jpeg", "scan.JPG", plainJPEG, nil},
		{"valid jpeg extension", "scan.jpeg", plainJPEG, nil},
		{"valid pdf", "report.pdf", plainPDF, nil},
		{"empty", "scan.png", nil, ErrEmptyFile},
		{"png that is pdf", "scan.png", plainPDF, ErrContentExtensionMismatch},
		{"pdf that is jpeg", "report.pdf", plainJPEG, ErrContentExtensionMismatch},
		{"jpg that is png", "scan.jpg", plainPNG, ErrContentExtensionMismatch},
		{"unknown extension", "scan.gif", plainPNG, ErrContentExtensionMismatch},
		{"text renamed png", "scan.png", []byte("not an image at all"), ErrCorruptedFile},
		{"png without IEND", "scan.png", pngWithoutEnd, ErrCorruptedFile},
		{"png bad IHDR", "scan.png", pngBadHeader, ErrCorruptedFile},
		{"png truncated", "scan.png", plainPNG[:20], ErrCorruptedFile},
		{"jpeg truncated", "scan.jpg", jpegTruncated, ErrCorruptedFile},
		{"jpeg only SOI", "scan.jpg", []byte{0xff, 0xd8, 0xff, 0xd9}, ErrCorruptedFile},
		{"pdf without EOF", "report.pdf", pdfWithoutEOF, ErrCorruptedFile},
		{"pdf without objects", "report.pdf", pdfWithoutObjects, ErrCorruptedFile},
		{"pdf EOF outside tail", "report.pdf", pdfEOFOutsideTail, ErrCorruptedFile},
		{"too large", "scan.png", append(append([]byte{}, plainPNG...), make([]byte, MaxFileBytes)...), ErrFileTooLarge},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(testingInstance *testing.T) {
			head := testCase.content
			if len(head) > HeadBytes {
				head = head[:HeadBytes]
			}
			validationErr := Validate(Detect(head), testCase.fileName, testCase.content)
			if testCase.expected == nil {
				if validationErr != nil {
					testingInstance.Fatalf("expected no error, got %v", validationErr)
				}
				return
			}
			if !errors.Is(validationErr, testCase.expected) {
				testingInstance.Fatalf("expected %v, got %v", testCase.expected, validationErr)
			}
		})
	}
}

func listPNGChunkTypes(testingInstance *testing.T, content []byte) []string {
	testingInstance.Helper()
	var chunkTypes []string
	offset := len(pngSignature)
	for offset+12 <= len(content) {
		chunkLength := int(binary.BigEndian.Uint32(content[offset:]))
		chunkTypes = append(chunkTypes, string(content[offset+4:offset+8]))
		offset += 12 + chunkLength
	}
	return chunkTypes
}

func TestStripMetadataPNG(testingInstance *testing.T) {
	original := buildPNGWithMetadata(testingInstance)

	cleaned, report, stripErr := StripMetadata(KindPNG, original)
	if stripErr != nil {
		testingInstance.Fatalf("strip: %v", stripErr)
	}

	if containsAny(cleaned, secretPatientName, secretDevice) {
		testingInstance.Fatal("secret strings survived in png output")
	}
	for _, chunkType := range listPNGChunkTypes(testingInstance, cleaned) {
		switch chunkType {
		case "tEXt", "zTXt", "iTXt", "eXIf", "tIME":
			testingInstance.Fatalf("ancillary chunk %s survived", chunkType)
		}
	}
	if validationErr := Validate(KindPNG, "scan.png", cleaned); validationErr != nil {
		testingInstance.Fatalf("cleaned png invalid: %v", validationErr)
	}

	decoded, _, decodeErr := image.Decode(bytes.NewReader(cleaned))
	if decodeErr != nil {
		testingInstance.Fatalf("cleaned png not decodable: %v", decodeErr)
	}
	reference, _, _ := image.Decode(bytes.NewReader(buildPlainPNG(testingInstance)))
	if decoded.At(5, 7) != reference.At(5, 7) || decoded.Bounds() != reference.Bounds() {
		testingInstance.Fatal("pixel data changed")
	}

	expectedPixelChunks := buildPlainPNG(testingInstance)
	if !bytes.Equal(cleaned, expectedPixelChunks) {
		testingInstance.Fatal("cleaned png should equal the metadata-free encoding byte for byte")
	}

	assertReport(testingInstance, report, []string{StepMetadataRemoved, StepTrailingDataRemoved}, []string{LimitationBurnedInText})
}

func TestStripMetadataPNGWithoutMetadata(testingInstance *testing.T) {
	plain := buildPlainPNG(testingInstance)
	cleaned, report, stripErr := StripMetadata(KindPNG, plain)
	if stripErr != nil {
		testingInstance.Fatalf("strip: %v", stripErr)
	}
	if !bytes.Equal(cleaned, plain) {
		testingInstance.Fatal("clean png must be unchanged")
	}
	assertReport(testingInstance, report, []string{StepNoMetadataFound}, []string{LimitationBurnedInText})
}

func listJPEGMarkers(testingInstance *testing.T, content []byte) []byte {
	testingInstance.Helper()
	var markers []byte
	offset := 2
	for offset+4 <= len(content) {
		if content[offset] != 0xff {
			testingInstance.Fatalf("expected marker at %d", offset)
		}
		marker := content[offset+1]
		markers = append(markers, marker)
		if marker == jpegMarkerSOS {
			break
		}
		offset += 2 + int(binary.BigEndian.Uint16(content[offset+2:]))
	}
	return markers
}

func TestStripMetadataJPEG(testingInstance *testing.T) {
	original := buildJPEGWithMetadata(testingInstance)
	plain := buildPlainJPEG(testingInstance)

	cleaned, report, stripErr := StripMetadata(KindJPEG, original)
	if stripErr != nil {
		testingInstance.Fatalf("strip: %v", stripErr)
	}

	if containsAny(cleaned, secretPatientName, secretDevice, "Exif", "Photoshop") {
		testingInstance.Fatal("secret strings survived in jpeg output")
	}
	if !bytes.Contains(cleaned, []byte("ICC_PROFILE")) {
		testingInstance.Fatal("ICC profile must be kept")
	}
	for _, marker := range listJPEGMarkers(testingInstance, cleaned) {
		if marker == 0xe1 || marker == 0xed || marker == 0xfe {
			testingInstance.Fatalf("marker %x survived", marker)
		}
	}
	if validationErr := Validate(KindJPEG, "scan.jpg", cleaned); validationErr != nil {
		testingInstance.Fatalf("cleaned jpeg invalid: %v", validationErr)
	}

	scanStart := bytes.Index(plain, []byte{0xff, jpegMarkerSOS})
	if !bytes.HasSuffix(cleaned, plain[scanStart:]) {
		testingInstance.Fatal("image scan data must be byte-identical (lossless)")
	}

	decoded, _, decodeErr := image.Decode(bytes.NewReader(cleaned))
	if decodeErr != nil {
		testingInstance.Fatalf("cleaned jpeg not decodable: %v", decodeErr)
	}
	reference, _, _ := image.Decode(bytes.NewReader(plain))
	if decoded.At(3, 9) != reference.At(3, 9) {
		testingInstance.Fatal("pixel data changed")
	}

	assertReport(testingInstance, report, []string{StepMetadataRemoved, StepTrailingDataRemoved}, []string{LimitationBurnedInText})
}

func TestStripMetadataJPEGWithoutMetadata(testingInstance *testing.T) {
	plain := buildPlainJPEG(testingInstance)
	cleaned, report, stripErr := StripMetadata(KindJPEG, plain)
	if stripErr != nil {
		testingInstance.Fatalf("strip: %v", stripErr)
	}
	if !bytes.Equal(cleaned, plain) {
		testingInstance.Fatal("clean jpeg must be unchanged")
	}
	assertReport(testingInstance, report, []string{StepNoMetadataFound}, []string{LimitationBurnedInText})
}

func TestStripMetadataPDF(testingInstance *testing.T) {
	original := buildPDF(true)

	cleaned, report, stripErr := StripMetadata(KindPDF, original)
	if stripErr != nil {
		testingInstance.Fatalf("strip: %v", stripErr)
	}

	if containsAny(cleaned, secretPatientName, secretDevice) {
		testingInstance.Fatal("secret strings survived in pdf output")
	}
	if validationErr := Validate(KindPDF, "report.pdf", cleaned); validationErr != nil {
		testingInstance.Fatalf("cleaned pdf invalid: %v", validationErr)
	}

	pdfContext, readErr := api.ReadContext(bytes.NewReader(cleaned), model.NewDefaultConfiguration())
	if readErr != nil {
		testingInstance.Fatalf("cleaned pdf cannot be parsed: %v", readErr)
	}
	rootDictionary, catalogErr := pdfContext.Catalog()
	if catalogErr != nil {
		testingInstance.Fatalf("catalog: %v", catalogErr)
	}
	if _, hasMetadata := rootDictionary.Find("Metadata"); hasMetadata {
		testingInstance.Fatal("catalog still references XMP metadata")
	}
	if pdfContext.Info != nil {
		infoDictionary, infoErr := pdfContext.DereferenceDict(*pdfContext.Info)
		if infoErr != nil {
			testingInstance.Fatalf("info dict: %v", infoErr)
		}
		for _, forbiddenKey := range []string{"Author", "Title", "Creator", "Subject", "Keywords"} {
			if _, present := infoDictionary.Find(forbiddenKey); present {
				testingInstance.Fatalf("info dictionary still has %s", forbiddenKey)
			}
		}
	}

	assertReport(testingInstance, report, []string{StepMetadataRemoved}, []string{LimitationBurnedInText})
}

func TestStripMetadataPDFFallbackReportsLimitation(testingInstance *testing.T) {
	brokenPDF := append([]byte("%PDF-1.4\n1 0 obj\n<< /Author (x) >>\nendobj\n"), []byte("%%EOF\n")...)

	returned, report, stripErr := StripMetadata(KindPDF, brokenPDF)
	if stripErr != nil {
		testingInstance.Fatalf("strip must not fail: %v", stripErr)
	}
	if !bytes.Equal(returned, brokenPDF) {
		testingInstance.Fatal("original bytes must be returned on fallback")
	}
	assertReport(testingInstance, report, []string{}, []string{LimitationBurnedInText, LimitationPDFMetadata})
}

func TestStripMetadataRejectsInvalidInput(testingInstance *testing.T) {
	testCases := []struct {
		name     string
		kind     Kind
		content  []byte
		expected error
	}{
		{"empty", KindPNG, nil, ErrEmptyFile},
		{"unknown kind", KindUnknown, []byte("abc"), ErrCorruptedFile},
		{"png garbage", KindPNG, []byte("garbage-garbage-garbage"), ErrCorruptedFile},
		{"png truncated chunk", KindPNG, append(append([]byte{}, pngSignature...), 0, 0, 0xff, 0xff, 'I', 'D', 'A', 'T', 1, 2), ErrCorruptedFile},
		{"jpeg garbage", KindJPEG, []byte("garbage-garbage-garbage"), ErrCorruptedFile},
		{"jpeg segment overflow", KindJPEG, []byte{0xff, 0xd8, 0xff, 0xe1, 0xff, 0xff, 0x00, 0x00}, ErrCorruptedFile},
		{"jpeg without scan", KindJPEG, []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x04, 'J', 'F'}, ErrCorruptedFile},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(testingInstance *testing.T) {
			_, _, stripErr := StripMetadata(testCase.kind, testCase.content)
			if !errors.Is(stripErr, testCase.expected) {
				testingInstance.Fatalf("expected %v, got %v", testCase.expected, stripErr)
			}
		})
	}
}

func TestNeutralFileName(testingInstance *testing.T) {
	testCases := map[Kind]string{
		KindPNG:     "exam.png",
		KindJPEG:    "exam.jpg",
		KindPDF:     "exam.pdf",
		KindUnknown: "exam",
	}
	for kind, expected := range testCases {
		if actual := NeutralFileName(kind); actual != expected {
			testingInstance.Fatalf("kind %q: expected %q, got %q", kind, expected, actual)
		}
	}
}

func assertReport(testingInstance *testing.T, report AnonymizationReport, expectedSteps []string, expectedLimitations []string) {
	testingInstance.Helper()
	if strings.Join(report.Steps, ",") != strings.Join(expectedSteps, ",") {
		testingInstance.Fatalf("steps: expected %v, got %v", expectedSteps, report.Steps)
	}
	if strings.Join(report.Limitations, ",") != strings.Join(expectedLimitations, ",") {
		testingInstance.Fatalf("limitations: expected %v, got %v", expectedLimitations, report.Limitations)
	}
	if report.Steps == nil || report.Limitations == nil {
		testingInstance.Fatal("report slices must be non-nil for stable JSON")
	}
}
