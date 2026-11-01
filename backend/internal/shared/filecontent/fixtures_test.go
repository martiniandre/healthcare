package filecontent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"
	"testing"
)

const (
	secretPatientName = "Maria Secreta Silva"
	secretDevice      = "ScannerXYZ-9000"
)

func buildTestImage() image.Image {
	canvas := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for row := 0; row < 16; row++ {
		for column := 0; column < 16; column++ {
			canvas.Set(column, row, color.RGBA{R: uint8(row * 16), G: uint8(column * 16), B: 128, A: 255})
		}
	}
	return canvas
}

func buildPNGChunk(chunkType string, data []byte) []byte {
	chunk := make([]byte, 0, len(data)+12)
	chunk = binary.BigEndian.AppendUint32(chunk, uint32(len(data)))
	chunk = append(chunk, chunkType...)
	chunk = append(chunk, data...)
	checksum := crc32.ChecksumIEEE(append([]byte(chunkType), data...))
	return binary.BigEndian.AppendUint32(chunk, checksum)
}

func buildPlainPNG(testingInstance *testing.T) []byte {
	testingInstance.Helper()
	var buffer bytes.Buffer
	if encodeErr := png.Encode(&buffer, buildTestImage()); encodeErr != nil {
		testingInstance.Fatalf("png encode: %v", encodeErr)
	}
	return buffer.Bytes()
}

func buildPNGWithMetadata(testingInstance *testing.T) []byte {
	testingInstance.Helper()
	plain := buildPlainPNG(testingInstance)
	endOffset := len(plain) - 12
	var enriched bytes.Buffer
	enriched.Write(plain[:endOffset])
	enriched.Write(buildPNGChunk("tEXt", []byte("Author\x00"+secretPatientName)))
	enriched.Write(buildPNGChunk("zTXt", []byte("Comment\x00\x00"+secretDevice)))
	enriched.Write(buildPNGChunk("iTXt", []byte("Description\x00\x00\x00\x00\x00"+secretPatientName)))
	enriched.Write(buildPNGChunk("eXIf", []byte("Exif\x00\x00"+secretDevice)))
	enriched.Write(buildPNGChunk("tIME", []byte{0x07, 0xe8, 1, 2, 3, 4, 5}))
	enriched.Write(plain[endOffset:])
	enriched.WriteString("TRAILING-" + secretPatientName)
	return enriched.Bytes()
}

func buildJPEGSegment(marker byte, payload []byte) []byte {
	segment := []byte{0xff, marker}
	segment = binary.BigEndian.AppendUint16(segment, uint16(len(payload)+2))
	return append(segment, payload...)
}

func buildPlainJPEG(testingInstance *testing.T) []byte {
	testingInstance.Helper()
	var buffer bytes.Buffer
	if encodeErr := jpeg.Encode(&buffer, buildTestImage(), &jpeg.Options{Quality: 90}); encodeErr != nil {
		testingInstance.Fatalf("jpeg encode: %v", encodeErr)
	}
	return buffer.Bytes()
}

func buildJPEGWithMetadata(testingInstance *testing.T) []byte {
	testingInstance.Helper()
	plain := buildPlainJPEG(testingInstance)
	var enriched bytes.Buffer
	enriched.Write(plain[:2])
	enriched.Write(buildJPEGSegment(0xe1, []byte("Exif\x00\x00"+secretPatientName)))
	enriched.Write(buildJPEGSegment(0xe1, []byte("http://ns.adobe.com/xap/1.0/\x00"+secretDevice)))
	enriched.Write(buildJPEGSegment(0xed, []byte("Photoshop 3.0\x008BIM"+secretPatientName)))
	enriched.Write(buildJPEGSegment(0xfe, []byte("comment "+secretPatientName)))
	enriched.Write(buildJPEGSegment(0xe2, []byte("ICC_PROFILE\x00\x01\x01profile-bytes")))
	enriched.Write(plain[2:])
	enriched.WriteString("APPENDED-" + secretPatientName)
	return enriched.Bytes()
}

func buildPDF(includeSecrets bool) []byte {
	infoEntry := "<< /Producer (plain) >>"
	catalogMetadata := ""
	metadataStream := "<< /Length 0 >>\nstream\n\nendstream"
	if includeSecrets {
		infoEntry = "<< /Author (" + secretPatientName + ") /Title (Exame de " + secretPatientName + ") /Creator (" + secretDevice + ") >>"
		catalogMetadata = " /Metadata 5 0 R"
		xmp := "<x:xmpmeta><dc:creator>" + secretPatientName + "</dc:creator></x:xmpmeta>"
		metadataStream = fmt.Sprintf("<< /Type /Metadata /Subtype /XML /Length %d >>\nstream\n%s\nendstream", len(xmp), xmp)
	}

	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R" + catalogMetadata + " >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 200 200] >>",
		infoEntry,
		metadataStream,
	}

	var document bytes.Buffer
	document.WriteString("%PDF-1.4\n")
	offsets := make([]int, len(objects))
	for index, objectBody := range objects {
		offsets[index] = document.Len()
		fmt.Fprintf(&document, "%d 0 obj\n%s\nendobj\n", index+1, objectBody)
	}
	crossReferenceOffset := document.Len()
	fmt.Fprintf(&document, "xref\n0 %d\n0000000000 65535 f \n", len(objects)+1)
	for _, offset := range offsets {
		fmt.Fprintf(&document, "%010d 00000 n \n", offset)
	}
	fmt.Fprintf(&document, "trailer\n<< /Size %d /Root 1 0 R /Info 4 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, crossReferenceOffset)
	return document.Bytes()
}

func containsAny(content []byte, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(string(content), needle) {
			return true
		}
	}
	return false
}
