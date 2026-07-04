import { http } from "../../shared/utils/http"
import type {
  PortalDashboard,
  PortalEncounter,
  PortalObservation,
  PortalCondition,
  PortalMedication,
  PortalReport,
  PortalImaging,
} from "./types"

export const portalApi = {
  getDashboard: async (): Promise<PortalDashboard> => {
    return http.get<PortalDashboard>("/portal/dashboard")
  },

  getEncounters: async (): Promise<PortalEncounter[]> => {
    return http.get<PortalEncounter[]>("/portal/encounters")
  },

  getObservations: async (): Promise<PortalObservation[]> => {
    return http.get<PortalObservation[]>("/portal/observations")
  },

  getConditions: async (): Promise<PortalCondition[]> => {
    return http.get<PortalCondition[]>("/portal/conditions")
  },

  getMedications: async (): Promise<PortalMedication[]> => {
    return http.get<PortalMedication[]>("/portal/medications")
  },

  getReports: async (): Promise<PortalReport[]> => {
    return http.get<PortalReport[]>("/portal/reports")
  },

  getImaging: async (): Promise<PortalImaging[]> => {
    return http.get<PortalImaging[]>("/portal/imaging")
  },
}
