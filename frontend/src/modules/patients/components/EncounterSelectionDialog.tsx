import { useNavigate } from "react-router-dom"
import { useTranslation } from "react-i18next"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../shared/components/ui/dialog"
import { Badge } from "../../../shared/components/ui/badge"
import { Button } from "../../../shared/components/ui/button"
import type { Encounter } from "../types"
import { formatDateTime } from "../../../shared/utils/dates"

interface EncounterSelectionDialogProps {
  isOpen: boolean
  onClose: () => void
  encounters: Encounter[]
  patientId: string
}

export function EncounterSelectionDialog({
  isOpen,
  onClose,
  encounters,
  patientId,
}: EncounterSelectionDialogProps) {
  const navigate = useNavigate()
  const { t } = useTranslation("patients")

  const handleSelect = (encounterFhirId: string) => {
    onClose()
    navigate(`/patients/${patientId}/encounters/${encounterFhirId}`)
  }

  if (!isOpen) {
    return null
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[520px] max-h-[70vh] flex flex-col">
        <DialogHeader>
          <DialogTitle className="text-left">
            {t("details.selectEncounterTitle")}
          </DialogTitle>
        </DialogHeader>
        <div className="flex flex-col gap-2 overflow-y-auto pr-1 mt-2">
          {encounters.length === 0 ? (
            <div className="text-center py-8 text-sm text-muted-foreground font-medium">
              {t("details.encountersCard.empty")}
            </div>
          ) : (
            encounters.map((encounter) => (
              <Button
                key={encounter.fhir_id}
                type="button"
                variant="outline"
                onClick={() => handleSelect(encounter.fhir_id)}
                className="h-auto w-full justify-between gap-3 px-4 py-3 text-left font-normal"
              >
                <div className="flex flex-col gap-0.5 min-w-0">
                  <span className="text-sm font-bold truncate">
                    {encounter.reason_display}
                  </span>
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary" className="uppercase tracking-wider">
                      {encounter.status}
                    </Badge>
                    <span className="text-[11px] text-muted-foreground font-semibold">
                      {formatDateTime(encounter.created_at)}
                    </span>
                  </div>
                </div>
              </Button>
            ))
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
