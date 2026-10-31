import type { FormEventHandler, ReactNode } from "react"
import { Button } from "../../../../shared/components/ui/button"
import {
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "../../../../shared/components/ui/dialog"
import { Spinner } from "../../../../shared/components/ui/spinner"
import { cn } from "../../../../shared/utils/cn"

interface FormDialogLayoutProperties {
  title: string
  description?: string
  widthClassName?: string
  onSubmit: FormEventHandler<HTMLFormElement>
  cancelLabel: string
  onCancel: () => void
  submitLabel: string
  isPending: boolean
  isSubmitDisabled?: boolean
  children: ReactNode
}

export const FormDialogLayout = ({
  title,
  description,
  widthClassName = "sm:max-w-[520px]",
  onSubmit,
  cancelLabel,
  onCancel,
  submitLabel,
  isPending,
  isSubmitDisabled = false,
  children,
}: FormDialogLayoutProperties) => {
  return (
    <DialogContent
      className={cn("flex max-h-[90dvh] flex-col gap-0 p-0 text-left", widthClassName)}
      {...(description ? {} : { "aria-describedby": undefined })}
    >
      <DialogHeader className="border-b px-6 py-4 pr-12">
        <DialogTitle className="text-left">{title}</DialogTitle>
        {description ? <DialogDescription>{description}</DialogDescription> : null}
      </DialogHeader>
      <form onSubmit={onSubmit} className="flex min-h-0 flex-1 flex-col" noValidate>
        <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto px-6 py-4">{children}</div>
        <DialogFooter className="border-t px-6 py-4">
          <Button variant="outline" type="button" className="min-h-11 md:min-h-9" onClick={onCancel}>
            {cancelLabel}
          </Button>
          <Button
            type="submit"
            className="min-h-11 md:min-h-9"
            disabled={isPending || isSubmitDisabled}
            aria-busy={isPending}
          >
            {isPending && <Spinner />}
            {submitLabel}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  )
}
