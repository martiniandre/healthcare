import { useState } from "react"
import { useParams, useNavigate, useSearchParams } from "react-router-dom"
import {
  usePatientQuery,
  useEncountersQuery,
  useCreateEncounterMutation,
  useObservationsQuery,
  useCreateObservationMutation,
  useDiagnosticReportsQuery,
  useCreateDiagnosticReportMutation,
  usePatientConditionsQuery,
  useCreateConditionMutation
} from "./queries"
import { useImagingStudiesQuery } from "../imaging/queries"
import { PatientHeader } from "./components/PatientHeader"
import { EncounterHistory } from "./components/EncounterHistory"
import { VitalSigns } from "./components/VitalSigns"
import { ClinicalReports } from "./components/ClinicalReports"
import { PACSStudies } from "./components/PACSStudies"
import { ClinicalConditions } from "./components/ClinicalConditions"
import { EncounterModal } from "./components/modals/EncounterModal"
import { ObservationModal } from "./components/modals/ObservationModal"
import { ReportModal } from "./components/modals/ReportModal"
import { ConditionModal } from "./components/modals/ConditionModal"
import { Card } from "../../shared/components/ui/Card"
import { 
  History, 
  Heart, 
  FileText, 
  Image as ImageIcon,
  AlertTriangle,
  FolderOpen,
  Activity
} from "lucide-react"

export const PatientDetails = () => {
  const { id = "" } = useParams<{ id: string }>()
  const navigate = useNavigate()

  const [searchParameters, setSearchParameters] = useSearchParams()
  const activeTab = (searchParameters.get("tab") || "encounters") as "encounters" | "vitals" | "reports" | "pacs" | "conditions"
  const setActiveTab = (tabName: "encounters" | "vitals" | "reports" | "pacs" | "conditions") => {
    setSearchParameters({ tab: tabName })
  }
  const [selectedEncounterId, setSelectedEncounterId] = useState<string | null>(null)

  const [isEncounterModalOpen, setIsEncounterModalOpen] = useState(false)
  const [isObservationModalOpen, setIsObservationModalOpen] = useState(false)
  const [isReportModalOpen, setIsReportModalOpen] = useState(false)
  const [isConditionModalOpen, setIsConditionModalOpen] = useState(false)

  const { data: patient, isLoading: isPatientLoading } = usePatientQuery(id)
  const { data: encounters = [] } = useEncountersQuery(id)
  
  const activeEncounterId = selectedEncounterId || (encounters.length > 0 ? encounters[encounters.length - 1].fhir_id : null)

  const { data: observations = [] } = useObservationsQuery(activeEncounterId || "")
  const { data: reports = [] } = useDiagnosticReportsQuery(activeEncounterId || "")
  const { data: studies = [] } = useImagingStudiesQuery(id)
  const { data: conditions = [] } = usePatientConditionsQuery(id)

  const createEncounterMutation = useCreateEncounterMutation()
  const createObservationMutation = useCreateObservationMutation()
  const createReportMutation = useCreateDiagnosticReportMutation()
  const createConditionMutation = useCreateConditionMutation()

  const selectedEncounter = encounters.find(e => e.fhir_id === activeEncounterId) || null

  const handleCreateEncounter = async (formData: { reasonDisplay: string }) => {
    try {
      const newEncounter = await createEncounterMutation.mutateAsync({
        patient_fhir_id: id,
        reason_display: formData.reasonDisplay,
        practitioner_id: "practitioner-1",
      })
      setIsEncounterModalOpen(false)
      setSelectedEncounterId(newEncounter.fhir_id)
      setActiveTab("encounters")
    } catch {
      alert("Falha ao registrar consulta.")
    }
  }

  const handleCreateObservation = async (formData: { loincCode: string; valueQuantity: number }) => {
    if (!selectedEncounter) {
      return
    }
    
    let display = "Frequência Cardíaca"
    let unit = "bpm"

    if (formData.loincCode === "8310-5") {
      display = "Temperatura Corporal"
      unit = "°C"
    } else if (formData.loincCode === "85354-9") {
      display = "Pressão Arterial Sistólica"
      unit = "mmHg"
    }

    try {
      await createObservationMutation.mutateAsync({
        encounter_fhir_id: selectedEncounter.fhir_id,
        patient_fhir_id: id,
        loinc_code: formData.loincCode,
        code_display: display,
        value_quantity: formData.valueQuantity,
        value_unit: unit,
      })
      setIsObservationModalOpen(false)
    } catch {
      alert("Erro ao registrar sinal vital.")
    }
  }

  const handleCreateReport = async (formData: { reportDisplay: string; conclusion: string }) => {
    if (!selectedEncounter) {
      return
    }
    try {
      await createReportMutation.mutateAsync({
        encounter_fhir_id: selectedEncounter.fhir_id,
        patient_fhir_id: id,
        report_display: formData.reportDisplay,
        conclusion: formData.conclusion,
      })
      setIsReportModalOpen(false)
    } catch {
      alert("Erro ao laudar exame.")
    }
  }

  const handleCreateCondition = async (formData: { icd10Code: string; codeDisplay: string }) => {
    try {
      await createConditionMutation.mutateAsync({
        patient_fhir_id: id,
        icd10_code: formData.icd10Code,
        code_display: formData.codeDisplay,
      })
      setIsConditionModalOpen(false)
    } catch {
      alert("Erro ao registrar diagnóstico.")
    }
  }

  if (isPatientLoading || !patient) {
    return (
      <div className="text-center py-16">
        <span className="text-sm text-muted">Carregando ficha clínica...</span>
      </div>
    )
  }

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full">
      <PatientHeader patient={patient} onBack={() => navigate("/")} />

      <div className="flex flex-col md:flex-row gap-6 items-start mt-2">
        <div className="w-full md:w-64 shrink-0 bg-white border border-border p-4 rounded-xl flex flex-col gap-4">
          <span className="text-[10px] font-black text-gray-500 uppercase tracking-widest px-3 text-left">
            Recursos Clínicos
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
              Atendimentos
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
              Sinais Vitais
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
              Laudos Clínicos
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
              Diagnósticos Ativos
              <span className="ml-auto text-[10px] bg-gray-100 text-gray-500 px-2 py-0.5 rounded font-black">
                {conditions.length}
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
              PACS DICOM
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
                    Consulta em Foco
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
                Alterar Consulta
              </button>
            </div>
          )}

          <div className="flex flex-col gap-6">
            {activeTab === "encounters" && (
              <EncounterHistory
                encounters={encounters}
                selectedEncounterId={activeEncounterId}
                onSelect={setSelectedEncounterId}
                onNew={() => setIsEncounterModalOpen(true)}
              />
            )}

            {activeTab === "vitals" && (
              selectedEncounter ? (
                <VitalSigns
                  observations={observations}
                  onAdd={() => setIsObservationModalOpen(true)}
                />
              ) : (
                <Card className="py-20 text-center">
                  <AlertTriangle className="w-12 h-12 text-gray-300 mx-auto mb-3" />
                  <h3 className="text-lg font-bold text-gray-800">Nenhum atendimento ativo</h3>
                  <p className="text-sm text-muted">
                    Por favor, selecione ou registre um atendimento na aba "Atendimentos".
                  </p>
                </Card>
              )
            )}

            {activeTab === "reports" && (
              selectedEncounter ? (
                <ClinicalReports
                  reports={reports}
                  onAdd={() => setIsReportModalOpen(true)}
                />
              ) : (
                <Card className="py-20 text-center">
                  <AlertTriangle className="w-12 h-12 text-gray-300 mx-auto mb-3" />
                  <h3 className="text-lg font-bold text-gray-800">Nenhum atendimento ativo</h3>
                  <p className="text-sm text-muted">
                    Por favor, selecione ou registre um atendimento na aba "Atendimentos".
                  </p>
                </Card>
              )
            )}

            {activeTab === "conditions" && (
              <ClinicalConditions
                conditions={conditions}
                onAdd={() => setIsConditionModalOpen(true)}
              />
            )}

            {activeTab === "pacs" && (
              <PACSStudies
                studies={studies}
                onOpen={(studyId) => navigate(`/imaging/${studyId}`)}
              />
            )}
          </div>
        </div>
      </div>

      <EncounterModal
        isOpen={isEncounterModalOpen}
        onClose={() => setIsEncounterModalOpen(false)}
        onSubmit={handleCreateEncounter}
        isPending={createEncounterMutation.isPending}
      />

      <ObservationModal
        isOpen={isObservationModalOpen}
        onClose={() => setIsObservationModalOpen(false)}
        onSubmit={handleCreateObservation}
        isPending={createObservationMutation.isPending}
      />

      <ReportModal
        isOpen={isReportModalOpen}
        onClose={() => setIsReportModalOpen(false)}
        onSubmit={handleCreateReport}
        isPending={createReportMutation.isPending}
      />

      <ConditionModal
        isOpen={isConditionModalOpen}
        onClose={() => setIsConditionModalOpen(false)}
        onSubmit={handleCreateCondition}
        isPending={createConditionMutation.isPending}
      />
    </div>
  )
}
