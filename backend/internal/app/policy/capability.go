package policy

import (
	"sort"

	"github.com/healthcare/backend/internal/shared/role"
)

type Capability string

const (
	CapabilityAuthPermissions Capability = "auth:permissions"

	CapabilityPatientRead   Capability = "patient:read"
	CapabilityPatientCreate Capability = "patient:create"

	CapabilityStaffRead         Capability = "staff:read"
	CapabilityStaffCreate       Capability = "staff:create"
	CapabilityStaffUpdateStatus Capability = "staff:update_status"
	CapabilityStaffDelete       Capability = "staff:delete"

	CapabilityEncounterRead   Capability = "encounter:read"
	CapabilityEncounterCreate Capability = "encounter:create"
	CapabilityEncounterUpdate Capability = "encounter:update"
	CapabilityEncounterDelete Capability = "encounter:delete"

	CapabilityObservationRead   Capability = "observation:read"
	CapabilityObservationCreate Capability = "observation:create"
	CapabilityObservationUpdate Capability = "observation:update"
	CapabilityObservationDelete Capability = "observation:delete"

	CapabilityConditionRead   Capability = "condition:read"
	CapabilityConditionCreate Capability = "condition:create"
	CapabilityConditionUpdate Capability = "condition:update"
	CapabilityConditionDelete Capability = "condition:delete"

	CapabilityAllergyRead   Capability = "allergy:read"
	CapabilityAllergyCreate Capability = "allergy:create"
	CapabilityAllergyUpdate Capability = "allergy:update"
	CapabilityAllergyDelete Capability = "allergy:delete"

	CapabilityMedicationRead   Capability = "medication:read"
	CapabilityMedicationCreate Capability = "medication:create"
	CapabilityMedicationCheck  Capability = "medication:check"

	CapabilityReportRead   Capability = "report:read"
	CapabilityReportCreate Capability = "report:create"
	CapabilityReportUpdate Capability = "report:update"

	CapabilityImagingRead   Capability = "imaging:read"
	CapabilityImagingCreate Capability = "imaging:create"

	CapabilityExamAnalysisRead   Capability = "exam_analysis:read"
	CapabilityExamAnalysisCreate Capability = "exam_analysis:create"
	CapabilityExamAnalysisDelete Capability = "exam_analysis:delete"

	CapabilityTelemetryRead       Capability = "telemetry:read"
	CapabilityTelemetryReadBeds   Capability = "telemetry:read_beds"
	CapabilityTelemetryUnlock     Capability = "telemetry:unlock"
	CapabilityTelemetryUpdateBeds Capability = "telemetry:update_bed"

	CapabilityAppointmentRead    Capability = "appointment:read"
	CapabilityAppointmentCreate  Capability = "appointment:create"
	CapabilityAppointmentUpdate  Capability = "appointment:update"
	CapabilityAppointmentCancel  Capability = "appointment:cancel"
	CapabilityAppointmentReadOwn Capability = "appointment:read_own"

	CapabilityUnavailabilityRead   Capability = "unavailability:read"
	CapabilityUnavailabilityCreate Capability = "unavailability:create"
	CapabilityUnavailabilityDelete Capability = "unavailability:delete"

	CapabilityAnalyticsRead Capability = "analytics:read"
	CapabilityTimelineRead  Capability = "timeline:read"

	CapabilityAuditLogRead   Capability = "audit_log:read"
	CapabilityAuditLogCreate Capability = "audit_log:create"

	CapabilityNotificationRead      Capability = "notification:read"
	CapabilityNotificationUpdate    Capability = "notification:update"
	CapabilityNotificationEventRead Capability = "notification_event:read"

	CapabilityPortalRead       Capability = "portal:read"
	CapabilityDesignSystemRead Capability = "design_system:read"
)

var (
	everyRole = []role.Role{role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception, role.RolePatient}

	staffRoles      = []role.Role{role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception}
	clinicalRoles   = []role.Role{role.RoleAdmin, role.RoleDoctor, role.RoleNurse}
	prescriberRoles = []role.Role{role.RoleAdmin, role.RoleDoctor}
	frontDeskRoles  = []role.Role{role.RoleAdmin, role.RoleReception}
	adminOnlyRoles  = []role.Role{role.RoleAdmin}
	patientOnly     = []role.Role{role.RolePatient}
)

var capabilityRoles = map[Capability][]role.Role{
	CapabilityAuthPermissions: everyRole,

	CapabilityPatientRead:   staffRoles,
	CapabilityPatientCreate: frontDeskRoles,

	CapabilityStaffRead:         staffRoles,
	CapabilityStaffCreate:       adminOnlyRoles,
	CapabilityStaffUpdateStatus: clinicalRoles,
	CapabilityStaffDelete:       adminOnlyRoles,

	CapabilityEncounterRead:   clinicalRoles,
	CapabilityEncounterCreate: staffRoles,
	CapabilityEncounterUpdate: clinicalRoles,
	CapabilityEncounterDelete: clinicalRoles,

	CapabilityObservationRead:   clinicalRoles,
	CapabilityObservationCreate: clinicalRoles,
	CapabilityObservationUpdate: clinicalRoles,
	CapabilityObservationDelete: clinicalRoles,

	CapabilityConditionRead:   clinicalRoles,
	CapabilityConditionCreate: clinicalRoles,
	CapabilityConditionUpdate: clinicalRoles,
	CapabilityConditionDelete: clinicalRoles,

	CapabilityAllergyRead:   clinicalRoles,
	CapabilityAllergyCreate: clinicalRoles,
	CapabilityAllergyUpdate: clinicalRoles,
	CapabilityAllergyDelete: clinicalRoles,

	CapabilityMedicationRead:   clinicalRoles,
	CapabilityMedicationCreate: prescriberRoles,
	CapabilityMedicationCheck:  prescriberRoles,

	CapabilityReportRead:   clinicalRoles,
	CapabilityReportCreate: clinicalRoles,
	CapabilityReportUpdate: clinicalRoles,

	CapabilityImagingRead:   clinicalRoles,
	CapabilityImagingCreate: clinicalRoles,

	CapabilityExamAnalysisRead:   clinicalRoles,
	CapabilityExamAnalysisCreate: clinicalRoles,
	CapabilityExamAnalysisDelete: prescriberRoles,

	CapabilityTelemetryRead:       staffRoles,
	CapabilityTelemetryReadBeds:   clinicalRoles,
	CapabilityTelemetryUnlock:     clinicalRoles,
	CapabilityTelemetryUpdateBeds: clinicalRoles,

	CapabilityAppointmentRead:    staffRoles,
	CapabilityAppointmentCreate:  frontDeskRoles,
	CapabilityAppointmentUpdate:  frontDeskRoles,
	CapabilityAppointmentCancel:  frontDeskRoles,
	CapabilityAppointmentReadOwn: patientOnly,

	CapabilityUnavailabilityRead:   prescriberRoles,
	CapabilityUnavailabilityCreate: prescriberRoles,
	CapabilityUnavailabilityDelete: prescriberRoles,

	CapabilityAnalyticsRead: clinicalRoles,
	CapabilityTimelineRead:  clinicalRoles,

	CapabilityAuditLogRead:   adminOnlyRoles,
	CapabilityAuditLogCreate: adminOnlyRoles,

	CapabilityNotificationRead:      everyRole,
	CapabilityNotificationUpdate:    everyRole,
	CapabilityNotificationEventRead: staffRoles,

	CapabilityPortalRead:       patientOnly,
	CapabilityDesignSystemRead: staffRoles,
}

var methodCapabilities = map[string]Capability{
	"/audit_logs.v1.AuditLogsService/CreateAuditLog": CapabilityAuditLogCreate,
	"/audit_logs.v1.AuditLogsService/ListAuditLogs":  CapabilityAuditLogRead,

	"/telemetry.v1.TelemetryService/GetRooms":           CapabilityTelemetryRead,
	"/telemetry.v1.TelemetryService/UnlockRoom":         CapabilityTelemetryUnlock,
	"/telemetry.v1.TelemetryService/GetBeds":            CapabilityTelemetryReadBeds,
	"/telemetry.v1.TelemetryService/UpdateBedCondition": CapabilityTelemetryUpdateBeds,

	"/exam_analyzer.v1.ExamAnalyzerService/ListAnalyses":   CapabilityExamAnalysisRead,
	"/exam_analyzer.v1.ExamAnalyzerService/GetAnalysis":    CapabilityExamAnalysisRead,
	"/exam_analyzer.v1.ExamAnalyzerService/DeleteAnalysis": CapabilityExamAnalysisDelete,

	"/staff.v1.StaffService/CreateEmployee":     CapabilityStaffCreate,
	"/staff.v1.StaffService/GetEmployee":        CapabilityStaffRead,
	"/staff.v1.StaffService/ListEmployees":      CapabilityStaffRead,
	"/staff.v1.StaffService/DeactivateEmployee": CapabilityStaffDelete,

	"/patients.v1.PatientService/CreatePatient":        CapabilityPatientCreate,
	"/patients.v1.PatientService/GetPatient":           CapabilityPatientRead,
	"/patients.v1.PatientService/GetPatientByDocument": CapabilityPatientRead,
	"/patients.v1.PatientService/ListPatients":         CapabilityPatientRead,

	"/encounter.v1.EncounterService/CreateEncounter": CapabilityEncounterCreate,
	"/encounter.v1.EncounterService/GetEncounter":    CapabilityEncounterRead,
	"/encounter.v1.EncounterService/GetEncounters":   CapabilityEncounterRead,

	"/observation.v1.ObservationService/CreateObservation":      CapabilityObservationCreate,
	"/observation.v1.ObservationService/CreateObservationBatch": CapabilityObservationCreate,
	"/observation.v1.ObservationService/GetObservations":        CapabilityObservationRead,

	"/condition.v1.ConditionService/CreateCondition": CapabilityConditionCreate,
	"/condition.v1.ConditionService/GetConditions":   CapabilityConditionRead,

	"/allergy.v1.AllergyService/CreateAllergyIntolerance": CapabilityAllergyCreate,
	"/allergy.v1.AllergyService/GetAllergyIntolerances":   CapabilityAllergyRead,

	"/medication.v1.MedicationService/CreateMedicationRequest": CapabilityMedicationCreate,
	"/medication.v1.MedicationService/GetMedicationRequests":   CapabilityMedicationRead,

	"/diagnostic_report.v1.DiagnosticReportService/CreateDiagnosticReport": CapabilityReportCreate,
	"/diagnostic_report.v1.DiagnosticReportService/GetDiagnosticReports":   CapabilityReportRead,

	"/imaging.v1.ImagingService/UploadDICOM":         CapabilityImagingCreate,
	"/imaging.v1.ImagingService/GetImagingStudy":     CapabilityImagingRead,
	"/imaging.v1.ImagingService/ListImagingStudies":  CapabilityImagingRead,
	"/imaging.v1.ImagingService/GetDICOMDownloadURL": CapabilityImagingRead,
}

var routeCapabilities = map[string]Capability{
	"GET /api/v1/auth/permissions": CapabilityAuthPermissions,

	"GET /api/v1/telemetry/rooms":                   CapabilityTelemetryRead,
	"POST /api/v1/telemetry/rooms/{roomId}/unlock":  CapabilityTelemetryUnlock,
	"GET /api/v1/telemetry/rooms/{roomId}/beds":     CapabilityTelemetryReadBeds,
	"POST /api/v1/telemetry/beds/{bedId}/condition": CapabilityTelemetryUpdateBeds,

	"GET /api/v1/staff/employees":                       CapabilityStaffRead,
	"POST /api/v1/staff/employees":                      CapabilityStaffCreate,
	"GET /api/v1/staff/departments":                     CapabilityStaffRead,
	"PATCH /api/v1/staff/employees/{employeeId}/status": CapabilityStaffUpdateStatus,
	"DELETE /api/v1/staff/employees/{employeeId}":       CapabilityStaffDelete,

	"POST /api/v1/appointments":                        CapabilityAppointmentCreate,
	"GET /api/v1/appointments":                         CapabilityAppointmentRead,
	"GET /api/v1/appointments/my":                      CapabilityAppointmentReadOwn,
	"GET /api/v1/appointments/{appointmentId}":         CapabilityAppointmentRead,
	"POST /api/v1/appointments/{appointmentId}/cancel": CapabilityAppointmentCancel,
	"PUT /api/v1/appointments/{appointmentId}":         CapabilityAppointmentUpdate,

	"POST /api/v1/schedule/unavailability":                      CapabilityUnavailabilityCreate,
	"GET /api/v1/schedule/unavailability":                       CapabilityUnavailabilityRead,
	"DELETE /api/v1/schedule/unavailability/{unavailabilityId}": CapabilityUnavailabilityDelete,

	"GET /api/v1/portal/dashboard":    CapabilityPortalRead,
	"GET /api/v1/portal/encounters":   CapabilityPortalRead,
	"GET /api/v1/portal/observations": CapabilityPortalRead,
	"GET /api/v1/portal/conditions":   CapabilityPortalRead,
	"GET /api/v1/portal/medications":  CapabilityPortalRead,
	"GET /api/v1/portal/reports":      CapabilityPortalRead,
	"GET /api/v1/portal/imaging":      CapabilityPortalRead,

	"GET /api/v1/patients":                 CapabilityPatientRead,
	"POST /api/v1/patients":                CapabilityPatientCreate,
	"GET /api/v1/patients/{patientFhirId}": CapabilityPatientRead,

	"GET /api/v1/patients/{patientFhirId}/observations":            CapabilityObservationRead,
	"GET /api/v1/encounters/{encounterFhirId}/observations":        CapabilityObservationRead,
	"POST /api/v1/encounters/{encounterFhirId}/observations":       CapabilityObservationCreate,
	"POST /api/v1/encounters/{encounterFhirId}/observations/batch": CapabilityObservationCreate,
	"PUT /api/v1/observations/{observationFhirId}":                 CapabilityObservationUpdate,
	"DELETE /api/v1/observations/{observationFhirId}":              CapabilityObservationDelete,

	"GET /api/v1/notifications":                        CapabilityNotificationRead,
	"POST /api/v1/notifications/{notificationId}/read": CapabilityNotificationUpdate,
	"GET /api/v1/notifications/unread-count":           CapabilityNotificationRead,
	"GET /api/v1/notifications/stream":                 CapabilityNotificationRead,
	"GET /api/v1/notification-events":                  CapabilityNotificationEventRead,

	"GET /api/v1/encounters/{encounterFhirId}/medications":    CapabilityMedicationRead,
	"POST /api/v1/encounters/{encounterFhirId}/medications":   CapabilityMedicationCreate,
	"POST /api/v1/patients/{patientFhirId}/medications/check": CapabilityMedicationCheck,

	"GET /api/v1/patients/{patientFhirId}/studies":  CapabilityImagingRead,
	"POST /api/v1/patients/{patientFhirId}/studies": CapabilityImagingCreate,
	"GET /api/v1/studies/{studyId}":                 CapabilityImagingRead,

	"GET /api/v1/exam-analyses":                 CapabilityExamAnalysisRead,
	"POST /api/v1/exam-analyses":                CapabilityExamAnalysisCreate,
	"GET /api/v1/exam-analyses/{analysisId}":    CapabilityExamAnalysisRead,
	"DELETE /api/v1/exam-analyses/{analysisId}": CapabilityExamAnalysisDelete,

	"GET /api/v1/patients/{patientFhirId}/encounters":  CapabilityEncounterRead,
	"POST /api/v1/patients/{patientFhirId}/encounters": CapabilityEncounterCreate,
	"GET /api/v1/encounters/{encounterFhirId}":         CapabilityEncounterRead,
	"PUT /api/v1/encounters/{encounterFhirId}":         CapabilityEncounterUpdate,
	"DELETE /api/v1/encounters/{encounterFhirId}":      CapabilityEncounterDelete,

	"GET /api/v1/encounters/{encounterFhirId}/reports":      CapabilityReportRead,
	"POST /api/v1/encounters/{encounterFhirId}/reports":     CapabilityReportCreate,
	"PUT /api/v1/reports/{reportFhirId}":                    CapabilityReportUpdate,
	"GET /api/v1/reports/{reportFhirId}/versions":           CapabilityReportRead,
	"GET /api/v1/reports/{reportFhirId}/versions/{version}": CapabilityReportRead,

	"GET /api/v1/patients/{patientFhirId}/conditions":                      CapabilityConditionRead,
	"POST /api/v1/patients/{patientFhirId}/conditions":                     CapabilityConditionCreate,
	"PUT /api/v1/patients/{patientFhirId}/conditions/{conditionFhirId}":    CapabilityConditionUpdate,
	"DELETE /api/v1/patients/{patientFhirId}/conditions/{conditionFhirId}": CapabilityConditionDelete,

	"GET /api/v1/audit-logs":  CapabilityAuditLogRead,
	"POST /api/v1/audit-logs": CapabilityAuditLogCreate,

	"GET /api/v1/patients/{patientFhirId}/allergies":                    CapabilityAllergyRead,
	"POST /api/v1/patients/{patientFhirId}/allergies":                   CapabilityAllergyCreate,
	"PUT /api/v1/patients/{patientFhirId}/allergies/{allergyFhirId}":    CapabilityAllergyUpdate,
	"DELETE /api/v1/patients/{patientFhirId}/allergies/{allergyFhirId}": CapabilityAllergyDelete,

	"GET /api/v1/patients/{patientFhirId}/timeline": CapabilityTimelineRead,

	"GET /api/v1/analytics":                                    CapabilityAnalyticsRead,
	"GET /api/v1/analytics/dashboard":                          CapabilityAnalyticsRead,
	"GET /api/v1/analytics/dashboard/consultations-per-doctor": CapabilityAnalyticsRead,
	"GET /api/v1/analytics/dashboard/occupancy-rate":           CapabilityAnalyticsRead,
	"GET /api/v1/analytics/dashboard/avg-wait-time":            CapabilityAnalyticsRead,
	"GET /api/v1/analytics/dashboard/top-diagnoses":            CapabilityAnalyticsRead,
}

func deriveRoles(capabilityByKey map[string]Capability) map[string][]role.Role {
	derivedRoles := make(map[string][]role.Role, len(capabilityByKey))
	for key, capability := range capabilityByKey {
		derivedRoles[key] = capabilityRoles[capability]
	}
	return derivedRoles
}

func CapabilitiesForRole(callerRole role.Role) []Capability {
	granted := make([]Capability, 0, len(capabilityRoles))
	for capability, allowedRoles := range capabilityRoles {
		for _, allowedRole := range allowedRoles {
			if allowedRole == callerRole {
				granted = append(granted, capability)
				break
			}
		}
	}
	sort.Slice(granted, func(firstIndex, secondIndex int) bool {
		return granted[firstIndex] < granted[secondIndex]
	})
	return granted
}
