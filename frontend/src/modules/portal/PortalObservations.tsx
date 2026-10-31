import { Heart } from "lucide-react"
import { useTranslation } from "react-i18next"
import { usePortalObservationsQuery } from "./queries"
import { PortalSectionState } from "./components/PortalSectionState"
import { PortalVitalValue } from "./components/PortalVitalValue"
import { Card, CardContent } from "../../shared/components/ui/card"
import { formatDateTime } from "../../shared/utils/dates"
import { useLocale } from "../../shared/hooks/useLocale"
import { findVitalSignDisplay } from "../patients/components/vitalSignDisplay"

export const PortalObservations = () => {
  const { t } = useTranslation(["portal", "patients"])
  const { data: observations, isLoading, isError, refetch } = usePortalObservationsQuery()
  const locale = useLocale()

  return (
    <PortalSectionState
      isLoading={isLoading}
      isError={isError}
      onRetry={() => refetch()}
      isEmpty={!observations || observations.length === 0}
      emptyTitle={t("portal:empty.observations")}
      emptyIcon={Heart}
      skeletonLayout="grid"
    >
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {observations?.map((observation) => {
          const displayMetadata = findVitalSignDisplay(observation.loinc_code)
          return (
            <Card key={observation.fhir_resource_id} className="py-4 sm:py-5">
              <CardContent className="px-4 sm:px-5">
                <p className="mb-1 break-words text-xs font-medium text-muted-foreground">
                  {displayMetadata?.labelKey ? t(`patients:${displayMetadata.labelKey}`) : observation.code_display}
                </p>
                <PortalVitalValue
                  notPerformed={observation.not_performed}
                  valueQuantity={observation.value_quantity}
                  valueUnit={observation.value_unit}
                  valueClassName="text-2xl font-bold text-foreground"
                />
                <p className="mt-2 text-xs tabular-nums text-muted-foreground">
                  {formatDateTime(observation.observed_at, locale)}
                </p>
              </CardContent>
            </Card>
          )
        })}
      </div>
    </PortalSectionState>
  )
}
