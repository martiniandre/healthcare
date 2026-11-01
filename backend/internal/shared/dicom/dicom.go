package dicom

import (
	"fmt"
	"strings"
)

const (
	TransferSyntaxImplicitVRLittleEndian = "1.2.840.10008.1.2"
	TransferSyntaxExplicitVRLittleEndian = "1.2.840.10008.1.2.1"
	TransferSyntaxDeflatedExplicitVRLE   = "1.2.840.10008.1.2.1.99"
	TransferSyntaxExplicitVRBigEndian    = "1.2.840.10008.1.2.2"
	TransferSyntaxJPEGBaseline           = "1.2.840.10008.1.2.4.50"
	TransferSyntaxJPEGExtended           = "1.2.840.10008.1.2.4.51"
	TransferSyntaxJPEGLossless           = "1.2.840.10008.1.2.4.57"
	TransferSyntaxJPEGLosslessP1         = "1.2.840.10008.1.2.4.70"
	TransferSyntaxJPEGLosslessSV1        = "1.2.840.10008.1.2.4.80"
	TransferSyntaxJPEGLSNearLossless     = "1.2.840.10008.1.2.4.81"
	TransferSyntaxJPEG2000Lossless       = "1.2.840.10008.1.2.4.90"
	TransferSyntaxJPEG2000               = "1.2.840.10008.1.2.4.91"
	TransferSyntaxRLELossless            = "1.2.840.10008.1.2.5"
)

const SOPClassSecondaryCapture = "1.2.840.10008.5.1.4.1.1.7"

var transferSyntaxDisplayNames = map[string]string{
	TransferSyntaxImplicitVRLittleEndian: "Implicit VR Little Endian",
	TransferSyntaxExplicitVRLittleEndian: "Explicit VR Little Endian",
	TransferSyntaxDeflatedExplicitVRLE:   "Deflated Explicit VR Little Endian",
	TransferSyntaxExplicitVRBigEndian:    "Explicit VR Big Endian",
	TransferSyntaxJPEGBaseline:           "JPEG Baseline (Process 1)",
	TransferSyntaxJPEGExtended:           "JPEG Extended (Process 2 and 4)",
	TransferSyntaxJPEGLossless:           "JPEG Lossless, Non-Hierarchical (Process 14)",
	TransferSyntaxJPEGLosslessP1:         "JPEG Lossless, Non-Hierarchical, First-Order Prediction",
	TransferSyntaxJPEGLosslessSV1:        "JPEG Lossless, Non-Hierarchical, SV1",
	TransferSyntaxJPEGLSNearLossless:     "JPEG-LS Near-Lossless",
	TransferSyntaxJPEG2000Lossless:       "JPEG 2000 Image Compression (Lossless Only)",
	TransferSyntaxJPEG2000:               "JPEG 2000 Image Compression",
	TransferSyntaxRLELossless:            "RLE Lossless",
}

var supportedTransferSyntaxes = map[string]bool{
	TransferSyntaxImplicitVRLittleEndian: true,
	TransferSyntaxExplicitVRLittleEndian: true,
}

func IsSupportedTransferSyntax(transferSyntaxUID string) bool {
	return supportedTransferSyntaxes[strings.TrimSpace(transferSyntaxUID)]
}

func DescribeTransferSyntax(transferSyntaxUID string) string {
	trimmedUID := strings.TrimSpace(transferSyntaxUID)
	if displayName, found := transferSyntaxDisplayNames[trimmedUID]; found {
		return fmt.Sprintf("%s (%s)", displayName, trimmedUID)
	}
	if trimmedUID == "" {
		return "unknown"
	}
	return trimmedUID
}
