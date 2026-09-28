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
	cardiologyDepartmentID := fetchDepartmentID(t, testServer.db, "Cardiologia")

	departmentsResponse := adminClient.Get(t, "/api/v1/staff/departments")
	requireStatusCode(t, departmentsResponse, http.StatusOK)
	var departmentsList []map[string]interface{}
	decodeJSONResponse(t, departmentsResponse, &departmentsList)
	if len(departmentsList) < 20 {
		t.Fatalf("expected at least 20 departments in catalog, got %d", len(departmentsList))
	}

	createResponse := adminClient.Post(t, "/api/v1/staff/employees", map[string]interface{}{
		"created_by":    adminUserID,
		"full_name":     "Dra. Marina Ribeiro",
		"email":         "marina.ribeiro@clinica.com",
		"role":          "DOCTOR",
		"crm_number":    "CRM 123456",
		"department_id": cardiologyDepartmentID,
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
		"created_by":    adminUserID,
		"full_name":     "Dra. Marina Ribeiro",
		"email":         "marina.ribeiro@clinica.com",
		"role":          "DOCTOR",
		"crm_number":    "CRM 123456",
		"department_id": cardiologyDepartmentID,
	})
	requireStatusCode(t, duplicateResponse, http.StatusConflict)

	withoutDepartmentResponse := adminClient.Post(t, "/api/v1/staff/employees", map[string]interface{}{
		"created_by": adminUserID,
		"full_name":  "Dr. Sem Departamento",
		"email":      "sem.departamento@clinica.com",
		"role":       "DOCTOR",
		"crm_number": "CRM 654321",
	})
	requireStatusCode(t, withoutDepartmentResponse, http.StatusBadRequest)

	listResponse := adminClient.Get(t, "/api/v1/staff/employees")
	requireStatusCode(t, listResponse, http.StatusOK)
	var employeesList []map[string]interface{}
	decodeJSONResponse(t, listResponse, &employeesList)
	foundEmployee := false
	for _, employee := range employeesList {
		if employee["id"] == employeeID {
			foundEmployee = true
			if employee["department_name"] != "Cardiologia" {
				t.Fatalf("expected department_name Cardiologia, got %v", employee["department_name"])
			}
			if employee["is_active"] != true {
				t.Fatal("expected newly created employee to be active")
			}
			break
		}
	}
	if !foundEmployee {
		t.Fatalf("expected employee %s in employees list", employeeID)
	}
}

func TestAdminTogglesAndDeletesEmployee(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")
	adminUserID := fetchUserID(t, testServer.db, "admin@hospital.com")
	nursingDepartmentID := fetchDepartmentID(t, testServer.db, "Pediatria")

	createResponse := adminClient.Post(t, "/api/v1/staff/employees", map[string]interface{}{
		"created_by":    adminUserID,
		"full_name":     "Enf. Carlos Dutra",
		"email":         "carlos.dutra@clinica.com",
		"role":          "NURSE",
		"crm_number":    "COREN 987654",
		"department_id": nursingDepartmentID,
	})
	requireStatusCode(t, createResponse, http.StatusCreated)
	var createdEmployee map[string]interface{}
	decodeJSONResponse(t, createResponse, &createdEmployee)
	employeeID := createdEmployee["employee_id"].(string)

	disableResponse := adminClient.Patch(t, "/api/v1/staff/employees/"+employeeID+"/status", map[string]interface{}{
		"is_active": false,
	})
	requireStatusCode(t, disableResponse, http.StatusOK)
	var disabledEmployee map[string]interface{}
	decodeJSONResponse(t, disableResponse, &disabledEmployee)
	if disabledEmployee["is_active"] != false {
		t.Fatal("expected employee to be disabled")
	}

	deleteResponse := adminClient.Delete(t, "/api/v1/staff/employees/"+employeeID)
	requireStatusCode(t, deleteResponse, http.StatusOK)

	unknownEmployeeResponse := adminClient.Patch(t, "/api/v1/staff/employees/"+employeeID+"/status", map[string]interface{}{
		"is_active": false,
	})
	requireStatusCode(t, unknownEmployeeResponse, http.StatusOK)

	notFoundResponse := adminClient.Patch(t, "/api/v1/staff/employees/nao-existe/status", map[string]interface{}{
		"is_active": true,
	})
	requireStatusCode(t, notFoundResponse, http.StatusBadRequest)
}

func TestDoctorCannotCreateEmployee(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")
	doctorUserID := fetchUserID(t, testServer.db, "medico@clinica.com")
	cardiologyDepartmentID := fetchDepartmentID(t, testServer.db, "Cardiologia")

	response := doctorClient.Post(t, "/api/v1/staff/employees", map[string]interface{}{
		"created_by":    doctorUserID,
		"full_name":     "Médico Sem Permissão",
		"email":         "sem.permissao@clinica.com",
		"role":          "DOCTOR",
		"crm_number":    "CRM 111222",
		"department_id": cardiologyDepartmentID,
	})
	requireStatusCode(t, response, http.StatusForbidden)
}
