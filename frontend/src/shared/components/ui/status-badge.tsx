import * as React from "react"
import { AlertTriangle, CheckCircle2, Info, MinusCircle, XOctagon } from "lucide-react"
import { Badge } from "@/shared/components/ui/badge"
import { cn } from "@/shared/utils/cn"

export type ClinicalStatus = "normal" | "attention" | "critical" | "info" | "neutral"

const statusPresentation: Record<
  ClinicalStatus,
  { icon: React.ComponentType<{ className?: string }>; className: string }
> = {
  normal: { icon: CheckCircle2, className: "border-success/30 bg-success/10 text-success" },
  attention: { icon: AlertTriangle, className: "border-warning/30 bg-warning/10 text-warning" },
  critical: { icon: XOctagon, className: "border-destructive/30 bg-destructive/10 text-destructive" },
  info: { icon: Info, className: "border-primary/30 bg-primary/10 text-primary" },
  neutral: { icon: MinusCircle, className: "border-border bg-muted text-muted-foreground" },
}

interface StatusBadgeProperties extends Omit<React.ComponentProps<typeof Badge>, "variant"> {
  status: ClinicalStatus
}

function StatusBadge({ status, className, children, ...badgeProperties }: StatusBadgeProperties) {
  const { icon: StatusIcon, className: statusClassName } = statusPresentation[status]

  return (
    <Badge
      variant="outline"
      data-status={status}
      className={cn("gap-1", statusClassName, className)}
      {...badgeProperties}
    >
      <StatusIcon aria-hidden="true" />
      {children}
    </Badge>
  )
}

export { StatusBadge }
