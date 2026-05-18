import { HashRouter, Routes, Route, Navigate } from "react-router-dom"
import { useAuthStore } from "../shared/store/auth_store"
import { Login } from "../modules/auth/Login"
import { Patients } from "../modules/patients/Patients"
import { PatientDetails } from "../modules/patients/PatientDetails"
import { ImagingWorkspace } from "../modules/imaging/ImagingWorkspace"
import { Telemetry } from "../modules/telemetry/Telemetry"
import { Stats } from "../modules/stats/Stats"
import { Staff } from "../modules/staff/Staff"
import { AppSidebar } from "../shared/components/AppSidebar"
import { AppHeader } from "../shared/components/AppHeader"

export const AppRoutes = () => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)

  return (
    <HashRouter>
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
                      <Route path="/imaging/:studyId" element={<ImagingWorkspace />} />
                      <Route path="/telemetry" element={<Telemetry />} />
                      <Route path="/stats" element={<Stats />} />
                      <Route path="/staff" element={<Staff />} />
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
    </HashRouter>
  )
}
