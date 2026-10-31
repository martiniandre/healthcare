import { Palette } from "lucide-react"
import { useTranslation } from "react-i18next"
import { PageContainer } from "@/shared/components/ui/PageContainer"
import { SECTION_IDS } from "./catalog"
import { ColorsSection } from "./components/ColorsSection"
import { ComponentsSection } from "./components/ComponentsSection"
import { GuidelinesSection } from "./components/GuidelinesSection"
import { TypographySection } from "./components/TypographySection"

export const DesignSystem = () => {
  const { t } = useTranslation("designSystem")

  return (
    <PageContainer className="gap-6">
      <header>
        <h1 className="flex items-center gap-2 text-2xl font-bold tracking-tight text-foreground">
          <Palette className="size-6 text-primary" aria-hidden="true" />
          {t("title")}
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">{t("description")}</p>
      </header>
      <div className="flex flex-col gap-6 lg:flex-row lg:items-start">
        <nav
          aria-label={t("nav.label")}
          className="sticky top-0 z-10 -mx-4 border-b bg-background/95 px-4 backdrop-blur sm:-mx-6 sm:px-6 lg:top-6 lg:mx-0 lg:w-48 lg:shrink-0 lg:border-b-0 lg:px-0"
        >
          <ul className="flex gap-1 overflow-x-auto py-2 lg:flex-col lg:overflow-visible">
            {SECTION_IDS.map((sectionId) => (
              <li key={sectionId} className="shrink-0">
                <a
                  href={`#${sectionId}`}
                  className="flex min-h-11 items-center rounded-md px-3 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground focus-visible:outline-2 focus-visible:outline-ring lg:min-h-9"
                >
                  {t(`nav.${sectionId}`)}
                </a>
              </li>
            ))}
          </ul>
        </nav>
        <div className="flex min-w-0 flex-1 flex-col gap-10">
          <ColorsSection />
          <TypographySection />
          <ComponentsSection />
          <GuidelinesSection />
        </div>
      </div>
    </PageContainer>
  )
}
