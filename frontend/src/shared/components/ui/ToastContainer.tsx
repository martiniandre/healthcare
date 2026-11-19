import { useEffect, useRef, useState } from "react"
import { useTranslation } from "react-i18next"
import { useToastStore, type ToastItem } from "../../store/toast_store"
import { CheckCircle2, AlertCircle, Info, X } from "lucide-react"
import { Button } from "@/shared/components/ui/button"
import { cn } from "../../utils/cn"

interface ToastCardProps {
  toastItem: ToastItem
  onDismiss: (id: string) => void
}

const ToastCard = ({ toastItem, onDismiss }: ToastCardProps) => {
  const { t } = useTranslation("common")
  const [isPaused, setIsPaused] = useState(false)
  const remainingTimeReference = useRef(toastItem.duration ?? null)
  const startedAtReference = useRef(0)

  useEffect(() => {
    if (remainingTimeReference.current === null || isPaused) return
    startedAtReference.current = Date.now()
    const dismissTimer = window.setTimeout(() => onDismiss(toastItem.id), remainingTimeReference.current)
    return () => {
      window.clearTimeout(dismissTimer)
      if (remainingTimeReference.current !== null) {
        remainingTimeReference.current = Math.max(0, remainingTimeReference.current - (Date.now() - startedAtReference.current))
      }
    }
  }, [isPaused, onDismiss, toastItem.id])

  const isError = toastItem.type === "error"

  return (
    <div
      role={isError ? "alert" : "status"}
      aria-live={isError ? "assertive" : "polite"}
      aria-atomic="true"
      onMouseEnter={() => setIsPaused(true)}
      onMouseLeave={() => setIsPaused(false)}
      onFocus={() => setIsPaused(true)}
      onBlur={() => setIsPaused(false)}
      className={cn(
        "pointer-events-auto flex w-full items-start gap-3 rounded-lg border border-border bg-card p-4 shadow-lg animate-toast-slide-in transition-shadow duration-150 hover:shadow-xl",
        isError && "border-destructive/40"
      )}
    >
      <div className="mt-0.5 flex-shrink-0" aria-hidden="true">
        {toastItem.type === "success" && <CheckCircle2 className="size-5 text-success" />}
        {isError && <AlertCircle className="size-5 text-danger" />}
        {toastItem.type === "info" && <Info className="size-5 text-primary" />}
      </div>

      <div className="min-w-0 flex-1">
        <span className="block text-sm font-semibold text-foreground">
          {toastItem.type === "success" && t("toast.success")}
          {isError && t("toast.error")}
          {toastItem.type === "info" && t("toast.info")}
        </span>
        <p className="mt-0.5 break-words text-xs leading-relaxed text-muted-foreground">{toastItem.message}</p>
      </div>

      <Button
        variant="ghost"
        size="icon-sm"
        onClick={() => onDismiss(toastItem.id)}
        aria-label={t("close")}
        className="shrink-0 text-muted-foreground"
      >
        <X className="size-4" />
      </Button>
    </div>
  )
}

export const ToastContainer = () => {
  const toasts = useToastStore((state) => state.toasts)
  const removeToast = useToastStore((state) => state.removeToast)

  return (
    <div className="pointer-events-none fixed bottom-4 right-4 z-[9999] flex w-full max-w-sm flex-col gap-2 px-4 sm:px-0">
      {toasts.map((toastItem) => (
        <ToastCard key={toastItem.id} toastItem={toastItem} onDismiss={removeToast} />
      ))}
    </div>
  )
}
