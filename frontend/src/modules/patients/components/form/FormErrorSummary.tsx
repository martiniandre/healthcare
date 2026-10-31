import { CircleAlert } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Alert, AlertDescription, AlertTitle } from "../../../../shared/components/ui/alert"

export interface FormErrorSummaryItem {
  fieldId: string
  label: string
  message: string
}

interface FormErrorSummaryProperties {
  items: FormErrorSummaryItem[]
}

export const FormErrorSummary = ({ items }: FormErrorSummaryProperties) => {
  const { t } = useTranslation("patients")

  const focusField = (fieldId: string) => {
    document.getElementById(fieldId)?.focus()
  }

  if (items.length === 0) {
    return null
  }

  return (
    <Alert variant="destructive" data-testid="form-error-summary">
      <CircleAlert aria-hidden="true" />
      <AlertTitle className="line-clamp-none">{t("form.errorSummaryTitle")}</AlertTitle>
      <AlertDescription>
        <ul className="flex flex-col gap-1">
          {items.map((item) => (
            <li key={item.fieldId}>
              <button
                type="button"
                className="min-h-8 text-left underline underline-offset-2 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
                onClick={() => focusField(item.fieldId)}
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
