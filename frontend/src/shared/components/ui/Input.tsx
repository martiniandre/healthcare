import * as React from "react"
import { cn } from "../../utils/cn"

interface ClinicalInputProperties extends React.InputHTMLAttributes<HTMLInputElement> {
  errorText?: string
}

export const Input = React.forwardRef<HTMLInputElement, ClinicalInputProperties>(
  ({ className, errorText, ...elementProperties }, reference) => {
    return (
      <div className="flex flex-col gap-1">
        <input
          ref={reference}
          className={cn(
            "flex w-full h-10 rounded-lg border border-input bg-surface px-3.5 text-sm text-foreground shadow-sm transition-colors placeholder:text-muted-foreground/70 focus-visible:outline-none focus-visible:border-primary focus-visible:ring-2 focus-visible:ring-primary/20 disabled:cursor-not-allowed disabled:opacity-50",
            errorText && "border-danger/40 focus-visible:border-danger focus-visible:ring-danger/20",
            className
          )}
          {...elementProperties}
        />
        {errorText && (
          <span className="text-xs text-danger font-medium px-1">
            {errorText}
          </span>
        )}
      </div>
    )
  }
)

Input.displayName = "Input"
