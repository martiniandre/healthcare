package interceptor

import "github.com/healthcare/backend/internal/shared/role"

var publicMethods = map[string]bool{
	"/auth.v1.AuthService/Login":    true,
	"/auth.v1.AuthService/Register": true,
	"/auth.v1.AuthService/Logout":   true,
	"/grpc.health.v1.Health/Check":  true,
	"/grpc.health.v1.Health/Watch":  true,
	"/audit_logs.v1.AuditLogsService/CreateAuditLog": true,
}

var methodPermissions = map[string][]role.Role{
	"/audit_logs.v1.AuditLogsService/ListAuditLogs":     {role.RoleAdmin},

	"/telemetry.v1.TelemetryService/GetRooms":           {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"/telemetry.v1.TelemetryService/UnlockRoom":         {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/telemetry.v1.TelemetryService/GetBeds":           {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/telemetry.v1.TelemetryService/UpdateBedCondition": {role.RoleDoctor, role.RoleNurse},

	"/staff.v1.StaffService/CreateEmployee":     {role.RoleAdmin},

	"/staff.v1.StaffService/GetEmployee":        {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"/staff.v1.StaffService/ListEmployees":      {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"/staff.v1.StaffService/DeactivateEmployee": {role.RoleAdmin},

	"/patients.v1.PatientService/CreatePatient":        {role.RoleAdmin, role.RoleReception},
	"/patients.v1.PatientService/GetPatient":           {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"/patients.v1.PatientService/GetPatientByDocument": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"/patients.v1.PatientService/ListPatients":         {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},

	"/encounter.v1.EncounterService/CreateEncounter":        {role.RoleDoctor},
	"/encounter.v1.EncounterService/GetEncounter":           {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/encounter.v1.EncounterService/GetEncounters":           {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/observation.v1.ObservationService/CreateObservation":  {role.RoleDoctor, role.RoleNurse},
	"/observation.v1.ObservationService/GetObservations":    {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/condition.v1.ConditionService/CreateCondition":        {role.RoleDoctor},
	"/condition.v1.ConditionService/GetConditions":          {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/allergy.v1.AllergyService/CreateAllergyIntolerance":   {role.RoleDoctor, role.RoleNurse},
	"/allergy.v1.AllergyService/GetAllergyIntolerances":     {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/medication.v1.MedicationService/CreateMedicationRequest":  {role.RoleDoctor},
	"/medication.v1.MedicationService/GetMedicationRequests":    {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/diagnostic_report.v1.DiagnosticReportService/CreateDiagnosticReport": {role.RoleDoctor},
	"/diagnostic_report.v1.DiagnosticReportService/GetDiagnosticReports":   {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/clinical.v1.ImagingService/UploadDICOM":         {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/clinical.v1.ImagingService/GetImagingStudy":     {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RolePatient},
	"/clinical.v1.ImagingService/ListImagingStudies":   {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RolePatient},
	"/clinical.v1.ImagingService/GetDICOMDownloadURL": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RolePatient},
}
