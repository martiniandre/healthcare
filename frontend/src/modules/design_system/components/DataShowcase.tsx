import { useState } from "react"
import { useTranslation } from "react-i18next"
import { Avatar, AvatarFallback } from "@/shared/components/ui/avatar"
import { Button } from "@/shared/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/components/ui/card"
import { StatusBadge } from "@/shared/components/ui/status-badge"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow, type TableDensity } from "@/shared/components/ui/table"
import { TABLE_ROWS } from "../catalog"
import { Showcase } from "./SectionShell"

const DENSITIES: readonly TableDensity[] = ["comfortable", "compact"]

export const DataShowcase = () => {
  const { t } = useTranslation("designSystem")
  const [density, setDensity] = useState<TableDensity>("comfortable")

  return (
    <>
      <Showcase title={t("components.card.title")}>
        <Card>
          <CardHeader>
            <CardTitle>{t("components.card.heading")}</CardTitle>
            <CardDescription>{t("components.card.description")}</CardDescription>
          </CardHeader>
          <CardContent className="text-sm">{t("components.card.content")}</CardContent>
        </Card>
      </Showcase>
      <Showcase title={t("components.avatar.title")}>
        <Avatar>
          <AvatarFallback>{t("components.avatar.initials")}</AvatarFallback>
        </Avatar>
      </Showcase>
      <Showcase title={t("components.table.title")}>
        <div className="flex flex-wrap gap-2" role="group" aria-label={t("components.table.density")}>
          {DENSITIES.map((option) => (
            <Button
              key={option}
              size="sm"
              variant={density === option ? "default" : "outline"}
              aria-pressed={density === option}
              onClick={() => setDensity(option)}
            >
              {t(`components.table.${option}`)}
            </Button>
          ))}
        </div>
        <Table density={density}>
          <TableHeader>
            <TableRow>
              <TableHead>{t("components.table.patient")}</TableHead>
              <TableHead>{t("components.table.heartRate")}</TableHead>
              <TableHead>{t("components.table.status")}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {TABLE_ROWS.map((row, index) => (
              <TableRow key={row}>
                <TableCell>{t(`components.table.${row}`)}</TableCell>
                <TableCell className="tabular-nums">{72 + index * 9}</TableCell>
                <TableCell>
                  <StatusBadge status={index === 2 ? "critical" : "normal"}>
                    {t(index === 2 ? "components.statusBadge.critical" : "components.statusBadge.normal")}
                  </StatusBadge>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Showcase>
    </>
  )
}
