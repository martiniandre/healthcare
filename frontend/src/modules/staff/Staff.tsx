import { useState } from "react"
import type { FormEvent } from "react"
import { Card } from "../../shared/components/ui/Card"
import { Button } from "../../shared/components/ui/Button"
import { Input } from "../../shared/components/ui/Input"
import { 
  Users, 
  Search, 
  UserPlus, 
  UserCheck, 
  Mail, 
  Clock,
  X
} from "lucide-react"

import { StaffRole, StaffStatus } from "../../shared/types"
import { useStaffListQuery, useCreateEmployeeMutation } from "./queries"
import { toast } from "../../shared/store/toast_store"

export const Staff = () => {
  const [filterRole, setFilterRole] = useState<string>("All")
  const [searchQuery, setSearchQuery] = useState<string>(" ")
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false)

  const { data: staffList = [], isLoading } = useStaffListQuery()
  const createEmployeeMutation = useCreateEmployeeMutation()

  const [newStaffName, setNewStaffName] = useState<string>("")
  const [newStaffRole, setNewStaffRole] = useState<StaffRole>(StaffRole.Doctor)
  const [newStaffLicense, setNewStaffLicense] = useState<string>("")
  const [newStaffEmail, setNewStaffEmail] = useState<string>("")
  const [newStaffDept, setNewStaffDept] = useState<string>("")

  const handleRegisterStaff = async (event: FormEvent) => {
    event.preventDefault()
    if (!newStaffName || !newStaffEmail) {
      return
    }

    try {
      const temporaryRandomUserId = crypto.randomUUID()
      await createEmployeeMutation.mutateAsync({
        userId: temporaryRandomUserId,
        fullName: newStaffName,
        email: newStaffEmail,
        role: newStaffRole,
        crmNumber: newStaffLicense || "N/A",
      })

      setIsModalOpen(false)
      setNewStaffName("")
      setNewStaffRole(StaffRole.Doctor)
      setNewStaffLicense("")
      setNewStaffEmail("")
      setNewStaffDept("")
      toast.success("Profissional de saúde cadastrado com sucesso!")
    } catch {
      toast.error("Falha ao registrar profissional de saúde.")
    }
  }

  const filteredStaff = staffList.filter((member) => {
    const matchesRole = filterRole === "All" || member.role === filterRole
    const matchesSearch = member.fullName.toLowerCase().includes(searchQuery.trim().toLowerCase()) ||
      member.email.toLowerCase().includes(searchQuery.trim().toLowerCase()) ||
      member.department.toLowerCase().includes(searchQuery.trim().toLowerCase())
    return matchesRole && matchesSearch
  })

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full select-none relative">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div className="text-left">
          <h2 className="text-xl font-black text-gray-900 leading-none flex items-center gap-2">
            <Users className="w-5 h-5 text-primary animate-pulse-glow" />
            Gestão de Equipes Hospitalares
          </h2>
          <span className="text-xs text-muted mt-1.5 block">
            Administração de profissionais de saúde, plantonistas ativos e credenciais operacionais
          </span>
        </div>

        <Button
          variantType="primary"
          onClick={() => setIsModalOpen(true)}
          className="px-4 gap-2 text-xs font-bold"
        >
          <UserPlus className="w-4 h-4" />
          Cadastrar Profissional
        </Button>
      </div>

      <Card className="p-4 flex flex-col gap-4">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="relative flex-1 max-w-md">
            <Search className="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              type="text"
              placeholder="Buscar por nome, e-mail ou especialidade..."
              value={searchQuery === " " ? "" : searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full bg-white border border-border rounded-lg pl-9 pr-4 py-2 text-xs text-gray-800 placeholder-gray-400 focus:outline-none focus:border-primary/50 transition-all duration-200"
            />
          </div>

          <div className="flex gap-2 flex-wrap">
            {["All", StaffRole.Doctor, StaffRole.Nurse, StaffRole.Receptionist, StaffRole.Admin].map((roleOption) => (
              <button
                key={roleOption}
                onClick={() => setFilterRole(roleOption)}
                className={`px-3 py-1.5 rounded-lg text-xs font-bold transition-all duration-200 border ${
                  filterRole === roleOption
                    ? "bg-primary/5 text-primary border-primary"
                    : "bg-white text-gray-500 border-border hover:bg-gray-50 hover:text-gray-900"
                }`}
              >
                {roleOption === "All" ? "Todos" : roleOption}
              </button>
            ))}
          </div>
        </div>

        {isLoading ? (
          <div className="text-center py-16">
            <span className="text-sm text-muted">Carregando corpo clínico...</span>
          </div>
        ) : (
          <div className="overflow-x-auto border border-border rounded-xl w-full">
            <table className="w-full text-left text-xs border-collapse min-w-[700px] md:min-w-0">
              <thead>
                <tr className="bg-gray-50/50 border-b border-border text-gray-500 font-bold uppercase tracking-wider">
                  <th className="py-3 px-4">Profissional</th>
                  <th className="py-3 px-4">Função / Categoria</th>
                  <th className="py-3 px-4">Registro (CRM/COREN)</th>
                  <th className="py-3 px-4">Departamento</th>
                  <th className="py-3 px-4">Contato</th>
                  <th className="py-3 px-4">Escala</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border text-gray-700 font-medium bg-white">
                {filteredStaff.map((member) => (
                  <tr key={member.id} className="hover:bg-gray-50/30 transition-all duration-150">
                    <td className="py-4 px-4">
                      <div className="flex items-center gap-3">
                        <div className="bg-primary/8 p-2 rounded-lg border border-primary/10">
                          <Users className="w-4 h-4 text-primary" />
                        </div>
                        <span className="font-extrabold text-gray-900">{member.fullName}</span>
                      </div>
                    </td>
                    <td className="py-4 px-4">
                      <span className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full font-bold text-[10px] uppercase border ${
                        member.role === StaffRole.Doctor 
                          ? "bg-blue-50 text-blue-600 border-blue-100" 
                          : member.role === StaffRole.Nurse
                            ? "bg-teal-50 text-teal-600 border-teal-100"
                            : member.role === StaffRole.Admin
                              ? "bg-purple-50 text-purple-600 border-purple-100"
                              : "bg-gray-50 text-gray-600 border-gray-100"
                      }`}>
                        {member.role}
                      </span>
                    </td>
                    <td className="py-4 px-4 font-semibold text-gray-600">{member.license}</td>
                    <td className="py-4 px-4 font-semibold text-gray-600">{member.department || "Geral"}</td>
                    <td className="py-4 px-4">
                      <span className="flex items-center gap-1.5 text-gray-500 font-semibold">
                        <Mail className="w-3.5 h-3.5 text-gray-400" />
                        {member.email}
                      </span>
                    </td>
                    <td className="py-4 px-4">
                      <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full font-bold text-[10px] uppercase ${
                        member.status === StaffStatus.OnDuty
                          ? "bg-emerald-50 text-emerald-600 border border-emerald-100"
                          : "bg-gray-50 text-gray-400 border border-gray-100"
                      }`}>
                        {member.status === StaffStatus.OnDuty ? (
                          <>
                            <UserCheck className="w-3 h-3" />
                            Plantonista
                          </>
                        ) : (
                          <>
                            <Clock className="w-3 h-3" />
                            Fora de Escala
                          </>
                        )}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>

      {isModalOpen && (
        <div className="fixed inset-0 z-50 bg-black/30 backdrop-blur-[2px] flex items-center justify-center p-4">
          <Card className="w-full max-w-[480px] p-6 text-left flex flex-col gap-5 border border-border">
            <div className="flex items-center justify-between border-b border-border pb-3">
              <h3 className="text-md font-bold text-gray-900 flex items-center gap-2">
                <UserPlus className="w-5 h-5 text-primary" />
                Cadastrar Novo Profissional
              </h3>
              <button
                onClick={() => setIsModalOpen(false)}
                className="text-gray-400 hover:text-gray-600 transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <form onSubmit={handleRegisterStaff} className="flex flex-col gap-4">
              <div className="flex flex-col gap-1">
                <label className="text-xs font-semibold text-gray-600">Nome Completo</label>
                <Input
                  type="text"
                  placeholder="Ex: Dr. André Silva de Araujo"
                  value={newStaffName}
                  onChange={(e) => setNewStaffName(e.target.value)}
                  required
                />
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="flex flex-col gap-1">
                  <label className="text-xs font-semibold text-gray-600">Categoria</label>
                  <select
                    value={newStaffRole}
                    onChange={(e) => setNewStaffRole(e.target.value as StaffRole)}
                    className="w-full bg-white border border-border rounded-lg px-3 py-2.5 text-xs text-gray-800 focus:outline-none focus:border-primary/50 transition-all duration-200"
                  >
                    <option value={StaffRole.Doctor}>{StaffRole.Doctor}</option>
                    <option value={StaffRole.Nurse}>{StaffRole.Nurse}</option>
                    <option value={StaffRole.Receptionist}>{StaffRole.Receptionist}</option>
                    <option value={StaffRole.Admin}>Administrativo</option>
                  </select>
                </div>

                <div className="flex flex-col gap-1">
                  <label className="text-xs font-semibold text-gray-600">Registro Profissional</label>
                  <Input
                    type="text"
                    placeholder="Ex: CRM-SP 12345"
                    value={newStaffLicense}
                    onChange={(e) => setNewStaffLicense(e.target.value)}
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="flex flex-col gap-1">
                  <label className="text-xs font-semibold text-gray-600">E-mail Corporativo</label>
                  <Input
                    type="email"
                    placeholder="Ex: nome@hospital.com"
                    value={newStaffEmail}
                    onChange={(e) => setNewStaffEmail(e.target.value)}
                    required
                  />
                </div>

                <div className="flex flex-col gap-1">
                  <label className="text-xs font-semibold text-gray-600">Departamento</label>
                  <Input
                    type="text"
                    placeholder="Ex: Cardiologia"
                    value={newStaffDept}
                    onChange={(e) => setNewStaffDept(e.target.value)}
                  />
                </div>
              </div>

              <div className="flex gap-3 justify-end border-t border-border pt-4 mt-2">
                <Button
                  type="button"
                  variantType="outline"
                  onClick={() => setIsModalOpen(false)}
                  className="px-4 py-2 text-xs font-bold"
                >
                  Cancelar
                </Button>
                <Button
                  type="submit"
                  disabled={createEmployeeMutation.isPending}
                  variantType="primary"
                  className="px-4 py-2 text-xs font-bold"
                >
                  Salvar Cadastro
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </div>
  )
}
