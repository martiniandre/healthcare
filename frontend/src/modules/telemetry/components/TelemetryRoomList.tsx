import { useTranslation } from "react-i18next"
import { Card, CardContent } from "@/shared/components/ui/card"
import { Badge } from "@/shared/components/ui/badge"
import { Button } from "@/shared/components/ui/button"
import { cn } from "@/shared/utils/cn"
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
            className={cn(
              "relative cursor-pointer gap-0 py-4 text-left transition-colors duration-200",
              isSelected ? "border-primary bg-primary/5" : "hover:border-ring/50"
            )}
          >
            <CardContent className="px-4">
              <div className="flex items-start justify-between gap-2">
                <h4 className="block text-sm font-semibold text-foreground">
                  {roomItem.name}
                </h4>
                {isUnlocked ? (
                  <Badge variant="success">
                    <Unlock />
                    {t("unlocked")}
                  </Badge>
                ) : (
                  <Badge variant="warning">
                    <Lock />
                    {t("protected")}
                  </Badge>
                )}
              </div>
              <span className="mt-1 block text-xs leading-normal text-muted-foreground">
                {roomItem.description}
              </span>

              {isUnlocked && (
                <Button
                  variant="ghost"
                  size="xs"
                  onClick={(event) => {
                    event.stopPropagation()
                    handleLockRoom(roomItem.id)
                  }}
                  className="absolute bottom-3 right-3 font-semibold text-destructive hover:text-destructive"
                >
                  {t("lockButton")}
                </Button>
              )}
            </CardContent>
          </Card>
        )
      })}
    </div>
  )
}
