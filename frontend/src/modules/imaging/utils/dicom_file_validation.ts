export const DICOM_MAX_FILE_GIGABYTES = 2
export const DICOM_MAX_FILE_BYTES = DICOM_MAX_FILE_GIGABYTES * 1024 * 1024 * 1024
export const DICOM_MAX_ARCHIVE_GIGABYTES = 2
export const DICOM_MAX_ARCHIVE_BYTES = DICOM_MAX_ARCHIVE_GIGABYTES * 1024 * 1024 * 1024
export const DICOM_MAX_ARCHIVE_ENTRIES = 2000
export const DICOM_ACCEPT_ATTRIBUTE = ".dcm,.zip"

const ZIP_LOCAL_HEADER_SIGNATURE = [0x50, 0x4b, 0x03, 0x04]

export type DicomFileValidationError = "empty" | "type" | "tooLarge" | "zipTooLarge" | "invalidZip"

const hasZipLocalHeader = async (file: File): Promise<boolean> => {
  const headerBytes = new Uint8Array(await file.slice(0, ZIP_LOCAL_HEADER_SIGNATURE.length).arrayBuffer())
  return ZIP_LOCAL_HEADER_SIGNATURE.every((signatureByte, index) => headerBytes[index] === signatureByte)
}

export const validateDicomFile = async (file: File): Promise<DicomFileValidationError | null> => {
  if (file.size === 0) {
    return "empty"
  }
  const lowerCaseName = file.name.toLowerCase()
  if (lowerCaseName.endsWith(".zip")) {
    if (file.size > DICOM_MAX_ARCHIVE_BYTES) {
      return "zipTooLarge"
    }
    return (await hasZipLocalHeader(file)) ? null : "invalidZip"
  }
  if (!lowerCaseName.endsWith(".dcm")) {
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
