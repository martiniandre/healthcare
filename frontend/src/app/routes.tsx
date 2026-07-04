import { lazy, Suspense } from "react"
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import { useAuthStore } from "../shared/store/auth_store"
import { AppSidebar } from "../shared/components/AppSidebar"
import { AppHeader } from "../shared/components/AppHeader"
import { Spinner } from "../shared/components/ui/Spinner"
import { usePageViewLogger } from "../shared/hooks/usePageViewLogger"

const Login = lazy(() => import("../modules/auth/Login").then((module) => ({ default: module.Login })))
const Patients = lazy(() => import("../modules/patients/Patients").then((module) => ({ default: module.Patients })))
const PatientDetails = lazy(() => import("../modules/patients/PatientDetails").then((module) => ({ default: module.PatientDetails })))
const Telemetry = lazy(() => import("../modules/telemetry/Telemetry").then((module) => ({ default: module.Telemetry })))
const Stats = lazy(() => import("../modules/analytics/Stats").then((module) => ({ default: module.Stats })))
const Staff = lazy(() => import("../modules/staff/Staff").then((module) => ({ default: module.Staff })))
const ExamAnalyzer = lazy(() => import("../modules/exam_analyzer/ExamAnalyzer").then((module) => ({ default: module.ExamAnalyzer })))
const AuditLogs = lazy(() => import("../modules/audit_logs/AuditLogs").then((module) => ({ default: module.AuditLogs })))
const PortalPage = lazy(() => import("../modules/portal/PortalPage").then((module) => ({ default: module.PortalPage })))
const ClinicalDashboard = lazy(() => import("../modules/analytics/ClinicalDashboard").then((module) => ({ default: module.ClinicalDashboard })))

const PageViewLogger = () => {
  usePageViewLogger()
  return null
}

export const AppRoutes = () => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  const role = useAuthStore((state) => state.role)

  return (
    <BrowserRouter>
      <PageViewLogger />
      <Suspense
        fallback={
          <div className="min-h-screen bg-background flex items-center justify-center">
            <Spinner className="w-8 h-8 text-primary" />
          </div>
        }
      >
        <Routes>
          <Route
            path="/login"
            element={!isAuthenticated ? <Login /> : <Navigate to="/" replace />}
          />

          <Route
            path="/*"
            element={
              isAuthenticated ? (
                <div className="min-h-screen bg-background flex">
                  <AppSidebar />
                  <div className="flex-1 flex flex-col min-w-0">
                    <AppHeader />
                    <main className="flex-1 flex flex-col">
                      <Routes>
                        <Route path="/" element={role === "PATIENT" ? <Navigate to="/portal" replace /> : <Patients />} />
                        <Route path="/portal/*" element={<PortalPage />} />
                        <Route path="/dashboard" element={<ClinicalDashboard />} />
                        <Route path="/patients/:id" element={<PatientDetails />} />
                        <Route path="/telemetry" element={<Telemetry />} />
                        <Route path="/analytics" element={<Stats />} />
                        <Route path="/staff" element={<Staff />} />
                        <Route path="/exam-analyzer" element={<ExamAnalyzer />} />
                        <Route
                          path="/audit-logs"
                          element={role === "ADMIN" ? <AuditLogs /> : <Navigate to="/" replace />}
                        />
                        <Route path="*" element={<Navigate to="/" replace />} />
                      </Routes>
                    </main>
                  </div>
                </div>
              ) : (
                <Navigate to="/login" replace />
              )
            }
          />
        </Routes>
      </Suspense>
    </BrowserRouter>
  )
}
