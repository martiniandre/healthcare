import { useQuery } from "@tanstack/react-query"
import { portalApi } from "./api"

export const portalQueryKeys = {
  all: ["portal"] as const,
  dashboard: () => [...portalQueryKeys.all, "dashboard"] as const,
  encounters: () => [...portalQueryKeys.all, "encounters"] as const,
  observations: () => [...portalQueryKeys.all, "observations"] as const,
  conditions: () => [...portalQueryKeys.all, "conditions"] as const,
  medications: () => [...portalQueryKeys.all, "medications"] as const,
  reports: () => [...portalQueryKeys.all, "reports"] as const,
  imaging: () => [...portalQueryKeys.all, "imaging"] as const,
}

export const usePortalDashboardQuery = () => {
  return useQuery({
    queryKey: portalQueryKeys.dashboard(),
    queryFn: portalApi.getDashboard,
  })
}

export const usePortalEncountersQuery = () => {
  return useQuery({
    queryKey: portalQueryKeys.encounters(),
    queryFn: portalApi.getEncounters,
  })
}

export const usePortalObservationsQuery = () => {
  return useQuery({
    queryKey: portalQueryKeys.observations(),
    queryFn: portalApi.getObservations,
  })
}

export const usePortalConditionsQuery = () => {
  return useQuery({
    queryKey: portalQueryKeys.conditions(),
    queryFn: portalApi.getConditions,
  })
}

export const usePortalMedicationsQuery = () => {
  return useQuery({
    queryKey: portalQueryKeys.medications(),
    queryFn: portalApi.getMedications,
  })
}

export const usePortalReportsQuery = () => {
  return useQuery({
    queryKey: portalQueryKeys.reports(),
    queryFn: portalApi.getReports,
  })
}

export const usePortalImagingQuery = () => {
  return useQuery({
    queryKey: portalQueryKeys.imaging(),
    queryFn: portalApi.getImaging,
  })
}
