import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import tseslint from 'typescript-eslint'
import { defineConfig, globalIgnores } from 'eslint/config'

import { noHardcodedUserText } from './eslint-rules/no-hardcoded-user-text.js'

const i18nUserTextPlugin = {
  rules: {
    'no-hardcoded-user-text': noHardcodedUserText,
  },
}

export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      js.configs.recommended,
      tseslint.configs.recommended,
      reactHooks.configs.flat.recommended,
      reactRefresh.configs.vite,
    ],
    languageOptions: {
      globals: globals.browser,
    },
  },
  {
    files: ['src/**/*.{ts,tsx}'],
    ignores: ['**/*.test.{ts,tsx}', '**/*.spec.{ts,tsx}', 'e2e/**'],
    plugins: {
      'i18n-user-text': i18nUserTextPlugin,
    },
    rules: {
      'i18n-user-text/no-hardcoded-user-text': 'error',
    },
  },
])
