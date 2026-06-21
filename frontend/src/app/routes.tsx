import { lazy, Suspense, useEffect } from "react"
import { HashRouter, Routes, Route, Navigate, useLocation, useParams, useNavigate } from "react-router-dom"
import { useAuthStore } from "../shared/store/auth_store"
import { AppSidebar } from "../shared/components/AppSidebar"
import { AppHeader } from "../shared/components/AppHeader"
import { Spinner } from "../shared/components/ui/Spinner"
import { auditLogsApi } from "../modules/audit_logs/api"

const Login = lazy(() => import("../modules/auth/Login").then((module) => ({ default: module.Login })))
const Patients = lazy(() => import("../modules/patients/Patients").then((module) => ({ default: module.Patients })))
const PatientDetails = lazy(() => import("../modules/patients/PatientDetails").then((module) => ({ default: module.PatientDetails })))
const ImagingWorkspace = lazy(() => import("../modules/imaging/ImagingWorkspace").then((module) => ({ default: module.ImagingWorkspace })))
const Telemetry = lazy(() => import("../modules/telemetry/Telemetry").then((module) => ({ default: module.Telemetry })))
const Stats = lazy(() => import("../modules/stats/Stats").then((module) => ({ default: module.Stats })))
const Staff = lazy(() => import("../modules/staff/Staff").then((module) => ({ default: module.Staff })))
const ExamAnalyzer = lazy(() => import("../modules/exam_analyzer/ExamAnalyzer").then((module) => ({ default: module.ExamAnalyzer })))
const AuditLogs = lazy(() => import("../modules/audit_logs/AuditLogs").then((module) => ({ default: module.AuditLogs })))

const PageViewLogger = () => {
  const location = useLocation()
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)

  useEffect(() => {
    if (isAuthenticated) {
      auditLogsApi.createAuditLog({
        action: "PAGE_VIEW",
        details: `Viewed page: ${location.pathname}${location.search}`,
        status: "SUCCESS",
      }).catch((logError) => {
        console.error("Failed to log page view", logError)
      })
    }
  }, [location.pathname, location.search, isAuthenticated])

  return null
}

const ImagingWorkspaceRouteWrapper = () => {
  const { studyId } = useParams()
  const navigate = useNavigate()
  return (
    <ImagingWorkspace
      studyId={studyId || ""}
      onBack={() => navigate(-1)}
    />
  )
}

export const AppRoutes = () => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  const role = useAuthStore((state) => state.role)

  return (
    <HashRouter>
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
                        <Route path="/" element={<Patients />} />
                        <Route path="/patients/:id" element={<PatientDetails />} />
                        <Route path="/imaging" element={<ImagingWorkspaceRouteWrapper />} />
                        <Route path="/imaging/:studyId" element={<ImagingWorkspaceRouteWrapper />} />
                        <Route path="/telemetry" element={<Telemetry />} />
                        <Route path="/stats" element={<Stats />} />
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
    </HashRouter>
  )
}
