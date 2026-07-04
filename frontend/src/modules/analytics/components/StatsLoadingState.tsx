import { Card } from "../../../shared/components/ui/Card"
import { Skeleton } from "../../../shared/components/ui/Skeleton"

export const StatsLoadingState = () => {
  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full">
      <div className="text-left">
        <Skeleton className="h-6 w-48" />
        <Skeleton className="h-3 w-72 mt-2.5" />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
        {[1, 2, 3, 4].map((indexValue) => (
          <Card key={indexValue} className="p-4 flex items-center justify-between border border-border h-24">
            <div className="flex-1 flex flex-col gap-2">
              <Skeleton className="h-3 w-20" />
              <Skeleton className="h-6 w-16" />
              <Skeleton className="h-3 w-28" />
            </div>
            <Skeleton className="w-12 h-12 rounded-xl" />
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card className="p-5 flex flex-col gap-5 border border-border h-80">
          <Skeleton className="h-4 w-32" />
          <Skeleton className="h-3 w-48" />
          <div className="flex items-center gap-6 mt-4">
            <Skeleton className="w-36 h-36 rounded-full" />
            <div className="flex-1 flex flex-col gap-3">
              {[1, 2, 3, 4].map((indexValue) => (
                <Skeleton key={indexValue} className="h-8 w-full" />
              ))}
            </div>
          </div>
        </Card>

        <Card className="p-5 flex flex-col gap-5 border border-border h-80">
          <Skeleton className="h-4 w-32" />
          <Skeleton className="h-3 w-48" />
          <div className="flex-1 flex items-end justify-between gap-3 h-40 pb-2">
            {[1, 2, 3, 4, 5, 6, 7].map((indexValue) => (
              <Skeleton key={indexValue} className="flex-1 h-32 rounded-t-md rounded-b-none" />
            ))}
          </div>
        </Card>
      </div>

      <Card className="p-5 flex flex-col gap-4 border border-border">
        <div className="flex justify-between items-center pb-3">
          <div className="flex flex-col gap-2">
            <Skeleton className="h-4 w-40" />
            <Skeleton className="h-3 w-60" />
          </div>
          <Skeleton className="w-24 h-8" />
        </div>
        <div className="flex flex-col gap-4 mt-2">
          {[1, 2, 3, 4].map((indexValue) => (
            <Skeleton key={indexValue} className="h-10 w-full" />
          ))}
        </div>
      </Card>
    </div>
  )
}
