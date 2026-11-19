import { useTranslation } from "react-i18next"

export const RequiredMark = () => {
  const { t } = useTranslation("patients")

  return (
    <>
      <span aria-hidden="true" className="text-destructive">
        *
      </span>
      <span className="sr-only"> {t("form.required")}</span>
    </>
  )
}
