import type { ReactNode } from "react"
import { useTranslation } from "react-i18next"
import { AlertCircle } from "lucide-react"
import { usePermissionsQuery } from "../../modules/auth/queries"
import { Alert, AlertDescription, AlertTitle } from "../components/ui/alert"
import { Button } from "../components/ui/button"
import { Skeleton } from "../components/ui/skeleton"
import { AbilityProvider } from "./AbilityContext"

interface PermissionsGateProps {
  children: ReactNode
}

export function PermissionsGate({ children }: PermissionsGateProps) {
  const { t } = useTranslation("common")
  const permissionsQuery = usePermissionsQuery()

  if (permissionsQuery.data) {
    return <AbilityProvider capabilities={permissionsQuery.data.capabilities}>{children}</AbilityProvider>
  }

  if (permissionsQuery.isError) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center p-6">
        <Alert variant="destructive" className="max-w-md">
          <AlertCircle />
          <AlertTitle>{t("permissions.errorTitle")}</AlertTitle>
          <AlertDescription>
            <p>{t("permissions.errorDescription")}</p>
            <Button
              variant="outline"
              size="sm"
              onClick={() => permissionsQuery.refetch()}
              disabled={permissionsQuery.isFetching}
            >
              {t("retry")}
            </Button>
          </AlertDescription>
        </Alert>
      </div>
    )
  }

  return (
    <div
      role="status"
      aria-busy="true"
      aria-label={t("permissions.loadingMessage")}
      className="min-h-screen bg-background flex"
    >
      <Skeleton className="hidden md:block w-[240px] shrink-0 h-screen rounded-none" />
      <div className="flex-1 flex flex-col gap-4 p-6">
        <Skeleton className="h-10 w-64" />
        <Skeleton className="h-40 w-full" />
        <Skeleton className="h-64 w-full" />
      </div>
    </div>
  )
}
