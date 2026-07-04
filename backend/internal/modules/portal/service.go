package portal

import (
	"context"
	"errors"
)

var ErrPatientNotFound = errors.New("patient not found")

type Service interface {
	GetDashboard(ctx context.Context, fhirResourceID string) (*PortalDashboard, error)
	GetEncounters(ctx context.Context, fhirResourceID string) ([]PortalEncounter, error)
	GetObservations(ctx context.Context, fhirResourceID string) ([]PortalObservation, error)
	GetConditions(ctx context.Context, fhirResourceID string) ([]PortalCondition, error)
	GetMedications(ctx context.Context, fhirResourceID string) ([]PortalMedication, error)
	GetReports(ctx context.Context, fhirResourceID string) ([]PortalReport, error)
	GetImaging(ctx context.Context, fhirResourceID string) ([]PortalImaging, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (portalService *service) GetDashboard(ctx context.Context, fhirResourceID string) (*PortalDashboard, error) {
	patientInfo, patientErr := portalService.repo.GetPatient(ctx, fhirResourceID)
	if patientErr != nil {
		return nil, ErrPatientNotFound
	}

	encounters, _ := portalService.repo.GetEncountersByPatient(ctx, fhirResourceID)
	observations, _ := portalService.repo.GetObservationsByPatient(ctx, fhirResourceID)
	conditions, _ := portalService.repo.GetConditionsByPatient(ctx, fhirResourceID)
	medications, _ := portalService.repo.GetMedicationsByPatient(ctx, fhirResourceID)
	reports, _ := portalService.repo.GetReportsByPatient(ctx, fhirResourceID)
	imaging, _ := portalService.repo.GetImagingByPatient(ctx, fhirResourceID)

	activeConditions := make([]PortalCondition, 0, len(conditions))
	for _, condition := range conditions {
		if condition.ClinicalStatus == "active" {
			activeConditions = append(activeConditions, condition)
		}
	}

	activeMedications := make([]PortalMedication, 0, len(medications))
	for _, medication := range medications {
		if medication.Status == "active" {
			activeMedications = append(activeMedications, medication)
		}
	}

	upcomingEncounters := make([]PortalEncounter, 0, len(encounters))
	for _, encounter := range encounters {
		if encounter.Status == "planned" || encounter.Status == "arrived" {
			upcomingEncounters = append(upcomingEncounters, encounter)
		}
	}

	recentObservations := observations
	if len(recentObservations) > 20 {
		recentObservations = recentObservations[:20]
	}

	recentReports := reports
	if len(recentReports) > 10 {
		recentReports = recentReports[:10]
	}

	recentImaging := imaging
	if len(recentImaging) > 10 {
		recentImaging = recentImaging[:10]
	}

	return &PortalDashboard{
		PatientInfo:        *patientInfo,
		UpcomingEncounters: upcomingEncounters,
		RecentObservations: recentObservations,
		ActiveConditions:   activeConditions,
		ActiveMedications:  activeMedications,
		RecentReports:      recentReports,
		RecentImaging:      recentImaging,
	}, nil
}

func (portalService *service) GetEncounters(ctx context.Context, fhirResourceID string) ([]PortalEncounter, error) {
	return portalService.repo.GetEncountersByPatient(ctx, fhirResourceID)
}

func (portalService *service) GetObservations(ctx context.Context, fhirResourceID string) ([]PortalObservation, error) {
	return portalService.repo.GetObservationsByPatient(ctx, fhirResourceID)
}

func (portalService *service) GetConditions(ctx context.Context, fhirResourceID string) ([]PortalCondition, error) {
	return portalService.repo.GetConditionsByPatient(ctx, fhirResourceID)
}

func (portalService *service) GetMedications(ctx context.Context, fhirResourceID string) ([]PortalMedication, error) {
	return portalService.repo.GetMedicationsByPatient(ctx, fhirResourceID)
}

func (portalService *service) GetReports(ctx context.Context, fhirResourceID string) ([]PortalReport, error) {
	return portalService.repo.GetReportsByPatient(ctx, fhirResourceID)
}

func (portalService *service) GetImaging(ctx context.Context, fhirResourceID string) ([]PortalImaging, error) {
	return portalService.repo.GetImagingByPatient(ctx, fhirResourceID)
}
