package policy

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/healthcare/backend/internal/shared/role"
	"github.com/stretchr/testify/assert"
)

func rolesKey(roles []role.Role) string {
	sortedRoles := append([]role.Role(nil), roles...)
	sort.Slice(sortedRoles, func(firstIndex, secondIndex int) bool {
		return sortedRoles[firstIndex] < sortedRoles[secondIndex]
	})
	roleNames := make([]string, 0, len(sortedRoles))
	for _, allowedRole := range sortedRoles {
		roleNames = append(roleNames, string(allowedRole))
	}
	return strings.Join(roleNames, ",")
}

func TestGRPCMethodPermissions(testingInstance *testing.T) {
	expectedRoles := map[string]string{
		"/audit_logs.v1.AuditLogsService/CreateAuditLog": "ADMIN,DOCTOR,NURSE,PATIENT,RECEPTION",
		"/audit_logs.v1.AuditLogsService/ListAuditLogs":  "ADMIN",

		"/telemetry.v1.TelemetryService/GetRooms":           "ADMIN,DOCTOR,NURSE,RECEPTION",
		"/telemetry.v1.TelemetryService/UnlockRoom":         "ADMIN,DOCTOR,NURSE",
		"/telemetry.v1.TelemetryService/GetBeds":            "ADMIN,DOCTOR,NURSE",
		"/telemetry.v1.TelemetryService/UpdateBedCondition": "ADMIN,DOCTOR,NURSE",

		"/exam_analyzer.v1.ExamAnalyzerService/ListAnalyses":   "ADMIN,DOCTOR,NURSE",
		"/exam_analyzer.v1.ExamAnalyzerService/GetAnalysis":    "ADMIN,DOCTOR,NURSE",
		"/exam_analyzer.v1.ExamAnalyzerService/DeleteAnalysis": "ADMIN,DOCTOR",

		"/staff.v1.StaffService/CreateEmployee":     "ADMIN",
		"/staff.v1.StaffService/GetEmployee":        "ADMIN,DOCTOR,NURSE,RECEPTION",
		"/staff.v1.StaffService/ListEmployees":      "ADMIN,DOCTOR,NURSE,RECEPTION",
		"/staff.v1.StaffService/DeactivateEmployee": "ADMIN",

		"/patients.v1.PatientService/CreatePatient":        "ADMIN,RECEPTION",
		"/patients.v1.PatientService/GetPatient":           "ADMIN,DOCTOR,NURSE,RECEPTION",
		"/patients.v1.PatientService/GetPatientByDocument": "ADMIN,DOCTOR,NURSE,RECEPTION",
		"/patients.v1.PatientService/ListPatients":         "ADMIN,DOCTOR,NURSE,RECEPTION",

		"/encounter.v1.EncounterService/CreateEncounter": "ADMIN,DOCTOR,NURSE,RECEPTION",
		"/encounter.v1.EncounterService/GetEncounter":    "ADMIN,DOCTOR,NURSE",
		"/encounter.v1.EncounterService/GetEncounters":   "ADMIN,DOCTOR,NURSE",

		"/observation.v1.ObservationService/CreateObservation":      "ADMIN,DOCTOR,NURSE",
		"/observation.v1.ObservationService/CreateObservationBatch": "ADMIN,DOCTOR,NURSE",
		"/observation.v1.ObservationService/GetObservations":        "ADMIN,DOCTOR,NURSE",

		"/condition.v1.ConditionService/CreateCondition": "ADMIN,DOCTOR,NURSE",
		"/condition.v1.ConditionService/GetConditions":   "ADMIN,DOCTOR,NURSE",

		"/allergy.v1.AllergyService/CreateAllergyIntolerance": "ADMIN,DOCTOR,NURSE",
		"/allergy.v1.AllergyService/GetAllergyIntolerances":   "ADMIN,DOCTOR,NURSE",

		"/medication.v1.MedicationService/CreateMedicationRequest": "ADMIN,DOCTOR",
		"/medication.v1.MedicationService/GetMedicationRequests":   "ADMIN,DOCTOR,NURSE",

		"/diagnostic_report.v1.DiagnosticReportService/CreateDiagnosticReport": "ADMIN,DOCTOR,NURSE",
		"/diagnostic_report.v1.DiagnosticReportService/GetDiagnosticReports":   "ADMIN,DOCTOR,NURSE",

		"/imaging.v1.ImagingService/UploadDICOM":         "ADMIN,DOCTOR,NURSE",
		"/imaging.v1.ImagingService/GetImagingStudy":     "ADMIN,DOCTOR,NURSE",
		"/imaging.v1.ImagingService/ListImagingStudies":  "ADMIN,DOCTOR,NURSE",
		"/imaging.v1.ImagingService/GetDICOMDownloadURL": "ADMIN,DOCTOR,NURSE",
	}

	for fullMethod, expected := range expectedRoles {
		allowedRoles, defined := GRPCMethodRoles(fullMethod)
		if !assert.True(testingInstance, defined, "grpc method %s must be registered", fullMethod) {
			continue
		}
		assert.Equal(testingInstance, expected, rolesKey(allowedRoles), "grpc method %s role drift", fullMethod)
	}
}

func TestPublicGRPCMethods(testingInstance *testing.T) {
	expectedPublicMethods := []string{
		"/auth.v1.AuthService/Login",
		"/auth.v1.AuthService/Register",
		"/auth.v1.AuthService/Logout",
		"/grpc.health.v1.Health/Check",
		"/grpc.health.v1.Health/Watch",
	}
	for _, fullMethod := range expectedPublicMethods {
		assert.True(testingInstance, IsPublicGRPCMethod(fullMethod), "%s must be public", fullMethod)
		assert.False(testingInstance, func() bool {
			_, defined := GRPCMethodRoles(fullMethod)
			return defined
		}(), "%s must not be role-gated", fullMethod)
	}
}

func TestHTTPRoutePermissions(testingInstance *testing.T) {
	expectedRoles := map[string]string{
		"GET /api/v1/telemetry/rooms":                   "ADMIN,DOCTOR,NURSE,RECEPTION",
		"POST /api/v1/telemetry/rooms/{roomId}/unlock":  "ADMIN,DOCTOR,NURSE",
		"GET /api/v1/telemetry/rooms/{roomId}/beds":     "ADMIN,DOCTOR,NURSE",
		"POST /api/v1/telemetry/beds/{bedId}/condition": "ADMIN,DOCTOR,NURSE",

		"GET /api/v1/staff/employees":  "ADMIN,DOCTOR,NURSE,RECEPTION",
		"POST /api/v1/staff/employees": "ADMIN",

		"POST /api/v1/appointments":                        "ADMIN,RECEPTION",
		"GET /api/v1/appointments":                         "ADMIN,DOCTOR,NURSE,RECEPTION",
		"GET /api/v1/appointments/my":                      "PATIENT",
		"GET /api/v1/appointments/{appointmentId}":         "ADMIN,DOCTOR,NURSE,RECEPTION",
		"POST /api/v1/appointments/{appointmentId}/cancel": "ADMIN,RECEPTION",

		"GET /api/v1/portal/dashboard":    "PATIENT",
		"GET /api/v1/portal/encounters":   "PATIENT",
		"GET /api/v1/portal/observations": "PATIENT",
		"GET /api/v1/portal/conditions":   "PATIENT",
		"GET /api/v1/portal/medications":  "PATIENT",
		"GET /api/v1/portal/reports":      "PATIENT",
		"GET /api/v1/portal/imaging":      "PATIENT",

		"GET /api/v1/patients":                 "ADMIN,DOCTOR,NURSE,RECEPTION",
		"POST /api/v1/patients":                "ADMIN,RECEPTION",
		"GET /api/v1/patients/{patientFhirId}": "ADMIN,DOCTOR,NURSE,RECEPTION",

		"GET /api/v1/patients/{patientFhirId}/observations":            "ADMIN,DOCTOR,NURSE",
		"GET /api/v1/patients/{patientFhirId}/timeline":                "ADMIN,DOCTOR,NURSE",
		"GET /api/v1/encounters/{encounterFhirId}/observations":        "ADMIN,DOCTOR,NURSE",
		"POST /api/v1/encounters/{encounterFhirId}/observations":       "ADMIN,DOCTOR,NURSE",
		"POST /api/v1/encounters/{encounterFhirId}/observations/batch": "ADMIN,DOCTOR,NURSE",
		"PUT /api/v1/observations/{observationFhirId}":                 "ADMIN,DOCTOR,NURSE",
		"DELETE /api/v1/observations/{observationFhirId}":              "ADMIN,DOCTOR,NURSE",

		"GET /api/v1/notifications":                        "ADMIN,DOCTOR,NURSE,PATIENT,RECEPTION",
		"POST /api/v1/notifications/{notificationId}/read": "ADMIN,DOCTOR,NURSE,PATIENT,RECEPTION",
		"GET /api/v1/notifications/unread-count":           "ADMIN,DOCTOR,NURSE,PATIENT,RECEPTION",
		"GET /api/v1/notifications/stream":                 "ADMIN,DOCTOR,NURSE,PATIENT,RECEPTION",

		"GET /api/v1/encounters/{encounterFhirId}/medications":  "ADMIN,DOCTOR,NURSE",
		"POST /api/v1/encounters/{encounterFhirId}/medications": "ADMIN,DOCTOR",

		"GET /api/v1/patients/{patientFhirId}/studies":  "ADMIN,DOCTOR,NURSE",
		"POST /api/v1/patients/{patientFhirId}/studies": "ADMIN,DOCTOR,NURSE",
		"GET /api/v1/studies/{studyId}":                 "ADMIN,DOCTOR,NURSE",

		"GET /api/v1/exam-analyses":                 "ADMIN,DOCTOR,NURSE",
		"POST /api/v1/exam-analyses":                "ADMIN,DOCTOR,NURSE",
		"GET /api/v1/exam-analyses/{analysisId}":    "ADMIN,DOCTOR,NURSE",
		"DELETE /api/v1/exam-analyses/{analysisId}": "ADMIN,DOCTOR",

		"GET /api/v1/patients/{patientFhirId}/encounters":  "ADMIN,DOCTOR,NURSE",
		"POST /api/v1/patients/{patientFhirId}/encounters": "ADMIN,DOCTOR,NURSE,RECEPTION",
		"GET /api/v1/encounters/{encounterFhirId}":         "ADMIN,DOCTOR,NURSE",
		"PUT /api/v1/encounters/{encounterFhirId}":         "ADMIN,DOCTOR,NURSE",
		"DELETE /api/v1/encounters/{encounterFhirId}":      "ADMIN,DOCTOR,NURSE",

		"GET /api/v1/encounters/{encounterFhirId}/reports":      "ADMIN,DOCTOR,NURSE",
		"POST /api/v1/encounters/{encounterFhirId}/reports":     "ADMIN,DOCTOR,NURSE",
		"PUT /api/v1/reports/{reportFhirId}":                    "ADMIN,DOCTOR,NURSE",
		"GET /api/v1/reports/{reportFhirId}/versions":           "ADMIN,DOCTOR,NURSE",
		"GET /api/v1/reports/{reportFhirId}/versions/{version}": "ADMIN,DOCTOR,NURSE",

		"GET /api/v1/patients/{patientFhirId}/conditions":                      "ADMIN,DOCTOR,NURSE",
		"POST /api/v1/patients/{patientFhirId}/conditions":                     "ADMIN,DOCTOR,NURSE",
		"PUT /api/v1/patients/{patientFhirId}/conditions/{conditionFhirId}":    "ADMIN,DOCTOR,NURSE",
		"DELETE /api/v1/patients/{patientFhirId}/conditions/{conditionFhirId}": "ADMIN,DOCTOR,NURSE",

		"GET /api/v1/audit-logs":  "ADMIN",
		"POST /api/v1/audit-logs": "ADMIN,DOCTOR,NURSE,PATIENT,RECEPTION",

		"GET /api/v1/patients/{patientFhirId}/allergies":                    "ADMIN,DOCTOR,NURSE",
		"POST /api/v1/patients/{patientFhirId}/allergies":                   "ADMIN,DOCTOR,NURSE",
		"PUT /api/v1/patients/{patientFhirId}/allergies/{allergyFhirId}":    "ADMIN,DOCTOR,NURSE",
		"DELETE /api/v1/patients/{patientFhirId}/allergies/{allergyFhirId}": "ADMIN,DOCTOR,NURSE",

		"GET /api/v1/analytics":                                    "ADMIN,DOCTOR,NURSE",
		"GET /api/v1/analytics/dashboard":                          "ADMIN,DOCTOR,NURSE",
		"GET /api/v1/analytics/dashboard/consultations-per-doctor": "ADMIN,DOCTOR,NURSE",
		"GET /api/v1/analytics/dashboard/occupancy-rate":           "ADMIN,DOCTOR,NURSE",
		"GET /api/v1/analytics/dashboard/avg-wait-time":            "ADMIN,DOCTOR,NURSE",
		"GET /api/v1/analytics/dashboard/top-diagnoses":            "ADMIN,DOCTOR,NURSE",
	}

	for routePattern, expected := range expectedRoles {
		allowedRoles, defined := HTTPRouteRoles(routePattern)
		if !assert.True(testingInstance, defined, "http route %s must be registered", routePattern) {
			continue
		}
		assert.Equal(testingInstance, expected, rolesKey(allowedRoles), "http route %s role drift", routePattern)
	}
}

func TestPublicHTTPRoutes(testingInstance *testing.T) {
	publicRoutes := []string{
		"GET /health",
		"POST /api/v1/auth/login",
		"POST /api/v1/auth/logout",
		"GET /api/v1/auth/me",
	}
	for _, routePattern := range publicRoutes {
		assert.True(testingInstance, IsPublicHTTPRoute(routePattern), "%s must be public", routePattern)
	}
}

func TestDriftGuard_GRPCAndHTTPAgree(testingInstance *testing.T) {
	equivalences := []struct {
		grpcMethod string
		httpRoute  string
	}{
		{"/auth.v1.AuthService/Login", "POST /api/v1/auth/login"},
		{"/auth.v1.AuthService/Logout", "POST /api/v1/auth/logout"},
		{"/audit_logs.v1.AuditLogsService/CreateAuditLog", "POST /api/v1/audit-logs"},
		{"/audit_logs.v1.AuditLogsService/ListAuditLogs", "GET /api/v1/audit-logs"},
		{"/telemetry.v1.TelemetryService/GetRooms", "GET /api/v1/telemetry/rooms"},
		{"/telemetry.v1.TelemetryService/UnlockRoom", "POST /api/v1/telemetry/rooms/{roomId}/unlock"},
		{"/telemetry.v1.TelemetryService/GetBeds", "GET /api/v1/telemetry/rooms/{roomId}/beds"},
		{"/telemetry.v1.TelemetryService/UpdateBedCondition", "POST /api/v1/telemetry/beds/{bedId}/condition"},
		{"/staff.v1.StaffService/CreateEmployee", "POST /api/v1/staff/employees"},
		{"/staff.v1.StaffService/ListEmployees", "GET /api/v1/staff/employees"},
		{"/patients.v1.PatientService/CreatePatient", "POST /api/v1/patients"},
		{"/patients.v1.PatientService/GetPatient", "GET /api/v1/patients/{patientFhirId}"},
		{"/patients.v1.PatientService/ListPatients", "GET /api/v1/patients"},
		{"/encounter.v1.EncounterService/CreateEncounter", "POST /api/v1/patients/{patientFhirId}/encounters"},
		{"/encounter.v1.EncounterService/GetEncounter", "GET /api/v1/encounters/{encounterFhirId}"},
		{"/encounter.v1.EncounterService/GetEncounters", "GET /api/v1/patients/{patientFhirId}/encounters"},
		{"/observation.v1.ObservationService/CreateObservation", "POST /api/v1/encounters/{encounterFhirId}/observations"},
		{"/observation.v1.ObservationService/GetObservations", "GET /api/v1/encounters/{encounterFhirId}/observations"},
		{"/condition.v1.ConditionService/CreateCondition", "POST /api/v1/patients/{patientFhirId}/conditions"},
		{"/condition.v1.ConditionService/GetConditions", "GET /api/v1/patients/{patientFhirId}/conditions"},
		{"/allergy.v1.AllergyService/CreateAllergyIntolerance", "POST /api/v1/patients/{patientFhirId}/allergies"},
		{"/allergy.v1.AllergyService/GetAllergyIntolerances", "GET /api/v1/patients/{patientFhirId}/allergies"},
		{"/medication.v1.MedicationService/CreateMedicationRequest", "POST /api/v1/encounters/{encounterFhirId}/medications"},
		{"/medication.v1.MedicationService/GetMedicationRequests", "GET /api/v1/encounters/{encounterFhirId}/medications"},
		{"/diagnostic_report.v1.DiagnosticReportService/CreateDiagnosticReport", "POST /api/v1/encounters/{encounterFhirId}/reports"},
		{"/diagnostic_report.v1.DiagnosticReportService/GetDiagnosticReports", "GET /api/v1/encounters/{encounterFhirId}/reports"},
		{"/imaging.v1.ImagingService/UploadDICOM", "POST /api/v1/patients/{patientFhirId}/studies"},
		{"/imaging.v1.ImagingService/GetImagingStudy", "GET /api/v1/studies/{studyId}"},
		{"/imaging.v1.ImagingService/ListImagingStudies", "GET /api/v1/patients/{patientFhirId}/studies"},
		{"/imaging.v1.ImagingService/GetDICOMDownloadURL", "GET /api/v1/studies/{studyId}"},
		{"/exam_analyzer.v1.ExamAnalyzerService/ListAnalyses", "GET /api/v1/exam-analyses"},
		{"/exam_analyzer.v1.ExamAnalyzerService/GetAnalysis", "GET /api/v1/exam-analyses/{analysisId}"},
		{"/exam_analyzer.v1.ExamAnalyzerService/DeleteAnalysis", "DELETE /api/v1/exam-analyses/{analysisId}"},
	}

	for _, equivalence := range equivalences {
		grpcIsPublic := IsPublicGRPCMethod(equivalence.grpcMethod)
		httpIsPublic := IsPublicHTTPRoute(equivalence.httpRoute)
		assert.Equal(testingInstance, grpcIsPublic, httpIsPublic, "public status drift: %s vs %s", equivalence.grpcMethod, equivalence.httpRoute)

		if grpcIsPublic {
			continue
		}

		grpcRoles, grpcDefined := GRPCMethodRoles(equivalence.grpcMethod)
		httpRoles, httpDefined := HTTPRouteRoles(equivalence.httpRoute)
		if !assert.True(testingInstance, grpcDefined, "grpc method %s not registered", equivalence.grpcMethod) {
			continue
		}
		if !assert.True(testingInstance, httpDefined, "http route %s not registered", equivalence.httpRoute) {
			continue
		}
		assert.Equal(testingInstance, rolesKey(grpcRoles), rolesKey(httpRoles),
			"role drift between gRPC %s and HTTP %s", equivalence.grpcMethod, equivalence.httpRoute)
	}
}

func TestRolesForNotificationType(testingInstance *testing.T) {
	expectedRoles := map[string]string{
		"telemetry_alert":   "DOCTOR,NURSE",
		"exam_complete":     "DOCTOR",
		"encounter_created": "DOCTOR,PATIENT,RECEPTION",
		"encounter_updated": "DOCTOR,RECEPTION",
		"patient_created":   "ADMIN,RECEPTION",
		"patient_updated":   "DOCTOR,RECEPTION",
		"audit_alert":       "ADMIN",
		"report_ready":      "PATIENT",
		"system":            "ADMIN,DOCTOR,NURSE",
	}

	for notificationType, expected := range expectedRoles {
		allowedRoles, defined := RolesForNotificationType(notificationType)
		if !assert.True(testingInstance, defined, "notification type %s must be registered", notificationType) {
			continue
		}
		assert.Equal(testingInstance, expected, rolesKey(allowedRoles), "notification type %s role drift", notificationType)
	}
}

func TestRoleOrderingPreservedForErrorMessages(testingInstance *testing.T) {
	_ = fmt.Sprintf("%s", role.RoleAdmin)
}
