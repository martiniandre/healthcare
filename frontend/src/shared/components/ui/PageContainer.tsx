import type { ReactNode } from "react"
import { cn } from "../../utils/cn"

interface PageContainerProps {
  className?: string
  children: ReactNode
}

export const PageContainer = ({ className, children }: PageContainerProps) => {
  return (
    <div
      className={cn(
        "flex-1 px-4 pb-10 pt-5 sm:px-6 md:px-8 md:pt-6 flex flex-col gap-5 md:gap-7 max-w-7xl mx-auto w-full animate-fade-in",
        className
      )}
    >
      {children}
    </div>
  )
}

interface PageHeaderProps {
  eyebrow?: string
  title: string
  description?: string
  actions?: ReactNode
  className?: string
}

export const PageHeader = ({ eyebrow, title, description, actions, className }: PageHeaderProps) => {
  return (
    <header className={cn("flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between", className)}>
      <div className="flex min-w-0 items-start gap-4">
        <div className="min-w-0">
          {eyebrow && (
            <p className="mb-1 text-xs font-medium text-muted-foreground">{eyebrow}</p>
          )}
          <h1 className="font-display text-2xl font-semibold leading-tight text-foreground">
            {title}
          </h1>
          {description && <p className="mt-1 max-w-2xl text-sm text-muted-foreground">{description}</p>}
        </div>
      </div>
      {actions && <div className="flex flex-wrap items-center gap-2 sm:shrink-0">{actions}</div>}
    </header>
  )
}
