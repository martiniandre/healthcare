import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
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
  const { t } = useTranslation()

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
                ? "border-primary bg-primary/5 shadow-sm" 
                : "border-border hover:border-gray-300 bg-white"
            }`}
          >
            <div className="flex items-center justify-between">
              <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">
                {t("telemetry.monitoredWing")}
              </span>
              <span className={`inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full border ${
                isUnlocked 
                  ? "bg-emerald-50 text-emerald-600 border-emerald-100" 
                  : "bg-amber-50 text-amber-600 border-amber-100"
              }`}>
                {isUnlocked ? (
                  <>
                    <Unlock className="w-3 h-3" />
                    {t("telemetry.unlocked")}
                  </>
                ) : (
                  <>
                    <Lock className="w-3 h-3" />
                    {t("telemetry.protected")}
                  </>
                )}
              </span>
            </div>

            <h4 className="text-sm font-extrabold text-gray-900 mt-2 block">
              {roomItem.name}
            </h4>
            <span className="text-[11px] text-gray-400 block mt-1 leading-normal">
              {roomItem.description}
            </span>

            {isUnlocked && (
              <button
                onClick={(event) => {
                  event.stopPropagation()
                  handleLockRoom(roomItem.id)
                }}
                className="absolute bottom-4 right-4 text-xs text-red-500 hover:text-red-700 transition-colors font-bold"
              >
                {t("telemetry.lockButton")}
              </button>
            )}
          </Card>
        )
      })}
    </div>
  )
}
