package integration

import (
	"net/http"
	"testing"
)

func TestAdminCreatesAndListsEmployees(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")
	adminUserID := fetchUserID(t, testServer.db, "admin@hospital.com")

	createResponse := adminClient.Post(t, "/api/v1/staff/employees", map[string]interface{}{
		"created_by": adminUserID,
		"full_name":  "Dra. Marina Ribeiro",
		"email":      "marina.ribeiro@clinica.com",
		"role":       "DOCTOR",
		"crm_number": "CRM 123456",
	})
	requireStatusCode(t, createResponse, http.StatusCreated)
	var createdEmployee map[string]interface{}
	decodeJSONResponse(t, createResponse, &createdEmployee)
	employeeID, hasEmployeeID := createdEmployee["employee_id"].(string)
	if !hasEmployeeID || employeeID == "" {
		t.Fatal("expected employee_id in create employee response")
	}
	staffFHIRID, hasStaffFHIRID := createdEmployee["fhir_resource_id"].(string)
	if !hasStaffFHIRID || staffFHIRID == "" {
		t.Fatal("expected fhir_resource_id in create employee response")
	}
	_ = staffFHIRID

	duplicateResponse := adminClient.Post(t, "/api/v1/staff/employees", map[string]interface{}{
		"created_by": adminUserID,
		"full_name":  "Dra. Marina Ribeiro",
		"email":      "marina.ribeiro@clinica.com",
		"role":       "DOCTOR",
		"crm_number": "CRM 123456",
	})
	requireStatusCode(t, duplicateResponse, http.StatusConflict)

	listResponse := adminClient.Get(t, "/api/v1/staff/employees")
	requireStatusCode(t, listResponse, http.StatusOK)
	var employeesList []map[string]interface{}
	decodeJSONResponse(t, listResponse, &employeesList)
	foundEmployee := false
	for _, employee := range employeesList {
		if employee["id"] == employeeID {
			foundEmployee = true
			break
		}
	}
	if !foundEmployee {
		t.Fatalf("expected employee %s in employees list", employeeID)
	}
}

func TestDoctorCannotCreateEmployee(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")
	doctorUserID := fetchUserID(t, testServer.db, "medico@clinica.com")

	response := doctorClient.Post(t, "/api/v1/staff/employees", map[string]interface{}{
		"created_by": doctorUserID,
		"full_name":  "Médico Sem Permissão",
		"email":      "sem.permissao@clinica.com",
		"role":       "DOCTOR",
	})
	requireStatusCode(t, response, http.StatusForbidden)
}
