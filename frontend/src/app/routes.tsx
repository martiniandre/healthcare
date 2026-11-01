import { lazy, Suspense } from "react"
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import { useAuthStore } from "../shared/store/auth_store"
import { PermissionsGate } from "../shared/auth/PermissionsGate"
import { useCapabilities } from "../shared/auth/capabilities_context"
import { AppSidebar } from "../shared/components/AppSidebar"
import { AppHeader } from "../shared/components/AppHeader"
import { RoleRoute } from "../shared/auth/RoleRoute"
import { getHomePathForRole, routeRequiredCapability } from "../shared/auth/route_access"
import { Spinner } from "../shared/components/ui/spinner"

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
const DesignSystem = lazy(() => import("../modules/design_system/DesignSystem").then((module) => ({ default: module.DesignSystem })))

const AuthenticatedLayout = () => {
  const { capabilities } = useCapabilities()

  return (
    <div className="min-h-screen bg-background flex">
      <AppSidebar />
      <div className="flex-1 flex flex-col min-w-0">
        <AppHeader />
        <main className="flex-1 flex flex-col">
          <Routes>
            <Route path="/" element={<RoleRoute requiredCapability={routeRequiredCapability.patients}><Patients /></RoleRoute>} />
            <Route path="/portal/*" element={<RoleRoute requiredCapability={routeRequiredCapability.portal}><PortalPage /></RoleRoute>} />
            <Route path="/dashboard" element={<RoleRoute requiredCapability={routeRequiredCapability.dashboard}><ClinicalDashboard /></RoleRoute>} />
            <Route path="/patients/:id" element={<RoleRoute requiredCapability={routeRequiredCapability.patients}><PatientDetails /></RoleRoute>} />
            <Route path="/patients/:id/encounters/:encounterId" element={<RoleRoute requiredCapability={routeRequiredCapability.patients}><EncounterDetail /></RoleRoute>} />
            <Route path="/telemetry" element={<RoleRoute requiredCapability={routeRequiredCapability.telemetry}><Telemetry /></RoleRoute>} />
            <Route path="/analytics" element={<RoleRoute requiredCapability={routeRequiredCapability.analytics}><Stats /></RoleRoute>} />
            <Route path="/staff" element={<RoleRoute requiredCapability={routeRequiredCapability.staff}><Staff /></RoleRoute>} />
            <Route path="/schedule" element={<RoleRoute requiredCapability={routeRequiredCapability.schedule}><Schedule /></RoleRoute>} />
            <Route path="/exam-analyzer" element={<RoleRoute requiredCapability={routeRequiredCapability.examAnalyzer}><ExamAnalyzer /></RoleRoute>} />
            <Route path="/audit-logs" element={<RoleRoute requiredCapability={routeRequiredCapability.auditLogs}><AuditLogs /></RoleRoute>} />
            <Route path="/design-system" element={<RoleRoute requiredCapability={routeRequiredCapability.designSystem}><DesignSystem /></RoleRoute>} />
            <Route path="*" element={<Navigate to={getHomePathForRole(capabilities)} replace />} />
          </Routes>
        </main>
      </div>
    </div>
  )
}

export const AppRoutes = () => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)

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
                <PermissionsGate>
                  <AuthenticatedLayout />
                </PermissionsGate>
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
