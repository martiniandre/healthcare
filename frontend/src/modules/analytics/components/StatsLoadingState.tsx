import { Card, CardContent } from "../../../shared/components/ui/card";
import { Skeleton } from "../../../shared/components/ui/skeleton";
import { PageContainer } from "../../../shared/components/ui/PageContainer";

export const StatsLoadingState = () => {
  return (
    <PageContainer>
      <div className="text-left">
        <Skeleton className="h-6 w-48" />
        <Skeleton className="mt-2.5 h-3 w-72" />
      </div>

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:gap-4 xl:grid-cols-4">
        {Array.from({ length: 4 }).map((_, indexValue) => (
          <Card key={String(indexValue)} className="gap-0 py-4">
            <CardContent className="flex flex-col gap-3 px-4">
              <div className="flex items-start justify-between gap-2">
                <Skeleton className="h-3 w-24" />
                <Skeleton className="size-8 rounded-lg" />
              </div>
              <Skeleton className="h-8 w-16" />
              <Skeleton className="h-4 w-4/5" />
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-1 gap-4 md:gap-6 xl:grid-cols-12">
        <Card className="gap-0 py-5 xl:col-span-5">
          <CardContent className="flex flex-col gap-5 px-5">
            <div className="flex items-start justify-between gap-4">
              <div className="flex flex-col gap-2">
                <Skeleton className="h-4 w-40" />
                <Skeleton className="h-3 w-56" />
              </div>
              <Skeleton className="h-11 w-20 rounded-lg" />
            </div>
            <div className="flex flex-col items-center gap-6 sm:flex-row">
              <Skeleton className="aspect-square w-full max-w-[200px] shrink-0 rounded-full" />
              <div className="flex w-full flex-col gap-2">
                {Array.from({ length: 4 }).map((_, indexValue) => (
                  <Skeleton
                    key={String(indexValue)}
                    className="h-10 w-full rounded-lg"
                  />
                ))}
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="gap-0 py-5 xl:col-span-7">
          <CardContent className="flex flex-col gap-5 px-5">
            <div className="flex flex-wrap items-start justify-between gap-4">
              <div className="flex flex-col gap-2">
                <Skeleton className="h-4 w-40" />
                <Skeleton className="h-3 w-56" />
              </div>
              <div className="flex gap-1.5">
                {Array.from({ length: 3 }).map((_, indexValue) => (
                  <Skeleton
                    key={String(indexValue)}
                    className="h-11 w-16 rounded-lg"
                  />
                ))}
              </div>
            </div>
            <Skeleton className="h-[260px] w-full rounded-lg" />
          </CardContent>
        </Card>
      </div>

      <Card className="gap-0 py-5">
        <CardContent className="flex flex-col gap-4 px-5">
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div className="flex flex-col gap-2">
              <Skeleton className="h-4 w-44" />
              <Skeleton className="h-3 w-60" />
            </div>
            <div className="flex gap-2">
              <Skeleton className="h-11 w-24 rounded-lg" />
              <Skeleton className="h-11 w-28 rounded-lg" />
            </div>
          </div>
          <div className="flex flex-col gap-3">
            {Array.from({ length: 4 }).map((_, indexValue) => (
              <Skeleton
                key={String(indexValue)}
                className="h-12 w-full rounded-lg"
              />
            ))}
          </div>
        </CardContent>
      </Card>
    </PageContainer>
  );
};
