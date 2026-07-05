import { useTranslation } from "react-i18next"
import { CheckCircle } from "lucide-react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../shared/components/ui/Dialog"
import type { Encounter } from "../types"

interface EncounterSelectionDialogProps {
  isOpen: boolean
  onClose: () => void
  encounters: Encounter[]
  selectedEncounterId: string | null
  onSelect: (id: string) => void
}

export function EncounterSelectionDialog({
  isOpen,
  onClose,
  encounters,
  selectedEncounterId,
  onSelect,
}: EncounterSelectionDialogProps) {
  const { t } = useTranslation("patients")

  const handleSelect = (id: string) => {
    onSelect(id)
    onClose()
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
            <div className="text-center py-8 text-sm text-gray-400 font-medium">
              {t("details.encountersCard.empty")}
            </div>
          ) : (
            encounters.map((encounter) => {
              const isActive = selectedEncounterId === encounter.fhir_id
              return (
                <button
                  key={encounter.fhir_id}
                  onClick={() => handleSelect(encounter.fhir_id)}
                  className={`w-full text-left flex items-center justify-between gap-3 px-4 py-3 rounded-lg border transition-all duration-200 ${
                    isActive
                      ? "bg-primary/8 border-primary/40 text-primary shadow-sm"
                      : "bg-white border-border hover:border-primary/30 hover:bg-gray-50 text-gray-700"
                  }`}
                >
                  <div className="flex flex-col gap-0.5 min-w-0">
                    <span className="text-sm font-bold truncate">
                      {encounter.reason_display}
                    </span>
                    <div className="flex items-center gap-2">
                      <span className="text-[10px] bg-gray-100 text-gray-500 px-2 py-0.5 rounded font-bold uppercase tracking-wider">
                        {encounter.status}
                      </span>
                      <span className="text-[11px] text-gray-400 font-semibold">
                        {new Date(encounter.created_at).toLocaleString()}
                      </span>
                    </div>
                  </div>
                  {isActive && (
                    <span className="text-xs font-bold text-primary flex items-center gap-1 shrink-0">
                      <CheckCircle className="w-4 h-4" />
                      {t("details.focus")}
                    </span>
                  )}
                </button>
              )
            })
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
