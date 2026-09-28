import type { LucideIcon } from "lucide-react"
import { Card } from "../../../shared/components/ui/Card"
import { cn } from "../../../shared/utils/cn"
import { toneClassNames, type MetricTone } from "../metric_tones"

export interface MetricStripMeta {
  text: string
  tone: string
  icon?: LucideIcon
}

export interface MetricStripItem {
  label: string
  value: string
  unit?: string
  icon: LucideIcon
  tone: MetricTone
  meta: MetricStripMeta
}

interface MetricStripProps {
  items: MetricStripItem[]
  className?: string
  columnsClassName?: string
}

export const MetricStrip = ({
  items,
  className,
  columnsClassName = "grid-cols-2 sm:grid-cols-3 lg:grid-cols-6",
}: MetricStripProps) => {
  return (
    <div className={cn("grid gap-3 lg:gap-4", columnsClassName, className)}>
      {items.map((item) => (
        <MetricStripCard key={item.label} item={item} />
      ))}
    </div>
  )
}

interface MetricStripCardProps {
  item: MetricStripItem
}

const MetricStripCard = ({ item }: MetricStripCardProps) => {
  const { label, value, unit, icon: Icon, tone, meta } = item
  const toneClassName = toneClassNames[tone]
  const MetaIcon = meta.icon

  return (
    <Card className="relative flex flex-col gap-3 overflow-hidden p-4">
      <span
        className={cn("absolute inset-x-0 bottom-0 h-[3px]", toneClassName.bar)}
        aria-hidden="true"
      />

      <div className="flex items-start justify-between gap-2">
        <span className="text-[11px] font-bold uppercase leading-tight tracking-wider text-muted-foreground">
          {label}
        </span>
        <span
          className={cn("grid size-8 shrink-0 place-items-center rounded-lg", toneClassName.surface)}
        >
          <Icon className={cn("size-4", toneClassName.text)} aria-hidden="true" />
        </span>
      </div>

      <div className="flex items-baseline gap-1.5">
        <span className="text-3xl font-black leading-none tracking-tight tabular-nums text-gray-900">
          {value}
        </span>
        {unit ? (
          <span className="text-xs font-bold text-muted-foreground">{unit}</span>
        ) : null}
      </div>

      <div className="mt-auto flex h-4 items-center">
        <span className={cn("flex items-center gap-1 text-[11px] font-bold", meta.tone)}>
          {MetaIcon ? <MetaIcon className="size-3 shrink-0" aria-hidden="true" /> : null}
          <span className="truncate">{meta.text}</span>
        </span>
      </div>
    </Card>
  )
}
