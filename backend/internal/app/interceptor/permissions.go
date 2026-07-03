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

	"/clinical.v1.ClinicalService/CreateEncounter":          {role.RoleDoctor},
	"/clinical.v1.ClinicalService/GetEncounters":             {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/clinical.v1.ClinicalService/CreateObservation":         {role.RoleDoctor, role.RoleNurse},
	"/clinical.v1.ClinicalService/GetObservations":           {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/clinical.v1.ClinicalService/CreateCondition":           {role.RoleDoctor},
	"/clinical.v1.ClinicalService/GetConditions":             {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/clinical.v1.ClinicalService/CreateAllergyIntolerance":  {role.RoleDoctor, role.RoleNurse},
	"/clinical.v1.ClinicalService/GetAllergyIntolerances":    {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/clinical.v1.ClinicalService/CreateMedicationRequest":   {role.RoleDoctor},
	"/clinical.v1.ClinicalService/GetMedicationRequests":     {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/clinical.v1.ClinicalService/CreateDiagnosticReport":    {role.RoleDoctor},
	"/clinical.v1.ClinicalService/GetDiagnosticReports":      {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/clinical.v1.ImagingService/UploadDICOM":         {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/clinical.v1.ImagingService/GetImagingStudy":     {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RolePatient},
	"/clinical.v1.ImagingService/ListImagingStudies":   {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RolePatient},
	"/clinical.v1.ImagingService/GetDICOMDownloadURL": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RolePatient},
}
