import type { LucideIcon } from "lucide-react";
import { Card, CardContent } from "../../../shared/components/ui/card";
import { StatusBadge, type ClinicalStatus } from "../../../shared/components/ui/status-badge";
import { cn } from "../../../shared/utils/cn";
import { toneClassNames, type MetricTone } from "../metric_tones";

export interface MetricStripMeta {
  text: string;
  tone: string;
  icon?: LucideIcon;
  status?: ClinicalStatus;
}

export interface MetricStripItem {
  label: string;
  value: string;
  unit?: string;
  icon: LucideIcon;
  tone: MetricTone;
  meta?: MetricStripMeta;
}

interface MetricStripProps {
  items: MetricStripItem[];
  className?: string;
  columnsClassName?: string;
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
  );
};

interface MetricStripCardProps {
  item: MetricStripItem;
}

const MetricStripCard = ({ item }: MetricStripCardProps) => {
  const { label, value, unit, tone, meta } = item;
  const toneClassName = toneClassNames[tone];

  return (
    <Card className="relative gap-0 overflow-hidden py-4">
      <span
        className={cn("absolute inset-x-0 bottom-0 h-[3px]", toneClassName.bar)}
        aria-hidden="true"
      />

      <CardContent className="flex flex-col gap-3 px-4">
        <span className="text-xs font-medium leading-tight text-muted-foreground">
          {label}
        </span>

        <div className="flex items-baseline gap-1.5">
          <span className="text-3xl font-semibold leading-none tracking-tight tabular-nums text-foreground">
            {value}
          </span>
          {unit ? (
            <span className="text-xs font-medium text-muted-foreground">
              {unit}
            </span>
          ) : null}
        </div>

        {meta ? (
        <div className="mt-auto flex min-h-5 items-center">
          {meta.status ? (
            <StatusBadge status={meta.status} className="max-w-full text-2xs font-semibold">
              <span className="truncate">{meta.text}</span>
            </StatusBadge>
          ) : (
          <span
            className={cn(
              "flex items-center gap-1 text-2xs font-semibold",
              meta.tone,
            )}
          >
            <span className="truncate">{meta.text}</span>
          </span>
          )}
        </div>
        ) : null}
      </CardContent>
    </Card>
  );
};
