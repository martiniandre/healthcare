import { useTranslation } from "react-i18next"
import { Button } from "@/shared/components/ui/button"
import { Card, CardContent } from "@/shared/components/ui/card"
import { Progress } from "@/shared/components/ui/progress"

interface ImagingUploadProgressProperties {
  percentage: number
  status: string
  fileName?: string
  fileSize?: string
  onCancel?: () => void
}

export const ImagingUploadProgress = ({
  percentage,
  status,
  fileName,
  fileSize,
  onCancel,
}: ImagingUploadProgressProperties) => {
  const { t } = useTranslation("imaging")

  return (
    <Card className="gap-0 border-primary/20 bg-primary/5 py-4 text-left">
      <CardContent className="flex flex-col gap-2.5 px-4">
        {fileName && (
          <p className="truncate text-xs text-muted-foreground">
            <span className="font-semibold text-foreground">{fileName}</span>
            {fileSize ? <span className="tabular-nums"> ({fileSize})</span> : null}
          </p>
        )}
        <div className="flex items-center justify-between gap-3 text-xs">
          <span role="status" className="font-bold text-primary">
            {status}
          </span>
          <span className="font-bold tabular-nums text-muted-foreground">{percentage}%</span>
        </div>
        <Progress value={percentage} aria-label={status} aria-valuetext={`${percentage}%`} />
        {onCancel && (
          <div className="flex justify-end">
            <Button type="button" variant="outline" size="sm" onClick={onCancel} className="max-sm:min-h-11 max-sm:w-full">
              {t("upload.cancel")}
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
