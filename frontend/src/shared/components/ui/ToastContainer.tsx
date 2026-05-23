import { useToastStore, type ToastItem } from "../../store/toast_store"
import { CheckCircle2, AlertCircle, Info, X } from "lucide-react"
import { cn } from "../../utils/cn"

export const ToastContainer = () => {
  const toastStore = useToastStore()

  return (
    <div className="fixed top-4 right-4 z-[9999] flex flex-col gap-2 max-w-sm w-full pointer-events-none px-4 sm:px-0">
      {toastStore.toasts.map((toastItem: ToastItem) => {
        const handleDismiss = () => {
          toastStore.removeToast(toastItem.id)
        }

        return (
          <div
            key={toastItem.id}
            className={cn(
              "pointer-events-auto flex items-start gap-3 w-full bg-white border border-gray-150 rounded-lg p-4 shadow-lg animate-toast-slide-in hover:shadow-xl transition-all duration-200"
            )}
          >
            <div className="flex-shrink-0 mt-0.5">
              {toastItem.type === "success" && (
                <CheckCircle2 className="w-5 h-5 text-success" />
              )}
              {toastItem.type === "error" && (
                <AlertCircle className="w-5 h-5 text-danger" />
              )}
              {toastItem.type === "info" && (
                <Info className="w-5 h-5 text-primary" />
              )}
            </div>

            <div className="flex-1 min-w-0">
              <span className="text-[13px] font-bold text-gray-900 block">
                {toastItem.type === "success" && "Sucesso"}
                {toastItem.type === "error" && "Erro"}
                {toastItem.type === "info" && "Aviso"}
              </span>
              <p className="text-xs text-gray-600 mt-0.5 leading-relaxed break-words">
                {toastItem.message}
              </p>
            </div>

            <button
              onClick={handleDismiss}
              className="flex-shrink-0 p-0.5 rounded-md text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors cursor-pointer"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        )
      })}
    </div>
  )
}
