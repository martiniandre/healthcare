import { useState, Suspense, lazy } from "react"
import { useParams, useNavigate, useSearchParams } from "react-router-dom"
import { useTranslation } from "react-i18next"
import {
  usePatientQuery,
  useEncountersQuery,
  usePatientConditionsQuery,
  usePatientAllergiesQuery,
} from "./queries"
import { useImagingStudiesQuery } from "../imaging/queries"
import { PatientHeader } from "./components/PatientHeader"

const EncounterHistory = lazy(() => import("./components/EncounterHistory"))
const VitalSigns = lazy(() => import("./components/VitalSigns"))
const ClinicalReports = lazy(() => import("./components/ClinicalReports"))
const ClinicalConditions = lazy(() => import("./components/ClinicalConditions"))
const ClinicalAllergies = lazy(() => import("./components/ClinicalAllergies"))
const ClinicalMedications = lazy(() => import("./components/ClinicalMedications"))

import { PACSStudies } from "./components/PACSStudies"
import { ExamAnalyzerModal } from "./components/modals/ExamAnalyzerModal"
import { Card } from "../../shared/components/ui/Card"
import { Button } from "../../shared/components/ui/Button"
import { 
  History, 
  Heart, 
  FileText, 
  Image as ImageIcon,
  AlertTriangle,
  FolderOpen,
  Activity,
  ShieldAlert,
  Pill,
  Sparkles,
  Loader2
} from "lucide-react"

export const PatientDetails = () => {
  const { id = "" } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { t } = useTranslation("patients")

  const [searchParameters, setSearchParameters] = useSearchParams()
  const activeTab = (searchParameters.get("tab") || "encounters") as "encounters" | "vitals" | "reports" | "pacs" | "conditions" | "allergies" | "medications"
  const setActiveTab = (tabName: "encounters" | "vitals" | "reports" | "pacs" | "conditions" | "allergies" | "medications") => {
    setSearchParameters({ tab: tabName })
  }
  const [selectedEncounterId, setSelectedEncounterId] = useState<string | null>(null)

  const [isExamModalOpen, setIsExamModalOpen] = useState(false)

  const { data: patient, isLoading: isPatientLoading } = usePatientQuery(id)
  const { data: encounters = [] } = useEncountersQuery(id)
  
  const activeEncounterId = selectedEncounterId || (encounters.length > 0 ? encounters[encounters.length - 1].fhir_id : null)

  const { data: conditions = [] } = usePatientConditionsQuery(id)
  const { data: allergies = [] } = usePatientAllergiesQuery(id)
  const { data: studies = [] } = useImagingStudiesQuery(id)

  const selectedEncounter = encounters.find((encounterItem) => encounterItem.fhir_id === activeEncounterId) || null

  if (isPatientLoading || !patient) {
    return (
      <div className="text-center py-16">
        <span className="text-sm text-muted">{t("details.loadingDetails")}</span>
      </div>
    )
  }

  const TabFallback = () => (
    <Card className="flex items-center justify-center min-h-[450px]">
      <div className="flex flex-col items-center gap-2">
        <Loader2 className="w-8 h-8 text-primary animate-spin" />
        <span className="text-sm text-gray-500 font-medium">Carregando componente...</span>
      </div>
    </Card>
  )

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full">
      <div className="flex flex-col xl:flex-row xl:items-start justify-between gap-4">
        <PatientHeader patient={patient} onBack={() => navigate("/")} />
        <Button 
          onClick={() => setIsExamModalOpen(true)} 
          className="gap-2 shrink-0 self-start xl:self-auto bg-primary/10 text-primary hover:bg-primary/20 border border-primary/20 hover:border-primary/40 font-bold"
        >
          <Sparkles className="w-4 h-4 text-primary" />
          Analisar Exame com IA
        </Button>
      </div>

      <div className="flex flex-col md:flex-row gap-6 items-start mt-2">
        <div className="w-full md:w-64 shrink-0 bg-white border border-border p-4 rounded-xl flex flex-col gap-4">
          <span className="text-[10px] font-black text-gray-500 uppercase tracking-widest px-3 text-left">
            {t("details.clinicalResources")}
          </span>
          <div className="flex flex-col gap-2">
            <button
              onClick={() => setActiveTab("encounters")}
              className={`w-full text-left flex items-center gap-3 px-4 py-3 rounded-lg text-xs font-extrabold transition-all duration-300 ${
                activeTab === "encounters"
                  ? "bg-primary/8 text-primary"
                  : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
              }`}
            >
              <History className="w-4 h-4 shrink-0" />
              {t("details.encounters")}
              <span className="ml-auto text-[10px] bg-gray-100 text-gray-500 px-2 py-0.5 rounded font-black">
                {encounters.length}
              </span>
            </button>

            <button
              onClick={() => setActiveTab("vitals")}
              className={`w-full text-left flex items-center gap-3 px-4 py-3 rounded-lg text-xs font-extrabold transition-all duration-300 ${
                activeTab === "vitals"
                  ? "bg-primary/8 text-primary"
                  : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
              }`}
            >
              <Heart className="w-4 h-4 shrink-0" />
              {t("details.vitals")}
            </button>

            <button
              onClick={() => setActiveTab("reports")}
              className={`w-full text-left flex items-center gap-3 px-4 py-3 rounded-lg text-xs font-extrabold transition-all duration-300 ${
                activeTab === "reports"
                  ? "bg-primary/8 text-primary"
                  : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
              }`}
            >
              <FileText className="w-4 h-4 shrink-0" />
              {t("details.reports")}
            </button>

            <button
              onClick={() => setActiveTab("conditions")}
              className={`w-full text-left flex items-center gap-3 px-4 py-3 rounded-lg text-xs font-extrabold transition-all duration-300 ${
                activeTab === "conditions"
                  ? "bg-primary/8 text-primary"
                  : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
              }`}
            >
              <Activity className="w-4 h-4 shrink-0" />
              {t("details.conditions")}
              <span className="ml-auto text-[10px] bg-gray-100 text-gray-500 px-2 py-0.5 rounded font-black">
                {conditions.length}
              </span>
            </button>

            <button
              onClick={() => setActiveTab("medications")}
              className={`w-full text-left flex items-center gap-3 px-4 py-3 rounded-lg text-xs font-extrabold transition-all duration-300 ${
                activeTab === "medications"
                  ? "bg-primary/8 text-primary"
                  : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
              }`}
            >
              <Pill className="w-4 h-4 shrink-0" />
              {t("details.medications")}
            </button>

            <button
              onClick={() => setActiveTab("allergies")}
              className={`w-full text-left flex items-center gap-3 px-4 py-3 rounded-lg text-xs font-extrabold transition-all duration-300 ${
                activeTab === "allergies"
                  ? "bg-primary/8 text-primary"
                  : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
              }`}
            >
              <ShieldAlert className="w-4 h-4 shrink-0" />
              {t("details.allergies")}
              <span className="ml-auto text-[10px] bg-gray-100 text-gray-500 px-2 py-0.5 rounded font-black">
                {allergies.length}
              </span>
            </button>

            <button
              onClick={() => setActiveTab("pacs")}
              className={`w-full text-left flex items-center gap-3 px-4 py-3 rounded-lg text-xs font-extrabold transition-all duration-300 ${
                activeTab === "pacs"
                  ? "bg-primary/8 text-primary"
                  : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
              }`}
            >
              <ImageIcon className="w-4 h-4 shrink-0" />
              {t("details.pacs")}
              <span className="ml-auto text-[10px] bg-gray-100 text-gray-500 px-2 py-0.5 rounded font-black">
                {studies.length}
              </span>
            </button>
          </div>
        </div>

        <div className="flex-1 flex flex-col gap-6 min-w-0 w-full">
          {selectedEncounter && activeTab !== "encounters" && activeTab !== "pacs" && (
            <div className="flex items-center justify-between bg-primary/5 border border-primary/20 p-4 rounded-xl text-left">
              <div className="flex items-center gap-3">
                <FolderOpen className="w-5 h-5 text-primary shrink-0" />
                <div>
                  <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">
                    {t("details.activeEncounterFocus")}
                  </span>
                  <span className="text-sm font-bold text-gray-800">
                    {selectedEncounter.reason_display}
                  </span>
                </div>
              </div>
              <button
                onClick={() => setActiveTab("encounters")}
                className="text-xs text-primary hover:underline font-bold shrink-0"
              >
                {t("details.changeEncounter")}
              </button>
            </div>
          )}

          <div className="flex flex-col gap-6">
            <Suspense fallback={<TabFallback />}>
              {activeTab === "encounters" && (
                <EncounterHistory
                  patientId={id}
                  selectedEncounterId={activeEncounterId}
                  onSelect={setSelectedEncounterId}
                />
              )}

              {activeTab === "vitals" && (
                selectedEncounter ? (
                  <VitalSigns
                    patientId={id}
                    encounterId={selectedEncounter.fhir_id}
                  />
                ) : (
                  <Card className="py-20 text-center">
                    <AlertTriangle className="w-12 h-12 text-gray-300 mx-auto mb-3" />
                    <h3 className="text-lg font-bold text-gray-800">
                      {t("details.noActiveEncounter")}
                    </h3>
                    <p className="text-sm text-muted">
                      {t("details.selectEncounterDesc")}
                    </p>
                  </Card>
                )
              )}

              {activeTab === "reports" && (
                selectedEncounter ? (
                  <ClinicalReports
                    patientId={id}
                    encounterId={selectedEncounter.fhir_id}
                  />
                ) : (
                  <Card className="py-20 text-center">
                    <AlertTriangle className="w-12 h-12 text-gray-300 mx-auto mb-3" />
                    <h3 className="text-lg font-bold text-gray-800">
                      {t("details.noActiveEncounter")}
                    </h3>
                    <p className="text-sm text-muted">
                      {t("details.selectEncounterDesc")}
                    </p>
                  </Card>
                )
              )}

              {activeTab === "conditions" && (
                <ClinicalConditions
                  patientId={id}
                />
              )}

              {activeTab === "medications" && (
                selectedEncounter ? (
                  <ClinicalMedications
                    patientId={id}
                    encounterId={selectedEncounter.fhir_id}
                  />
                ) : (
                  <Card className="py-20 text-center">
                    <AlertTriangle className="w-12 h-12 text-gray-300 mx-auto mb-3" />
                    <h3 className="text-lg font-bold text-gray-800">
                      {t("details.noActiveEncounter")}
                    </h3>
                    <p className="text-sm text-muted">
                      {t("details.selectEncounterDesc")}
                    </p>
                  </Card>
                )
              )}

              {activeTab === "allergies" && (
                <ClinicalAllergies
                  patientId={id}
                />
              )}

              {activeTab === "pacs" && (
                <PACSStudies
                  studies={studies}
                  onOpen={(studyId) => navigate(`/imaging/${studyId}`)}
                />
              )}
            </Suspense>
          </div>
        </div>
      </div>

      <ExamAnalyzerModal
        isOpen={isExamModalOpen}
        onClose={() => setIsExamModalOpen(false)}
        patientFhirId={id}
      />
    </div>
  )
}

