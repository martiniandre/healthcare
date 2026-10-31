import { useTranslation } from "react-i18next"
import { Users, Activity, Database } from "lucide-react"
import { Card, CardContent } from "../../../shared/components/ui/card"

interface PatientsMetricsGridProps {
  totalPatients: number
}

export const PatientsMetricsGrid = ({ totalPatients }: PatientsMetricsGridProps) => {
  const { t } = useTranslation("patients")

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      <Card className="py-4 gap-0">
        <CardContent className="flex items-center gap-4 px-4">
          <div className="w-10 h-10 rounded-lg bg-primary/8 flex items-center justify-center shrink-0">
            <Users className="w-5 h-5 text-primary" />
          </div>
          <div>
            <span className="text-[10px] text-muted-foreground font-semibold uppercase tracking-wider block">
              {t("cards.total")}
            </span>
            <span className="text-2xl font-black text-foreground leading-none mt-0.5 block">{totalPatients}</span>
          </div>
        </CardContent>
      </Card>

      <Card className="py-4 gap-0">
        <CardContent className="flex items-center gap-4 px-4">
          <div className="w-10 h-10 rounded-lg bg-secondary/8 flex items-center justify-center shrink-0">
            <Activity className="w-5 h-5 text-secondary" />
          </div>
          <div>
            <span className="text-[10px] text-muted-foreground font-semibold uppercase tracking-wider block">
              {t("cards.standard")}
            </span>
            <span className="text-sm font-bold text-foreground mt-0.5 block">{t("cards.fhirCompliant")}</span>
          </div>
        </CardContent>
      </Card>

      <Card className="py-4 gap-0">
        <CardContent className="flex items-center gap-4 px-4">
          <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center shrink-0">
            <Database className="w-5 h-5 text-primary" />
          </div>
          <div>
            <span className="text-[10px] text-muted-foreground font-semibold uppercase tracking-wider block">
              {t("cards.integration")}
            </span>
            <span className="text-sm font-bold text-foreground mt-0.5 block">{t("cards.cloudIntegration")}</span>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
