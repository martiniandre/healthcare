import * as React from "react"
import { cn } from "../../utils/cn"

interface ClinicalButtonProperties extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variantType?: "primary" | "secondary" | "outline" | "danger"
}

export const Button = React.forwardRef<HTMLButtonElement, ClinicalButtonProperties>(
  ({ className, variantType = "primary", children, ...elementProperties }, reference) => {
    return (
      <button
        ref={reference}
        className={cn(
          "px-4 py-2 rounded-lg font-semibold transition-all duration-200 active:scale-[0.97] disabled:opacity-50 cursor-pointer text-sm flex items-center justify-center gap-2",
          variantType === "primary" && "bg-primary text-white hover:bg-primary/90",
          variantType === "secondary" && "bg-secondary text-white hover:bg-secondary/90",
          variantType === "outline" && "border border-border text-gray-600 hover:bg-gray-50 hover:border-gray-300",
          variantType === "danger" && "bg-red-50 border border-red-200 text-red-600 hover:bg-red-100",
          className
        )}
        {...elementProperties}
      >
        {children}
      </button>
    )
  }
)

Button.displayName = "Button"
