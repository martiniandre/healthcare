interface ChartTooltipProps {
  label: string
  detail?: string
  value: string
  unit?: string
  barFill: string
}

export const ChartTooltip = ({ label, detail, value, unit, barFill }: ChartTooltipProps) => {
  return (
    <div className="rounded-lg border border-border bg-popover px-3 py-2 shadow-lg">
      <p className="text-2xs font-semibold text-foreground">{label}</p>
      {detail ? <p className="mt-0.5 text-2xs font-semibold text-muted-foreground">{detail}</p> : null}
      <p className="mt-1.5 flex items-center gap-1.5">
        <span
          className="size-2 shrink-0 rounded-full"
          style={{ backgroundColor: barFill }}
          aria-hidden="true"
        />
        <span className="text-sm font-semibold tabular-nums text-foreground">{value}</span>
        {unit ? <span className="text-2xs font-semibold text-muted-foreground">{unit}</span> : null}
      </p>
    </div>
  )
}
