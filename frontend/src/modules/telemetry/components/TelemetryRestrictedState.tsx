import { type FormEvent } from "react"
import { useTranslation } from "react-i18next"
import { Card } from "@/shared/components/ui/card"
import { Button } from "@/shared/components/ui/button"
import { Input } from "@/shared/components/ui/input"
import { Field, FieldError, FieldLabel } from "@/shared/components/ui/field"
import { Separator } from "@/shared/components/ui/separator"
import { Spinner } from "@/shared/components/ui/spinner"
import { Lock, KeyRound, Hospital, Unlock } from "lucide-react"

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

      <form onSubmit={handleUnlockRoom} className="w-full max-w-[320px] flex flex-col gap-3">
        <Field data-invalid={passcodeError ? true : undefined} className="text-left">
          <FieldLabel htmlFor="telemetry-passcode" className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
            {t("passcodeLabel")}
          </FieldLabel>
          <div className="relative">
            <KeyRound className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              id="telemetry-passcode"
              type="password"
              placeholder={t("passcodePlaceholder")}
              value={passcodeInput}
              onChange={(event) => setPasscodeInput(event.target.value)}
              aria-invalid={passcodeError ? true : undefined}
              className="pl-9 font-mono text-xs"
              required
            />
          </div>
          {passcodeError && <FieldError className="text-[10px] font-bold">{passcodeError}</FieldError>}
        </Field>

        <Button
          type="submit"
          disabled={isPending}
          className="w-full text-xs font-bold"
        >
          {isPending ? <Spinner /> : <Unlock />}
          {t("unlockButton")}
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
