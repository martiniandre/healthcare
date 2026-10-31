import { useTranslation } from "react-i18next"
import { COLOR_TOKENS } from "../catalog"
import { ColorSwatch } from "./ColorSwatch"
import { SectionShell } from "./SectionShell"

export const ColorsSection = () => {
  const { t } = useTranslation("designSystem")

  return (
    <SectionShell id="colors" title={t("colors.title")} description={t("colors.description")}>
      <ul className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-3">
        {COLOR_TOKENS.map((token) => (
          <ColorSwatch key={token.name} token={token} />
        ))}
      </ul>
    </SectionShell>
  )
}
