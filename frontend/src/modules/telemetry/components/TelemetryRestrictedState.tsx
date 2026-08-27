import { type FormEvent } from "react"
import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
import { Button } from "../../../shared/components/ui/Button"
import { Input } from "../../../shared/components/ui/Input"
import { Label } from "../../../shared/components/ui/Label"
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
    <Card className="flex-1 p-8 border border-border bg-muted-soft flex flex-col items-center justify-center text-center gap-5 min-h-[400px]">
      <div className="bg-warning-soft p-4 rounded-full border border-warning/20 text-warning animate-pulse-glow">
        <Lock className="w-10 h-10" />
      </div>

      <div className="max-w-md flex flex-col gap-1">
        <h3 className="text-md font-extrabold text-foreground">
          {t("restrictedRoom")}
        </h3>
        <p className="text-xs text-muted-foreground leading-normal">
          {t("restrictedRoomDesc", { roomName: activeRoomName })}
        </p>
      </div>

      <form onSubmit={handleUnlockRoom} className="w-full max-w-[320px] flex flex-col gap-3">
        <div className="flex flex-col gap-1.5 text-left">
          <Label className="text-[10px] font-bold text-muted-foreground uppercase tracking-wider">{t("passcodeLabel")}</Label>
          <div className="relative">
            <KeyRound className="w-4 h-4 text-muted-foreground absolute left-3 top-1/2 -translate-y-1/2" />
            <Input
              type="password"
              placeholder={t("passcodePlaceholder")}
              value={passcodeInput}
              onChange={(event) => setPasscodeInput(event.target.value)}
              className="pl-9 pr-4 text-xs font-mono"
              required
            />
          </div>
          {passcodeError && (
            <span className="text-[10px] text-danger font-bold mt-1 block">
              {passcodeError}
            </span>
          )}
        </div>

        <Button
          type="submit"
          variantType="primary"
          disabled={isPending}
          className="w-full py-2.5 text-xs font-bold gap-2"
        >
          <Unlock className="w-4 h-4" />
          {t("unlockButton")}
        </Button>
      </form>

      <div className="text-[10px] text-muted-foreground font-semibold border-t border-border/80 pt-4 w-full max-w-sm mt-2 flex items-center justify-center gap-1.5">
        <Hospital className="w-3.5 h-3.5" />
        {t("clinicalCouncil")}
      </div>
    </Card>
  )
}
