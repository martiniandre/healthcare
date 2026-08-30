import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "../../utils/cn"

const badgeVariants = cva(
  "inline-flex items-center rounded-full border border-border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
  {
    variants: {
      variant: {
        default:
          "border-transparent bg-primary text-primary-foreground hover:bg-primary/80",
        secondary:
          "border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80",
        destructive:
          "border-transparent bg-red-50 text-red-600 border-red-200 hover:bg-red-100 dark:bg-red-500/10 dark:text-red-400 dark:border-red-500/30 dark:hover:bg-red-500/20",
        outline: "text-foreground",
        success:
          "border-transparent bg-green-50 text-green-600 border-green-200 hover:bg-green-100 dark:bg-green-500/10 dark:text-green-400 dark:border-green-500/30 dark:hover:bg-green-500/20",
        info:
          "border-transparent bg-blue-50 text-blue-600 border-blue-200 hover:bg-blue-100 dark:bg-blue-500/10 dark:text-blue-400 dark:border-blue-500/30 dark:hover:bg-blue-500/20",
        warning:
          "border-transparent bg-yellow-50 text-yellow-600 border-yellow-200 hover:bg-yellow-100 dark:bg-yellow-500/10 dark:text-yellow-400 dark:border-yellow-500/30 dark:hover:bg-yellow-500/20",
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
