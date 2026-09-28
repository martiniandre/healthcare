import type { ReactNode } from "react"
import { Card } from "../../../shared/components/ui/Card"
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
    <Card className={cn("flex flex-col overflow-hidden p-0", className)}>
      <header className="flex items-start justify-between gap-4 border-b border-border px-5 py-4">
        <div className="min-w-0">
          <h3 className="text-sm font-extrabold tracking-tight text-gray-900">{title}</h3>
          {subtitle ? <p className="mt-1 text-xs text-muted">{subtitle}</p> : null}
        </div>
        {trailing ? <div className="shrink-0 pt-0.5">{trailing}</div> : null}
      </header>
      <div className={cn("flex flex-1 flex-col px-5 py-5", bodyClassName)}>{children}</div>
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
    tone === "primary" ? "text-primary" : tone === "secondary" ? "text-secondary" : "text-gray-900"

  return (
    <div
      className={cn(
        "flex flex-col items-end rounded-lg border border-border bg-gray-50 px-2.5 py-1.5",
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
