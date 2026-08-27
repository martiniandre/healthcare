import { useState, Suspense, lazy } from "react"
import { useParams, useNavigate, useSearchParams } from "react-router-dom"
import { useTranslation } from "react-i18next"
import {
  usePatientQuery,
  useEncountersQuery,
} from "./queries"
import { useImagingStudiesQuery } from "../imaging/queries"
import { ImagingWorkspace } from "../imaging/ImagingWorkspace"
import { Action, Feature } from "../../shared/auth/AbilityContext"
import { useAbility } from "../../shared/auth/useAbility"
import { EncounterSelectionDialog } from "./components/EncounterSelectionDialog"
import { ExamAnalyzerModal } from "./components/modals/ExamAnalyzerModal"
import { Can } from "../../shared/auth/AbilityContext"
import { PatientPageLayout } from "./components/PatientPageLayout"
import { Card, Button } from "../../shared/components/ui"
import {
  History,
  Image as ImageIcon,
  Activity,
  ShieldAlert,
  Sparkles,
  Loader2
} from "lucide-react"

const PatientTab = {
  Encounters: "encounters",
  Conditions: "conditions",
  Allergies: "allergies",
  Pacs: "pacs",
} as const

type PatientTab = typeof PatientTab[keyof typeof PatientTab]

const EncounterHistory = lazy(() => import("./components/EncounterHistory"))
const ClinicalConditions = lazy(() => import("./components/ClinicalConditions"))
const ClinicalAllergies = lazy(() => import("./components/ClinicalAllergies"))

import { PACSStudies } from "./components/PACSStudies"

const TabFallback = () => (
  <Card className="flex items-center justify-center min-h-[450px]">
    <div className="flex flex-col items-center gap-2">
      <Loader2 className="w-8 h-8 text-primary animate-spin" />
      <span className="text-sm text-muted-foreground font-medium">Carregando componente...</span>
    </div>
  </Card>
)

export const PatientDetails = () => {
  const { id = "" } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { t } = useTranslation("patients")

  const ability = useAbility()

  const [searchParameters, setSearchParameters] = useSearchParams()
  const activeTab = (searchParameters.get("tab") || PatientTab.Encounters) as PatientTab
  const [selectedStudyId, setSelectedStudyId] = useState<string | null>(null)
  const setActiveTab = (tabName: PatientTab) => {
    setSearchParameters({ tab: tabName })
    setSelectedStudyId(null)
  }
  const [isEncounterSelectionOpen, setIsEncounterSelectionOpen] = useState(false)

  const [isExamModalOpen, setIsExamModalOpen] = useState(false)

  const { data: patient, isLoading: isPatientLoading, isError: isPatientError } = usePatientQuery(id)
  const { data: encounters = [] } = useEncountersQuery(id)

  const canReadConditions = ability.can(Action.Read, Feature.Condition)
  const canReadAllergies = ability.can(Action.Read, Feature.Allergy)
  const canReadStudies = ability.can(Action.Read, Feature.ImagingStudy)

  const availableTabs: Record<string, boolean> = {
    [PatientTab.Encounters]: true,
    [PatientTab.Conditions]: canReadConditions,
    [PatientTab.Allergies]: canReadAllergies,
    [PatientTab.Pacs]: canReadStudies,
  }
  const resolvedActiveTab = availableTabs[activeTab] ? activeTab : PatientTab.Encounters

  const { data: studies = [] } = useImagingStudiesQuery(id, canReadStudies)

  if (isPatientError) {
    return (
      <div className="text-center py-16">
        <span className="text-sm text-danger font-medium">{t("details.loadError")}</span>
      </div>
    )
  }

  if (isPatientLoading || !patient) {
    return (
      <div className="text-center py-16">
        <span className="text-sm text-muted-foreground">{t("details.loadingDetails")}</span>
      </div>
    )
  }

  const sidebarTabs = (
    <>
      <Button
        variantType="ghost"
        onClick={() => setActiveTab(PatientTab.Encounters)}
        className={`w-full justify-start gap-3 px-4 py-3 h-auto rounded-lg text-xs font-extrabold transition-all duration-300 ${
          resolvedActiveTab === PatientTab.Encounters
            ? "bg-primary-soft text-primary hover:bg-primary-soft"
            : "text-muted-foreground hover:text-foreground hover:bg-muted-soft"
        }`}
      >
        <History className="w-4 h-4 shrink-0" />
        {t("details.encounters")}
        <span className="ml-auto text-[10px] bg-muted-soft text-muted-foreground px-2 py-0.5 rounded font-black">
          {encounters.length}
        </span>
      </Button>

      {canReadConditions && (
        <Button
          variantType="ghost"
          onClick={() => setActiveTab(PatientTab.Conditions)}
          className={`w-full justify-start gap-3 px-4 py-3 h-auto rounded-lg text-xs font-extrabold transition-all duration-300 ${
            resolvedActiveTab === PatientTab.Conditions
              ? "bg-primary-soft text-primary hover:bg-primary-soft"
              : "text-muted-foreground hover:text-foreground hover:bg-muted-soft"
          }`}
        >
          <Activity className="w-4 h-4 shrink-0" />
          {t("details.conditions")}
        </Button>
      )}

      {canReadAllergies && (
        <Button
          variantType="ghost"
          onClick={() => setActiveTab(PatientTab.Allergies)}
          className={`w-full justify-start gap-3 px-4 py-3 h-auto rounded-lg text-xs font-extrabold transition-all duration-300 ${
            resolvedActiveTab === PatientTab.Allergies
              ? "bg-primary-soft text-primary hover:bg-primary-soft"
              : "text-muted-foreground hover:text-foreground hover:bg-muted-soft"
          }`}
        >
          <ShieldAlert className="w-4 h-4 shrink-0" />
          {t("details.allergies")}
        </Button>
      )}

      {canReadStudies && (
        <Button
          variantType="ghost"
          onClick={() => setActiveTab(PatientTab.Pacs)}
          className={`w-full justify-start gap-3 px-4 py-3 h-auto rounded-lg text-xs font-extrabold transition-all duration-300 ${
            resolvedActiveTab === PatientTab.Pacs
              ? "bg-primary-soft text-primary hover:bg-primary-soft"
              : "text-muted-foreground hover:text-foreground hover:bg-muted-soft"
          }`}
        >
          <ImageIcon className="w-4 h-4 shrink-0" />
          {t("details.pacs")}
          <span className="ml-auto text-[10px] bg-muted-soft text-muted-foreground px-2 py-0.5 rounded font-black">
            {studies.length}
          </span>
        </Button>
      )}
    </>
  )

  return (
    <>
      <PatientPageLayout
        patient={patient}
        onBack={() => navigate("/")}
        sidebarTop={
          <>
            <span className="text-[10px] font-black text-muted-foreground uppercase tracking-widest px-3 text-left">
              {t("details.clinicalResources")}
            </span>
            {sidebarTabs}
          </>
        }
        headerActions={
          <Can I={Action.Create} a={Feature.ExamAnalysis}>
            <Button
              onClick={() => setIsExamModalOpen(true)}
              className="gap-2 shrink-0 self-start xl:self-auto bg-primary-soft text-primary hover:bg-primary/20 border border-primary/20 hover:border-primary/40 font-bold"
            >
              <Sparkles className="w-4 h-4 text-primary" />
              Analisar Exame com IA
            </Button>
          </Can>
        }
      >
        <Suspense fallback={<TabFallback />}>
          {resolvedActiveTab === PatientTab.Encounters && (
            <EncounterHistory
              patientId={id}
            />
          )}

          {resolvedActiveTab === PatientTab.Conditions && (
            <ClinicalConditions
              patientId={id}
            />
          )}

          {resolvedActiveTab === PatientTab.Allergies && (
            <ClinicalAllergies
              patientId={id}
            />
          )}

          {resolvedActiveTab === PatientTab.Pacs && (
            selectedStudyId ? (
              <ImagingWorkspace
                studyId={selectedStudyId}
                onBack={() => setSelectedStudyId(null)}
              />
            ) : (
              <PACSStudies
                studies={studies}
                onOpen={(studyId) => setSelectedStudyId(studyId)}
              />
            )
          )}
        </Suspense>
      </PatientPageLayout>

      <EncounterSelectionDialog
        isOpen={isEncounterSelectionOpen}
        onClose={() => setIsEncounterSelectionOpen(false)}
        encounters={encounters}
        patientId={id}
      />

      <ExamAnalyzerModal
        isOpen={isExamModalOpen}
        onClose={() => setIsExamModalOpen(false)}
        patientFhirId={id}
      />
    </>
  )
}
