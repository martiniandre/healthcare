import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
import { EmptyState } from "../../../shared/components/ui/EmptyState"
import type { DepartmentWaitTime } from "../dashboard_types"

interface WaitTimeChartProps {
  waitTimeByDepartment: DepartmentWaitTime[]
}

export const WaitTimeChart = ({ waitTimeByDepartment }: WaitTimeChartProps) => {
  const { t } = useTranslation("analytics")
  const maxMinutes = waitTimeByDepartment.length > 0
    ? Math.max(...waitTimeByDepartment.map((department) => department.minutes))
    : 1

  return (
    <Card className="p-5 flex flex-col gap-5 text-left border border-border">
      <div>
        <h3 className="font-extrabold text-gray-900 text-md">{t("dashboard.waitTime.title")}</h3>
        <span className="text-xs text-muted block mt-1">{t("dashboard.waitTime.subtitle")}</span>
      </div>

      <div className="overflow-x-auto w-full">
        {waitTimeByDepartment.length > 0 ? (
          <div className="flex flex-col gap-3 min-w-[300px]">
            {waitTimeByDepartment.map((department) => {
              const barWidth = (department.minutes / maxMinutes) * 100
              return (
                <div key={department.department} className="flex items-center gap-3">
                  <div className="w-32 shrink-0 text-left">
                    <span className="text-xs font-bold text-gray-700 block truncate">{department.department}</span>
                  </div>
                  <div className="flex-1 flex items-center gap-2">
                    <div className="flex-1 bg-gray-100 rounded-full h-5 overflow-hidden">
                      <div
                        className="bg-amber-500 rounded-full h-full transition-all duration-500"
                        style={{ width: `${barWidth}%` }}
                      />
                    </div>
                    <span className="text-xs font-black text-gray-900 w-10 text-right">{department.minutes.toFixed(0)} min</span>
                  </div>
                </div>
              )
            })}
          </div>
        ) : (
          <EmptyState
            title={t("dashboard.waitTime.noDataTitle")}
            description={t("dashboard.waitTime.noDataDescription")}
            className="h-40"
          />
        )}
      </div>
    </Card>
  )
}
