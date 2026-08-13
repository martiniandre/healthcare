import { Suspense, lazy } from "react"
import { useParams, useNavigate, useSearchParams } from "react-router-dom"
import { useTranslation } from "react-i18next"
import {
  usePatientQuery,
} from "./queries"
import { Action, Feature } from "../../shared/auth/AbilityContext"
import { useAbility } from "../../shared/auth/useAbility"
import { PatientPageLayout } from "./components/PatientPageLayout"
import { Card } from "../../shared/components/ui/Card"
import { Heart, FileText, Pill, Loader2 } from "lucide-react"

const EncounterTab = {
  Vitals: "vitals",
  Reports: "reports",
  Medications: "medications",
} as const

type EncounterTab = typeof EncounterTab[keyof typeof EncounterTab]

const VitalSigns = lazy(() => import("./components/VitalSigns"))
const ClinicalReports = lazy(() => import("./components/ClinicalReports"))
const ClinicalMedications = lazy(() => import("./components/ClinicalMedications"))

const TabFallback = () => (
  <Card className="flex items-center justify-center min-h-[450px]">
    <div className="flex flex-col items-center gap-2">
      <Loader2 className="w-8 h-8 text-primary animate-spin" />
      <span className="text-sm text-gray-500 font-medium">Carregando componente...</span>
    </div>
  </Card>
)

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

  const { data: patient, isLoading: isPatientLoading, isError: isPatientError } = usePatientQuery(id)

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
      <div className="text-center py-16">
        <span className="text-sm text-red-500 font-medium">{t("details.loadError")}</span>
      </div>
    )
  }

  if (isPatientLoading || !patient) {
    return (
      <div className="text-center py-16">
        <span className="text-sm text-muted">{t("details.loadingDetails")}</span>
      </div>
    )
  }

  const sidebarTabs = (
    <>
      {canReadObservations && (
        <button
          onClick={() => setActiveTab(EncounterTab.Vitals)}
          className={`w-full text-left flex items-center gap-3 px-4 py-3 rounded-lg text-xs font-extrabold transition-all duration-300 ${
            resolvedActiveTab === EncounterTab.Vitals
              ? "bg-primary/8 text-primary"
              : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
          }`}
        >
          <Heart className="w-4 h-4 shrink-0" />
          {t("details.vitals")}
        </button>
      )}

      {canReadReports && (
        <button
          onClick={() => setActiveTab(EncounterTab.Reports)}
          className={`w-full text-left flex items-center gap-3 px-4 py-3 rounded-lg text-xs font-extrabold transition-all duration-300 ${
            resolvedActiveTab === EncounterTab.Reports
              ? "bg-primary/8 text-primary"
              : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
          }`}
        >
          <FileText className="w-4 h-4 shrink-0" />
          {t("details.reports")}
        </button>
      )}

      {canReadMedications && (
        <button
          onClick={() => setActiveTab(EncounterTab.Medications)}
          className={`w-full text-left flex items-center gap-3 px-4 py-3 rounded-lg text-xs font-extrabold transition-all duration-300 ${
            resolvedActiveTab === EncounterTab.Medications
              ? "bg-primary/8 text-primary"
              : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
          }`}
        >
          <Pill className="w-4 h-4 shrink-0" />
          {t("details.medications")}
        </button>
      )}
    </>
  )

  return (
    <PatientPageLayout
      patient={patient}
      onBack={() => navigate(`/patients/${id}`)}
      sidebarTop={
        <>
          <span className="text-[10px] font-black text-gray-500 uppercase tracking-widest px-3 text-left">
            {t("details.clinicalResources")}
          </span>
          {sidebarTabs}
        </>
      }
    >
      <Suspense fallback={<TabFallback />}>
        {{
          [EncounterTab.Vitals]: <VitalSigns patientId={id} encounterId={encounterId} />,
          [EncounterTab.Reports]: <ClinicalReports patientId={id} encounterId={encounterId} />,
          [EncounterTab.Medications]: <ClinicalMedications patientId={id} encounterId={encounterId} />,
        }[resolvedActiveTab]}
      </Suspense>
    </PatientPageLayout>
  )
}
