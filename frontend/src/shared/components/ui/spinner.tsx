import { useTranslation } from "react-i18next"
import { cn } from "@/shared/utils/cn"
import { Loader2Icon } from "lucide-react"

function Spinner({ className, ...props }: React.ComponentProps<"svg">) {
  const { t } = useTranslation("common")

  return (
    <Loader2Icon
      role="status"
      aria-label={t("loading")}
      className={cn("size-4 animate-spin", className)}
      {...props}
    />
  )
}

export { Spinner }
