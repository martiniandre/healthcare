import { create } from "zustand"

export type ToastType = "success" | "error" | "info"

export interface ToastItem {
  id: string
  message: string
  type: ToastType
  duration?: number
}

interface ToastState {
  toasts: ToastItem[]
  addToast: (message: string, type: ToastType, duration?: number) => void
  removeToast: (id: string) => void
}

export const useToastStore = create<ToastState>((set) => ({
  toasts: [],
  addToast: (message, type, duration = 4000) => {
    const generatedId = Date.now().toString() + Math.random().toString(36).substring(2, 9)
    set((currentState) => ({
      toasts: [...currentState.toasts, { id: generatedId, message, type, duration }],
    }))
    setTimeout(() => {
      set((currentState) => ({
        toasts: currentState.toasts.filter((toastItem) => toastItem.id !== generatedId),
      }))
    }, duration)
  },
  removeToast: (id) =>
    set((currentState) => ({
      toasts: currentState.toasts.filter((toastItem) => toastItem.id !== id),
    })),
}))

export const toast = {
  success: (message: string, duration?: number) =>
    useToastStore.getState().addToast(message, "success", duration),
  error: (message: string, duration?: number) =>
    useToastStore.getState().addToast(message, "error", duration),
  info: (message: string, duration?: number) =>
    useToastStore.getState().addToast(message, "info", duration),
}
