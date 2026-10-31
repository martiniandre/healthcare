import { useTranslation } from "react-i18next"
import { CalendarCheck } from "lucide-react"
import { Alert, AlertDescription, AlertTitle } from "@/shared/components/ui/alert"

export const RequiredIndicator = () => {
  const { t } = useTranslation("schedule")

  return (
    <>
      <span aria-hidden="true" className="text-destructive">*</span>
      <span className="sr-only">{t("form.required")}</span>
    </>
  )
}

export const OptionalHint = () => {
  const { t } = useTranslation("schedule")

  return <span className="font-normal text-muted-foreground">({t("form.optional")})</span>
}

interface SelectionSummaryProperties {
  id: string
  text: string
  detail?: string
}

export const SelectionSummary = ({ id, text, detail }: SelectionSummaryProperties) => {
  const { t } = useTranslation("schedule")

  return (
    <div
      id={id}
      role="status"
      aria-live="polite"
      aria-atomic="true"
      data-testid="selection-summary"
      className="flex items-start gap-2 rounded-md border bg-muted/50 px-3 py-2 text-sm"
    >
      <CalendarCheck aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-primary" />
      <div className="flex flex-col">
        <span className="text-xs font-medium text-muted-foreground">{t("form.summary.label")}</span>
        <span className="font-medium text-foreground">{text}</span>
        {detail && <span className="text-xs text-muted-foreground">{detail}</span>}
      </div>
    </div>
  )
}

export interface FormErrorSummaryItem {
  name: string
  label: string
  message: string
}

interface FormErrorSummaryProperties {
  items: FormErrorSummaryItem[]
  onSelectItem: (name: string) => void
}

export const FormErrorSummary = ({ items, onSelectItem }: FormErrorSummaryProperties) => {
  const { t } = useTranslation("schedule")

  if (items.length === 0) {
    return null
  }

  return (
    <Alert variant="destructive" data-testid="form-error-summary">
      <AlertTitle>{t("form.errorSummary.title", { count: items.length })}</AlertTitle>
      <AlertDescription>
        <ul className="flex list-disc flex-col gap-1 pl-4">
          {items.map((item) => (
            <li key={item.name}>
              <button
                type="button"
                className="min-h-8 text-left underline underline-offset-2 focus-visible:outline-2 focus-visible:outline-offset-2 md:min-h-0"
                onClick={() => onSelectItem(item.name)}
              >
                {item.label}: {item.message}
              </button>
            </li>
          ))}
        </ul>
      </AlertDescription>
    </Alert>
  )
}
