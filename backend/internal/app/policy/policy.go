package policy

import "github.com/healthcare/backend/internal/shared/role"

var publicGRPCMethods = map[string]bool{
	"/auth.v1.AuthService/Login":    true,
	"/auth.v1.AuthService/Register": true,
	"/auth.v1.AuthService/Logout":   true,
	"/grpc.health.v1.Health/Check":  true,
	"/grpc.health.v1.Health/Watch":  true,
}

var grpcMethodRoles = map[string][]role.Role{
	"/audit_logs.v1.AuditLogsService/CreateAuditLog": {role.RoleAdmin},
	"/audit_logs.v1.AuditLogsService/ListAuditLogs":  {role.RoleAdmin},

	"/telemetry.v1.TelemetryService/GetRooms":           {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"/telemetry.v1.TelemetryService/UnlockRoom":         {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/telemetry.v1.TelemetryService/GetBeds":            {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/telemetry.v1.TelemetryService/UpdateBedCondition": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/exam_analyzer.v1.ExamAnalyzerService/ListAnalyses":   {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/exam_analyzer.v1.ExamAnalyzerService/GetAnalysis":    {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/exam_analyzer.v1.ExamAnalyzerService/DeleteAnalysis": {role.RoleAdmin, role.RoleDoctor},

	"/staff.v1.StaffService/CreateEmployee":     {role.RoleAdmin},
	"/staff.v1.StaffService/GetEmployee":        {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"/staff.v1.StaffService/ListEmployees":      {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"/staff.v1.StaffService/DeactivateEmployee": {role.RoleAdmin},

	"/patients.v1.PatientService/CreatePatient":        {role.RoleAdmin, role.RoleReception},
	"/patients.v1.PatientService/GetPatient":           {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"/patients.v1.PatientService/GetPatientByDocument": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"/patients.v1.PatientService/ListPatients":         {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},

	"/encounter.v1.EncounterService/CreateEncounter": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"/encounter.v1.EncounterService/GetEncounter":    {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/encounter.v1.EncounterService/GetEncounters":   {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/observation.v1.ObservationService/CreateObservation":      {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/observation.v1.ObservationService/CreateObservationBatch": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/observation.v1.ObservationService/GetObservations":        {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/condition.v1.ConditionService/CreateCondition": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/condition.v1.ConditionService/GetConditions":   {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/allergy.v1.AllergyService/CreateAllergyIntolerance": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/allergy.v1.AllergyService/GetAllergyIntolerances":   {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/medication.v1.MedicationService/CreateMedicationRequest": {role.RoleAdmin, role.RoleDoctor},
	"/medication.v1.MedicationService/GetMedicationRequests":   {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/diagnostic_report.v1.DiagnosticReportService/CreateDiagnosticReport": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/diagnostic_report.v1.DiagnosticReportService/GetDiagnosticReports":   {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"/imaging.v1.ImagingService/UploadDICOM":         {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/imaging.v1.ImagingService/GetImagingStudy":     {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/imaging.v1.ImagingService/ListImagingStudies":  {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"/imaging.v1.ImagingService/GetDICOMDownloadURL": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
}

var publicHTTPRoutes = map[string]bool{
	"GET /health":              true,
	"POST /api/v1/auth/login":  true,
	"POST /api/v1/auth/logout": true,
	"GET /api/v1/auth/me":      true,
}

var httpRouteRoles = map[string][]role.Role{
	"GET /api/v1/telemetry/rooms":                   {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"POST /api/v1/telemetry/rooms/{roomId}/unlock":  {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"GET /api/v1/telemetry/rooms/{roomId}/beds":     {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"POST /api/v1/telemetry/beds/{bedId}/condition": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"GET /api/v1/staff/employees":  {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"POST /api/v1/staff/employees": {role.RoleAdmin},

	"POST /api/v1/appointments":                        {role.RoleAdmin, role.RoleReception},
	"GET /api/v1/appointments":                         {role.RoleAdmin, role.RoleReception, role.RoleDoctor, role.RoleNurse},
	"GET /api/v1/appointments/my":                      {role.RolePatient},
	"GET /api/v1/appointments/{appointmentId}":         {role.RoleAdmin, role.RoleReception, role.RoleDoctor, role.RoleNurse},
	"POST /api/v1/appointments/{appointmentId}/cancel": {role.RoleAdmin, role.RoleReception},
	"PUT /api/v1/appointments/{appointmentId}":         {role.RoleAdmin, role.RoleReception},

	"POST /api/v1/schedule/unavailability":                      {role.RoleAdmin, role.RoleDoctor},
	"GET /api/v1/schedule/unavailability":                       {role.RoleAdmin, role.RoleDoctor},
	"DELETE /api/v1/schedule/unavailability/{unavailabilityId}": {role.RoleAdmin, role.RoleDoctor},

	"GET /api/v1/portal/dashboard":    {role.RolePatient},
	"GET /api/v1/portal/encounters":   {role.RolePatient},
	"GET /api/v1/portal/observations": {role.RolePatient},
	"GET /api/v1/portal/conditions":   {role.RolePatient},
	"GET /api/v1/portal/medications":  {role.RolePatient},
	"GET /api/v1/portal/reports":      {role.RolePatient},
	"GET /api/v1/portal/imaging":      {role.RolePatient},

	"GET /api/v1/patients":                 {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"POST /api/v1/patients":                {role.RoleAdmin, role.RoleReception},
	"GET /api/v1/patients/{patientFhirId}": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},

	"GET /api/v1/patients/{patientFhirId}/observations":            {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"GET /api/v1/encounters/{encounterFhirId}/observations":        {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"POST /api/v1/encounters/{encounterFhirId}/observations":       {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"POST /api/v1/encounters/{encounterFhirId}/observations/batch": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"PUT /api/v1/observations/{observationFhirId}":                 {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"DELETE /api/v1/observations/{observationFhirId}":              {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"GET /api/v1/notifications":                        {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception, role.RolePatient},
	"POST /api/v1/notifications/{notificationId}/read": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception, role.RolePatient},
	"GET /api/v1/notifications/unread-count":           {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception, role.RolePatient},
	"GET /api/v1/notifications/stream":                 {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception, role.RolePatient},
	"GET /api/v1/notification-events":                  {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},

	"GET /api/v1/encounters/{encounterFhirId}/medications":  {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"POST /api/v1/encounters/{encounterFhirId}/medications": {role.RoleAdmin, role.RoleDoctor},

	"GET /api/v1/patients/{patientFhirId}/studies":  {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"POST /api/v1/patients/{patientFhirId}/studies": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"GET /api/v1/studies/{studyId}":                 {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"GET /api/v1/exam-analyses":                 {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"POST /api/v1/exam-analyses":                {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"GET /api/v1/exam-analyses/{analysisId}":    {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"DELETE /api/v1/exam-analyses/{analysisId}": {role.RoleAdmin, role.RoleDoctor},

	"GET /api/v1/patients/{patientFhirId}/encounters":  {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"POST /api/v1/patients/{patientFhirId}/encounters": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception},
	"GET /api/v1/encounters/{encounterFhirId}":         {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"PUT /api/v1/encounters/{encounterFhirId}":         {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"DELETE /api/v1/encounters/{encounterFhirId}":      {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"GET /api/v1/encounters/{encounterFhirId}/reports":      {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"POST /api/v1/encounters/{encounterFhirId}/reports":     {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"PUT /api/v1/reports/{reportFhirId}":                    {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"GET /api/v1/reports/{reportFhirId}/versions":           {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"GET /api/v1/reports/{reportFhirId}/versions/{version}": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"GET /api/v1/patients/{patientFhirId}/conditions":                      {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"POST /api/v1/patients/{patientFhirId}/conditions":                     {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"PUT /api/v1/patients/{patientFhirId}/conditions/{conditionFhirId}":    {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"DELETE /api/v1/patients/{patientFhirId}/conditions/{conditionFhirId}": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"GET /api/v1/audit-logs":  {role.RoleAdmin},
	"POST /api/v1/audit-logs": {role.RoleAdmin},

	"GET /api/v1/patients/{patientFhirId}/allergies":                    {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"POST /api/v1/patients/{patientFhirId}/allergies":                   {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"PUT /api/v1/patients/{patientFhirId}/allergies/{allergyFhirId}":    {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"DELETE /api/v1/patients/{patientFhirId}/allergies/{allergyFhirId}": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},

	"GET /api/v1/analytics":                                    {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"GET /api/v1/analytics/dashboard":                          {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"GET /api/v1/analytics/dashboard/consultations-per-doctor": {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"GET /api/v1/analytics/dashboard/occupancy-rate":           {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"GET /api/v1/analytics/dashboard/avg-wait-time":            {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
	"GET /api/v1/analytics/dashboard/top-diagnoses":            {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
}

var notificationTypeRoles = map[string][]role.Role{
	"telemetry_alert":   {role.RoleNurse, role.RoleDoctor},
	"exam_complete":     {role.RoleDoctor},
	"encounter_created": {role.RoleDoctor, role.RoleReception, role.RolePatient},
	"encounter_updated": {role.RoleDoctor, role.RoleReception},
	"patient_created":   {role.RoleReception, role.RoleAdmin},
	"patient_updated":   {role.RoleReception, role.RoleDoctor},
	"audit_alert":       {role.RoleAdmin},
	"report_ready":      {role.RolePatient},
	"system":            {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
}

func IsPublicGRPCMethod(fullMethod string) bool {
	return publicGRPCMethods[fullMethod]
}

func GRPCMethodRoles(fullMethod string) ([]role.Role, bool) {
	allowedRoles, defined := grpcMethodRoles[fullMethod]
	return allowedRoles, defined
}

func IsPublicHTTPRoute(routePattern string) bool {
	return publicHTTPRoutes[routePattern]
}

func HTTPRouteRoles(routePattern string) ([]role.Role, bool) {
	allowedRoles, defined := httpRouteRoles[routePattern]
	return allowedRoles, defined
}

func RolesForNotificationType(notificationType string) ([]role.Role, bool) {
	allowedRoles, defined := notificationTypeRoles[notificationType]
	return allowedRoles, defined
}
