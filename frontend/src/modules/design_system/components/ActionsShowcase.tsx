import { Plus } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Badge } from "@/shared/components/ui/badge"
import { Button } from "@/shared/components/ui/button"
import { Spinner } from "@/shared/components/ui/spinner"
import { StatusBadge } from "@/shared/components/ui/status-badge"
import { BADGE_VARIANTS, BUTTON_SIZES, BUTTON_VARIANTS, CLINICAL_STATUSES } from "../catalog"
import { Showcase } from "./SectionShell"

export const ActionsShowcase = () => {
  const { t } = useTranslation("designSystem")

  return (
    <>
      <Showcase title={t("components.button.variants")}>
        <div className="flex flex-wrap gap-2">
          {BUTTON_VARIANTS.map((variant) => (
            <Button key={variant} variant={variant}>
              {variant}
            </Button>
          ))}
        </div>
      </Showcase>
      <Showcase title={t("components.button.sizes")}>
        <div className="flex flex-wrap items-center gap-2">
          {BUTTON_SIZES.map((size) =>
            size === "icon" ? (
              <Button key={size} size={size} aria-label={t("components.button.iconLabel")}>
                <Plus aria-hidden="true" />
              </Button>
            ) : (
              <Button key={size} size={size}>
                {size}
              </Button>
            ),
          )}
        </div>
      </Showcase>
      <Showcase title={t("components.button.states")}>
        <div className="flex flex-wrap gap-2">
          <Button disabled>{t("components.button.disabled")}</Button>
          <Button disabled aria-busy="true">
            <Spinner />
            {t("components.button.loading")}
          </Button>
        </div>
      </Showcase>
      <Showcase title={t("components.badge.title")}>
        <div className="flex flex-wrap gap-2">
          {BADGE_VARIANTS.map((variant) => (
            <Badge key={variant} variant={variant}>
              {variant}
            </Badge>
          ))}
        </div>
      </Showcase>
      <Showcase title={t("components.statusBadge.title")}>
        <div className="flex flex-wrap gap-2">
          {CLINICAL_STATUSES.map((status) => (
            <StatusBadge key={status} status={status}>
              {t(`components.statusBadge.${status}`)}
            </StatusBadge>
          ))}
        </div>
      </Showcase>
    </>
  )
}
