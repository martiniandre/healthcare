---
name: go-module-gen
description: Creates a new independent backend Go domain module (gRPC) in the Healthcare monorepo, following Hexagonal Architecture and the AGENTS.MD blueprint (proto, pb stubs, model, repository, service, grpc_handler, register, RBAC policy, migration when operational, unit tests). Use when the user says "crie o módulo backend [Nome] com os campos [Campos]", "novo módulo Go", "gerar módulo backend", "backend module", or triggers the Antigravity AG-01 flow.
---

# Go Backend Module Generator (AG-01)

Create a new independent domain module under `backend/internal/modules/{domain}/` and plug it into the core (`backend/cmd/api/main.go`).

## Trigger

User request: `"Antigravity, crie o módulo backend [Nome] com os campos [Campos]."`

## Step 0 — Clarification (MANDATORY)

Identify ambiguities before writing any code. If anything is unclear about fields, business rules, persistence target, integrations, or expected behavior, ask the user and wait for the answer. Never assume silently.

## Step 0.5 — Worktree isolation (MANDATORY when other features may be in flight)

If this module is created inside a dedicated worktree, follow the `parallel-worktrees` skill: it is already claimed by the orchestrator. If the current working directory is the main repo, load `parallel-worktrees` Step 0 and activate a worktree first — never scaffold on `main`. All `policy.go` / `main.go` edits must follow its shared-file additive-anchored protocol.

## Step 1 — Persistence decision (MANDATORY)

Decide where the data lives BEFORE scaffolding:

- **Clinical (Google Cloud Healthcare API / FHIR):** `Patient`, `Observation`, `Encounter`, `Condition`, `DiagnosticReport`, `ImagingStudy`, `AllergyIntolerance`, `MedicationRequest`. The repository talks to `shared/healthcare.FHIRClient`. **No SQL migration.**
- **Operational (PostgreSQL local):** `auth`, `staff`, `telemetry`, `audit_logs`, `notifications`, `schedule`. The repository uses `pgxpool`. **Requires a migration** in `backend/migrations/`.

Migration naming convention: `{NNN}_{snake_case_name}.up.sql` and `{NNN}_{snake_case_name}.down.sql`. The `NNN` is **reserved from the worktree registry**, never computed as "local max + 1" (sibling worktrees compute the same number): `node scripts/worktrees.mjs reserve-migration --slug <slug> --name <snake_case_name>`. In a single-worktree session with no registry entry, fall back to the highest existing number + 1 and verify uniqueness.

## Step 2 — Proto contract

Create `backend/proto/{domain}.proto`:

```proto
syntax = "proto3";

package {domain}.v1;

option go_package = "github.com/healthcare/backend/internal/modules/{domain}/pb";

service {Domain}Service {
  rpc Create{Domain}(Create{Domain}Request) returns (Create{Domain}Response);
  rpc Get{Domains}(Get{Domains}Request) returns (Get{Domains}Response);
}
```

Use `snake_case` field names and plain proto types (string, int32, repeated). Wire the `go_package` exactly as above.

## Step 3 — pb stubs (hand-written, no protoc)

This repository does NOT run `protoc`/`buf`. Create `backend/internal/modules/{domain}/pb/{domain}.pb.go` as **hand-written stubs** that mirror the proto contract (see `condition/pb/condition.pb.go` as reference):

```go
package pb

import "context"

type {Domain}ServiceServer interface {
	Create{Domain}(ctx context.Context, request *Create{Domain}Request) (*Create{Domain}Response, error)
	Get{Domains}(ctx context.Context, request *Get{Domains}Request) (*Get{Domains}Response, error)
}

type Create{Domain}Request struct {
	FieldA string
}

type Create{Domain}Response struct {
	{Domain}FhirId string
}

type Get{Domains}Request struct {
	PatientFhirId string
}

type {Domain} struct {
	FhirId string
	FieldA string
}

type Get{Domains}Response struct {
	{Domains} []*{Domain}
}

func Register{Domain}ServiceServer(_ interface{}, server {Domain}ServiceServer) {}
```

Field names in Go use CamelCase with the `Fhir` acronym capitalized as `Fhir` (e.g., `PatientFhirId`, `Icd10Code`).

## Step 4 — Model and DTOs

Create `backend/internal/modules/{domain}/model.go` with the pure domain entity (plain struct with exported fields, no tags).

If the module needs partial-update inputs, create `backend/internal/modules/{domain}/inputs.go` with `Update{Domain}Input` structs using `*string` / `*int32` pointer fields to distinguish omitted fields from zero values.

Use the generated pb structs as DTOs. Never expose persistence internals through the service interface.

## Step 5 — Repository

Create `backend/internal/modules/{domain}/repository.go`:

```go
type Repository interface {
	Create{Domain}(ctx context.Context, entity *{Domain}) (*{Domain}, error)
	Get{Domain}ByID(ctx context.Context, fhirResourceID string) (*{Domain}, error)
	Get{Domains}ByPatient(ctx context.Context, patientFHIRID string) ([]*{Domain}, error)
	Update{Domain}(ctx context.Context, fhirResourceID string, entity *{Domain}) (*{Domain}, error)
	Delete{Domain}(ctx context.Context, fhirResourceID string) error
}
```

- **FHIR modules:** inject `healthcare.FHIRClient`. Use `fhir.New{Domain}Resource(...)` builders, `fhirClient.CreateResource` / `GetResource` / `SearchResources` / `UpdateResource` / `DeleteResource`, and decode with `fhir.DecodeResource[fhir.{Domain}]` or `fhir.DecodeBundle[fhir.{Domain}]`. Map FHIR -> domain with a `mapFHIR{Domain}ToDomain` helper using `fhir.CodeableConceptParts`, `fhir.SplitReferenceID`, `fhir.ParseRFC3339`. Translate not-found into `apperrors.Err{Domain}NotFound` when `healthcare.IsNotFound(err)`.
- **Operational modules:** inject `*pgxpool.Pool`. Write SQL queries against the tables created in the migration.

## Step 6 — Service

Create `backend/internal/modules/{domain}/service.go` with the business rules:

```go
type Service interface {
	Create{Domain}(ctx context.Context, input Create{Domain}Input) (*{Domain}, error)
	Get{Domains}ByPatient(ctx context.Context, patientFHIRID string) ([]*{Domain}, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
```

Validation patterns used across the codebase:
- Return `apperrors.InvalidArgument("invalid {domain} input", fieldViolations)` where `fieldViolations` is `map[string]string` of field name -> message.
- Use `shared/validator` helpers (`IsValidICD10`, `IsValidClinicalStatus`, CPF, date checks, etc.).
- Apply defaults via a `normalize{Domain}Field` helper (e.g., `normalizeClinicalStatus` defaults to `"active"`).
- Every `context.Context` must be the first parameter.

## Step 7 — gRPC handler

Create `backend/internal/modules/{domain}/grpc_handler.go` implementing the pb interface, acting as a thin controller that calls the Service. Convert all errors with `apperrors.ToGRPCStatus(err)`. Do not put business logic here.

## Step 8 — RBAC permissions (MANDATORY)

Register every new endpoint in `backend/internal/app/policy/policy.go`:
- gRPC methods in the `grpcMethodRoles` map (full path `"/{domain}.v1.{Domain}Service/{Method}"`).
- HTTP routes in the `httpRouteRoles` map (pattern `"METHOD /api/v1/..."`) when the module also exposes HTTP endpoints.

**Endpoints not registered are blocked by default.** Available roles: `role.RoleAdmin`, `role.RoleDoctor`, `role.RoleNurse`, `role.RoleReception`, `role.RolePatient`. Reads (`Get*`) usually allow doctor/nurse; writes (`Create*`) usually allow doctor (sometimes nurse); staff/audit admin-only. Follow the pattern of the nearest existing module.

Insert new map entries at their **sorted key position** (never append before a map's closing brace) so parallel worktrees merging into `main` land on non-overlapping lines. Additive-only: never restructure existing entries.

## Step 9 — Core injection

Create `backend/internal/modules/{domain}/register.go`:

```go
package {domain}

type Dependency struct {
	FHIRClient healthcare.FHIRClient // or DB *pgxpool.Pool
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.FHIRClient)
	svc := NewService(repo)
	handler := NewGRPCHandler(svc)
	{domain}pb.Register{Domain}ServiceServer(grpcServer, handler)
	return svc
}
```

Then add only `{domain}.Register(applicationServer.GRPCServer, {domain}.Dependency{...})` in `backend/cmd/api/main.go`, mirroring the existing registrations. Insert the Go import at its **alphabetical position** and the `Register(...)` / HTTP-handler / router-argument lines near the related module block — additive and anchored, so parallel worktrees merge cleanly.

## Step 10 — Unit tests (MANDATORY)

Create:
- `backend/internal/modules/{domain}/mocks/repository.go` — a `MockRepository` struct with call counters and canned responses.
- `backend/internal/modules/{domain}/tests/service_test.go` — same package tests covering: input validation failures, defaults applied, happy path, and repository failure propagation.

Follow the structure of `condition/tests/service_test.go`.

## Validation commands

```bash
cd backend
go build ./...
go test -v ./internal/modules/{domain}/...
go vet ./...
```

## Definition of Done

- [ ] Proto + pb stubs hand-written (no protoc) and consistent.
- [ ] Model, repository, service, grpc_handler, register created.
- [ ] Every endpoint registered in `policy.go` (RBAC) — else it is blocked.
- [ ] `main.go` registers the module with one line.
- [ ] Migration created ONLY for operational (PostgreSQL) modules.
- [ ] Unit tests pass with MockRepository.
- [ ] `go build ./...` and `go vet ./...` clean.
- [ ] Zero comments in all generated code; descriptive multi-word variable names (no single-letter identifiers).
