import { useEffect, useRef } from "react"
import { useSearchParams } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { usePortalDashboardQuery } from "./queries"
import { PortalDashboardOverview } from "./PortalDashboardOverview"
import { PortalEncounters } from "./PortalEncounters"
import { PortalObservations } from "./PortalObservations"
import { PortalConditions } from "./PortalConditions"
import { PortalMedications } from "./PortalMedications"
import { PortalReports } from "./PortalReports"
import { PortalImaging } from "./PortalImaging"
import { PortalAppointments } from "./PortalAppointments"
import { Avatar, AvatarFallback } from "../../shared/components/ui/avatar"
import { Card, CardContent } from "../../shared/components/ui/card"
import { Alert, AlertDescription, AlertTitle } from "../../shared/components/ui/alert"
import { Button } from "../../shared/components/ui/button"
import { Skeleton } from "../../shared/components/ui/skeleton"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../../shared/components/ui/tabs"
import {
  LayoutDashboard,
  History,
  Heart,
  Activity,
  Pill,
  FileText,
  Image,
  CalendarClock,
  AlertCircle,
} from "lucide-react"

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

const sidebarItems: { key: PortalTab; labelKey: string; icon: React.ReactNode }[] = [
  { key: "dashboard", labelKey: "tabs.dashboard", icon: <LayoutDashboard aria-hidden="true" /> },
  { key: "appointments", labelKey: "tabs.appointments", icon: <CalendarClock aria-hidden="true" /> },
  { key: "encounters", labelKey: "tabs.encounters", icon: <History aria-hidden="true" /> },
  { key: "observations", labelKey: "tabs.observations", icon: <Heart aria-hidden="true" /> },
  { key: "conditions", labelKey: "tabs.conditions", icon: <Activity aria-hidden="true" /> },
  { key: "medications", labelKey: "tabs.medications", icon: <Pill aria-hidden="true" /> },
  { key: "reports", labelKey: "tabs.reports", icon: <FileText aria-hidden="true" /> },
  { key: "imaging", labelKey: "tabs.imaging", icon: <Image aria-hidden="true" /> },
]

const tabTriggerClassName =
  "h-auto min-h-11 min-w-11 flex-none gap-2 px-3 py-2 text-xs font-bold data-[state=active]:bg-primary/10 data-[state=active]:text-primary md:justify-start md:gap-3 md:px-4 md:py-3 md:font-extrabold"

const getInitials = (fullName: string): string =>
  fullName
    .split(" ")
    .filter(Boolean)
    .map((namePart) => namePart[0])
    .join("")
    .slice(0, 2)
    .toUpperCase()

export const PortalPage = () => {
  const { t } = useTranslation("portal")
  const [searchParameters, setSearchParameters] = useSearchParams()
  const activeTab = (searchParameters.get("tab") || PortalTab.Dashboard) as PortalTab
  const mobileTabListReference = useRef<HTMLDivElement>(null)
  const setActiveTab = (tabName: PortalTab) => {
    setSearchParameters({ tab: tabName })
  }

  const { data: dashboardData, isLoading: isDashboardLoading, isError: hasDashboardError, refetch: refetchDashboard } =
    usePortalDashboardQuery()

  useEffect(() => {
    const activeTrigger = mobileTabListReference.current?.querySelector<HTMLElement>('[data-state="active"]')
    activeTrigger?.scrollIntoView?.({ block: "nearest", inline: "center" })
  }, [activeTab, dashboardData])

  if (hasDashboardError && !dashboardData) {
    return (
      <div className="mx-auto flex w-full max-w-7xl flex-1 flex-col p-4 md:p-8">
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
      </div>
    )
  }

  if (isDashboardLoading || !dashboardData) {
    return (
      <div
        className="mx-auto flex w-full max-w-7xl flex-1 flex-col gap-4 p-4 md:p-8"
        role="status"
        aria-busy="true"
        aria-label={t("loading")}
      >
        <div className="flex items-center gap-3">
          <Skeleton className="size-10 rounded-full md:size-14" />
          <div className="flex flex-1 flex-col gap-2">
            <Skeleton className="h-5 w-1/2" />
            <Skeleton className="h-3 w-1/3" />
          </div>
        </div>
        <Skeleton className="h-11 w-full" />
        <Skeleton className="h-32 w-full" />
        <Skeleton className="h-32 w-full" />
      </div>
    )
  }

  const patientName = dashboardData.patient_info.full_name

  return (
    <div className="mx-auto flex w-full min-w-0 max-w-7xl flex-1 flex-col gap-3 px-4 pt-3 pb-[max(1rem,env(safe-area-inset-bottom))] md:gap-6 md:p-8">
      <div className="flex items-center gap-3 md:gap-4 md:rounded-xl md:border md:bg-card md:p-6">
        <Avatar className="size-10 md:size-14">
          <AvatarFallback className="bg-primary/10 text-base font-bold text-primary md:text-xl">
            {getInitials(patientName)}
          </AvatarFallback>
        </Avatar>
        <div className="min-w-0 flex-1">
          <h1 className="truncate text-lg font-display font-bold text-foreground md:text-xl" title={patientName}>
            {patientName}
          </h1>
          <p className="truncate text-xs text-muted-foreground md:text-sm">{t("title")}</p>
        </div>
      </div>

      <Tabs value={activeTab} onValueChange={(tabValue) => setActiveTab(tabValue as PortalTab)} className="gap-0">
        <div className="flex flex-col gap-3 md:flex-row md:items-start md:gap-6">
          <div className="sticky top-0 z-20 -mx-4 border-b bg-background/95 px-4 backdrop-blur md:hidden">
            <TabsList
              ref={mobileTabListReference}
              variant="line"
              aria-label={t("navigationLabel")}
              className="group-data-[orientation=horizontal]/tabs:h-auto w-full justify-start gap-1 overflow-x-auto py-1 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
            >
              {sidebarItems.map((item) => (
                <TabsTrigger key={item.key} value={item.key} className={tabTriggerClassName}>
                  {item.icon}
                  {t(item.labelKey)}
                </TabsTrigger>
              ))}
            </TabsList>
          </div>

          <Card className="hidden w-56 shrink-0 gap-4 py-4 md:flex">
            <CardContent className="flex flex-col gap-4 px-4">
              <span className="px-3 text-left text-[10px] font-black uppercase tracking-widest text-muted-foreground">
                {t("navigation")}
              </span>
              <TabsList
                variant="line"
                aria-label={t("navigationLabel")}
                className="group-data-[orientation=horizontal]/tabs:h-auto h-auto w-full flex-col items-stretch gap-2 p-0"
              >
                {sidebarItems.map((item) => (
                  <TabsTrigger key={item.key} value={item.key} className={tabTriggerClassName}>
                    {item.icon}
                    <span className="truncate">{t(item.labelKey)}</span>
                  </TabsTrigger>
                ))}
              </TabsList>
            </CardContent>
          </Card>

          <div className="flex w-full min-w-0 flex-1 flex-col gap-6">
            <TabsContent value={PortalTab.Dashboard}>
              <PortalDashboardOverview dashboard={dashboardData} />
            </TabsContent>
            <TabsContent value={PortalTab.Appointments}>
              <PortalAppointments />
            </TabsContent>
            <TabsContent value={PortalTab.Encounters}>
              <PortalEncounters />
            </TabsContent>
            <TabsContent value={PortalTab.Observations}>
              <PortalObservations />
            </TabsContent>
            <TabsContent value={PortalTab.Conditions}>
              <PortalConditions />
            </TabsContent>
            <TabsContent value={PortalTab.Medications}>
              <PortalMedications />
            </TabsContent>
            <TabsContent value={PortalTab.Reports}>
              <PortalReports />
            </TabsContent>
            <TabsContent value={PortalTab.Imaging}>
              <PortalImaging />
            </TabsContent>
          </div>
        </div>
      </Tabs>
    </div>
  )
}
