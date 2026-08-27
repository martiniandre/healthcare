import * as React from "react"
import { cn } from "../../utils/cn"

const Textarea = React.forwardRef<
  HTMLTextAreaElement,
  React.TextareaHTMLAttributes<HTMLTextAreaElement>
>(({ className, ...properties }, reference) => (
  <textarea
    ref={reference}
    className={cn(
      "flex min-h-[80px] w-full rounded-lg border border-input bg-surface px-3.5 py-2.5 text-sm text-foreground shadow-sm transition-colors placeholder:text-muted-foreground/70 focus-visible:border-primary focus-visible:ring-2 focus-visible:ring-primary/20 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50 resize-y",
      className
    )}
    {...properties}
  />
))
Textarea.displayName = "Textarea"

export { Textarea }
