import { useTranslation } from "react-i18next"
import { TABULAR_SAMPLE, TYPOGRAPHY_SCALE } from "../catalog"
import { SectionShell, Showcase } from "./SectionShell"

export const TypographySection = () => {
  const { t } = useTranslation("designSystem")

  return (
    <SectionShell id="typography" title={t("typography.title")} description={t("typography.description")}>
      <Showcase title={t("typography.scale")}>
        <ul className="flex flex-col divide-y">
          {TYPOGRAPHY_SCALE.map((step) => (
            <li key={step.token} className="flex flex-col gap-1 py-2 sm:flex-row sm:items-baseline sm:gap-4">
              <span className="w-24 shrink-0 font-mono text-xs text-muted-foreground">{step.token}</span>
              <span className={step.className}>{t("typography.sample")}</span>
            </li>
          ))}
        </ul>
      </Showcase>
      <Showcase title={t("typography.tabular")}>
        <div className="grid gap-3 sm:grid-cols-2">
          <p className="text-lg">{TABULAR_SAMPLE}</p>
          <p className="text-lg tabular-nums">{TABULAR_SAMPLE}</p>
        </div>
        <p className="text-xs text-muted-foreground">{t("typography.tabularHint")}</p>
      </Showcase>
    </SectionShell>
  )
}
