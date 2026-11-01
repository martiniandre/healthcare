import { ImagingStudyStatus } from "../../../shared/types"
import type { ClinicalStatus } from "@/shared/components/ui/status-badge"

export const createValidDicomBlob = (): Blob => {
  const dicomBufferArray = new ArrayBuffer(136)
  const dataViewWriter = new DataView(dicomBufferArray)
  dataViewWriter.setUint8(128, 68)
  dataViewWriter.setUint8(129, 73)
  dataViewWriter.setUint8(130, 67)
  dataViewWriter.setUint8(131, 77)

  return new Blob([dicomBufferArray], { type: "application/dicom" })
}

export const waitForUploadFrame = async (millisecondsToWait: number): Promise<void> => {
  await new Promise((resolve) => setTimeout(resolve, millisecondsToWait))
}

export const resolveStudyClinicalStatus = (status: ImagingStudyStatus): ClinicalStatus => {
  if (status === ImagingStudyStatus.Available || status === ImagingStudyStatus.Ready) {
    return "normal"
  }
  if (
    status === ImagingStudyStatus.Registered ||
    status === ImagingStudyStatus.Receiving ||
    status === ImagingStudyStatus.Processing
  ) {
    return "info"
  }
  if (status === ImagingStudyStatus.Partial) {
    return "attention"
  }
  if (
    status === ImagingStudyStatus.Cancelled ||
    status === ImagingStudyStatus.EnteredInError ||
    status === ImagingStudyStatus.Failed
  ) {
    return "critical"
  }
  return "neutral"
}

const DICOM_DATE_PATTERN = /^(\d{4})(\d{2})(\d{2})$/
const SHORT_UID_THRESHOLD = 24
const SHORT_UID_EDGE_LENGTH = 10

export const normalizeDicomStudyDate = (studyDate: string): string => {
  const dicomDateMatch = DICOM_DATE_PATTERN.exec(studyDate)
  return dicomDateMatch ? `${dicomDateMatch[1]}-${dicomDateMatch[2]}-${dicomDateMatch[3]}` : studyDate
}

export const shortenDicomUid = (uid: string): string => {
  if (uid.length <= SHORT_UID_THRESHOLD) {
    return uid
  }
  return `${uid.slice(0, SHORT_UID_EDGE_LENGTH)}…${uid.slice(-SHORT_UID_EDGE_LENGTH)}`
}
