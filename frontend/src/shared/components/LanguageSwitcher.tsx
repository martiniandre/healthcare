import { useState, useRef, useEffect } from "react"
import { useTranslation } from "react-i18next"
import { Globe, Check, ChevronDown } from "lucide-react"

interface LanguageOption {
  code: string
  label: string
  flag: string
}

const languageOptions: LanguageOption[] = [
  { code: "pt-BR", label: "Português", flag: "🇧🇷" },
  { code: "en-US", label: "English", flag: "🇺🇸" },
  { code: "es-ES", label: "Español", flag: "🇪🇸" }
]

export const LanguageSwitcher = () => {
  const { i18n } = useTranslation()
  const [isOpen, setIsOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  const activeLanguage = languageOptions.find((option) => option.code === i18n.language) || languageOptions[0]

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    document.addEventListener("mousedown", handleClickOutside)
    return () => {
      document.removeEventListener("mousedown", handleClickOutside)
    }
  }, [])

  const handleLanguageChange = (code: string) => {
    i18n.changeLanguage(code)
    setIsOpen(false)
  }

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-2 px-3 py-2 rounded-lg bg-gray-50 border border-border hover:bg-gray-100 hover:border-gray-300 transition-all duration-200 text-xs font-semibold text-gray-700 hover:text-gray-900 select-none shadow-sm cursor-pointer"
      >
        <Globe className="w-3.5 h-3.5 text-gray-400" />
        <span className="text-[14px] leading-none shrink-0">{activeLanguage.flag}</span>
        <span className="hidden sm:inline">{activeLanguage.label}</span>
        <ChevronDown className={`w-3 h-3 text-gray-400 transition-transform duration-200 ${isOpen ? "rotate-180" : ""}`} />
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-1.5 w-40 rounded-xl border border-border/80 bg-white/90 backdrop-blur-md shadow-lg py-1.5 z-50 animate-fade-in flex flex-col">
          {languageOptions.map((option) => {
            const isSelected = option.code === i18n.language

            return (
              <button
                key={option.code}
                onClick={() => handleLanguageChange(option.code)}
                className={`w-full flex items-center justify-between px-3 py-2 text-left text-xs font-medium transition-all duration-150 hover:bg-primary/5 cursor-pointer ${
                  isSelected ? "text-primary bg-primary/4 font-semibold" : "text-gray-600 hover:text-gray-900"
                }`}
              >
                <div className="flex items-center gap-2">
                  <span className="text-sm leading-none">{option.flag}</span>
                  <span>{option.label}</span>
                </div>
                {isSelected && <Check className="w-3.5 h-3.5 text-primary shrink-0" />}
              </button>
            )
          })}
        </div>
      )}
    </div>
  )
}
