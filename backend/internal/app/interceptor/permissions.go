package interceptor

import "github.com/healthcare/backend/internal/modules/auth"

var publicMethods = map[string]bool{
	"/auth.v1.AuthService/Login":    true,
	"/auth.v1.AuthService/Register": true,
	"/auth.v1.AuthService/Logout":   true,
	"/grpc.health.v1.Health/Check":  true,
	"/grpc.health.v1.Health/Watch":  true,
	"/audit_logs.v1.AuditLogsService/CreateAuditLog": true,
}

var methodPermissions = map[string][]auth.Role{
	"/audit_logs.v1.AuditLogsService/ListAuditLogs":     {auth.RoleAdmin},

	"/telemetry.v1.TelemetryService/GetRooms":           {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception},
	"/telemetry.v1.TelemetryService/UnlockRoom":         {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse},
	"/telemetry.v1.TelemetryService/GetBeds":           {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse},
	"/telemetry.v1.TelemetryService/UpdateBedCondition": {auth.RoleDoctor, auth.RoleNurse},

	"/staff.v1.StaffService/CreateEmployee":     {auth.RoleAdmin},

	"/staff.v1.StaffService/GetEmployee":        {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception},
	"/staff.v1.StaffService/ListEmployees":      {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception},
	"/staff.v1.StaffService/DeactivateEmployee": {auth.RoleAdmin},

	"/patients.v1.PatientService/CreatePatient":        {auth.RoleAdmin, auth.RoleReception},
	"/patients.v1.PatientService/GetPatient":           {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception},
	"/patients.v1.PatientService/GetPatientByDocument": {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception},
	"/patients.v1.PatientService/ListPatients":         {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception},

	"/clinical.v1.ClinicalService/CreateEncounter":          {auth.RoleDoctor},
	"/clinical.v1.ClinicalService/GetEncounters":             {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse},
	"/clinical.v1.ClinicalService/CreateObservation":         {auth.RoleDoctor, auth.RoleNurse},
	"/clinical.v1.ClinicalService/GetObservations":           {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse},
	"/clinical.v1.ClinicalService/CreateCondition":           {auth.RoleDoctor},
	"/clinical.v1.ClinicalService/GetConditions":             {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse},
	"/clinical.v1.ClinicalService/CreateAllergyIntolerance":  {auth.RoleDoctor, auth.RoleNurse},
	"/clinical.v1.ClinicalService/GetAllergyIntolerances":    {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse},
	"/clinical.v1.ClinicalService/CreateMedicationRequest":   {auth.RoleDoctor},
	"/clinical.v1.ClinicalService/GetMedicationRequests":     {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse},
	"/clinical.v1.ClinicalService/CreateDiagnosticReport":    {auth.RoleDoctor},
	"/clinical.v1.ClinicalService/GetDiagnosticReports":      {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse},

	"/clinical.v1.ImagingService/UploadDICOM":         {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse},
	"/clinical.v1.ImagingService/GetImagingStudy":     {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RolePatient},
	"/clinical.v1.ImagingService/ListImagingStudies":   {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RolePatient},
	"/clinical.v1.ImagingService/GetDICOMDownloadURL": {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RolePatient},
}
