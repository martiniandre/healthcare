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
        "flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full animate-fade-in",
        className
      )}
    >
      {children}
    </div>
  )
}

interface PageTitleProps {
  icon?: ReactNode
  title: string
  description?: string
  actions?: ReactNode
}

export const PageTitle = ({ icon, title, description, actions }: PageTitleProps) => {
  return (
    <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 className="text-xl font-display font-bold text-foreground tracking-tight leading-none flex items-center gap-2">
          {icon}
          {title}
        </h2>
        {description && (
          <span className="text-xs text-muted mt-1.5 block">{description}</span>
        )}
      </div>
      {actions && <div className="flex items-center gap-2 self-start sm:self-auto">{actions}</div>}
    </div>
  )
}
