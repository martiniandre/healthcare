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
            "w-full bg-gray-50 border border-border rounded-lg px-3.5 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/20 transition-colors",
            errorText && "border-red-300 focus:border-red-500 focus:ring-red-200",
            className
          )}
          {...elementProperties}
        />
        {errorText && (
          <span className="text-xs text-red-500 font-medium px-1">
            {errorText}
          </span>
        )}
      </div>
    )
  }
)

Input.displayName = "Input"
