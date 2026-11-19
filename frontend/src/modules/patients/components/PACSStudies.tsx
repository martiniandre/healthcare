import { Copy, Image as ImageIcon } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Badge } from "../../../shared/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "../../../shared/components/ui/card"
import { Button } from "../../../shared/components/ui/button"
import { ClinicalQuerySection } from "./ClinicalQuerySection"
import { Tooltip, TooltipContent, TooltipTrigger } from "../../../shared/components/ui/tooltip"
import { toast } from "../../../shared/store/toast_store"
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
      isLoading={isLoading}
      isError={isError}
      onRetry={() => onRetry?.()}
      isEmpty={studies.length === 0}
      emptyTitle={t("details.pacsCard.empty")}
      emptyIcon={ImageIcon}
    >
      <Card className="min-h-[450px] gap-0">
        <CardHeader className="border-b border-border">
          <CardTitle className="text-base font-semibold">
            {t("details.pacsCard.title")}
          </CardTitle>
        </CardHeader>
        <CardContent className="w-full overflow-x-auto px-0">
            <Table>
              <TableHeader>
                <TableRow className="hover:bg-transparent">
                  <TableHead className="pl-6">
                    {t("details.pacsCard.display")}
                  </TableHead>
                  <TableHead>
                    {t("details.pacsCard.modality")}
                  </TableHead>
                  <TableHead>
                    {t("details.pacsCard.uid")}
                  </TableHead>
                  <TableHead>
                    {t("details.pacsCard.date")}
                  </TableHead>
                  <TableHead className="pr-6 text-right">
                    {t("details.pacsCard.action")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {studies.map((study) => (
                  <TableRow key={study.id}>
                    <TableCell className="pl-6">
                      <div className="flex items-center gap-3">
                        <div className="rounded-md border border-border bg-muted p-2 text-muted-foreground">
                          <ImageIcon className="w-4 h-4" />
                        </div>
                        <span className="block text-sm font-semibold text-foreground">
                          {study.title}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline" className="font-mono uppercase">
                        {study.modality}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <span tabIndex={0} className="block max-w-[180px] truncate rounded-sm font-mono text-xs text-muted-foreground outline-none focus-visible:ring-2 focus-visible:ring-ring">
                              {study.study_instance_uid}
                            </span>
                          </TooltipTrigger>
                          <TooltipContent className="font-mono">{study.study_instance_uid}</TooltipContent>
                        </Tooltip>
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon-xs"
                          aria-label={t("details.pacsCard.copyUid")}
                          onClick={() => {
                            void navigator.clipboard
                              ?.writeText(study.study_instance_uid)
                              .then(() => toast.success(t("details.pacsCard.uidCopied")))
                          }}
                        >
                          <Copy />
                        </Button>
                      </div>
                    </TableCell>
                    <TableCell>
                      <span className="text-sm tabular-nums text-muted-foreground">
                        {formatDateTime(study.created_at)}
                      </span>
                    </TableCell>
                    <TableCell className="pr-6 text-right">
                      <Button
                        variant="outline"
                        size="sm"
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
