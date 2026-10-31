import { Image as ImageIcon } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Badge } from "../../../shared/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "../../../shared/components/ui/card"
import { Button } from "../../../shared/components/ui/button"
import { ClinicalQuerySection } from "./ClinicalQuerySection"
import { formatDateTime } from "../../../shared/utils/dates"
import type { ImagingStudy } from "../../imaging/types"
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from "../../../shared/components/ui/table"

interface PACSStudiesProps {
  studies: ImagingStudy[]
  onOpen: (id: string) => void
  isLoading?: boolean
  isError?: boolean
  onRetry?: () => void
}

export const PACSStudies = ({ studies, onOpen, isLoading = false, isError = false, onRetry }: PACSStudiesProps) => {
  const { t } = useTranslation("patients")

  return (
    <ClinicalQuerySection
      title={t("details.pacsCard.title")}
      icon={<ImageIcon className="w-4 h-4 text-secondary" />}
      isLoading={isLoading}
      isError={isError}
      onRetry={() => onRetry?.()}
      isEmpty={studies.length === 0}
      emptyTitle={t("details.pacsCard.empty")}
      emptyIcon={ImageIcon}
    >
      <Card className="min-h-[450px] gap-0">
        <CardHeader className="border-b border-border">
          <CardTitle className="flex items-center gap-2 font-extrabold">
            <ImageIcon className="w-4 h-4 text-secondary" />
            {t("details.pacsCard.title")}
          </CardTitle>
        </CardHeader>
        <CardContent className="overflow-x-auto w-full">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50">
                  <TableHead className="py-3.5 px-4 text-xs font-black text-muted-foreground uppercase tracking-wider">
                    {t("details.pacsCard.display")}
                  </TableHead>
                  <TableHead className="py-3.5 px-4 text-xs font-black text-muted-foreground uppercase tracking-wider">
                    {t("details.pacsCard.modality")}
                  </TableHead>
                  <TableHead className="py-3.5 px-4 text-xs font-black text-muted-foreground uppercase tracking-wider">
                    {t("details.pacsCard.uid")}
                  </TableHead>
                  <TableHead className="py-3.5 px-4 text-xs font-black text-muted-foreground uppercase tracking-wider">
                    {t("details.pacsCard.date")}
                  </TableHead>
                  <TableHead className="py-3.5 px-4 text-right text-xs font-black text-muted-foreground uppercase tracking-wider pr-6">
                    {t("details.pacsCard.action")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {studies.map((study) => (
                  <TableRow key={study.id}>
                    <TableCell className="py-4 px-4">
                      <div className="flex items-center gap-3">
                        <div className="bg-secondary/10 p-2 rounded-lg border border-secondary/20 text-secondary">
                          <ImageIcon className="w-4 h-4" />
                        </div>
                        <span className="text-sm font-bold text-foreground block">
                          {study.title}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell className="py-4 px-4">
                      <Badge variant="secondary" className="uppercase">
                        {study.modality}
                      </Badge>
                    </TableCell>
                    <TableCell className="py-4 px-4">
                      <span className="text-xs font-mono text-muted-foreground max-w-[150px] block truncate">
                        {study.study_instance_uid}
                      </span>
                    </TableCell>
                    <TableCell className="py-4 px-4">
                      <span className="text-xs text-muted-foreground font-semibold">
                        {formatDateTime(study.created_at)}
                      </span>
                    </TableCell>
                    <TableCell className="py-4 px-4 text-right pr-6">
                      <Button
                        variant="outline"
                        size="xs"
                        onClick={() => onOpen(study.id)}
                      >
                        {t("details.pacsCard.action")}
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
      </Card>
    </ClinicalQuerySection>
  )
}
