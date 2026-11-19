import { useTranslation } from "react-i18next"
import { Globe, ChevronDown } from "lucide-react"
import { Button } from "@/shared/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/shared/components/ui/dropdown-menu"
import { cn } from "@/shared/utils/cn"

interface LanguageOption {
  code: string
  label: string
  shortCode: string
}

interface LanguageSwitcherProperties {
  sidebarLayout?: boolean
  compact?: boolean
}

const languageOptions: LanguageOption[] = [
  { code: "pt-BR", label: "Português", shortCode: "PT" },
  { code: "en-US", label: "English", shortCode: "EN" },
  { code: "es-ES", label: "Español", shortCode: "ES" }
]

export const LanguageSwitcher = ({ sidebarLayout = false, compact = false }: LanguageSwitcherProperties) => {
  const { i18n } = useTranslation()

  const activeLanguage = languageOptions.find((option) => option.code === i18n.language) || languageOptions[0]

  const handleLanguageChange = (code: string) => {
    i18n.changeLanguage(code)
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant={sidebarLayout ? "ghost" : "outline"}
          size="sm"
          aria-label={activeLanguage.label}
          className={cn(
            "group select-none",
            sidebarLayout
              ? cn(
                  "h-11 w-full justify-start gap-3 rounded-md px-3 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-foreground md:h-9",
                  compact && "md:size-10 md:justify-center md:px-0"
                )
              : "text-xs font-semibold"
          )}
        >
          {sidebarLayout ? (
            <>
              <Globe className={cn("size-[18px]", compact && "md:hidden")} />
              {compact && <span className="hidden text-xs font-semibold md:inline" aria-hidden="true">{activeLanguage.shortCode}</span>}
              <span className={cn("flex items-center gap-2 truncate", compact && "md:sr-only")}>
                <span>{activeLanguage.label}</span>
              </span>
              <ChevronDown className={cn("ml-auto size-3.5 transition-transform duration-200 group-data-[state=open]:rotate-180", compact && "md:hidden")} />
            </>
          ) : (
            <>
              <Globe className="h-3.5 w-3.5 text-muted-foreground" />
              <span className="hidden sm:inline">{activeLanguage.label}</span>
              <ChevronDown className="h-3 w-3 transition-transform duration-200 group-data-[state=open]:rotate-180" />
            </>
          )}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align={sidebarLayout ? "start" : "end"} side={sidebarLayout ? "right" : "bottom"} className="w-44">
        <DropdownMenuRadioGroup value={activeLanguage.code} onValueChange={handleLanguageChange}>
          {languageOptions.map((option) => (
            <DropdownMenuRadioItem key={option.code} value={option.code} className="gap-2 text-xs font-medium">
              <span>{option.label}</span>
            </DropdownMenuRadioItem>
          ))}
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
