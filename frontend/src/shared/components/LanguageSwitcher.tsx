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
  flag: string
}

interface LanguageSwitcherProperties {
  sidebarLayout?: boolean
}

const languageOptions: LanguageOption[] = [
  { code: "pt-BR", label: "Português", flag: "🇧🇷" },
  { code: "en-US", label: "English", flag: "🇺🇸" },
  { code: "es-ES", label: "Español", flag: "🇪🇸" }
]

export const LanguageSwitcher = ({ sidebarLayout = false }: LanguageSwitcherProperties) => {
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
            sidebarLayout ? "w-full justify-start gap-3 px-3 py-2 text-[13px] font-medium text-muted-foreground" : "text-xs font-semibold"
          )}
        >
          {sidebarLayout ? (
            <>
              <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-muted text-muted-foreground">
                <Globe className="h-4 w-4" />
              </span>
              <span className="flex items-center gap-2 truncate">
                <span className="text-[14px] leading-none shrink-0">{activeLanguage.flag}</span>
                <span>{activeLanguage.label}</span>
              </span>
              <ChevronDown className="ml-auto h-3 w-3 transition-transform duration-200 group-data-[state=open]:rotate-180" />
            </>
          ) : (
            <>
              <Globe className="h-3.5 w-3.5 text-muted-foreground" />
              <span className="text-[14px] leading-none shrink-0">{activeLanguage.flag}</span>
              <span className="hidden sm:inline">{activeLanguage.label}</span>
              <ChevronDown className="h-3 w-3 transition-transform duration-200 group-data-[state=open]:rotate-180" />
            </>
          )}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align={sidebarLayout ? "start" : "end"} className="w-44">
        <DropdownMenuRadioGroup value={activeLanguage.code} onValueChange={handleLanguageChange}>
          {languageOptions.map((option) => (
            <DropdownMenuRadioItem key={option.code} value={option.code} className="gap-2 text-xs font-medium">
              <span className="text-sm leading-none">{option.flag}</span>
              <span>{option.label}</span>
            </DropdownMenuRadioItem>
          ))}
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
