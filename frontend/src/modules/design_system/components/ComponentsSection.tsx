import { useTranslation } from "react-i18next"
import { COMPONENT_GROUP_IDS } from "../catalog"
import { ActionsShowcase } from "./ActionsShowcase"
import { DataShowcase } from "./DataShowcase"
import { FeedbackShowcase } from "./FeedbackShowcase"
import { FormsShowcase } from "./FormsShowcase"
import { NavigationShowcase } from "./NavigationShowcase"
import { OverlaysShowcase } from "./OverlaysShowcase"
import { SectionShell } from "./SectionShell"

const GROUPS = {
  actions: ActionsShowcase,
  feedback: FeedbackShowcase,
  forms: FormsShowcase,
  navigation: NavigationShowcase,
  overlays: OverlaysShowcase,
  data: DataShowcase,
} as const

export const ComponentsSection = () => {
  const { t } = useTranslation("designSystem")

  return (
    <SectionShell id="components" title={t("components.title")} description={t("components.description")}>
      {COMPONENT_GROUP_IDS.map((groupId) => {
        const Group = GROUPS[groupId]
        return (
          <div key={groupId} className="flex flex-col gap-3">
            <h3 className="text-base font-semibold text-foreground">{t(`components.groups.${groupId}`)}</h3>
            <div className="grid gap-3 lg:grid-cols-2">
              <Group />
            </div>
          </div>
        )
      })}
    </SectionShell>
  )
}
