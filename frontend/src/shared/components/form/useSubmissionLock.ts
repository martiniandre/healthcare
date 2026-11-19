import { useCallback, useRef } from "react"

export const useSubmissionLock = () => {
  const inFlightReference = useRef(false)

  return useCallback(async (task: () => Promise<unknown>) => {
    if (inFlightReference.current) {
      return
    }
    inFlightReference.current = true
    try {
      await task()
    } finally {
      inFlightReference.current = false
    }
  }, [])
}
