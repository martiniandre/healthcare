import type { ReactNode } from "react"
import { PatientHeader } from "./PatientHeader"
import { Card } from "../../../shared/components/ui/card"
import { PageContainer } from "../../../shared/components/ui/PageContainer"

interface PatientPageLayoutProps {
  patient: {
    patient_id: string
    fhir_resource_id: string
    full_name: string
    birth_date: string
    document_id: string
    phone_number: string
  }
  onBack: () => void
  sidebarTop: ReactNode
  headerActions?: ReactNode
  encounterBanner?: ReactNode
  children: ReactNode
}

export function PatientPageLayout({
  patient,
  onBack,
  sidebarTop,
  headerActions,
  encounterBanner,
  children,
}: PatientPageLayoutProps) {
  return (
    <PageContainer>
      <PatientHeader patient={patient} onBack={onBack} />
      {headerActions && <div className="flex justify-end">{headerActions}</div>}

      <div className="flex flex-col md:flex-row gap-6 items-start mt-2">
        <Card className="w-full md:w-64 shrink-0 gap-4 p-4">
          {sidebarTop}
        </Card>

        <div className="flex-1 flex flex-col gap-6 min-w-0 w-full">
          {encounterBanner}
          {children}
        </div>
      </div>
    </PageContainer>
  )
}
