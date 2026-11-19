import type { ReactNode } from "react"
import { Button } from "@/shared/components/ui/button"
import { Spinner } from "@/shared/components/ui/spinner"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/components/ui/dialog"

interface ConfirmActionDialogProps {
  isOpen: boolean
  title: string
  description: string
  details?: ReactNode
  hint?: string
  confirmLabel: string
  cancelLabel: string
  isDestructive?: boolean
  isPending?: boolean
  onConfirm: () => void
  onClose: () => void
}

export const ConfirmActionDialog = ({
  isOpen,
  title,
  description,
  details,
  hint,
  confirmLabel,
  cancelLabel,
  isDestructive = false,
  isPending = false,
  onConfirm,
  onClose,
}: ConfirmActionDialogProps) => (
  <Dialog open={isOpen} onOpenChange={(open) => !open && !isPending && onClose()}>
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle className="text-left">{title}</DialogTitle>
        <DialogDescription className="text-left">{description}</DialogDescription>
      </DialogHeader>
      {details && (
        <div className="flex flex-col gap-1 rounded-md border bg-muted/50 px-3 py-2 text-left text-sm">{details}</div>
      )}
      {hint && <p className="text-left text-xs text-muted-foreground">{hint}</p>}
      <DialogFooter>
        <Button variant="outline" type="button" autoFocus className="h-11 md:h-9" disabled={isPending} onClick={onClose}>
          {cancelLabel}
        </Button>
        <Button
          variant={isDestructive ? "destructive" : "default"}
          type="button"
          className="h-11 md:h-9"
          disabled={isPending}
          aria-busy={isPending}
          onClick={onConfirm}
        >
          {isPending && <Spinner />}
          {confirmLabel}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
)
