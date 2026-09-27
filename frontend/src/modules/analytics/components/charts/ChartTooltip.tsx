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
      <p className="text-[11px] font-bold text-gray-900">{label}</p>
      {detail ? <p className="mt-0.5 text-[10px] font-semibold text-muted-foreground">{detail}</p> : null}
      <p className="mt-1.5 flex items-center gap-1.5">
        <span
          className="size-2 shrink-0 rounded-full"
          style={{ backgroundColor: barFill }}
          aria-hidden="true"
        />
        <span className="text-sm font-black tabular-nums text-gray-900">{value}</span>
        {unit ? <span className="text-[10px] font-bold text-muted-foreground">{unit}</span> : null}
      </p>
    </div>
  )
}
