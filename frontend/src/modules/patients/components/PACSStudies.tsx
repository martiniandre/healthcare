import { Image as ImageIcon } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
import { Button } from "../../../shared/components/ui/Button"
import type { ImagingStudy } from "../../imaging/types"
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from "../../../shared/components/ui/Table"

interface PACSStudiesProps {
  studies: ImagingStudy[]
  onOpen: (id: string) => void
}

export const PACSStudies = ({ studies, onOpen }: PACSStudiesProps) => {
  const { t } = useTranslation("patients")

  return (
    <Card className="flex flex-col gap-5 min-h-[450px]">
      <div className="flex items-center justify-between border-b border-border pb-4">
        <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2">
          <ImageIcon className="w-4 h-4 text-secondary" />
          {t("details.pacsCard.title")}
        </h3>
      </div>

      {studies.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center gap-2 py-16">
          <ImageIcon className="w-8 h-8 text-gray-300" />
          <span className="text-xs text-muted">
            {t("details.pacsCard.empty")}
          </span>
        </div>
      ) : (
        <div className="overflow-x-auto w-full">
          <Table className="w-full text-left border-collapse">
            <TableHeader>
              <TableRow className="border-b border-border bg-gray-50/80">
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.pacsCard.display")}
                </TableHead>
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.pacsCard.modality")}
                </TableHead>
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.pacsCard.uid")}
                </TableHead>
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.pacsCard.date")}
                </TableHead>
                <TableHead className="py-3.5 px-4 text-right text-xs font-black text-gray-400 uppercase tracking-wider pr-6">
                  {t("details.pacsCard.action")}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {studies.map((study) => (
                <TableRow key={study.id} className="border-b border-border/60 hover:bg-gray-50 transition-colors duration-300">
                  <TableCell className="py-4 px-4">
                    <div className="flex items-center gap-3">
                      <div className="bg-secondary/10 p-2 rounded-lg border border-secondary/20 text-secondary">
                        <ImageIcon className="w-4 h-4" />
                      </div>
                      <span className="text-sm font-bold text-gray-800 block">
                        {study.title}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell className="py-4 px-4">
                    <span className="text-[10px] bg-secondary/15 text-secondary px-2.5 py-1 rounded font-black uppercase">
                      {study.modality}
                    </span>
                  </TableCell>
                  <TableCell className="py-4 px-4">
                    <span className="text-xs font-mono text-gray-500 max-w-[150px] block truncate">
                      {study.study_instance_uid}
                    </span>
                  </TableCell>
                  <TableCell className="py-4 px-4">
                    <span className="text-xs text-gray-500 font-semibold">
                      {new Date(study.created_at).toLocaleString()}
                    </span>
                  </TableCell>
                  <TableCell className="py-4 px-4 text-right pr-6">
                    <Button
                      variantType="outline"
                      onClick={() => onOpen(study.id)}
                      className="px-2.5 py-1 text-[10px] font-bold border-secondary/20 hover:bg-secondary/10 text-secondary"
                    >
                      {t("details.pacsCard.action")}
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </Card>
  )
}
