export const DICOM_MAX_FILE_GIGABYTES = 2
export const DICOM_MAX_FILE_BYTES = DICOM_MAX_FILE_GIGABYTES * 1024 * 1024 * 1024
export const DICOM_ACCEPT_ATTRIBUTE = ".dcm"

export type DicomFileValidationError = "empty" | "type" | "tooLarge"

export const validateDicomFile = (file: File): DicomFileValidationError | null => {
  if (file.size === 0) {
    return "empty"
  }
  if (!file.name.toLowerCase().endsWith(".dcm")) {
    return "type"
  }
  if (file.size > DICOM_MAX_FILE_BYTES) {
    return "tooLarge"
  }
  return null
}

export const formatDicomFileSize = (sizeInBytes: number): string => {
  if (sizeInBytes < 1024) {
    return `${sizeInBytes} B`
  }
  if (sizeInBytes < 1024 * 1024) {
    return `${(sizeInBytes / 1024).toFixed(1)} KB`
  }
  if (sizeInBytes < 1024 * 1024 * 1024) {
    return `${(sizeInBytes / (1024 * 1024)).toFixed(1)} MB`
  }
  return `${(sizeInBytes / (1024 * 1024 * 1024)).toFixed(2)} GB`
}
