import i18n from "i18next"
import { initReactI18next } from "react-i18next"
import LanguageDetector from "i18next-browser-languagedetector"

import ptBRResource from "./locales/pt-BR"
import enUSResource from "./locales/en-US"
import esESResource from "./locales/es-ES"

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources: {
      "pt-BR": { ...ptBRResource, translation: ptBRResource },
      "en-US": { ...enUSResource, translation: enUSResource },
      "es-ES": { ...esESResource, translation: esESResource }
    },
    fallbackLng: "pt-BR",
    interpolation: {
      escapeValue: false
    }
  })

export default i18n
