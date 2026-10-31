import * as React from "react"
import { Input } from "@/shared/components/ui/input"

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
      if (mask.includes("/")) {
        const truncated = digits.slice(0, 8)
        if (truncated.length <= 2) {
          return truncated
        }
        if (truncated.length <= 4) {
          return `${truncated.slice(0, 2)}/${truncated.slice(2)}`
        }
        return `${truncated.slice(0, 2)}/${truncated.slice(2, 4)}/${truncated.slice(4)}`
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
        <Input
          ref={reference}
          aria-invalid={errorText ? true : undefined}
          className={className}
          onChange={handleInputChange}
          {...elementProperties}
        />
        {errorText && (
          <span className="px-1 text-xs font-medium text-destructive">
            {errorText}
          </span>
        )}
      </div>
    )
  }
)

MaskedInput.displayName = "MaskedInput"
