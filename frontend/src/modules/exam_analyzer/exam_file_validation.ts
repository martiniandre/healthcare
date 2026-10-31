export const EXAM_MAX_FILE_MEGABYTES = 15
export const EXAM_MAX_FILE_BYTES = EXAM_MAX_FILE_MEGABYTES * 1024 * 1024
export const EXAM_ALLOWED_EXTENSIONS = [".png", ".jpg", ".jpeg", ".pdf"]
export const EXAM_ACCEPT_ATTRIBUTE = [
  "image/png",
  "image/jpeg",
  "application/pdf",
  ...EXAM_ALLOWED_EXTENSIONS,
].join(",")

export type ExamFileValidationError = "empty" | "type" | "tooLarge"

export const validateExamFile = (file: File): ExamFileValidationError | null => {
  if (file.size === 0) {
    return "empty"
  }
  const lowerCaseName = file.name.toLowerCase()
  const hasAllowedExtension = EXAM_ALLOWED_EXTENSIONS.some((extension) => lowerCaseName.endsWith(extension))
  if (!hasAllowedExtension) {
    return "type"
  }
  if (file.size > EXAM_MAX_FILE_BYTES) {
    return "tooLarge"
  }
  return null
}

export const formatFileSize = (sizeInBytes: number): string => {
  if (sizeInBytes < 1024) {
    return `${sizeInBytes} B`
  }
  if (sizeInBytes < 1024 * 1024) {
    return `${(sizeInBytes / 1024).toFixed(1)} KB`
  }
  return `${(sizeInBytes / (1024 * 1024)).toFixed(2)} MB`
}
