import { useTranslation } from "react-i18next"

const fallbackLocale = "pt-BR"

export const useLocale = (): string => {
  const { i18n } = useTranslation()
  return i18n.language || fallbackLocale
}
