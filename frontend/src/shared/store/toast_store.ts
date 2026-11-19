import { create } from "zustand"

export type ToastType = "success" | "error" | "info"

export const DEFAULT_TOAST_DURATION_MS = 6000

export interface ToastItem {
  id: string
  message: string
  type: ToastType
  duration?: number | null
}

interface ToastState {
  toasts: ToastItem[]
  addToast: (message: string, type: ToastType, duration?: number | null) => void
  removeToast: (id: string) => void
}

export const resolveToastDuration = (type: ToastType, duration?: number | null): number | null => {
  if (duration !== undefined) return duration
  return type === "error" ? null : DEFAULT_TOAST_DURATION_MS
}

export const useToastStore = create<ToastState>((set) => ({
  toasts: [],
  addToast: (message, type, duration) => {
    const generatedId = Date.now().toString() + Math.random().toString(36).substring(2, 9)
    set((currentState) => ({
      toasts: [...currentState.toasts, { id: generatedId, message, type, duration: resolveToastDuration(type, duration) }],
    }))
  },
  removeToast: (id) =>
    set((currentState) => ({
      toasts: currentState.toasts.filter((toastItem) => toastItem.id !== id),
    })),
}))

export const toast = {
  success: (message: string, duration?: number | null) =>
    useToastStore.getState().addToast(message, "success", duration),
  error: (message: string, duration?: number | null) =>
    useToastStore.getState().addToast(message, "error", duration),
  info: (message: string, duration?: number | null) =>
    useToastStore.getState().addToast(message, "info", duration),
}
