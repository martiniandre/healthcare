import { useState } from "react"
import { Card } from "../../shared/components/ui/Card"
import { Button } from "../../shared/components/ui/Button"
import { 
  BarChart3, 
  Users, 
  Clock, 
  CheckSquare, 
  ArrowUpRight, 
  FileSpreadsheet, 
  Activity 
} from "lucide-react"

interface ModalityData {
  modality: string
  percentage: number
  count: number
  color: string
}

interface ConsultationsDayData {
  dayName: string
  count: number
}

export const Stats = () => {
  const [selectedModality, setSelectedModality] = useState<string | null>(null)
  const [hoveredBarIndex, setHoveredBarIndex] = useState<number | null>(null)

  const examModalitiesData: ModalityData[] = [
    { modality: "CT (Tomografia)", percentage: 45, count: 180, color: "#2563eb" },
    { modality: "MR (Ressonância)", percentage: 30, count: 120, color: "#0d9488" },
    { modality: "CR (Raio-X)", percentage: 15, count: 60, color: "#8b5cf6" },
    { modality: "US (Ultrassom)", percentage: 10, count: 40, color: "#f59e0b" }
  ]

  const consultationsWeeklyData: ConsultationsDayData[] = [
    { dayName: "Seg", count: 24 },
    { dayName: "Ter", count: 38 },
    { dayName: "Qua", count: 42 },
    { dayName: "Qui", count: 35 },
    { dayName: "Sex", count: 48 },
    { dayName: "Sáb", count: 18 },
    { dayName: "Dom", count: 8 }
  ]

  const totalRegisteredPatients = 340
  const fhirComplianceRate = 99.4
  const averageServiceDurationMinutes = 14.5
  const activeConsultationsTotal = 56

  return (
    <div className="flex-1 p-8 flex flex-col gap-6 max-w-7xl mx-auto w-full select-none">
      <div className="text-left">
        <h2 className="text-xl font-black text-gray-900 leading-none flex items-center gap-2">
          <BarChart3 className="w-5 h-5 text-primary animate-pulse-glow" />
          Estatísticas Clínicas & Analytics
        </h2>
        <span className="text-xs text-muted mt-1.5 block">
          Visão epidemiológica agregada de prontuários FHIR e laudos radiológicos PACS
        </span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
        <Card className="p-4 flex items-center justify-between border border-border">
          <div className="text-left">
            <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Pacientes Ativos</span>
            <span className="text-2xl font-black text-gray-900 mt-1 block">{totalRegisteredPatients}</span>
            <span className="text-[10px] text-emerald-600 font-bold flex items-center gap-1 mt-1.5">
              <ArrowUpRight className="w-3.5 h-3.5" />
              +12% este mês
            </span>
          </div>
          <div className="bg-primary/8 p-3 rounded-xl">
            <Users className="w-6 h-6 text-primary" />
          </div>
        </Card>

        <Card className="p-4 flex items-center justify-between border border-border">
          <div className="text-left">
            <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Conformidade FHIR</span>
            <span className="text-2xl font-black text-gray-900 mt-1 block">{fhirComplianceRate}%</span>
            <span className="text-[10px] text-emerald-600 font-bold flex items-center gap-1 mt-1.5">
              <CheckSquare className="w-3.5 h-3.5" />
              R4 Core compliant
            </span>
          </div>
          <div className="bg-emerald-50 p-3 rounded-xl">
            <CheckSquare className="w-6 h-6 text-emerald-600" />
          </div>
        </Card>

        <Card className="p-4 flex items-center justify-between border border-border">
          <div className="text-left">
            <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">T. Médio Consulta</span>
            <span className="text-2xl font-black text-gray-900 mt-1 block">{averageServiceDurationMinutes} min</span>
            <span className="text-[10px] text-emerald-600 font-bold flex items-center gap-1 mt-1.5">
              <Clock className="w-3.5 h-3.5" />
              -1.2m vs mês anterior
            </span>
          </div>
          <div className="bg-purple-50 p-3 rounded-xl">
            <Clock className="w-6 h-6 text-purple-600" />
          </div>
        </Card>

        <Card className="p-4 flex items-center justify-between border border-border">
          <div className="text-left">
            <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Atendimentos Semanal</span>
            <span className="text-2xl font-black text-gray-900 mt-1 block">{activeConsultationsTotal}</span>
            <span className="text-[10px] text-amber-600 font-bold flex items-center gap-1 mt-1.5">
              <Activity className="w-3.5 h-3.5" />
              5 leitos de UTI ativos
            </span>
          </div>
          <div className="bg-amber-50 p-3 rounded-xl">
            <Activity className="w-6 h-6 text-amber-500" />
          </div>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card className="p-5 flex flex-col gap-5 text-left border border-border">
          <div>
            <h3 className="font-extrabold text-gray-900 text-md">Distribuição de Exames (PACS)</h3>
            <span className="text-xs text-muted block mt-1">Modalidade dos estudos DICOM integrados</span>
          </div>

          <div className="flex flex-col sm:flex-row items-center justify-around gap-6 py-4">
            <div className="relative w-44 h-44 flex items-center justify-center shrink-0">
              <svg className="w-full h-full transform -rotate-90" viewBox="0 0 100 100">
                <circle cx="50" cy="50" r="38" fill="transparent" stroke="#f3f4f6" strokeWidth="8" />
                
                <circle
                  cx="50"
                  cy="50"
                  r="38"
                  fill="transparent"
                  stroke="#2563eb"
                  strokeWidth="8.5"
                  strokeDasharray="238.76"
                  strokeDashoffset={238.76 - (238.76 * 45) / 100}
                />
                
                <circle
                  cx="50"
                  cy="50"
                  r="38"
                  fill="transparent"
                  stroke="#0d9488"
                  strokeWidth="8.5"
                  strokeDasharray="238.76"
                  strokeDashoffset={238.76 - (238.76 * 30) / 100}
                  className="transform origin-center rotate-[162deg]"
                />

                <circle
                  cx="50"
                  cy="50"
                  r="38"
                  fill="transparent"
                  stroke="#8b5cf6"
                  strokeWidth="8.5"
                  strokeDasharray="238.76"
                  strokeDashoffset={238.76 - (238.76 * 15) / 100}
                  className="transform origin-center rotate-[270deg]"
                />

                <circle
                  cx="50"
                  cy="50"
                  r="38"
                  fill="transparent"
                  stroke="#f59e0b"
                  strokeWidth="8.5"
                  strokeDasharray="238.76"
                  strokeDashoffset={238.76 - (238.76 * 10) / 100}
                  className="transform origin-center rotate-[324deg]"
                />
              </svg>

              <div className="absolute flex flex-col items-center justify-center text-center">
                <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider">Total</span>
                <span className="text-xl font-black text-gray-900">400</span>
                <span className="text-[9px] text-muted font-semibold mt-0.5">Estudos</span>
              </div>
            </div>

            <div className="flex flex-col gap-3 w-full">
              {examModalitiesData.map((item) => (
                <div
                  key={item.modality}
                  onMouseEnter={() => setSelectedModality(item.modality)}
                  onMouseLeave={() => setSelectedModality(null)}
                  className={`flex items-center justify-between p-2.5 rounded-lg border transition-all duration-200 ${
                    selectedModality === item.modality 
                      ? "bg-gray-50 border-gray-300" 
                      : "bg-white border-transparent"
                  }`}
                >
                  <div className="flex items-center gap-2.5">
                    <div className="w-3.5 h-3.5 rounded-full shrink-0" style={{ backgroundColor: item.color }} />
                    <span className="text-xs font-bold text-gray-700">{item.modality}</span>
                  </div>
                  <div className="text-right">
                    <span className="text-xs font-black text-gray-900 block">{item.count} exames</span>
                    <span className="text-[10px] text-gray-500 font-semibold">{item.percentage}%</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </Card>

        <Card className="p-5 flex flex-col gap-5 text-left border border-border">
          <div>
            <h3 className="font-extrabold text-gray-900 text-md">Volume de Atendimentos</h3>
            <span className="text-xs text-muted block mt-1">Evolução diária de atendimentos médicos e triagem</span>
          </div>

          <div className="flex items-end justify-between gap-2.5 h-48 border-b border-border pb-2 pt-6 px-4">
            {consultationsWeeklyData.map((item, index) => {
              const isHovered = hoveredBarIndex === index
              const maxScaleCount = 50
              const percentageHeight = (item.count / maxScaleCount) * 100

              return (
                <div
                  key={item.dayName}
                  className="flex-1 flex flex-col items-center gap-2 group relative"
                  onMouseEnter={() => setHoveredBarIndex(index)}
                  onMouseLeave={() => setHoveredBarIndex(null)}
                >
                  {isHovered && (
                    <div className="absolute -top-10 bg-gray-900 text-white text-[10px] font-bold px-2 py-1 rounded shadow-md z-10 whitespace-nowrap">
                      {item.count} Consultas
                    </div>
                  )}

                  <div
                    className={`w-full max-w-[28px] rounded-t-md transition-all duration-300 ${
                      isHovered ? "bg-primary" : "bg-primary/20"
                    }`}
                    style={{ height: `${percentageHeight}%` }}
                  />

                  <span className="text-[10px] font-bold text-gray-500 uppercase">
                    {item.dayName}
                  </span>
                </div>
              )
            })}
          </div>

          <div className="flex justify-between items-center text-xs text-gray-500 px-2.5">
            <span>Menor: 8 (Dom)</span>
            <span>Média: 30 / Dia</span>
            <span>Pico: 48 (Sex)</span>
          </div>
        </Card>
      </div>

      <Card className="p-5 flex flex-col gap-4 text-left border border-border">
        <div className="flex items-center justify-between border-b border-border pb-3 flex-wrap gap-2">
          <div>
            <h3 className="font-extrabold text-gray-900 text-md">Epidemiologia e Diagnósticos (FHIR Core)</h3>
            <span className="text-xs text-muted block mt-1">Casos clínicos ativos mapeados na base de dados FHIR</span>
          </div>
          <Button variantType="outline" className="px-3 gap-1.5 text-xs">
            <FileSpreadsheet className="w-4 h-4" />
            Exportar CSV
          </Button>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs border-collapse">
            <thead>
              <tr className="border-b border-border text-gray-500 font-bold uppercase tracking-wider">
                <th className="py-3 px-3">Código CID</th>
                <th className="py-3 px-3">Descrição da Patologia</th>
                <th className="py-3 px-3">Categoria FHIR</th>
                <th className="py-3 px-3">Casos Ativos</th>
                <th className="py-3 px-3 text-right">Tendência</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border text-gray-700 font-medium">
              <tr className="hover:bg-gray-50/50">
                <td className="py-3 px-3 font-mono font-bold text-primary">J45.9</td>
                <td className="py-3 px-3">Asma não especificada</td>
                <td className="py-3 px-3">Respiratory</td>
                <td className="py-3 px-3 font-bold text-gray-900">45</td>
                <td className="py-3 px-3 text-right text-emerald-600 font-bold">+5%</td>
              </tr>
              <tr className="hover:bg-gray-50/50">
                <td className="py-3 px-3 font-mono font-bold text-primary">I10</td>
                <td className="py-3 px-3">Hipertensão essencial primária</td>
                <td className="py-3 px-3">Cardiovascular</td>
                <td className="py-3 px-3 font-bold text-gray-900">120</td>
                <td className="py-3 px-3 text-right text-gray-400 font-bold">Estável</td>
              </tr>
              <tr className="hover:bg-gray-50/50">
                <td className="py-3 px-3 font-mono font-bold text-primary">E11.9</td>
                <td className="py-3 px-3">Diabetes mellitus tipo 2</td>
                <td className="py-3 px-3">Endocrine</td>
                <td className="py-3 px-3 font-bold text-gray-900">85</td>
                <td className="py-3 px-3 text-right text-red-500 font-bold">+12%</td>
              </tr>
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  )
}
