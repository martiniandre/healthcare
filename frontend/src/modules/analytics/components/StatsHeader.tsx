import { useTranslation } from "react-i18next"
import { PageHeader } from "@/shared/components/ui/PageContainer"

export const StatsHeader = () => {
  const { t } = useTranslation("analytics")

  return <PageHeader title={t("title")} />
}
