import { Suspense, lazy } from "react"
import { useParams, useNavigate, useSearchParams } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { usePatientQuery } from "./queries"
import { Action, Feature } from "../../shared/auth/AbilityContext"
import { useAbility } from "../../shared/auth/useAbility"
import { PatientPageLayout } from "./components/PatientPageLayout"
import { PatientPageSkeleton } from "./components/PatientPageSkeleton"
import { QueryErrorAlert } from "./components/QueryErrorAlert"
import { PageContainer } from "../../shared/components/ui/PageContainer"
import { Card } from "../../shared/components/ui/card"
import { Spinner } from "../../shared/components/ui/spinner"
import { Tabs, TabsList, TabsTrigger } from "../../shared/components/ui/tabs"
import { Heart, FileText, Pill } from "lucide-react"

const EncounterTab = {
  Vitals: "vitals",
  Reports: "reports",
  Medications: "medications",
} as const

type EncounterTab = (typeof EncounterTab)[keyof typeof EncounterTab]

const VitalSigns = lazy(() => import("./components/VitalSigns"))
const ClinicalReports = lazy(() => import("./components/ClinicalReports"))
const ClinicalMedications = lazy(() => import("./components/ClinicalMedications"))

const TabFallback = () => {
  const { t } = useTranslation("patients")

  return (
    <Card className="flex items-center justify-center min-h-[450px]">
      <div className="flex flex-col items-center gap-2">
        <Spinner className="size-8 text-primary" />
        <span className="text-sm text-muted-foreground font-medium">{t("loadingComponent")}</span>
      </div>
    </Card>
  )
}

export function EncounterDetail() {
  const { id = "", encounterId = "" } = useParams<{ id: string; encounterId: string }>()
  const navigate = useNavigate()
  const { t } = useTranslation("patients")

  const ability = useAbility()

  const [searchParameters, setSearchParameters] = useSearchParams()
  const activeTab = (searchParameters.get("tab") || EncounterTab.Vitals) as EncounterTab
  const setActiveTab = (tabName: EncounterTab) => {
    setSearchParameters({ tab: tabName })
  }

  const { data: patient, isLoading: isPatientLoading, isError: isPatientError, refetch: refetchPatient } = usePatientQuery(id)

  const canReadObservations = ability.can(Action.Read, Feature.Observation)
  const canReadReports = ability.can(Action.Read, Feature.DiagnosticReport)
  const canReadMedications = ability.can(Action.Read, Feature.MedicationRequest)

  const availableTabs: Record<EncounterTab, boolean> = {
    [EncounterTab.Vitals]: canReadObservations,
    [EncounterTab.Reports]: canReadReports,
    [EncounterTab.Medications]: canReadMedications,
  }
  const resolvedActiveTab = availableTabs[activeTab] ? activeTab : EncounterTab.Vitals

  if (isPatientError) {
    return (
      <PageContainer>
        <QueryErrorAlert title={t("details.loadError")} onRetry={() => void refetchPatient()} />
      </PageContainer>
    )
  }

  if (isPatientLoading || !patient) {
    return <PatientPageSkeleton />
  }

  const sidebarTabs = (
    <Tabs
      orientation="vertical"
      value={resolvedActiveTab}
      onValueChange={(tabName) => setActiveTab(tabName as EncounterTab)}
      className="w-full"
    >
      <TabsList variant="line" className="w-full gap-1 p-0">
        {canReadObservations && (
          <TabsTrigger
            value={EncounterTab.Vitals}
            className="h-auto flex-none justify-start gap-3 px-4 py-3 text-xs font-extrabold text-muted-foreground data-[state=active]:bg-primary/10 data-[state=active]:text-primary data-[state=active]:shadow-none"
          >
            <Heart className="w-4 h-4 shrink-0" />
            {t("details.vitals")}
          </TabsTrigger>
        )}

        {canReadReports && (
          <TabsTrigger
            value={EncounterTab.Reports}
            className="h-auto flex-none justify-start gap-3 px-4 py-3 text-xs font-extrabold text-muted-foreground data-[state=active]:bg-primary/10 data-[state=active]:text-primary data-[state=active]:shadow-none"
          >
            <FileText className="w-4 h-4 shrink-0" />
            {t("details.reports")}
          </TabsTrigger>
        )}

        {canReadMedications && (
          <TabsTrigger
            value={EncounterTab.Medications}
            className="h-auto flex-none justify-start gap-3 px-4 py-3 text-xs font-extrabold text-muted-foreground data-[state=active]:bg-primary/10 data-[state=active]:text-primary data-[state=active]:shadow-none"
          >
            <Pill className="w-4 h-4 shrink-0" />
            {t("details.medications")}
          </TabsTrigger>
        )}
      </TabsList>
    </Tabs>
  )

  return (
    <PatientPageLayout
      patient={patient}
      onBack={() => navigate(`/patients/${id}`)}
      sidebarTop={
        <>
          <span className="text-[10px] font-black text-muted-foreground uppercase tracking-widest px-3 text-left">
            {t("details.clinicalResources")}
          </span>
          {sidebarTabs}
        </>
      }
    >
      <Suspense fallback={<TabFallback />}>
        {
          {
            [EncounterTab.Vitals]: <VitalSigns patientId={id} encounterId={encounterId} />,
            [EncounterTab.Reports]: <ClinicalReports patientId={id} encounterId={encounterId} />,
            [EncounterTab.Medications]: <ClinicalMedications patientId={id} encounterId={encounterId} />,
          }[resolvedActiveTab]
        }
      </Suspense>
    </PatientPageLayout>
  )
}
