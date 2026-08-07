import { useRef, useCallback } from "react"

const generateIdempotencyKey = (): string => {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return crypto.randomUUID()
  }
  const randomSuffix = Math.random().toString(36).substring(2, 15) + Date.now().toString(36)
  return `booking-${randomSuffix}`
}

export function useIdempotencyKey() {
  const idempotencyKeyRef = useRef<string | null>(null)

  const getOrCreateKey = useCallback((): string => {
    if (!idempotencyKeyRef.current) {
      idempotencyKeyRef.current = generateIdempotencyKey()
    }
    return idempotencyKeyRef.current
  }, [])

  const resetKey = useCallback(() => {
    idempotencyKeyRef.current = null
  }, [])

  return { getOrCreateKey, resetKey }
}
