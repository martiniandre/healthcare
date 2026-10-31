import { useEffect, useRef, useState, type FormEvent } from "react"
import { useTranslation } from "react-i18next"
import { Card } from "@/shared/components/ui/card"
import { Button } from "@/shared/components/ui/button"
import { Input } from "@/shared/components/ui/input"
import { Field, FieldDescription, FieldError, FieldLabel } from "@/shared/components/ui/field"
import { Separator } from "@/shared/components/ui/separator"
import { Spinner } from "@/shared/components/ui/spinner"
import { Lock, KeyRound, Hospital, Unlock, Eye, EyeOff } from "lucide-react"

interface TelemetryRestrictedStateProps {
  activeRoomName?: string
  passcodeInput: string
  setPasscodeInput: (value: string) => void
  passcodeError: string
  isPending: boolean
  handleUnlockRoom: (event: FormEvent) => void
}

export const TelemetryRestrictedState = ({
  activeRoomName,
  passcodeInput,
  setPasscodeInput,
  passcodeError,
  isPending,
  handleUnlockRoom
}: TelemetryRestrictedStateProps) => {
  const { t } = useTranslation("telemetry")
  const inputReference = useRef<HTMLInputElement | null>(null)
  const [isPasscodeVisible, setIsPasscodeVisible] = useState(false)
  const [isEmptyAttempt, setIsEmptyAttempt] = useState(false)
  const displayedError = isEmptyAttempt ? t("passcodeRequired") : passcodeError
  const describedBy = displayedError
    ? "telemetry-passcode-hint telemetry-passcode-error"
    : "telemetry-passcode-hint"

  useEffect(() => {
    if (passcodeError && !isPending) {
      inputReference.current?.focus()
      inputReference.current?.select()
    }
  }, [passcodeError, isPending])

  const handleSubmit = (event: FormEvent) => {
    if (!passcodeInput.trim()) {
      event.preventDefault()
      setIsEmptyAttempt(true)
      inputReference.current?.focus()
      return
    }
    setIsEmptyAttempt(false)
    handleUnlockRoom(event)
  }

  const handlePasscodeChange = (value: string) => {
    setIsEmptyAttempt(false)
    setPasscodeInput(value)
  }

  return (
    <Card className="flex-1 min-h-[400px] items-center justify-center gap-5 bg-muted/50 p-8 text-center">
      <div className="rounded-full border border-warning/30 bg-warning/10 p-4 text-warning animate-pulse-glow">
        <Lock className="size-10" />
      </div>

      <div className="max-w-md flex flex-col gap-1">
        <h3 className="text-base font-extrabold text-foreground">
          {t("restrictedRoom")}
        </h3>
        <p className="text-xs leading-normal text-muted-foreground">
          {t("restrictedRoomDesc", { roomName: activeRoomName })}
        </p>
      </div>

      <form
        onSubmit={handleSubmit}
        noValidate
        aria-busy={isPending}
        className="w-full max-w-[320px] flex flex-col gap-3"
      >
        <Field data-invalid={displayedError ? true : undefined} className="text-left">
          <FieldLabel htmlFor="telemetry-passcode" className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
            {t("passcodeLabel")}
            <span aria-hidden="true" className="text-destructive">*</span>
            <span className="sr-only">{t("requiredLabel")}</span>
          </FieldLabel>
          <div className="relative">
            <KeyRound aria-hidden="true" className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              id="telemetry-passcode"
              ref={inputReference}
              type={isPasscodeVisible ? "text" : "password"}
              inputMode="numeric"
              autoComplete="off"
              autoCapitalize="none"
              spellCheck={false}
              maxLength={32}
              autoFocus
              placeholder={t("passcodePlaceholder")}
              value={passcodeInput}
              onChange={(event) => handlePasscodeChange(event.target.value)}
              readOnly={isPending}
              aria-required="true"
              aria-invalid={displayedError ? true : undefined}
              aria-describedby={describedBy}
              className="h-11 pl-9 pr-28 font-mono text-base md:h-9 md:text-xs"
            />
            <Button
              type="button"
              variant="ghost"
              size="sm"
              aria-pressed={isPasscodeVisible}
              onClick={() => setIsPasscodeVisible((visible) => !visible)}
              className="absolute right-1 top-1/2 h-9 -translate-y-1/2 gap-1 px-2 text-[11px] text-muted-foreground md:h-7"
            >
              {isPasscodeVisible ? <EyeOff /> : <Eye />}
              {t("togglePasscode")}
            </Button>
          </div>
          <FieldDescription id="telemetry-passcode-hint" className="text-[11px]">
            {t("passcodeHint")}
          </FieldDescription>
          {displayedError && (
            <FieldError id="telemetry-passcode-error" className="text-[11px] font-bold">
              {displayedError}
            </FieldError>
          )}
        </Field>

        <Button
          type="submit"
          disabled={isPending}
          className="h-11 w-full text-xs font-bold md:h-9"
        >
          {isPending ? <Spinner /> : <Unlock />}
          {isPending ? t("unlocking") : t("unlockButton")}
        </Button>
      </form>

      <div className="mt-2 flex w-full max-w-sm flex-col items-center gap-4">
        <Separator />
        <div className="flex items-center justify-center gap-1.5 text-[10px] font-semibold text-muted-foreground">
          <Hospital className="size-3.5" />
          {t("clinicalCouncil")}
        </div>
      </div>
    </Card>
  )
}
