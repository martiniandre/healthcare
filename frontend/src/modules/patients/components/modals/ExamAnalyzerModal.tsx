import { useTranslation } from "react-i18next"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../../shared/components/ui/Dialog"
import { ExamAnalyzer } from "../../../exam_analyzer/ExamAnalyzer"

interface ExamAnalyzerModalProps {
  isOpen: boolean
  onClose: () => void
  patientFhirId: string
}

export const ExamAnalyzerModal = ({
  isOpen,
  onClose,
  patientFhirId,
}: ExamAnalyzerModalProps) => {
  const { t } = useTranslation()

  if (!isOpen) {
    return null
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-6xl w-[95vw] h-[90vh] overflow-y-auto bg-gray-50/50 flex flex-col p-0">
        <DialogHeader className="p-4 sm:p-6 pb-0 shrink-0">
          <DialogTitle className="text-left sr-only">
            {t("examAnalyzer.title")}
          </DialogTitle>
        </DialogHeader>
        <div className="flex-1 overflow-auto">
          <ExamAnalyzer patientFhirId={patientFhirId} />
        </div>
      </DialogContent>
    </Dialog>
  )
}
