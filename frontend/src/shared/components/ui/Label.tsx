import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import * as LabelPrimitive from "@radix-ui/react-label"
import { cn } from "../../utils/cn"

const labelVariants = cva(
  "text-[13px] font-semibold leading-none text-foreground mb-1.5 block peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
)

const Label = React.forwardRef<
  React.ElementRef<typeof LabelPrimitive.Root>,
  React.ComponentPropsWithoutRef<typeof LabelPrimitive.Root> &
    VariantProps<typeof labelVariants>
>(({ className, ...properties }, reference) => (
  <LabelPrimitive.Root
    ref={reference}
    className={cn(labelVariants(), className)}
    {...properties}
  />
))
Label.displayName = LabelPrimitive.Root.displayName

export { Label }
