import { Card, CardContent } from "@/shared/components/ui/card"
import { Progress } from "@/shared/components/ui/progress"

interface ImagingUploadProgressProperties {
  percentage: number
  status: string
}

export const ImagingUploadProgress = ({ percentage, status }: ImagingUploadProgressProperties) => {
  return (
    <Card className="gap-0 border-primary/20 bg-primary/5 py-4 text-left">
      <CardContent className="flex flex-col gap-2.5 px-4">
        <div className="flex items-center justify-between text-xs">
          <span className="font-bold text-primary">{status}</span>
          <span className="font-bold tabular-nums text-muted-foreground">{percentage}%</span>
        </div>
        <Progress value={percentage} aria-label={status} />
      </CardContent>
    </Card>
  )
}
