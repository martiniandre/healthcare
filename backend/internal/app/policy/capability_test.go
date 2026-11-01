package policy

import (
	"sort"
	"testing"

	"github.com/healthcare/backend/internal/shared/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const permissionsRoute = "GET /api/v1/auth/permissions"

var capabilitiesWithoutBackendRoute = map[Capability]bool{
	CapabilityDesignSystemRead: true,
}

func snapshotRoles(derived map[string][]role.Role, ignoredKeys ...string) map[string]string {
	ignored := make(map[string]bool, len(ignoredKeys))
	for _, ignoredKey := range ignoredKeys {
		ignored[ignoredKey] = true
	}
	snapshot := make(map[string]string, len(derived))
	for key, roles := range derived {
		if ignored[key] {
			continue
		}
		snapshot[key] = rolesKey(roles)
	}
	return snapshot
}

func TestDerivedGRPCMethodRolesMatchGolden(testingInstance *testing.T) {
	assert.Equal(testingInstance, goldenGRPCMethodRoles, snapshotRoles(grpcMethodRoles))
}

func TestDerivedHTTPRouteRolesMatchGolden(testingInstance *testing.T) {
	assert.Equal(testingInstance, goldenHTTPRouteRoles, snapshotRoles(httpRouteRoles, permissionsRoute))
}

func TestPermissionsRouteIsRegisteredForEveryRole(testingInstance *testing.T) {
	allowedRoles, defined := HTTPRouteRoles(permissionsRoute)
	require.True(testingInstance, defined)
	assert.Equal(testingInstance, "ADMIN,DOCTOR,NURSE,PATIENT,RECEPTION", rolesKey(allowedRoles))
	assert.False(testingInstance, IsPublicHTTPRoute(permissionsRoute))
}

func TestEveryMappedCapabilityHasRoles(testingInstance *testing.T) {
	for route, capability := range routeCapabilities {
		assert.NotEmpty(testingInstance, capabilityRoles[capability], "route %s uses capability %s without roles", route, capability)
	}
	for method, capability := range methodCapabilities {
		assert.NotEmpty(testingInstance, capabilityRoles[capability], "method %s uses capability %s without roles", method, capability)
	}
}

func TestEveryCapabilityGuardsARouteOrIsListed(testingInstance *testing.T) {
	used := make(map[Capability]bool)
	for _, capability := range routeCapabilities {
		used[capability] = true
	}
	for _, capability := range methodCapabilities {
		used[capability] = true
	}
	for capability := range capabilityRoles {
		if capabilitiesWithoutBackendRoute[capability] {
			continue
		}
		assert.True(testingInstance, used[capability], "capability %s guards no route or method", capability)
	}
}

func TestCapabilitiesForRoleAreSortedAndDeterministic(testingInstance *testing.T) {
	for _, callerRole := range everyRole {
		first := CapabilitiesForRole(callerRole)
		second := CapabilitiesForRole(callerRole)
		assert.Equal(testingInstance, first, second)
		assert.True(testingInstance, sort.SliceIsSorted(first, func(firstIndex, secondIndex int) bool {
			return first[firstIndex] < first[secondIndex]
		}))
	}
}

func TestCapabilitiesForRole(testingInstance *testing.T) {
	assert.Contains(testingInstance, CapabilitiesForRole(role.RoleNurse), CapabilityMedicationRead)
	assert.NotContains(testingInstance, CapabilitiesForRole(role.RoleNurse), CapabilityMedicationCreate)
	assert.Contains(testingInstance, CapabilitiesForRole(role.RoleAdmin), CapabilityMedicationCreate)
	assert.Contains(testingInstance, CapabilitiesForRole(role.RoleDoctor), CapabilityMedicationCreate)
	assert.NotContains(testingInstance, CapabilitiesForRole(role.RoleReception), CapabilityMedicationCreate)
	assert.Contains(testingInstance, CapabilitiesForRole(role.RoleReception), CapabilityDesignSystemRead)
	assert.NotContains(testingInstance, CapabilitiesForRole(role.RolePatient), CapabilityDesignSystemRead)
	assert.Empty(testingInstance, CapabilitiesForRole(role.Role("UNKNOWN")))
}

func TestPatientCapabilitiesAreLimitedToSelfService(testingInstance *testing.T) {
	assert.Equal(testingInstance, []Capability{
		CapabilityAppointmentReadOwn,
		CapabilityAuthPermissions,
		CapabilityNotificationRead,
		CapabilityNotificationUpdate,
		CapabilityPortalRead,
	}, CapabilitiesForRole(role.RolePatient))
}
