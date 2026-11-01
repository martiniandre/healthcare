export interface ApiErrorDetails {
  code?: string
  message?: string
  reason?: string
}

const readText = (value: unknown): string | undefined => {
  return typeof value === "string" && value.trim() !== "" ? value : undefined
}

export const extractApiErrorDetails = (error: unknown): ApiErrorDetails => {
  if (typeof error !== "object" || error === null || !("response" in error)) {
    return {}
  }
  const responseBody = (error as { response?: { data?: unknown } }).response?.data
  if (typeof responseBody !== "object" || responseBody === null) {
    return {}
  }
  const errorBody = responseBody as Record<string, unknown>
  return {
    code: readText(errorBody.code),
    message: readText(errorBody.error),
    reason: readText(errorBody.reason),
  }
}
