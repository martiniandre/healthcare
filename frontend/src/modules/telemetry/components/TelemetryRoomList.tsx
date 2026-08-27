import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
import { Button } from "../../../shared/components/ui/Button"
import { Lock, Unlock } from "lucide-react"

interface Room {
  id: string
  name: string
  description: string
}

interface TelemetryRoomListProps {
  rooms: Room[]
  selectedRoomId: string | null
  unlockedRoomIds: string[]
  handleSelectRoom: (roomId: string) => void
  handleLockRoom: (roomId: string) => void
}

export const TelemetryRoomList = ({
  rooms,
  selectedRoomId,
  unlockedRoomIds,
  handleSelectRoom,
  handleLockRoom
}: TelemetryRoomListProps) => {
  const { t } = useTranslation("telemetry")

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      {rooms.map((roomItem) => {
        const isSelected = roomItem.id === selectedRoomId
        const isUnlocked = unlockedRoomIds.includes(roomItem.id)

        return (
          <Card
            key={roomItem.id}
            onClick={() => handleSelectRoom(roomItem.id)}
            className={`p-4 cursor-pointer text-left transition-all duration-200 border relative ${
              isSelected 
                ? "border-primary bg-primary-soft shadow-sm" 
                : "border-border hover:border-border-strong bg-surface"
            }`}
          >
            <div className="flex items-center justify-between">
              <span className="text-[10px] text-muted-foreground font-bold uppercase tracking-wider block">
                {t("monitoredWing")}
              </span>
              <span className={`inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full border ${
                isUnlocked 
                  ? "bg-success-soft text-success border-success/20" 
                  : "bg-warning-soft text-warning border-warning/20"
              }`}>
                {isUnlocked ? (
                  <>
                    <Unlock className="w-3 h-3" />
                    {t("unlocked")}
                  </>
                ) : (
                  <>
                    <Lock className="w-3 h-3" />
                    {t("protected")}
                  </>
                )}
              </span>
            </div>

            <h4 className="text-sm font-extrabold text-foreground mt-2 block">
              {roomItem.name}
            </h4>
            <span className="text-[11px] text-muted-foreground block mt-1 leading-normal">
              {roomItem.description}
            </span>

            {isUnlocked && (
              <Button
                variantType="ghost"
                size="sm"
                onClick={(event) => {
                  event.stopPropagation()
                  handleLockRoom(roomItem.id)
                }}
                className="absolute bottom-4 right-4 text-xs text-danger hover:text-danger hover:bg-danger-soft transition-colors font-bold"
              >
                {t("lockButton")}
              </Button>
            )}
          </Card>
        )
      })}
    </div>
  )
}
