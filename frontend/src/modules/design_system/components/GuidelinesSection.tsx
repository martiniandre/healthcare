import { Check, X } from "lucide-react"
import { useTranslation } from "react-i18next"
import { GUIDELINE_IDS } from "../catalog"
import { SectionShell } from "./SectionShell"

export const GuidelinesSection = () => {
  const { t } = useTranslation("designSystem")

  return (
    <SectionShell id="guidelines" title={t("guidelines.title")} description={t("guidelines.description")}>
      <div className="grid gap-4 md:grid-cols-2">
        <div className="rounded-xl border border-success/30 bg-success/5 p-4">
          <h3 className="flex items-center gap-2 text-sm font-semibold text-success">
            <Check className="size-4" aria-hidden="true" />
            {t("guidelines.doTitle")}
          </h3>
          <ul className="mt-3 flex list-disc flex-col gap-2 pl-5 text-sm text-foreground">
            {GUIDELINE_IDS.map((guidelineId) => (
              <li key={guidelineId}>{t(`guidelines.do.${guidelineId}`)}</li>
            ))}
          </ul>
        </div>
        <div className="rounded-xl border border-destructive/30 bg-destructive/5 p-4">
          <h3 className="flex items-center gap-2 text-sm font-semibold text-destructive">
            <X className="size-4" aria-hidden="true" />
            {t("guidelines.dontTitle")}
          </h3>
          <ul className="mt-3 flex list-disc flex-col gap-2 pl-5 text-sm text-foreground">
            {GUIDELINE_IDS.map((guidelineId) => (
              <li key={guidelineId}>{t(`guidelines.dont.${guidelineId}`)}</li>
            ))}
          </ul>
        </div>
      </div>
    </SectionShell>
  )
}
