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
