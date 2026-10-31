import type { ReactNode } from "react"
import { useTranslation } from "react-i18next"
import { Field, FieldDescription, FieldError, FieldLabel } from "../../../../shared/components/ui/field"
import { cn } from "../../../../shared/utils/cn"
import { RequiredMark } from "./RequiredMark"

export interface FormFieldControlProperties {
  id: string
  "aria-invalid": boolean
  "aria-required": boolean | undefined
  "aria-describedby": string | undefined
}

interface FormFieldProperties {
  id: string
  label: string
  required?: boolean
  showOptionalHint?: boolean
  description?: string
  error?: string
  counter?: ReactNode
  className?: string
  children: (controlProperties: FormFieldControlProperties) => ReactNode
}

export const FormField = ({
  id,
  label,
  required = true,
  showOptionalHint = false,
  description,
  error,
  counter,
  className,
  children,
}: FormFieldProperties) => {
  const { t } = useTranslation("patients")
  const descriptionId = `${id}-description`
  const errorId = `${id}-error`
  const describedBy = [description ? descriptionId : null, error ? errorId : null].filter(Boolean).join(" ")

  return (
    <Field className={cn("gap-1.5", className)} data-invalid={Boolean(error)}>
      <div className="flex flex-wrap items-baseline gap-x-1">
        <FieldLabel htmlFor={id} className="text-sm font-medium text-foreground">
          {label}
        </FieldLabel>
        {required ? <RequiredMark /> : null}
        {!required && showOptionalHint ? (
          <span className="text-xs font-normal text-muted-foreground">{t("form.optional")}</span>
        ) : null}
      </div>
      {children({
        id,
        "aria-invalid": Boolean(error),
        "aria-required": required ? true : undefined,
        "aria-describedby": describedBy || undefined,
      })}
      {description || counter ? (
        <div className="flex items-start justify-between gap-3">
          {description ? (
            <FieldDescription id={descriptionId} className="text-xs">
              {description}
            </FieldDescription>
          ) : (
            <span />
          )}
          {counter}
        </div>
      ) : null}
      {error ? <FieldError id={errorId}>{error}</FieldError> : null}
    </Field>
  )
}
