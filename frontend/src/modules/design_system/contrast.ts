export interface RgbColor {
  red: number
  green: number
  blue: number
}

export type ContrastKind = "text" | "ui"

export interface ContrastEvaluation {
  ratio: number
  passesAa: boolean
  passesAaa: boolean
}

const CONTRAST_THRESHOLDS: Record<ContrastKind, { aa: number; aaa: number }> = {
  text: { aa: 4.5, aaa: 7 },
  ui: { aa: 3, aaa: 4.5 },
}

const HEX_PATTERN = /^#([0-9a-f]{3}|[0-9a-f]{6})$/i
const RGB_PATTERN = /^rgba?\(\s*([\d.]+)[\s,]+([\d.]+)[\s,]+([\d.]+)/i
const SRGB_PATTERN = /^color\(\s*srgb\s+([\d.]+)\s+([\d.]+)\s+([\d.]+)/i

const expandHex = (hex: string): string =>
  hex.length === 3 ? hex.split("").map((digit) => digit + digit).join("") : hex

export const parseColor = (value: string): RgbColor | null => {
  const trimmed = value.trim()
  const hexMatch = HEX_PATTERN.exec(trimmed)
  if (hexMatch) {
    const hex = expandHex(hexMatch[1])
    return {
      red: Number.parseInt(hex.slice(0, 2), 16),
      green: Number.parseInt(hex.slice(2, 4), 16),
      blue: Number.parseInt(hex.slice(4, 6), 16),
    }
  }
  const rgbMatch = RGB_PATTERN.exec(trimmed)
  if (rgbMatch) {
    return {
      red: Number(rgbMatch[1]),
      green: Number(rgbMatch[2]),
      blue: Number(rgbMatch[3]),
    }
  }
  const srgbMatch = SRGB_PATTERN.exec(trimmed)
  if (srgbMatch) {
    return {
      red: Math.round(Number(srgbMatch[1]) * 255),
      green: Math.round(Number(srgbMatch[2]) * 255),
      blue: Math.round(Number(srgbMatch[3]) * 255),
    }
  }
  return null
}

const linearizeChannel = (channel: number): number => {
  const normalized = channel / 255
  return normalized <= 0.03928 ? normalized / 12.92 : ((normalized + 0.055) / 1.055) ** 2.4
}

export const relativeLuminance = (color: RgbColor): number =>
  0.2126 * linearizeChannel(color.red) +
  0.7152 * linearizeChannel(color.green) +
  0.0722 * linearizeChannel(color.blue)

export const contrastRatio = (first: RgbColor, second: RgbColor): number => {
  const firstLuminance = relativeLuminance(first)
  const secondLuminance = relativeLuminance(second)
  const lighter = Math.max(firstLuminance, secondLuminance)
  const darker = Math.min(firstLuminance, secondLuminance)
  return (lighter + 0.05) / (darker + 0.05)
}

export const evaluateContrast = (
  foreground: string,
  background: string,
  kind: ContrastKind = "text",
): ContrastEvaluation | null => {
  const foregroundColor = parseColor(foreground)
  const backgroundColor = parseColor(background)
  if (!foregroundColor || !backgroundColor) {
    return null
  }
  const ratio = contrastRatio(foregroundColor, backgroundColor)
  const thresholds = CONTRAST_THRESHOLDS[kind]
  return {
    ratio,
    passesAa: ratio >= thresholds.aa,
    passesAaa: ratio >= thresholds.aaa,
  }
}
