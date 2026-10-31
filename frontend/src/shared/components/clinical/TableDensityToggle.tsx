import { Rows3, Rows4 } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Button } from "@/shared/components/ui/button"
import { useLayoutStore } from "@/shared/store/layout_store"

export const TableDensityToggle = () => {
  const { t } = useTranslation("common")
  const tableDensity = useLayoutStore((state) => state.tableDensity)
  const toggleTableDensity = useLayoutStore((state) => state.toggleTableDensity)
  const isCompact = tableDensity === "compact"
  const densityLabel = isCompact ? t("tableDensity.comfortable") : t("tableDensity.compact")

  return (
    <Button
      type="button"
      variant="ghost"
      size="icon-sm"
      onClick={toggleTableDensity}
      aria-label={densityLabel}
      aria-pressed={isCompact}
      title={densityLabel}
    >
      {isCompact ? <Rows3 aria-hidden="true" /> : <Rows4 aria-hidden="true" />}
    </Button>
  )
}
