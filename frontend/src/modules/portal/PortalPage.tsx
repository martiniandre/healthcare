import type { ReactNode } from "react"
import { useSearchParams } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { AlertCircle } from "lucide-react"
import { usePortalDashboardQuery } from "./queries"
import { PortalDashboardOverview } from "./PortalDashboardOverview"
import { PortalEncounters } from "./PortalEncounters"
import { PortalObservations } from "./PortalObservations"
import { PortalConditions } from "./PortalConditions"
import { PortalMedications } from "./PortalMedications"
import { PortalReports } from "./PortalReports"
import { PortalImaging } from "./PortalImaging"
import { PortalAppointments } from "./PortalAppointments"
import { Alert, AlertDescription, AlertTitle } from "../../shared/components/ui/alert"
import { Button } from "../../shared/components/ui/button"
import { Skeleton } from "../../shared/components/ui/skeleton"
import { PageContainer, PageHeader } from "../../shared/components/ui/PageContainer"

const PortalTab = {
  Dashboard: "dashboard",
  Appointments: "appointments",
  Encounters: "encounters",
  Observations: "observations",
  Conditions: "conditions",
  Medications: "medications",
  Reports: "reports",
  Imaging: "imaging",
} as const

type PortalTab = (typeof PortalTab)[keyof typeof PortalTab]

const portalTabValues = Object.values(PortalTab) as string[]

const resolvePortalTab = (requestedTab: string | null): PortalTab =>
  requestedTab && portalTabValues.includes(requestedTab) ? (requestedTab as PortalTab) : PortalTab.Dashboard

export const PortalPage = () => {
  const { t } = useTranslation("portal")
  const [searchParameters] = useSearchParams()
  const activeTab = resolvePortalTab(searchParameters.get("tab"))

  const { data: dashboardData, isLoading: isDashboardLoading, isError: hasDashboardError, refetch: refetchDashboard } =
    usePortalDashboardQuery()

  if (hasDashboardError && !dashboardData) {
    return (
      <PageContainer>
        <Alert variant="destructive">
          <AlertCircle aria-hidden="true" />
          <AlertTitle>{t("error.title")}</AlertTitle>
          <AlertDescription>
            <p>{t("error.description")}</p>
            <Button type="button" variant="outline" className="mt-2 min-h-11 min-w-11" onClick={() => refetchDashboard()}>
              {t("error.retry")}
            </Button>
          </AlertDescription>
        </Alert>
      </PageContainer>
    )
  }

  if (isDashboardLoading || !dashboardData) {
    return (
      <PageContainer>
        <div className="flex flex-col gap-4" role="status" aria-busy="true" aria-label={t("loading")}>
          <div className="flex flex-col gap-2">
            <Skeleton className="h-7 w-48" />
            <Skeleton className="h-4 w-64" />
          </div>
          <Skeleton className="h-32 w-full" />
          <Skeleton className="h-32 w-full" />
        </div>
      </PageContainer>
    )
  }

  const portalSections: Record<PortalTab, ReactNode> = {
    [PortalTab.Dashboard]: <PortalDashboardOverview dashboard={dashboardData} />,
    [PortalTab.Appointments]: <PortalAppointments />,
    [PortalTab.Encounters]: <PortalEncounters />,
    [PortalTab.Observations]: <PortalObservations />,
    [PortalTab.Conditions]: <PortalConditions />,
    [PortalTab.Medications]: <PortalMedications />,
    [PortalTab.Reports]: <PortalReports />,
    [PortalTab.Imaging]: <PortalImaging />,
  }

  return (
    <PageContainer>
      <PageHeader
        title={
          activeTab === PortalTab.Dashboard
            ? t("greeting", { firstName: dashboardData.patient_info.full_name.split(" ")[0] })
            : t(`tabs.${activeTab}`)
        }
      />
      <div className="flex min-w-0 flex-col gap-6">{portalSections[activeTab]}</div>
    </PageContainer>
  )
}
