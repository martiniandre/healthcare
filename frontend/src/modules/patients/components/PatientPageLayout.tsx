import type { ReactNode } from "react"
import { PatientHeader } from "./PatientHeader"

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
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full">
      <div className="flex flex-col xl:flex-row xl:items-start justify-between gap-4">
        <PatientHeader patient={patient} onBack={onBack} />
        {headerActions}
      </div>

      <div className="flex flex-col md:flex-row gap-6 items-start mt-2">
        <div className="w-full md:w-64 shrink-0 bg-white border border-border p-4 rounded-xl flex flex-col gap-4">
          {sidebarTop}
        </div>

        <div className="flex-1 flex flex-col gap-6 min-w-0 w-full">
          {encounterBanner}
          {children}
        </div>
      </div>
    </div>
  )
}
