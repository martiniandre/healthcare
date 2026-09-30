import { readdirSync, readFileSync } from "node:fs"
import { join } from "node:path"

const localesDirectory = join(process.cwd(), "src/shared/i18n/locales")
const locales = readdirSync(localesDirectory, { withFileTypes: true })
  .filter((entry) => entry.isDirectory())
  .map((entry) => entry.name)
  .sort()

const referenceLocale = locales[0]

const collectKeys = (locale, fileName) => {
  const parsed = JSON.parse(readFileSync(join(localesDirectory, locale, fileName), "utf8"))
  const keys = []
  const pending = Object.entries(parsed).map(([key, value]) => [key, value])
  while (pending.length > 0) {
    const [key, value] = pending.pop()
    if (value && typeof value === "object" && !Array.isArray(value)) {
      for (const [childKey, childValue] of Object.entries(value)) {
        pending.push([`${key}.${childKey}`, childValue])
      }
    } else {
      keys.push(key)
    }
  }
  return keys.sort()
}

const fileNames = readdirSync(join(localesDirectory, referenceLocale))
  .filter((fileName) => fileName.endsWith(".json"))
  .sort()

let missingTotal = 0

for (const fileName of fileNames) {
  const referenceKeys = new Set(collectKeys(referenceLocale, fileName))
  for (const locale of locales) {
    if (locale === referenceLocale) {
      continue
    }
    const localeKeys = new Set(collectKeys(locale, fileName))
    const missing = [...referenceKeys].filter((key) => !localeKeys.has(key))
    const extra = [...localeKeys].filter((key) => !referenceKeys.has(key))
    if (missing.length > 0 || extra.length > 0) {
      console.log(`${fileName} [${locale}]`)
      if (missing.length > 0) {
        console.log(`  missing: ${missing.join(", ")}`)
      }
      if (extra.length > 0) {
        console.log(`  extra: ${extra.join(", ")}`)
      }
      missingTotal += missing.length + extra.length
    }
  }
}

console.log(
  missingTotal === 0
    ? `PARITY OK across ${locales.length} locales for ${fileNames.length} namespaces`
    : `PARITY FAILURES: ${missingTotal}`,
)
process.exit(missingTotal === 0 ? 0 : 1)
