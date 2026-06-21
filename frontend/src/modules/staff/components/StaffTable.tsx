import { useTranslation } from "react-i18next"
import { Users } from "lucide-react"
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from "../../../shared/components/ui/Table"
import { Skeleton } from "../../../shared/components/ui/Skeleton"
import { StaffRole } from "../../../shared/types"
import type { StaffMember } from "../types"

interface StaffTableProps {
  isLoading: boolean
  filteredStaff: StaffMember[]
}

export const StaffTable = ({ isLoading, filteredStaff }: StaffTableProps) => {
  console.log({ filteredStaff })
  const { t } = useTranslation()

  const getRoleLabel = (role: string) => {
    switch (role) {
      case StaffRole.Doctor: return t("staff.table.roles.doctor", "Médico")
      case StaffRole.Nurse: return t("staff.table.roles.nurse", "Enfermeiro")
      case StaffRole.Receptionist: return t("staff.table.roles.receptionist", "Recepção")
      case StaffRole.Admin: return t("staff.table.roles.admin", "Admin")
      default: return role
    }
  }

  return (
    <div className="overflow-x-auto border border-border rounded-xl w-full bg-white">
      <Table className="min-w-[700px] md:min-w-0">
        <TableHeader>
          <TableRow className="hover:bg-transparent">
            <TableHead>{t("staff.table.professional")}</TableHead>
            <TableHead>{t("staff.table.role")}</TableHead>
            <TableHead>{t("staff.table.license")}</TableHead>
            <TableHead>{t("staff.table.department")}</TableHead>
            <TableHead>{t("staff.table.status")}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading ? (
            Array.from({ length: 5 }).map((_, i) => (
              <TableRow key={`skeleton-${i}`}>
                <TableCell>
                  <div className="flex items-center gap-3">
                    <Skeleton className="w-8 h-8 rounded-lg" />
                    <div className="flex flex-col gap-1.5">
                      <Skeleton className="h-4 w-32" />
                      <Skeleton className="h-3 w-24" />
                    </div>
                  </div>
                </TableCell>
                <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                <TableCell><Skeleton className="h-4 w-20" /></TableCell>
              </TableRow>
            ))
          ) : (
            filteredStaff.map((member) => (
              <TableRow key={member.id} className="group">
                <TableCell>
                  <div className="flex items-center gap-3">
                    <div className="bg-primary/8 p-2 rounded-lg border border-primary/10">
                      <Users className="w-4 h-4 text-primary" />
                    </div>
                    <div className="flex flex-col min-w-0">
                      <span className="font-extrabold text-gray-900 truncate">{member.fullName}</span>
                      <span className="text-[10px] text-gray-500 truncate mt-0.5">{member.email}</span>
                    </div>
                  </div>
                </TableCell>
                <TableCell>
                  <span className="font-semibold text-gray-700 text-xs">{getRoleLabel(member.role)}</span>
                </TableCell>
                <TableCell className="font-mono text-xs text-gray-600">
                  {member.license}
                </TableCell>
                <TableCell className="text-xs font-medium text-gray-600">
                  {member.department}
                </TableCell>
                <TableCell>
                  <span className="text-xs font-semibold text-gray-700">{member.status}</span>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  )
}
