import { lazy, Suspense } from "react"
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import { useAuthStore } from "../shared/store/auth_store"
import { AbilityProvider } from "../shared/auth/AbilityContext"
import { AppSidebar } from "../shared/components/AppSidebar"
import { AppHeader } from "../shared/components/AppHeader"
import { RoleRoute } from "../shared/auth/RoleRoute"
import { getHomePathForRole, routeAllowedRoles } from "../shared/auth/route_access"
import { Spinner } from "../shared/components/ui/Spinner"

const Login = lazy(() => import("../modules/auth/Login").then((module) => ({ default: module.Login })))
const Patients = lazy(() => import("../modules/patients/Patients").then((module) => ({ default: module.Patients })))
const PatientDetails = lazy(() => import("../modules/patients/PatientDetails").then((module) => ({ default: module.PatientDetails })))
const EncounterDetail = lazy(() => import("../modules/patients/EncounterDetail").then((module) => ({ default: module.EncounterDetail })))
const Telemetry = lazy(() => import("../modules/telemetry/Telemetry").then((module) => ({ default: module.Telemetry })))
const Stats = lazy(() => import("../modules/analytics/Stats").then((module) => ({ default: module.Stats })))
const Staff = lazy(() => import("../modules/staff/Staff").then((module) => ({ default: module.Staff })))
const ExamAnalyzer = lazy(() => import("../modules/exam_analyzer/ExamAnalyzer").then((module) => ({ default: module.ExamAnalyzer })))
const AuditLogs = lazy(() => import("../modules/audit_logs/AuditLogs").then((module) => ({ default: module.AuditLogs })))
const PortalPage = lazy(() => import("../modules/portal/PortalPage").then((module) => ({ default: module.PortalPage })))
const Schedule = lazy(() => import("../modules/schedule/Schedule").then((module) => ({ default: module.Schedule })))
const ClinicalDashboard = lazy(() => import("../modules/analytics/ClinicalDashboard").then((module) => ({ default: module.ClinicalDashboard })))

export const AppRoutes = () => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  const role = useAuthStore((state) => state.role)

  return (
    <BrowserRouter>
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
                <AbilityProvider role={role}>
                  <div className="min-h-screen bg-background flex">
                    <AppSidebar />
                    <div className="flex-1 flex flex-col min-w-0">
                      <AppHeader />
                      <main className="flex-1 flex flex-col">
                        <Routes>
                          <Route path="/" element={<RoleRoute allowedRoles={routeAllowedRoles.patients}><Patients /></RoleRoute>} />
                          <Route path="/portal/*" element={<RoleRoute allowedRoles={routeAllowedRoles.portal}><PortalPage /></RoleRoute>} />
                          <Route path="/dashboard" element={<RoleRoute allowedRoles={routeAllowedRoles.dashboard}><ClinicalDashboard /></RoleRoute>} />
                          <Route path="/patients/:id" element={<RoleRoute allowedRoles={routeAllowedRoles.patients}><PatientDetails /></RoleRoute>} />
                          <Route path="/patients/:id/encounters/:encounterId" element={<RoleRoute allowedRoles={routeAllowedRoles.patients}><EncounterDetail /></RoleRoute>} />
                          <Route path="/telemetry" element={<RoleRoute allowedRoles={routeAllowedRoles.telemetry}><Telemetry /></RoleRoute>} />
                          <Route path="/analytics" element={<RoleRoute allowedRoles={routeAllowedRoles.analytics}><Stats /></RoleRoute>} />
                          <Route path="/staff" element={<RoleRoute allowedRoles={routeAllowedRoles.staff}><Staff /></RoleRoute>} />
                          <Route path="/schedule" element={<RoleRoute allowedRoles={routeAllowedRoles.schedule}><Schedule /></RoleRoute>} />
                          <Route path="/exam-analyzer" element={<RoleRoute allowedRoles={routeAllowedRoles.examAnalyzer}><ExamAnalyzer /></RoleRoute>} />
                          <Route path="/audit-logs" element={<RoleRoute allowedRoles={routeAllowedRoles.auditLogs}><AuditLogs /></RoleRoute>} />
                          <Route path="*" element={<Navigate to={getHomePathForRole(role)} replace />} />
                        </Routes>
                      </main>
                    </div>
                  </div>
                </AbilityProvider>
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
