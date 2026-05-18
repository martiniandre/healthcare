import * as React from "react"
import { cn } from "../../utils/cn"

interface PanelContainerProperties extends React.HTMLAttributes<HTMLDivElement> {
  glowingType?: "cyan" | "amethyst" | "none"
}

export const Card = React.forwardRef<HTMLDivElement, PanelContainerProperties>(
  ({ className, glowingType = "none", children, ...elementProperties }, reference) => {
    return (
      <div
        ref={reference}
        className={cn(
          "clinical-glass rounded-xl p-6 transition-all duration-300",
          glowingType === "cyan" && "glow-cyan border-primary/20",
          glowingType === "amethyst" && "glow-amethyst border-secondary/20",
          className
        )}
        {...elementProperties}
      >
        {children}
      </div>
    )
  }
)

Card.displayName = "Card"
