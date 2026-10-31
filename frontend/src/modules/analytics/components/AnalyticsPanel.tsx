import type { ReactNode } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../../shared/components/ui/card"
import { cn } from "../../../shared/utils/cn"

interface AnalyticsPanelProps {
  title: string
  subtitle?: string
  trailing?: ReactNode
  bodyClassName?: string
  className?: string
  children: ReactNode
}

export const AnalyticsPanel = ({
  title,
  subtitle,
  trailing,
  bodyClassName,
  className,
  children,
}: AnalyticsPanelProps) => {
  return (
    <Card className={cn("gap-0 overflow-hidden py-0", className)}>
      <CardHeader className="flex flex-row items-start justify-between gap-4 border-b border-border px-5 py-4 [.border-b]:pb-4">
        <div className="min-w-0">
          <CardTitle className="text-sm font-extrabold tracking-tight">{title}</CardTitle>
          {subtitle ? <CardDescription className="mt-1 text-xs">{subtitle}</CardDescription> : null}
        </div>
        {trailing ? <div className="shrink-0 pt-0.5">{trailing}</div> : null}
      </CardHeader>
      <CardContent className={cn("flex flex-1 flex-col px-5 py-5", bodyClassName)}>{children}</CardContent>
    </Card>
  )
}

interface PanelMetricProps {
  label: string
  value: string
  tone?: "neutral" | "primary" | "secondary"
  className?: string
}

export const PanelMetric = ({
  label,
  value,
  tone = "neutral",
  className,
}: PanelMetricProps) => {
  const toneTextClassName =
    tone === "primary" ? "text-primary" : tone === "secondary" ? "text-secondary" : "text-foreground"

  return (
    <div
      className={cn(
        "flex flex-col items-end rounded-lg border border-border bg-muted px-2.5 py-1.5",
        className
      )}
    >
      <span className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
        {label}
      </span>
      <span className={cn("text-xs font-black tabular-nums", toneTextClassName)}>{value}</span>
    </div>
  )
}
