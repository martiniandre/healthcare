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
import { patientTabListClassName, patientTabTriggerClassName } from "./components/patientTabStyles"
import { PatientPageSkeleton } from "./components/PatientPageSkeleton"
import { QueryErrorAlert } from "@/shared/components/ui/QueryErrorAlert"
import { PageContainer } from "../../shared/components/ui/PageContainer"
import { Card } from "../../shared/components/ui/card"
import { Spinner } from "../../shared/components/ui/spinner"
import { Tabs, TabsList, TabsTrigger } from "../../shared/components/ui/tabs"
import { Badge } from "../../shared/components/ui/badge"
import { Button } from "../../shared/components/ui/button"
import { Sparkles } from "lucide-react"

const PatientTab = {
  Encounters: "encounters",
  Timeline: "timeline",
  Conditions: "conditions",
  Allergies: "allergies",
  Pacs: "pacs",
} as const

type PatientTab = (typeof PatientTab)[keyof typeof PatientTab]

const EncounterHistory = lazy(() => import("./components/EncounterHistory"))
const ClinicalTimeline = lazy(() => import("./components/timeline/ClinicalTimeline"))
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
  const canReadTimeline = ability.can(Action.Read, Feature.Timeline)

  const availableTabs: Record<string, boolean> = {
    [PatientTab.Encounters]: true,
    [PatientTab.Timeline]: canReadTimeline,
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
      orientation="horizontal"
      value={resolvedActiveTab}
      onValueChange={(tabName) => setActiveTab(tabName as PatientTab)}
      className="w-full"
    >
      <TabsList variant="line" className={patientTabListClassName}>
        <TabsTrigger
          value={PatientTab.Encounters}
          className={patientTabTriggerClassName}
        >
          {t("details.encounters")}
          <Badge variant="secondary" className="rounded bg-muted px-1.5 text-2xs tabular-nums text-muted-foreground group-data-[state=active]/trigger:bg-sage-100 group-data-[state=active]/trigger:text-sage-800">
            {encounters.length}
          </Badge>
        </TabsTrigger>

        {canReadTimeline && (
          <TabsTrigger
            value={PatientTab.Timeline}
            className={patientTabTriggerClassName}
          >
            {t("details.timeline")}
          </TabsTrigger>
        )}

        {canReadConditions && (
          <TabsTrigger
            value={PatientTab.Conditions}
            className={patientTabTriggerClassName}
          >
            {t("details.conditions")}
          </TabsTrigger>
        )}

        {canReadAllergies && (
          <TabsTrigger
            value={PatientTab.Allergies}
            className={patientTabTriggerClassName}
          >
            {t("details.allergies")}
          </TabsTrigger>
        )}

        {canReadStudies && (
          <TabsTrigger
            value={PatientTab.Pacs}
            className={patientTabTriggerClassName}
          >
            {t("details.pacs")}
            <Badge variant="secondary" className="rounded bg-muted px-1.5 text-2xs tabular-nums text-muted-foreground group-data-[state=active]/trigger:bg-sage-100 group-data-[state=active]/trigger:text-sage-800">
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
            <span className="sr-only">
              {t("details.clinicalResources")}
            </span>
            {sidebarTabs}
          </>
        }
        headerActions={
          <Can I={Action.Create} a={Feature.ExamAnalysis}>
            <Button
              onClick={() => setIsExamModalOpen(true)}
              variant="outline"
              className="w-full sm:w-auto"
            >
              <Sparkles />
              {t("analyzeExamWithAi")}
            </Button>
          </Can>
        }
      >
        <Suspense fallback={<TabFallback />}>
          {resolvedActiveTab === PatientTab.Encounters && <EncounterHistory patientId={id} />}

          {resolvedActiveTab === PatientTab.Timeline && <ClinicalTimeline patientId={id} />}

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
