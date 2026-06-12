import { type FormEvent } from "react"
import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
import { Button } from "../../../shared/components/ui/Button"
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
  const { t } = useTranslation()

  return (
    <Card className="flex-1 p-8 border border-border bg-gray-50/50 flex flex-col items-center justify-center text-center gap-5 min-h-[400px]">
      <div className="bg-amber-50 p-4 rounded-full border border-amber-100 text-amber-500 animate-pulse-glow">
        <Lock className="w-10 h-10" />
      </div>

      <div className="max-w-md flex flex-col gap-1">
        <h3 className="text-md font-extrabold text-gray-900">
          {t("telemetry.restrictedRoom")}
        </h3>
        <p className="text-xs text-gray-500 leading-normal">
          {t("telemetry.restrictedRoomDesc", { roomName: activeRoomName })}
        </p>
      </div>

      <form onSubmit={handleUnlockRoom} className="w-full max-w-[320px] flex flex-col gap-3">
        <div className="flex flex-col gap-1.5 text-left">
          <label className="text-[10px] font-bold text-gray-500 uppercase tracking-wider">{t("telemetry.passcodeLabel")}</label>
          <div className="relative">
            <KeyRound className="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              type="password"
              placeholder={t("telemetry.passcodePlaceholder")}
              value={passcodeInput}
              onChange={(event) => setPasscodeInput(event.target.value)}
              className="w-full bg-white border border-border rounded-lg pl-9 pr-4 py-2.5 text-xs text-gray-800 focus:outline-none focus:border-primary/50 transition-all duration-200 font-mono"
              required
            />
          </div>
          {passcodeError && (
            <span className="text-[10px] text-red-500 font-bold mt-1 block">
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
          {t("telemetry.unlockButton")}
        </Button>
      </form>

      <div className="text-[10px] text-gray-400 font-semibold border-t border-border/80 pt-4 w-full max-w-sm mt-2 flex items-center justify-center gap-1.5">
        <Hospital className="w-3.5 h-3.5" />
        {t("telemetry.clinicalCouncil")}
      </div>
    </Card>
  )
}
