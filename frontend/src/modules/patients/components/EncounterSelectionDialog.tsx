import { useNavigate } from "react-router-dom"
import { useTranslation } from "react-i18next"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  Button,
} from "../../../shared/components/ui"
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
            <div className="text-center py-8 text-sm text-muted-foreground/70 font-medium">
              {t("details.encountersCard.empty")}
            </div>
          ) : (
            encounters.map((encounter) => (
              <Button
                key={encounter.fhir_id}
                type="button"
                variantType="outline"
                onClick={() => handleSelect(encounter.fhir_id)}
                className="w-full justify-between gap-3 px-4 py-3 h-auto rounded-lg hover:border-primary/30 hover:bg-muted-soft text-foreground transition-all duration-200"
              >
                <div className="flex flex-col gap-0.5 min-w-0">
                  <span className="text-sm font-bold truncate">
                    {encounter.reason_display}
                  </span>
                  <div className="flex items-center gap-2">
                    <span className="text-[10px] bg-muted-soft text-muted-foreground px-2 py-0.5 rounded font-bold uppercase tracking-wider">
                      {encounter.status}
                    </span>
                    <span className="text-[11px] text-muted-foreground/70 font-semibold">
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