package filecontent

import (
	"bytes"
	"path/filepath"
	"strings"
)

type Kind string

const (
	KindUnknown Kind = ""
	KindPNG     Kind = "png"
	KindJPEG    Kind = "jpeg"
	KindPDF     Kind = "pdf"
)

const MaxFileBytes = 15 << 20

const HeadBytes = 16

var (
	pngSignature  = []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	jpegSignature = []byte{0xff, 0xd8, 0xff}
	pdfSignature  = []byte("%PDF-")
)

func Detect(headBytes []byte) Kind {
	switch {
	case bytes.HasPrefix(headBytes, pngSignature):
		return KindPNG
	case bytes.HasPrefix(headBytes, jpegSignature):
		return KindJPEG
	case bytes.HasPrefix(headBytes, pdfSignature):
		return KindPDF
	default:
		return KindUnknown
	}
}

func KindFromFileName(fileName string) Kind {
	switch strings.ToLower(filepath.Ext(fileName)) {
	case ".png":
		return KindPNG
	case ".jpg", ".jpeg":
		return KindJPEG
	case ".pdf":
		return KindPDF
	default:
		return KindUnknown
	}
}

func (kind Kind) Extension() string {
	switch kind {
	case KindPNG:
		return ".png"
	case KindJPEG:
		return ".jpg"
	case KindPDF:
		return ".pdf"
	default:
		return ""
	}
}

func NeutralFileName(kind Kind) string {
	return "exam" + kind.Extension()
}
