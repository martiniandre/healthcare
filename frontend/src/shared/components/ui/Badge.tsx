import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "../../utils/cn"

const badgeVariants = cva(
  "inline-flex items-center rounded-full border border-border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
  {
    variants: {
      variant: {
        default:
          "border-transparent bg-primary text-primary-foreground hover:bg-primary/90",
        secondary:
          "border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/90",
        destructive:
          "border-transparent bg-danger-soft text-danger border-danger/15 hover:bg-danger/10",
        outline: "text-foreground",
        success:
          "border-transparent bg-success-soft text-success border-success/15 hover:bg-success/10",
        info:
          "border-transparent bg-info-soft text-info border-info/15 hover:bg-info/10",
        warning:
          "border-transparent bg-warning-soft text-warning border-warning/15 hover:bg-warning/10",
        muted:
          "border-transparent bg-muted-soft text-muted-foreground hover:bg-muted-soft/70",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
)

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <div className={cn(badgeVariants({ variant }), className)} {...props} />
  )
}

export { Badge }
