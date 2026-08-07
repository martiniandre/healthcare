# Implementation Plan — Deepening the healthcare codebase

Date: Aug 7, 2026
Source: architecture-review-20260807-181900.html
Scope: all six candidates, phased, dependency-ordered.

Each phase is a self-contained milestone: code compiles, unit tests pass,
E2E (frontend) pass, then `git add . && commit && push` per AGENTS.md.

---

## Phase 0 — Create `shared/fhir` decoder module (candidate 1, part A)

The highest-leverage deepening. Collapses 8 duplicated FHIR bundle parsers
into one deep, directly-testable module.

### 0.1 New module: `backend/internal/shared/fhir/`
- `model.go` — typed read-side structs per resource:
  `Patient`, `Encounter`, `Observation`, `Condition`, `AllergyIntolerance`,
  `MedicationRequest`, `DiagnosticReport`, `ImagingStudy`. These are the
  domain DTOs produced by decoding; they match what the module repos consume.
- `decoder.go` — the deep module behind a small interface. Public surface:
  one `DecodeBundle(raw json.RawMessage) (*Bundle, error)` (or a
  `BundleDecoder` with `Decode*` per resource type) plus the generic helpers
  the repos need. Internal implementation hides: entry extraction,
  CodeableConcept fallback chain, reference splitting
  (`strings.SplitN(ref, "/", 2)`), `time.Parse(RFC3339)` with visibility of
  parse failures (no more silent zero-times), `clinicalStatus` policy.
- Private helpers mirror the current duplicated logic exactly first (see
  "behavior freeze" below) so nothing changes semantics on day one.

### 0.2 Behavior freeze (critical)
Before rewriting call sites, snapshot current parsing behavior with tests
so the consolidation is a pure refactor:
- `decoder_test.go` — fixtures (raw JSON bundles) captured from the current
  parsers' input/output shape; assert field-level equality against the new
  decoder. Tests live in the module (`package fhir`), per repo test convention.
- Include the documented `clinicalStatus` divergence as a **known-difference
  test**: condition parses it as a string, portal as a CodeableConcept.
  Decide the single canonical representation here (recommendation: keep the
  type the frontend actually consumes; confirm at grilling time) and encode
  it as the decoder's contract.

### 0.3 Migrate the 8 call sites
- `modules/{patients,encounter,observation,condition,allergy,medication,diagnostic_report,portal}/repository.go`
  — replace each local `extractBundleEntries` / inline walks / bespoke parsers
  with a call into `shared/fhir`. Each repo keeps only its thin
  `fhir.X → module struct` mapping.
- Delete the local parse helpers entirely (deletion test).

### 0.4 Tests
- `backend/internal/shared/fhir/decoder_test.go` — field-level fixtures.
- Update existing `repository.go` tests: they now exercise the module mapping,
  not the parsing; keep coverage green.

### Validation
```bash
go build ./...
go test ./backend/internal/shared/fhir/... ./backend/internal/modules/...
```
Commit: `refactor(fhir): consolidate FHIR read-side decoding into shared module`

---

## Phase 1 — Roll ADR-0004 thin adapters to remaining modules (candidate 2)

**Status: DONE** (all 7 modules committed; `go test ./internal/...` green).

Apply the proven pilot pattern (encounter, diagnostic_report, observation) to
the 7 modules that still validate in both transports.

### 1.1 Per module `{allergy, condition, medication, patients, staff, telemetry, imaging}`
- `service.go` — Service interface owns required-field validation, defaults,
  patch-merge on update, typed errors via `apperrors`.
- `grpc_handler.go` — thin translation only; no validation/defaults.
- `http_handler.go` — thin translation only; render errors via
  `apperrors.HTTPCode` (kills silent 500s).
- Ensure the shared FHIR client surfaces 404 so the module returns real
  not-found (mirror pilot).

### 1.2 Tests
- Per module `service_test.go` with `MockRepository` (same package),
  covering validation, defaults, and repository failure paths.
- Handler tests (if present in pilot) mirrored; otherwise ensure service
  tests fully cover the moved rules.

### Validation
```go
go test ./backend/internal/modules/...
```
Commit: `refactor(modules): move validation behind service interface (ADR-0004 rollout)`

---

## Phase 2 — Single RBAC policy module (candidate 3)

Unify the three independent permission stores into one deep module both
transports read from. Deny-by-default stays the single mental model.

### 2.0 Confirm drifts are bugs first (grilling decision)
The 3 documented drifts (imaging Patient gRPC yes/HTTP no; medication
Admin+Doctor gRPC vs Doctor HTTP; audit-logs public gRPC vs authed HTTP)
must be confirmed as unintended before unify. If any is intentional,
document it in the policy, don't "fix" it silently.

### 2.1 New module `backend/internal/app/policy/`
- `policy.go` — the single source of truth: one map of
  `{service/method | route pattern} → allowed roles`, plus `public` set.
  Keep the existing `RoleAdmin/Doctor/Nurse/Reception/Patient` vocabulary.
- `policy_test.go` — table tests asserting every endpoint's roles; a test
  that gRPC and HTTP maps agree (drift guard).

### 2.2 Rewire both transports
- `backend/internal/app/interceptor/permissions.go` — read from policy module
  (keep deny-by-default + 6 public methods behavior identical).
- HTTP: replace per-route `middleware.RequireRoles` across the 17
  `http_handler.go` files with a single route→roles lookup against the same
  policy module. Route registration stays in handlers; the roles source moves.
- Fold `notificationRoleRoutes` (`modules/notifications/service.go`) into the
  policy module.

### 2.3 Tests
- Policy table tests (above) become the drift guard going forward.
- Interceptor + middleware tests updated to the new source.

### Validation
```go
go test ./backend/internal/...
```
Commit: `refactor(security): unify RBAC policy across gRPC and HTTP`

---

## Phase 3 — Single i18n accessor (candidate 5)

Fix the dual-convention i18n leak (raw `patients.validation.*` keys rendered
in UI) before touching the modals, since both modals and schemas depend on it.

### 3.1 `frontend/src/shared/i18n/i18n.ts`
- Collapse to one convention: a single module-bound accessor factory, e.g.
  `createModuleTranslator(namespace)` returning `t(key)` that resolves within
  the module namespace — kills `useTranslation("patients")` +
  `resolvePatientsKey` manual prefixing.
- Fix resource registration so nested keys resolve under the module namespace
  (remove the whole-tree-as-`translation` double registration).

### 3.2 Update consumers
- `frontend/src/modules/patients/patient_schemas.ts` — replace
  `resolvePatientsKey(t)` with the module accessor; update unit tests that
  currently pin the literal leaked keys to assert translated output.
- All `useTranslation()` call sites updated to the accessor convention.
- Add Vitest unit test for the accessor (translation resolution + fallback).

### Validation
```bash
cd frontend && npm run lint && npm run test
```
Commit: `refactor(i18n): single module-bound translation accessor`

---

## Phase 4 — Collapse the clinical form modals (candidate 4)

Build on Phase 3 (schemas no longer carry the i18n workaround).

### 4.1 One deep form module
- `frontend/src/modules/patients/components/ClinicalFormModal/`
  - `ClinicalFormModal.tsx` — owns `useForm` + `zodResolver` + Dialog
    header/footer + submit handling + camelCase↔snake_case mapping.
  - Config-driven: props carry schema, field option lists (LOINC/ICD-10),
    transform flags (uppercase ICD-10), and which optional fields to render
    (e.g. practitioner select). The 6 existing modals' deltas become config.
- Move business mapping currently in parents
  (`VitalSigns.tsx:26-52`, `ClinicalReports.tsx:29-43`,
  `EncounterHistory.tsx:31-44`) into the module so the leak closes.
- Remove schema twin tiers: drop `get*Schema` indirection, keep single schema
  factory taking the module translator.

### 4.2 Migrate call sites
- Replace `{Report,Encounter,Observation,Condition,Allergy,Medication}Modal.tsx`
  with config instantiations of `ClinicalFormModal`; delete the 6 skeletons.
  `ExamAnalyzerModal`/`ReportVersionsModal` stay (different surfaces) — verify
  they don't share the skeleton before touching them.

### 4.3 Tests
- Vitest unit test for the form module core (submit mapping, transform flags).
- Update `EncounterModal.test.tsx` / `ReportModal.test.tsx` to the new surface.
- **E2E (Playwright)**: update `frontend/e2e/` patient-module specs to the new
  modal interaction; all network/auth calls routed via `helpers.ts`
  (`page.route()`), fully offline. Test titles in English.
- i18n unit tests from Phase 3 stay green.

### Validation
```bash
cd frontend && npm run lint && npm run test && npm run e2e
```
Commit: `refactor(patients): consolidate clinical form modals into one module`

---

## Phase 5 — Typed eventbus contract (candidate 6, speculative)

Lowest priority; only after 0–4 are stable. Bus is small today.

### 5.1 `backend/internal/shared/eventbus/eventbus.go`
- Keep the bus, change the envelope: payloads become typed event structs
  registered per topic (e.g. `AppointmentCreated` with typed fields), replacing
  `map[string]any` + string-key type assertions.
- `notifications` subscriber consumes the typed events; remove hand-rolled
  `"title"/"body"/"actor_id"` parsing.
- Decide lifecycle of the 2 dead events (`appointment.created`,
  `appointment.cancelled`): either wire a subscriber or drop the publish —
  do not leave convention-only publishers.

### 5.2 Tests
- `eventbus_test.go` — typed publish/subscribe round-trip, compile-time
  contract, dead-topic detection.

### Validation
```go
go build ./... && go test ./backend/internal/shared/eventbus/... ./backend/internal/modules/notifications/...
```
Commit: `refactor(eventbus): type event payloads end-to-end`

---

## Cross-cutting rules (all phases)
- **Zero comments**; descriptive names only (AGENTS.md).
- Backend hexagonal layout: model/service/repository/handler roles preserved.
- Every service/handler logic change carries a unit test (MockRepository,
  same package).
- Every frontend flow change carries Vitest unit + Playwright E2E (offline,
  `page.route()`, English titles).
- Endpoints: keep `methodPermissions`/policy deny-by-default — nothing new is
  callable until registered in the policy.
- Commit + push after each phase (Conventional Commits), working tree clean.

## Sequencing rationale
- Phase 0 first: 8× duplication is highest leverage and unsticks 8
  repository.go files used by Phase 1.
- Phase 1 before Phase 2: both touch the same http_handler.go files; moving
  validation first avoids reworking role wiring twice.
- Phase 3 before Phase 4: schemas (i18n workaround) feed the modals.
- Phase 5 last: smallest payoff, speculative.
