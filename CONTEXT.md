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

---

## Architecture

### Backend (Go — 16 modules)

```
backend/
├── cmd/api/main.go          — Composition root
├── proto/                   — 9 service proto files
├── internal/
│   ├── api/                 — HTTP router
│   ├── app/                 — gRPC server + interceptors
│   │   └── interceptor/     — Auth, RBAC, rate-limit, logging (AOP)
│   ├── modules/
│   │   ├── allergy/         — FHIR AllergyIntolerance CRUD
│   │   ├── analytics/       — Aggregated metrics (was stats/)
│   │   ├── audit_logs/      — Immutable HIPAA/LGPD records
│   │   ├── auth/            — JWT + RBAC
│   │   ├── condition/       — FHIR Condition CRUD
│   │   ├── diagnostic_report/ — FHIR DiagnosticReport CRUD
│   │   ├── encounter/       — FHIR Encounter CRUD
│   │   ├── exam_analyzer/   — AI analysis (Vertex AI)
│   │   ├── health/          — Healthcheck + readiness
│   │   ├── imaging/         — DICOM + FHIR ImagingStudy
│   │   ├── medication/      — FHIR MedicationRequest CRUD
│   │   ├── observation/     — FHIR Observation CRUD
│   │   ├── patients/        — FHIR Patient CRUD
│   │   ├── staff/           — Employee management (PostgreSQL)
│   │   └── telemetry/       — Real-time vitals monitoring
│   └── shared/
│       ├── apperrors/       — Standardized error types
│       ├── cache/           — Redis client
│       ├── config/          — Env-based config
│       ├── ctxkeys/         — Context key constants
│       ├── database/        — pgxpool wrapper
│       ├── fhir/            — 8 typed FHIR resource builders
│       ├── healthcare/      — GCP Healthcare API client
│       ├── logger/          — Structured logging
│       ├── migrations/      — SQL migrations
│       ├── role/            — Shared Role type (extracted from auth)
│       ├── storage/         — GCS client
│       └── validator/       — CPF, dates, ICD-10
└── migrations/              — SQL migration files
```

### Frontend (React + Vite — 8 modules)

```
frontend/src/
├── app/                     — Router + layout
├── modules/
│   ├── analytics/           — Charts + metrics dashboard
│   ├── audit_logs/          — Compliance log viewer
│   ├── auth/                — Login, registration
│   ├── exam_analyzer/       — AI exam upload + results
│   ├── imaging/             — DICOM viewer + upload
│   ├── patients/            — Patient CRUD + clinical tabs
│   ├── staff/               — Employee CRUD
│   └── telemetry/           — Room vitals dashboard
└── shared/
    ├── components/ui/       — Button, Card, Dialog, Select, etc.
    ├── hooks/               — useAuthInit, useDebounce, usePageViewLogger
    ├── i18n/                — pt-BR, en-US, es-ES
    ├── services/            — Legacy API services
    ├── store/               — Zustand auth store
    └── utils/               — http, cn, validators
```

### Testing

| Layer     | Framework | Tests | Command |
|-----------|-----------|-------|---------|
| Backend   | Go `testing` | 40+ service/handler tests | `go test ./internal/...` |
| Frontend  | Vitest     | 22 unit tests (hooks, schemas, validators, Button) | `npm run test` |
| E2E       | Playwright | Imaging + Telemetry flows | `npm run e2e` |

### Key architectural decisions

- **Hexagonal (Ports & Adapters)**: Every backend module has `repository.go` (port), `service.go` (business logic), `grpc_handler.go` (inbound adapter), `register.go` (DI wiring with `Dependency` struct)
- **FHIR-first for clinical data**: 7 of 15 backend modules use GCP Healthcare API (patients, encounter, observation, condition, allergy, medication, diagnostic_report). Only auth, staff, audit_logs, imaging metadata, and analytics use PostgreSQL.
- **AOP security**: All gRPC endpoints protected by shared interceptors (JWT auth + RBAC + rate-limit). Endpoints not in `permissions.go` are blocked by default.
- **Single `Register(grpcServer, dep)` pattern**: Every module follows the same wiring signature, keeping `main.go` declarative.
- **Pre-push hook**: Runs only `go vet` + `npm run lint` (fast). Heavy checks (tests, build) run in CI.
