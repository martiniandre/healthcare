import { Card } from "../../../shared/components/ui/Card"

interface OccupancyGaugeProps {
  occupancyRate: number
  totalBeds: number
  occupiedBeds: number
}

export const OccupancyGauge = ({ occupancyRate, totalBeds, occupiedBeds }: OccupancyGaugeProps) => {
  const clampedRate = Math.min(occupancyRate, 100)
  const circumference = 2 * Math.PI * 70
  const arcLength = circumference * 0.75
  const offset = arcLength - (clampedRate / 100) * arcLength

  const getGaugeColor = (rate: number) => {
    if (rate >= 85) return "#dc2626"
    if (rate >= 65) return "#f59e0b"
    return "#10b981"
  }

  const gaugeColor = getGaugeColor(clampedRate)

  return (
    <Card className="p-5 flex flex-col gap-4 text-left border border-border">
      <div>
        <h3 className="font-extrabold text-gray-900 text-md">Taxa de Ocupação</h3>
        <span className="text-xs text-muted block mt-1">Leitos ocupados vs disponíveis</span>
      </div>

      <div className="flex flex-col items-center justify-center py-4">
        <div className="relative w-48 h-28 overflow-hidden">
          <svg className="w-full h-full" viewBox="0 0 200 120">
            <path
              d="M 30 100 A 70 70 0 0 1 170 100"
              fill="none"
              stroke="#f3f4f6"
              strokeWidth="14"
              strokeLinecap="round"
            />
            <path
              d="M 30 100 A 70 70 0 0 1 170 100"
              fill="none"
              stroke={gaugeColor}
              strokeWidth="14"
              strokeLinecap="round"
              strokeDasharray={`${arcLength} ${circumference}`}
              strokeDashoffset={offset}
              className="transition-all duration-700"
            />
          </svg>

          <div className="absolute inset-0 flex flex-col items-center justify-center pt-4">
            <span className="text-3xl font-black text-gray-900">{clampedRate.toFixed(0)}%</span>
            <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider mt-1">Ocupação</span>
          </div>
        </div>

        <div className="flex items-center justify-center gap-6 mt-4 text-xs">
          <div className="text-center">
            <span className="font-black text-gray-900 text-sm block">{occupiedBeds}</span>
            <span className="text-gray-500 font-semibold">Ocupados</span>
          </div>
          <div className="w-px h-8 bg-gray-200" />
          <div className="text-center">
            <span className="font-black text-gray-900 text-sm block">{totalBeds}</span>
            <span className="text-gray-500 font-semibold">Totais</span>
          </div>
          <div className="w-px h-8 bg-gray-200" />
          <div className="text-center">
            <span className="font-black text-gray-900 text-sm block">{totalBeds - occupiedBeds}</span>
            <span className="text-gray-500 font-semibold">Disponíveis</span>
          </div>
        </div>
      </div>
    </Card>
  )
}
