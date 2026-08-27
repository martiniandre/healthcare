import * as React from "react"
import { cn } from "../../utils/cn"

interface MaskedInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  mask: string
  errorText?: string
}

export const MaskedInput = React.forwardRef<HTMLInputElement, MaskedInputProps>(
  ({ className, mask, errorText, onChange, ...elementProperties }, reference) => {
    const formatValue = (rawValue: string): string => {
      const digits = rawValue.replace(/\D/g, "")
      if (mask.includes(".") || mask === "cpf") {
        const truncated = digits.slice(0, 11)
        if (truncated.length <= 3) {
          return truncated
        }
        if (truncated.length <= 6) {
          return `${truncated.slice(0, 3)}.${truncated.slice(3)}`
        }
        if (truncated.length <= 9) {
          return `${truncated.slice(0, 3)}.${truncated.slice(3, 6)}.${truncated.slice(6)}`
        }
        return `${truncated.slice(0, 3)}.${truncated.slice(3, 6)}.${truncated.slice(6, 9)}-${truncated.slice(9)}`
      }
      if (mask.includes("(") || mask === "phone") {
        const truncated = digits.slice(0, 11)
        if (truncated.length <= 2) {
          return truncated.length > 0 ? `(${truncated}` : ""
        }
        if (truncated.length <= 6) {
          return `(${truncated.slice(0, 2)}) ${truncated.slice(2)}`
        }
        if (truncated.length <= 10) {
          return `(${truncated.slice(0, 2)}) ${truncated.slice(2, 6)}-${truncated.slice(6)}`
        }
        return `(${truncated.slice(0, 2)}) ${truncated.slice(2, 7)}-${truncated.slice(7)}`
      }
      if (mask.includes("-") || mask === "date") {
        const truncated = digits.slice(0, 8)
        if (truncated.length <= 4) {
          return truncated
        }
        if (truncated.length <= 6) {
          return `${truncated.slice(0, 4)}-${truncated.slice(4)}`
        }
        return `${truncated.slice(0, 4)}-${truncated.slice(4, 6)}-${truncated.slice(6)}`
      }
      return rawValue
    }

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
      const formatted = formatValue(event.target.value)
      event.target.value = formatted
      if (onChange) {
        onChange(event)
      }
    }

    return (
      <div className="flex flex-col gap-1">
        <input
          ref={reference}
          className={cn(
            "flex w-full h-10 rounded-lg border border-input bg-surface px-3.5 text-sm text-foreground shadow-sm transition-colors placeholder:text-muted-foreground/70 focus-visible:outline-none focus-visible:border-primary focus-visible:ring-2 focus-visible:ring-primary/20 disabled:cursor-not-allowed disabled:opacity-50",
            errorText && "border-danger/40 focus-visible:border-danger focus-visible:ring-danger/20",
            className
          )}
          onChange={handleInputChange}
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

MaskedInput.displayName = "MaskedInput"
