import { useEffect, useRef, useState } from "react"
import { useTranslation } from "react-i18next"
import { StatusBadge } from "@/shared/components/ui/status-badge"
import { evaluateContrast, type ContrastEvaluation } from "../contrast"
import type { ColorToken } from "../catalog"

interface ColorSwatchProperties {
  token: ColorToken
}

const useResolvedContrast = (token: ColorToken) => {
  const swatchReference = useRef<HTMLDivElement>(null)
  const [evaluation, setEvaluation] = useState<ContrastEvaluation | null>(null)

  useEffect(() => {
    const measure = () => {
      const swatch = swatchReference.current
      if (!swatch) {
        return
      }
      const computedStyle = getComputedStyle(swatch)
      setEvaluation(evaluateContrast(computedStyle.color, computedStyle.backgroundColor, token.kind))
    }
    const frame = requestAnimationFrame(measure)
    const observer = new MutationObserver(measure)
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ["class", "data-theme"] })
    return () => {
      cancelAnimationFrame(frame)
      observer.disconnect()
    }
  }, [token.kind])

  return { swatchReference, evaluation }
}

export const ColorSwatch = ({ token }: ColorSwatchProperties) => {
  const { t } = useTranslation("designSystem")
  const { swatchReference, evaluation } = useResolvedContrast(token)

  return (
    <li className="flex flex-col overflow-hidden rounded-xl border bg-card text-card-foreground">
      <div
        ref={swatchReference}
        data-testid={`swatch-${token.name}`}
        className="flex h-20 items-end justify-between border-b p-3 text-lg font-semibold"
        style={{
          backgroundColor: `var(--color-${token.name})`,
          color: `var(--color-${token.against})`,
        }}
      >
        <span aria-hidden="true">{t("colors.sample")}</span>
      </div>
      <div className="flex flex-col gap-2 p-3">
        <div>
          <p className="text-sm font-medium">{token.name}</p>
          <p className="font-mono text-xs text-muted-foreground">{`--color-${token.name}`}</p>
          <p className="text-xs text-muted-foreground">
            {t("colors.against", { token: token.against })}
          </p>
        </div>
        {evaluation ? (
          <>
            <p className="text-sm font-medium tabular-nums">
              {t("colors.ratio", { ratio: evaluation.ratio.toFixed(2) })}
            </p>
            <div className="flex flex-wrap gap-1.5">
              <StatusBadge status={evaluation.passesAa ? "normal" : "critical"}>
                {evaluation.passesAa ? t("colors.aaPass") : t("colors.aaFail")}
              </StatusBadge>
              <StatusBadge status={evaluation.passesAaa ? "normal" : "attention"}>
                {evaluation.passesAaa ? t("colors.aaaPass") : t("colors.aaaFail")}
              </StatusBadge>
            </div>
          </>
        ) : (
          <p className="text-xs text-muted-foreground">{t("colors.unavailable")}</p>
        )}
      </div>
    </li>
  )
}
