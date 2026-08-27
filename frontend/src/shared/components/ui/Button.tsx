import * as React from "react"
import { Loader2 } from "lucide-react"
import { cn } from "../../utils/cn"

interface ClinicalButtonProperties extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variantType?: "primary" | "secondary" | "outline" | "danger" | "ghost"
  size?: "sm" | "md" | "lg"
  isLoading?: boolean
}

export const Button = React.forwardRef<HTMLButtonElement, ClinicalButtonProperties>(
  ({ className, variantType = "primary", size = "md", isLoading = false, children, disabled, ...elementProperties }, reference) => {
    return (
      <button
        ref={reference}
        disabled={disabled || isLoading}
        className={cn(
          "inline-flex items-center justify-center gap-2 rounded-lg font-semibold transition-all duration-200 active:scale-[0.97] disabled:opacity-50 disabled:pointer-events-none cursor-pointer text-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring shadow-sm",
          size === "sm" && "h-8 px-3 text-xs",
          size === "md" && "h-10 px-4",
          size === "lg" && "h-11 px-6",
          variantType === "primary" && "bg-primary text-primary-foreground hover:bg-primary-hover hover:shadow-card",
          variantType === "secondary" && "bg-secondary text-secondary-foreground hover:bg-secondary-hover",
          variantType === "outline" && "border border-border-strong bg-surface text-foreground hover:bg-muted-soft hover:border-border-strong",
          variantType === "danger" && "bg-danger-soft border border-danger/20 text-danger hover:bg-danger/10",
          variantType === "ghost" && "text-muted-foreground shadow-none hover:bg-muted-soft hover:text-foreground",
          className
        )}
        aria-busy={isLoading}
        {...elementProperties}
      >
        {isLoading && <Loader2 className="w-4 h-4 animate-spin" />}
        {children}
      </button>
    )
  }
)

Button.displayName = "Button"
