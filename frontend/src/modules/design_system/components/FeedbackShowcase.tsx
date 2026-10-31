import { Inbox } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Alert, AlertDescription, AlertTitle } from "@/shared/components/ui/alert"
import { Button } from "@/shared/components/ui/button"
import { EmptyState } from "@/shared/components/ui/EmptyState"
import { Progress } from "@/shared/components/ui/progress"
import { Skeleton } from "@/shared/components/ui/skeleton"
import { Spinner } from "@/shared/components/ui/spinner"
import { ALERT_VARIANTS } from "../catalog"
import { Showcase } from "./SectionShell"

export const FeedbackShowcase = () => {
  const { t } = useTranslation("designSystem")

  return (
    <>
      <Showcase title={t("components.alert.title")}>
        <div className="flex flex-col gap-3">
          {ALERT_VARIANTS.map((variant) => (
            <Alert key={variant} variant={variant}>
              <AlertTitle>{variant}</AlertTitle>
              <AlertDescription>{t("components.alert.description")}</AlertDescription>
            </Alert>
          ))}
        </div>
      </Showcase>
      <Showcase title={t("components.progress.title")}>
        <Progress value={65} aria-label={t("components.progress.label")} />
      </Showcase>
      <Showcase title={t("components.spinner.title")}>
        <Spinner className="size-6 text-primary" />
      </Showcase>
      <Showcase title={t("components.skeleton.title")}>
        <div className="flex flex-col gap-2">
          <Skeleton className="h-4 w-2/3" />
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-1/2" />
        </div>
      </Showcase>
      <Showcase title={t("components.emptyState.title")}>
        <EmptyState
          icon={Inbox}
          title={t("components.emptyState.heading")}
          description={t("components.emptyState.description")}
          action={<Button variant="outline">{t("components.emptyState.action")}</Button>}
        />
      </Showcase>
    </>
  )
}
