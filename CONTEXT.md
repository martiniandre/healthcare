# Healthcare Platform

A clinical engine and administrative interface for healthcare facilities, designed for HIPAA/LGPD compliance with FHIR R4 interoperability and DICOM medical imaging support.

## Language

### Patient:
A person receiving healthcare services, whose clinical data lives in GCP Healthcare API as a FHIR R4 Patient resource.
_Avoid_: Client, customer, person, individual

### Staff:
A healthcare professional or administrative worker with system access via a User account.
_Avoid_: Employee, worker, collaborator, team member

### User:
An authenticated identity with a role-based access profile (Admin, Doctor, Nurse, Reception, Patient).
_Avoid_: Account, login, credential

### Encounter:
A clinical interaction between a Patient and one or more Staff providers.
_Avoid_: Appointment, visit, consultation, attendance

### Observation:
A measured value or clinical finding about a Patient, identified by a LOINC code with associated value, unit, and reference range.
_Avoid_: Vital sign, measurement, reading, vital

### Condition:
A diagnosis, problem, or medical condition identified for a Patient.
_Avoid_: Diagnosis, problem, issue, complaint, disorder

### DiagnosticReport:
A structured report of a diagnostic procedure, combining clinical findings and interpretation by a professional.
_Avoid_: Lab report, exam result, study report, test report

### AllergyIntolerance:
An adverse reaction or intolerance to a substance, recorded per Patient.
_Avoid_: Allergy, reaction, intolerance, sensitivity

### MedicationRequest:
A prescription or order for a medication to be administered to a Patient.
_Avoid_: Prescription, drug order, medication order, Rx

### ImagingStudy:
A DICOM imaging study whose metadata is stored in PostgreSQL and pixel data in GCS, with a FHIR ImagingStudy resource in the Healthcare API.
_Avoid_: DICOM study, image, scan, radiology study

### Telemetry:
Real-time Patient vital signs monitoring organized by Room and Bed, with configurable thresholds and status indicators.
_Avoid_: Monitoring, vitals, patient tracking, remote monitoring

### ExamAnalysis:
An AI-powered analysis of a medical exam, processed asynchronously via Vertex AI (Gemini), with consent and anonymization controls.
_Avoid_: AI analysis, exam review, automated diagnosis, AI diagnosis

### AuditLog:
An immutable record of system access and operations, persisted for HIPAA/LGPD compliance.
_Avoid_: Log, trail, history entry, activity record

### Analytics:
Aggregated clinical and operational metrics derived from FHIR and local data stores, presented as charts and summaries.
_Avoid_: Stats, dashboard, reports, metrics, BI

### Room:
A monitored physical space containing Beds, secured by a passcode, within the Telemetry module.
_Avoid_: Ward, unit, hall, chamber

### Bed:
A monitored Bed within a Room, tracking a Patient's vital signs (BPM, SpO2, temperature) and status/condition.
_Avoid_: Station, spot, unit, slot

### Auth:
Authentication and authorization subsystem that issues JWT tokens, manages sessions, and enforces role-based access control across all API surfaces.
_Avoid_: Security, login, identity, access control
