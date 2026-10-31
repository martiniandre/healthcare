import type { ReactNode } from "react"

interface SectionShellProperties {
  id: string
  title: string
  description: string
  children: ReactNode
}

export const SectionShell = ({ id, title, description, children }: SectionShellProperties) => (
  <section id={id} aria-labelledby={`${id}-title`} className="flex scroll-mt-28 flex-col gap-4">
    <div>
      <h2 id={`${id}-title`} className="text-xl font-bold tracking-tight text-foreground">
        {title}
      </h2>
      <p className="mt-1 text-sm text-muted-foreground">{description}</p>
    </div>
    {children}
  </section>
)

interface ShowcaseProperties {
  title: string
  children: ReactNode
}

export const Showcase = ({ title, children }: ShowcaseProperties) => (
  <div className="flex flex-col gap-3 rounded-xl border bg-card p-4 text-card-foreground">
    <h3 className="text-sm font-semibold text-foreground">{title}</h3>
    {children}
  </div>
)
