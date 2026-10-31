import { useCallback, useState } from "react"

const STORAGE_KEY = "telemetry.alarmsMuted"

const readStoredMuted = (): boolean => {
  try {
    const storedValue = window.localStorage.getItem(STORAGE_KEY)
    return storedValue === null ? true : storedValue === "true"
  } catch {
    return true
  }
}

export const useAlarmMute = () => {
  const [isMuted, setIsMutedState] = useState<boolean>(readStoredMuted)

  const setIsMuted = useCallback((nextMuted: boolean) => {
    setIsMutedState(nextMuted)
    try {
      window.localStorage.setItem(STORAGE_KEY, String(nextMuted))
    } catch {
      return
    }
  }, [])

  return { isMuted, setIsMuted }
}
