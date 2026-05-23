import { useState, useMemo } from "react"
import { useNavigate } from "react-router-dom"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Card } from "../../shared/components/ui/Card"
import { Input } from "../../shared/components/ui/Input"
import { Button } from "../../shared/components/ui/Button"
import { newPatientSchema, type NewPatientFormData } from "./patient_schemas"
import { usePatientsQuery, useCreatePatientMutation } from "./queries"
import {
  Search,
  UserPlus,
  Users,
  User,
  ArrowRight,
  Activity,
  Database,
  Clock,
  Filter,
  X
} from "lucide-react"

type SortDirection = "asc" | "desc"
type SortField = "full_name" | "birth_date" | "document_id"

export const Patients = () => {
  const navigate = useNavigate()
  const [searchTerm, setSearchTerm] = useState("")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [sortField, setSortField] = useState<SortField>("full_name")
  const [sortDirection, setSortDirection] = useState<SortDirection>("asc")
  const [currentPage, setCurrentPage] = useState(1)
  const itemsPerPage = 5

  const { data: patients = [], isLoading } = usePatientsQuery(searchTerm)
  const createPatientMutation = useCreatePatientMutation()

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<NewPatientFormData>({
    resolver: zodResolver(newPatientSchema),
  })

  const onSubmit = async (formData: NewPatientFormData) => {
    try {
      await createPatientMutation.mutateAsync({
        full_name: formData.fullName,
        birth_date: formData.birthDate,
        document_id: formData.documentId,
        phone_number: formData.phoneNumber,
      })
      reset()
      setIsModalOpen(false)
    } catch {
      alert("Falha ao registrar paciente.")
    }
  }

  const handleSortToggle = (field: SortField) => {
    if (sortField === field) {
      setSortDirection(sortDirection === "asc" ? "desc" : "asc")
    } else {
      setSortField(field)
      setSortDirection("asc")
    }
    setCurrentPage(1)
  }

  const handleSearchChange = (termValue: string) => {
    setSearchTerm(termValue)
    setCurrentPage(1)
  }

  const sortedPatients = useMemo(() => {
    return [...patients].sort((patientA, patientB) => {
      const valueA = patientA[sortField] || ""
      const valueB = patientB[sortField] || ""
      const comparison = valueA.localeCompare(valueB)
      return sortDirection === "asc" ? comparison : -comparison
    })
  }, [patients, sortField, sortDirection])

  const paginatedPatients = useMemo(() => {
    const startIndex = (currentPage - 1) * itemsPerPage
    return sortedPatients.slice(startIndex, startIndex + itemsPerPage)
  }, [sortedPatients, currentPage])

  const totalPages = Math.ceil(sortedPatients.length / itemsPerPage)

  const sortIndicator = (field: SortField) => {
    if (sortField !== field) return ""
    return sortDirection === "asc" ? " ↑" : " ↓"
  }

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-5 animate-fade-in">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h2 className="text-xl font-black text-gray-900 tracking-tight leading-none">
            Pacientes
          </h2>
          <span className="text-xs text-muted mt-1 block">
            Gestão de prontuários e dados clínicos FHIR
          </span>
        </div>
        <Button onClick={() => setIsModalOpen(true)} className="py-2 px-4 self-start sm:self-auto gap-2">
          <UserPlus className="w-4 h-4" />
          Novo Paciente
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white border border-border rounded-xl p-4 flex items-center gap-4">
          <div className="w-10 h-10 rounded-lg bg-primary/8 flex items-center justify-center shrink-0">
            <Users className="w-5 h-5 text-primary" />
          </div>
          <div>
            <span className="text-[10px] text-muted font-semibold uppercase tracking-wider block">
              Total Cadastrados
            </span>
            <span className="text-2xl font-black text-gray-900 leading-none mt-0.5 block">
              {patients.length}
            </span>
          </div>
        </div>

        <div className="bg-white border border-border rounded-xl p-4 flex items-center gap-4">
          <div className="w-10 h-10 rounded-lg bg-secondary/8 flex items-center justify-center shrink-0">
            <Activity className="w-5 h-5 text-secondary" />
          </div>
          <div>
            <span className="text-[10px] text-muted font-semibold uppercase tracking-wider block">
              Padrão de Dados
            </span>
            <span className="text-sm font-bold text-gray-800 mt-0.5 block">FHIR R4 Compliant</span>
          </div>
        </div>

        <div className="bg-white border border-border rounded-xl p-4 flex items-center gap-4">
          <div className="w-10 h-10 rounded-lg bg-blue-50 flex items-center justify-center shrink-0">
            <Database className="w-5 h-5 text-blue-500" />
          </div>
          <div>
            <span className="text-[10px] text-muted font-semibold uppercase tracking-wider block">
              Integração
            </span>
            <span className="text-sm font-bold text-gray-800 mt-0.5 block">Cloud Healthcare API</span>
          </div>
        </div>
      </div>

      <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
        <div className="flex-1 flex items-center gap-2.5 bg-white border border-border rounded-lg px-4 py-2.5">
          <Search className="w-4 h-4 text-gray-400 shrink-0" />
          <input
            type="text"
            placeholder="Buscar por nome, CPF ou telefone..."
            value={searchTerm}
            onChange={(event) => handleSearchChange(event.target.value)}
            className="w-full bg-transparent text-sm text-gray-800 placeholder-gray-400 focus:outline-none"
          />
          {searchTerm && (
            <button
              onClick={() => handleSearchChange("")}
              className="p-0.5 rounded hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-colors"
            >
              <X className="w-3.5 h-3.5" />
            </button>
          )}
        </div>

        <div className="flex items-center justify-center sm:justify-start gap-1.5 bg-white border border-border rounded-lg px-3 py-2.5 shrink-0">
          <Filter className="w-3.5 h-3.5 text-gray-400" />
          <span className="text-[11px] text-muted font-medium">
            {sortedPatients.length} resultado{sortedPatients.length !== 1 ? "s" : ""}
          </span>
        </div>
      </div>

      {isLoading ? (
        <div className="flex-1 flex items-center justify-center py-20">
          <div className="flex items-center gap-3 text-muted">
            <Clock className="w-5 h-5 animate-spin" />
            <span className="text-sm font-medium">Carregando registros clínicos...</span>
          </div>
        </div>
      ) : sortedPatients.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center gap-3 py-20">
          <div className="w-14 h-14 rounded-2xl bg-gray-50 border border-border flex items-center justify-center">
            <Users className="w-7 h-7 text-gray-300" />
          </div>
          <h3 className="text-base font-bold text-gray-800">
            {searchTerm ? "Nenhum resultado encontrado" : "Nenhum paciente cadastrado"}
          </h3>
          <p className="text-xs text-muted max-w-xs text-center">
            {searchTerm
              ? `Nenhum paciente corresponde a "${searchTerm}". Tente outro termo.`
              : "Registre o primeiro paciente para começar a gestão clínica."
            }
          </p>
        </div>
      ) : (
        <div className="bg-white border border-border rounded-xl overflow-hidden">
          <div className="overflow-x-auto w-full">
            <table className="w-full text-left border-collapse min-w-[650px] md:min-w-0">
            <thead>
              <tr className="border-b border-border bg-gray-50/80">
                <th
                  className="py-3 px-5 text-[10px] font-bold text-muted uppercase tracking-widest cursor-pointer hover:text-gray-700 transition-colors select-none"
                  onClick={() => handleSortToggle("full_name")}
                >
                  Paciente{sortIndicator("full_name")}
                </th>
                <th
                  className="py-3 px-5 text-[10px] font-bold text-muted uppercase tracking-widest cursor-pointer hover:text-gray-700 transition-colors select-none"
                  onClick={() => handleSortToggle("document_id")}
                >
                  CPF / Documento{sortIndicator("document_id")}
                </th>
                <th
                  className="py-3 px-5 text-[10px] font-bold text-muted uppercase tracking-widest cursor-pointer hover:text-gray-700 transition-colors select-none"
                  onClick={() => handleSortToggle("birth_date")}
                >
                  Data de Nascimento{sortIndicator("birth_date")}
                </th>
                <th className="py-3 px-5 text-[10px] font-bold text-muted uppercase tracking-widest">
                  Telefone
                </th>
                <th className="py-3 px-5 text-[10px] font-bold text-muted uppercase tracking-widest text-right pr-5">
                  Ação
                </th>
              </tr>
            </thead>
            <tbody>
              {paginatedPatients.map((patient, rowIndex) => (
                <tr
                  key={patient.patient_id}
                  className={`border-b border-border/60 hover:bg-blue-50/30 transition-colors duration-150 group ${
                    rowIndex % 2 === 1 ? "bg-gray-50/30" : ""
                  }`}
                >
                  <td className="py-3 px-5">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-lg bg-primary/6 flex items-center justify-center shrink-0 group-hover:bg-primary/10 transition-colors">
                        <User className="w-4 h-4 text-primary/60 group-hover:text-primary transition-colors" />
                      </div>
                      <div className="min-w-0">
                        <span className="text-[13px] font-semibold text-gray-800 block truncate group-hover:text-primary transition-colors">
                          {patient.full_name}
                        </span>
                        <span className="text-[10px] text-gray-400 font-mono block mt-0.5">
                          {patient.fhir_resource_id}
                        </span>
                      </div>
                    </div>
                  </td>
                  <td className="py-3 px-5">
                    <span className="text-xs font-mono text-gray-600">
                      {patient.document_id}
                    </span>
                  </td>
                  <td className="py-3 px-5">
                    <span className="text-xs text-gray-500">
                      {patient.birth_date}
                    </span>
                  </td>
                  <td className="py-3 px-5">
                    <span className="text-xs text-gray-500">
                      {patient.phone_number}
                    </span>
                  </td>
                  <td className="py-3 px-5 text-right">
                    <button
                      onClick={() => navigate(`/patients/${patient.fhir_resource_id}`)}
                      className="inline-flex items-center gap-1.5 text-[11px] font-semibold text-primary hover:text-primary/80 px-3 py-1.5 rounded-md hover:bg-primary/5 transition-all"
                    >
                      Prontuário
                      <ArrowRight className="w-3 h-3 transition-transform group-hover:translate-x-0.5" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          </div>
          {totalPages > 1 && (
            <div className="flex items-center justify-between border-t border-border px-6 py-4 bg-gray-50/50">
              <span className="text-xs text-gray-500 font-semibold">
                Página {currentPage} de {totalPages} ({sortedPatients.length} pacientes)
              </span>
              <div className="flex gap-2">
                <Button
                  variantType="outline"
                  className="px-3 py-1.5 text-xs font-bold"
                  disabled={currentPage === 1}
                  onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                >
                  Anterior
                </Button>
                <Button
                  variantType="outline"
                  className="px-3 py-1.5 text-xs font-bold"
                  disabled={currentPage === totalPages}
                  onClick={() => setCurrentPage((prev) => Math.min(prev + 1, totalPages))}
                >
                  Próxima
                </Button>
              </div>
            </div>
          )}
        </div>
      )}

      {isModalOpen && (
        <div className="fixed inset-0 z-50 bg-black/30 backdrop-blur-[2px] flex items-center justify-center p-4">
          <Card glowingType="cyan" className="w-full max-w-[480px] p-7 relative animate-fade-in">
            <h3 className="text-lg font-bold text-gray-900 mb-5 text-left">Registrar Novo Paciente</h3>

            <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 text-left">
              <div className="flex flex-col gap-1">
                <label className="text-xs font-semibold text-gray-600">Nome Completo</label>
                <Input
                  type="text"
                  placeholder="Nome Completo do Paciente"
                  errorText={errors.fullName?.message}
                  {...register("fullName")}
                />
              </div>

              <div className="flex flex-col gap-1">
                <label className="text-xs font-semibold text-gray-600">Data de Nascimento</label>
                <Input
                  type="text"
                  placeholder="AAAA-MM-DD"
                  errorText={errors.birthDate?.message}
                  {...register("birthDate")}
                />
              </div>

              <div className="flex flex-col gap-1">
                <label className="text-xs font-semibold text-gray-600">CPF ou Documento de Identidade</label>
                <Input
                  type="text"
                  placeholder="123.456.789-00"
                  errorText={errors.documentId?.message}
                  {...register("documentId")}
                />
              </div>

              <div className="flex flex-col gap-1">
                <label className="text-xs font-semibold text-gray-600">Telefone de Contato</label>
                <Input
                  type="text"
                  placeholder="(11) 98765-4321"
                  errorText={errors.phoneNumber?.message}
                  {...register("phoneNumber")}
                />
              </div>

              <div className="flex gap-3 justify-end mt-3">
                <Button variantType="outline" type="button" onClick={() => setIsModalOpen(false)}>
                  Cancelar
                </Button>
                <Button type="submit" disabled={createPatientMutation.isPending}>
                  Confirmar Cadastro
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </div>
  )
}
