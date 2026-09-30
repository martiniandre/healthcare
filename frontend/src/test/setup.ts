import '@testing-library/jest-dom'
import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

import enUSResource from '../shared/i18n/locales/en-US'

Element.prototype.scrollIntoView = () => {}

if (!i18n.isInitialized) {
  void i18n.use(initReactI18next).init({
    resources: {
      'en-US': enUSResource,
    },
    lng: 'en-US',
    fallbackLng: 'en-US',
    interpolation: {
      escapeValue: false,
    },
  })
}
