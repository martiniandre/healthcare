package policy

import "github.com/healthcare/backend/internal/shared/role"

var publicGRPCMethods = map[string]bool{
	"/auth.v1.AuthService/Login":    true,
	"/auth.v1.AuthService/Register": true,
	"/auth.v1.AuthService/Logout":   true,
	"/grpc.health.v1.Health/Check":  true,
	"/grpc.health.v1.Health/Watch":  true,
}

var grpcMethodRoles = deriveRoles(methodCapabilities)

var publicHTTPRoutes = map[string]bool{
	"GET /health":               true,
	"POST /api/v1/auth/login":   true,
	"POST /api/v1/auth/refresh": true,
	"POST /api/v1/auth/logout":  true,
	"GET /api/v1/auth/me":       true,
}

var httpRouteRoles = deriveRoles(routeCapabilities)

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
