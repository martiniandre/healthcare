import { useState, Suspense, lazy } from "react"
import { useParams, useNavigate, useSearchParams } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { usePatientQuery, useEncountersQuery } from "./queries"
import { useImagingStudiesQuery } from "../imaging/queries"
import { ImagingWorkspace } from "../imaging/ImagingWorkspace"
import { Action, Feature } from "../../shared/auth/AbilityContext"
import { useAbility } from "../../shared/auth/useAbility"
import { EncounterSelectionDialog } from "./components/EncounterSelectionDialog"
import { ExamAnalyzerModal } from "./components/modals/ExamAnalyzerModal"
import { Can } from "../../shared/auth/AbilityContext"
import { PatientPageLayout } from "./components/PatientPageLayout"
import { PatientPageSkeleton } from "./components/PatientPageSkeleton"
import { QueryErrorAlert } from "./components/QueryErrorAlert"
import { PageContainer } from "../../shared/components/ui/PageContainer"
import { Card } from "../../shared/components/ui/card"
import { Spinner } from "../../shared/components/ui/spinner"
import { Tabs, TabsList, TabsTrigger } from "../../shared/components/ui/tabs"
import { Badge } from "../../shared/components/ui/badge"
import { Button } from "../../shared/components/ui/button"
import { History, Image as ImageIcon, Activity, ShieldAlert, Sparkles } from "lucide-react"

const PatientTab = {
  Encounters: "encounters",
  Conditions: "conditions",
  Allergies: "allergies",
  Pacs: "pacs",
} as const

type PatientTab = (typeof PatientTab)[keyof typeof PatientTab]

const EncounterHistory = lazy(() => import("./components/EncounterHistory"))
const ClinicalConditions = lazy(() => import("./components/ClinicalConditions"))
const ClinicalAllergies = lazy(() => import("./components/ClinicalAllergies"))

import { PACSStudies } from "./components/PACSStudies"

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

  const { data: patient, isLoading: isPatientLoading, isError: isPatientError, refetch: refetchPatient } = usePatientQuery(id)
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

  const {
    data: studies = [],
    isLoading: isStudiesLoading,
    isError: isStudiesError,
    refetch: refetchStudies,
  } = useImagingStudiesQuery(id, canReadStudies)

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
      onValueChange={(tabName) => setActiveTab(tabName as PatientTab)}
      className="w-full"
    >
      <TabsList variant="line" className="w-full gap-1 p-0">
        <TabsTrigger
          value={PatientTab.Encounters}
          className="h-auto flex-none justify-start gap-3 px-4 py-3 text-xs font-extrabold text-muted-foreground data-[state=active]:bg-primary/10 data-[state=active]:text-primary data-[state=active]:shadow-none"
        >
          <History className="w-4 h-4 shrink-0" />
          {t("details.encounters")}
          <Badge variant="secondary" className="ml-auto">
            {encounters.length}
          </Badge>
        </TabsTrigger>

        {canReadConditions && (
          <TabsTrigger
            value={PatientTab.Conditions}
            className="h-auto flex-none justify-start gap-3 px-4 py-3 text-xs font-extrabold text-muted-foreground data-[state=active]:bg-primary/10 data-[state=active]:text-primary data-[state=active]:shadow-none"
          >
            <Activity className="w-4 h-4 shrink-0" />
            {t("details.conditions")}
          </TabsTrigger>
        )}

        {canReadAllergies && (
          <TabsTrigger
            value={PatientTab.Allergies}
            className="h-auto flex-none justify-start gap-3 px-4 py-3 text-xs font-extrabold text-muted-foreground data-[state=active]:bg-primary/10 data-[state=active]:text-primary data-[state=active]:shadow-none"
          >
            <ShieldAlert className="w-4 h-4 shrink-0" />
            {t("details.allergies")}
          </TabsTrigger>
        )}

        {canReadStudies && (
          <TabsTrigger
            value={PatientTab.Pacs}
            className="h-auto flex-none justify-start gap-3 px-4 py-3 text-xs font-extrabold text-muted-foreground data-[state=active]:bg-primary/10 data-[state=active]:text-primary data-[state=active]:shadow-none"
          >
            <ImageIcon className="w-4 h-4 shrink-0" />
            {t("details.pacs")}
            <Badge variant="secondary" className="ml-auto">
              {studies.length}
            </Badge>
          </TabsTrigger>
        )}
      </TabsList>
    </Tabs>
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
              className="gap-2 shrink-0 self-start xl:self-auto bg-primary/10 text-primary hover:bg-primary/20 border border-primary/20 hover:border-primary/40 font-bold"
            >
              <Sparkles className="w-4 h-4 text-primary" />
              {t("analyzeExamWithAi")}
            </Button>
          </Can>
        }
      >
        <Suspense fallback={<TabFallback />}>
          {resolvedActiveTab === PatientTab.Encounters && <EncounterHistory patientId={id} />}

          {resolvedActiveTab === PatientTab.Conditions && <ClinicalConditions patientId={id} />}

          {resolvedActiveTab === PatientTab.Allergies && <ClinicalAllergies patientId={id} />}

          {resolvedActiveTab === PatientTab.Pacs &&
            (selectedStudyId ? (
              <ImagingWorkspace studyId={selectedStudyId} onBack={() => setSelectedStudyId(null)} />
            ) : (
              <PACSStudies
                studies={studies}
                onOpen={(studyId) => setSelectedStudyId(studyId)}
                isLoading={isStudiesLoading}
                isError={isStudiesError}
                onRetry={() => void refetchStudies()}
              />
            ))}
        </Suspense>
      </PatientPageLayout>

      <EncounterSelectionDialog
        isOpen={isEncounterSelectionOpen}
        onClose={() => setIsEncounterSelectionOpen(false)}
        encounters={encounters}
        patientId={id}
      />

      <ExamAnalyzerModal isOpen={isExamModalOpen} onClose={() => setIsExamModalOpen(false)} patientFhirId={id} />
    </>
  )
}
