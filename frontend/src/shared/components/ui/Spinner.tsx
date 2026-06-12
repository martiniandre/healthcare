import * as React from "react"
import { Loader2 } from "lucide-react"
import { cn } from "../../utils/cn"

export type SpinnerProps = React.SVGProps<SVGSVGElement>

const Spinner = React.forwardRef<SVGSVGElement, SpinnerProps>(
  ({ className, ...props }, ref) => {
    return (
      <Loader2
        ref={ref}
        className={cn("animate-spin text-muted", className)}
        {...props}
      />
    )
  }
)
Spinner.displayName = "Spinner"

export { Spinner }
