import type { ReactNode } from "react"
import { PatientHeader } from "./PatientHeader"
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
    <PageContainer className="gap-5 md:gap-6">
      <PatientHeader patient={patient} onBack={onBack} actions={headerActions} />

      <nav className="-mx-4 overflow-x-auto px-4 sm:mx-0 sm:px-0">{sidebarTop}</nav>

      <div className="flex min-w-0 flex-col gap-6">
        {encounterBanner}
        {children}
      </div>
    </PageContainer>
  )
}
