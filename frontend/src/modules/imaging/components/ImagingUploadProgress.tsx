import { Card } from "../../../shared/components/ui/Card"

interface ImagingUploadProgressProperties {
  percentage: number
  status: string
}

export const ImagingUploadProgress = ({ percentage, status }: ImagingUploadProgressProperties) => {
  return (
    <Card className="p-4 bg-primary/5 border border-primary/20 flex flex-col gap-2.5 text-left">
      <div className="flex justify-between items-center text-xs">
        <span className="text-primary font-bold">{status}</span>
        <span className="text-gray-500 font-bold">{percentage}%</span>
      </div>
      <div className="w-full bg-gray-100 rounded-full h-2">
        <div
          className="bg-primary h-2 rounded-full transition-all duration-300"
          style={{ width: `${percentage}%` }}
        />
      </div>
    </Card>
  )
}
