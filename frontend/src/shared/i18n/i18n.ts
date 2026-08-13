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
      "pt-BR": ptBRResource,
      "en-US": enUSResource,
      "es-ES": esESResource
    },
    fallbackLng: "pt-BR",
    interpolation: {
      escapeValue: false
    }
  })

export const createModuleTranslator =
  (namespace: string) =>
  (key: string, defaultValue?: string): string => {
    let namespacedKey: string
    if (key.startsWith(`${namespace}:`) || key.startsWith(`${namespace}.`)) {
      namespacedKey = key.startsWith(`${namespace}:`)
        ? key
        : `${namespace}:${key.slice(namespace.length + 1)}`
    } else {
      namespacedKey = `${namespace}:${key}`
    }
    return defaultValue
      ? i18n.t(namespacedKey, { defaultValue })
      : i18n.t(namespacedKey)
  }

export default i18n
