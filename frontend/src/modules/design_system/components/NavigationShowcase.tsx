import { ChevronDown } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from "@/shared/components/ui/breadcrumb"
import { Button } from "@/shared/components/ui/button"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/shared/components/ui/dropdown-menu"
import { Pagination, PaginationContent, PaginationItem, PaginationLink, PaginationNext, PaginationPrevious } from "@/shared/components/ui/pagination"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/shared/components/ui/tabs"
import { PAGINATION_PAGES } from "../catalog"
import { Showcase } from "./SectionShell"

export const NavigationShowcase = () => {
  const { t } = useTranslation("designSystem")

  return (
    <>
      <Showcase title={t("components.tabs.title")}>
        <Tabs defaultValue="summary">
          <TabsList>
            <TabsTrigger value="summary">{t("components.tabs.summary")}</TabsTrigger>
            <TabsTrigger value="history">{t("components.tabs.history")}</TabsTrigger>
          </TabsList>
          <TabsContent value="summary" className="text-sm text-muted-foreground">
            {t("components.tabs.summaryContent")}
          </TabsContent>
          <TabsContent value="history" className="text-sm text-muted-foreground">
            {t("components.tabs.historyContent")}
          </TabsContent>
        </Tabs>
      </Showcase>
      <Showcase title={t("components.breadcrumb.title")}>
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink href="#components">{t("components.breadcrumb.home")}</BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbLink href="#components">{t("components.breadcrumb.patients")}</BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbPage>{t("components.breadcrumb.detail")}</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </Showcase>
      <Showcase title={t("components.pagination.title")}>
        <Pagination>
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious href="#components" />
            </PaginationItem>
            {PAGINATION_PAGES.map((page) => (
              <PaginationItem key={page}>
                <PaginationLink href="#components" isActive={page === 2}>
                  {page}
                </PaginationLink>
              </PaginationItem>
            ))}
            <PaginationItem>
              <PaginationNext href="#components" />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      </Showcase>
      <Showcase title={t("components.dropdown.title")}>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline">
              {t("components.dropdown.trigger")}
              <ChevronDown aria-hidden="true" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start">
            <DropdownMenuLabel>{t("components.dropdown.label")}</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>{t("components.dropdown.edit")}</DropdownMenuItem>
            <DropdownMenuItem>{t("components.dropdown.duplicate")}</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </Showcase>
    </>
  )
}
