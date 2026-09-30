# Hardcoded UI Text and Raw-String Enum Candidates

**Status:** research note — no source files were modified.
**Scope:** full monorepo (`backend/` Go + gRPC/proto, `frontend/` React + TypeScript).
**Method:** `rg`/`Select-String` inventory followed by line-by-line reads of every candidate. Every finding below is anchored to `file:line`. Counts are reproducible with the commands in [Appendix C](#appendix-c--reproduction-commands).

---

## Executive summary

The i18n groundwork is genuinely good: 3 locales × 13 namespaces, a factory-based translator, and 25/25 `ColumnDef` headers correctly routed through `t()`. The problems are concentrated, not diffuse — and the two biggest ones are not in the React tree at all.

| Signal | Count |
| --- | --- |
| Tier 1 findings (i18n bypass in scope) | 58 sites / 20 files |
| Tier 2 findings (user-visible, out of Tier 1 scope, or duplicate-key) | 9 findings / 8 areas |
| Tier 3 findings (intentionally not translatable) | 12 sites |
| Hardcoded pt-BR backend `render.Error` messages surfaced to the UI | **71 across 16 files** |
| `enum` declarations in `backend/proto` | **0** (across 13 files, 224 `string` fields) |
| Enums whose *values* are Portuguese display text | 2 of 12 in `shared/types/index.ts` |

Three findings dominate:

1. **The axios interceptor pipes raw backend pt-BR error strings straight into the UI.** `frontend/src/shared/services/api.ts:44-47` overwrites `requestError.message` with the response body's `error` field. 71 `render.Error(...)` calls across 16 Go handlers hardcode Portuguese. Switching the UI to `en-US` does nothing for any of them.
2. **`backend/proto` has no enums at all.** The Go layer *does* define named types (`role.Role`, `notifications.NotificationType`, `schedule.AppointmentStatus`) — the closed sets exist in-process and evaporate on the wire.
3. **The `portal` module is half-migrated.** A locale namespace exists with 2 keys, five of nine files already destructure `t`, and `PortalAppointments.tsx` is fully clean — yet the other eight files hardcode 34 strings and render 5 raw status values. The pattern to follow already exists one directory over.

---

## Part A — Hardcoded user-facing text

### Tier definitions

| Tier | Criterion |
| --- | --- |
| **1** | Text in a file that already calls `useTranslation` / `t()`, **or** any user-visible text under `frontend/src/modules/**` or `frontend/src/shared/components/**` |
| **2** | User-visible text with no i18n path, outside the Tier 1 location scope, **or** text duplicating / colliding with an existing locale key |
| **3** | Proper nouns, medical codes, units, brands, protocol IDs, machine-facing and log strings — intentionally not translatable |

### Tier 1 — 58 sites across 20 files

#### A.1 `portal` module — 8 files, 34 hardcoded sites + 5 raw status renders (39 total)

The `portal` locale namespace **exists in all three locales** (`*/portal.json`) but contains only two keys, `emptyEncounters` and `emptyConditions`. Five of the nine portal files already import a translator. The module is therefore not "un-i18n'd" — it is **half-migrated**, which is what makes it dangerous: the infrastructure is there and unused.

**Files that call `useTranslation` and still hardcode text (4 files, 13 sites):**

| File | Translator | Working `t()` calls | Hardcoded |
| --- | --- | --- | --- |
| `PortalDashboardOverview.tsx` | `:15` `useTranslation("patients")` | `:107` | `:20` `Total de Consultas`, `:26` `Condições Ativas`, `:32` `Medicamentos Ativos`, `:38` `Exames Recentes`, `:48` `Olá, {name}!`, `:51` welcome paragraph, `:71` `Próximas Consultas`, `:97` `Últimos Sinais Vitais` |
| `PortalEncounters.tsx` | `:9` `useTranslation("portal")` | `:24` | `:35,36,37` `Data`, `Motivo`, `Status` |
| `PortalConditions.tsx` | `:27` `useTranslation("portal")` | `:42` | `:58` `Início:` + `"N/I"` |
| `PortalObservations.tsx` | `:11` `useTranslation("patients")` | `:41` | `:26` `Nenhum sinal vital registrado.` |

`PortalDashboardOverview.tsx:48` is the sharpest case: it interpolates a translated-adjacent greeting around a real patient name while ignoring the `t` it already destructured two lines earlier.

**Files with no translator at all (4 files, 21 sites):**

| File | Hardcoded |
| --- | --- |
| `PortalPage.tsx` | `:37-44` eight sidebar labels (`Visão Geral`, `Agendamentos`, `Consultas`, `Sinais Vitais`, `Condições`, `Medicamentos`, `Exames`, `Imagens`), `:61` `Carregando portal...`, `:80` `Portal do Paciente`, `:88` `Navegação` |
| `PortalMedications.tsx` | `:38` `Nenhum medicamento prescrito.`, `:52` `"Medicamento"`, `:54` `Prescrito em:`, `:57` `"N/I"` |
| `PortalReports.tsx` | `:40` `Nenhum exame encontrado.`, `:54` `"Exame"`, `:58` `"N/I"` |
| `PortalImaging.tsx` | `:22` `Nenhum exame de imagem encontrado.`, `:36` `"N/I"`, `:45` `"Estudo de Imagem"` |

**Raw status rendered instead of translated — the same defect, five times:**

| File:line | Renders |
| --- | --- |
| `frontend/src/modules/portal/PortalDashboardOverview.tsx:87` | `{encounter.status}` |
| `frontend/src/modules/portal/PortalEncounters.tsx:49` | `{encounter.status}` |
| `frontend/src/modules/portal/PortalMedications.tsx:63` | `{medication.status}` |
| `frontend/src/modules/portal/PortalReports.tsx:70` | `{report.status}` |
| `frontend/src/modules/portal/PortalImaging.tsx:48` | `{study.status}` |

> **Contrast — the fix already exists, twice, in the same module.**
> `frontend/src/modules/portal/PortalAppointments.tsx:58` does `t(\`status.${appointment.status}\`)` against `schedule.json:9-14`, and its four table headers at `:40-43` resolve against the `portal` subtree at `schedule.json:82-88`. That file is the only one of the nine with **zero** hardcoded user-facing text. The other eight need to be brought in line with it — no new key-naming invention, no new namespace, no architectural change.
>
> `PortalAppointments.tsx:7` is nonetheless a Part B item: `Record<string, string>` with no union, so a typo in a status key fails silently.


#### A.2 `shared/components/**` — files that already have a translator (3 sites)

| File:line | Hardcoded text | Why Tier 1 |
| --- | --- | --- |
| `frontend/src/shared/components/ui/Dialog.tsx:47` | `Close` (`sr-only`) | Under `shared/components/**`; screen-reader copy |
| `frontend/src/shared/components/ui/ToastContainer.tsx:36-38` | `Sucesso`, `Erro`, `Aviso` | Toast variant labels; the component is the notification surface |
| `frontend/src/shared/components/ErrorBoundary.tsx:51` | `"Unknown Error"` | **Same file already calls `i18next.t` at lines 44, 47, 55** |

#### A.3 Module files that already use `useTranslation` (8 sites)

| File:line | Hardcoded text | Note |
| --- | --- | --- |
| `frontend/src/modules/staff/schemas/staff_schemas.ts:5` | `O nome deve ter no mínimo 3 caracteres` | zod message, rendered by `StaffModal.tsx:97` |
| `frontend/src/modules/staff/schemas/staff_schemas.ts:7` | `Selecione uma categoria válida` | rendered by `StaffModal.tsx:124` |
| `frontend/src/modules/staff/schemas/staff_schemas.ts:19` | `Formato inválido. Ex: CRM-SP 12345` | rendered by `StaffModal.tsx:134` |
| `frontend/src/modules/staff/schemas/staff_schemas.ts:21` | `E-mail inválido` | rendered by `StaffModal.tsx:146` |
| `frontend/src/modules/staff/components/StaffModal.tsx:114-117` | raw `DOCTOR` / `NURSE` / `RECEPTION` / `ADMIN` as `<option>` text | machine enum values rendered as visible labels |
| `frontend/src/modules/schedule/components/ScheduleViewToggle.tsx:13` | `aria-label="Calendar view"` | the only user-visible string-prop hit in the entire app |
| `frontend/src/modules/schedule/components/AppointmentDetailsModal.tsx:64` | `` `${fullName} — ${role}` `` | appends raw `DOCTOR` after the translated name |
| `frontend/src/modules/imaging/ImagingWorkspace.tsx:47` | `"Nova Ressonância Magnética"` | fallback study title; same file uses `t()` at 40, 52, 54, 56 |

> **Contrast — three of four schema files use the correct factory pattern.**
> `frontend/src/modules/patients/patient_schemas.ts:77-150`, `frontend/src/modules/schedule/schedule_schemas.ts:19,103`, and `frontend/src/modules/auth/auth_schemas.ts:3` all take a `translateFunction` parameter and call `translateFunction("validation.…")`. `staff_schemas.ts` is the sole holdout. The migration is mechanical: convert to a factory and pass `createModuleTranslator("staff")` at the call site.

#### A.4 Module components with no translator at all (8 sites / 6 files)

| File:line | Hardcoded text |
| --- | --- |
| `frontend/src/modules/imaging/components/DicomToolControls.tsx:21` | `Zoom (Arrastar)` |
| `frontend/src/modules/imaging/components/DicomToolControls.tsx:29` | `Luminosidade` |
| `frontend/src/modules/imaging/components/DicomToolControls.tsx:37` | `Régua (Medir)` |
| `frontend/src/modules/imaging/components/DicomToolControls.tsx:43,46,49` | `Tecido Mole`, `Osso`, `Pulmão` |
| `frontend/src/modules/telemetry/components/TelemetryBedList.tsx:68` | `FC: ` (Portuguese-only abbreviation for heart rate) |
| `frontend/src/modules/patients/PatientDetails.tsx:46` | `Carregando componente...` |
| `frontend/src/modules/patients/PatientDetails.tsx:186` | `Analisar Exame com IA` |
| `frontend/src/modules/patients/EncounterDetail.tsx:29` | `Carregando componente...` |

---

### Tier 2 — 9 findings

#### B.1 Backend error messages reach the UI untranslated (highest impact)

`frontend/src/shared/services/api.ts:44-47`:

```ts
const errorPayload = requestError.response.data as { error?: string }
if (errorPayload?.error) {
  requestError.message = errorPayload.error
}
```

The interceptor deliberately promotes the server's `error` string to the user-facing message. **71 `render.Error` calls across 16 Go files hardcode Portuguese**, out of 94 total `render.Error` calls:

`allergy`, `analytics`, `audit_logs`, `auth`, `condition`, `diagnostic_report`, `encounter`, `exam_analyzer`, `medication`, `observation`, `patients`, `portal`, `schedule` (×2 files), `staff`, `telemetry` — all under `backend/internal/modules/*/http_handler.go`.

Representative anchors:

| File:line | Message |
| --- | --- |
| `backend/internal/modules/auth/http_handler.go:60` | `Credenciais inválidas.` |
| `backend/internal/modules/auth/http_handler.go:181` | `Usuário não encontrado.` |
| `backend/internal/modules/portal/http_handler.go:41` | `Nenhum prontuário vinculado a este usuário.` |
| `backend/internal/modules/schedule/http_handler.go:179` | `Usuário não autenticado.` |
| `backend/internal/modules/exam_analyzer/http_handler.go:121` | `Formato de arquivo não suportado. Envie PNG, JPG ou PDF.` |
| `backend/internal/modules/exam_analyzer/http_handler.go:288` | `Análise e arquivo excluídos com sucesso.` (in a `200 OK` body) |

Two structural problems compound this:
- The remaining 23 `render.Error` calls are *English* or machine-oriented, so the API surface is not even internally consistent.
- `render.JSON(..., {"success": "Análise e arquivo excluídos com sucesso."})` puts user-facing prose in a success payload, which no i18n layer can intercept.

**Fix direction:** return a stable machine code (e.g. `{"error":"FILE_TYPE_UNSUPPORTED","params":{"ext":["png","jpg","pdf"]}}`) and let `api.ts` resolve it against a `common` namespace. This is a wire-format change and must be versioned.

#### B.2 Backend-generated notifications are rendered verbatim

`frontend/src/modules/notifications/components/NotificationItem.tsx:40,42` renders `notification.title` and `notification.body` as-is. Seven Go sites build that prose:

| File:line | Title / body |
| --- | --- |
| `backend/internal/modules/telemetry/service.go:102-103` | `Alerta Clínico - Leito <n>` / `Paciente %s apresenta condição %s (BPM: %d, SpO2: %d%%).` |
| `backend/internal/modules/telemetry/simulator.go:82-83` | duplicate of the above |
| `backend/internal/modules/auth/service.go:110-111` | `Login Realizado` / `Login realizado com sucesso` |
| `backend/internal/modules/encounter/service.go:63-64` | `Novo Atendimento Criado` / encounter prose |
| `backend/internal/modules/diagnostic_report/service.go:184-185` | `Laudo Disponível` / `O laudo do paciente %s está pronto para consulta.` |
| `backend/internal/modules/exam_analyzer/worker.go:177-178` | `Análise de Exame Concluída` / `Laudo disponível para revisão: ` |
| `backend/internal/modules/exam_analyzer/worker.go:131` | `Não foi possível gerar análise confiável devido à qualidade ou ilegibilidade do arquivo enviado.` |

Note `telemetry/service.go:102-103` and `telemetry/simulator.go:82-83` are byte-identical duplicates.

#### B.3 `AppHeader` role keys can never resolve

`frontend/src/shared/components/AppHeader.tsx:31` builds `` t(`roles.${userRole}`, { defaultValue: t("roles.RoleDefault") }) ``. The wire value is a raw string (`ADMIN`, `DOCTOR`, …) from `frontend/src/shared/store/auth_store.ts:8,23-29`. The locale files define `RoleAdmin`, `RoleDoctor`, `RoleNurse`, `RoleReception` under `roles` (each `*/header.json:6-12`). **No key can ever match**, and the explicit `defaultValue` makes every role render as `RoleDefault` (`Profissional` / `Professional`) — including doctors and admins. The bug is invisible precisely *because* the fallback is specified.

#### B.4 `staff` role labels use missing keys with pt-BR defaults

`frontend/src/modules/staff/components/StaffTable.tsx:26-29` and `StaffFilters.tsx:38-41` call `t("table.roles.doctor", "Médico")` and similar. **No `roles` block exists in any `staff.json`** (`pt-BR`, `en-US`, `es-ES`). The pt-BR `defaultValue` masks the missing key in the default locale, so the bug is invisible in development and produces Portuguese text in `en-US`/`es-ES`.

`StaffTable.tsx:30` compounds it: the `default:` branch of the same switch is `return role`, so any unmapped role renders the raw wire token rather than even the Portuguese default.

#### B.5 `StaffStatus` derives a display string from a boolean

`frontend/src/modules/staff/api.ts:31` maps `is_active` to `StaffStatus.OnDuty` / `OffDuty`, whose values are `Plantonista` / `Fora de Escala` (`frontend/src/shared/types/index.ts:27-30`). `StaffTable.tsx:89` renders `member.status` directly. The correct translations already exist, authored and unused: `staff:table.onDuty` and `staff:table.offDuty` (`*/staff.json:15-16`).

#### B.6 Fabricated and duplicated data

- `frontend/src/modules/staff/api.ts:32` — `department: "Geral"` is invented client-side.
- `frontend/src/modules/staff/api.ts:5-13` — a local `mapRole` duplicates the wire-role list and uses an unchecked `role as StaffRole` cast.

#### B.7 Static shell and PWA metadata

| File:line | Text |
| --- | --- |
| `frontend/index.html:2` | `lang="en"` — contradicts the `pt-BR` fallback in `i18n.ts:9-22` |
| `frontend/index.html:16` | `<title>HealthCare Portal</title>` |
| `frontend/public/manifest.json` | `name`, `short_name`, and a Portuguese `description` |

There are **zero** `document.title` assignments anywhere in app source, E2E, or `index.html` — the document title is never localized per route.

#### B.8 Clinical data written to FHIR as Portuguese

`backend/internal/modules/observation/metrics.go:11-18` hardcodes all 8 `CodeDisplay` values (`Frequência Cardíaca`, `Pressão Arterial Sistólica`, `Saturação de Oxigênio`, `Frequência Respiratória`, `Temperatura Corporal`, `Peso Corporal`, `Altura Corporal`, …). These become `Coding.display` in the Healthcare API. Once written, the store only ever holds pt-BR, so the defect is **persistent, not cosmetic**.

`backend/internal/modules/analytics/service.go:96-101` ships a `displayNameMap` whose four values are pt-BR: `CT (Tomografia)`, `MR (Ressonância)`, `CR (Raio-X)`, `US (Ultrassom)`. This one *is* fixable purely in the frontend, since the DICOM modality code is already transmitted alongside.

#### B.9 `exam_analyzer` fallback report corpus is entirely pt-BR

`backend/internal/modules/exam_analyzer/service.go:228-386` — `runHeuristicSimulation` returns one of four complete clinical reports (radiograph, cranial MRI, general photograph, haematology), with every `Finding`, `NextSteps`, `Limitations`, `ExamType`, and `Disclaimer` hardcoded in Portuguese. These render directly in the analysis UI.

The LLM prompt at `service.go:122-127` is a related design decision worth flagging rather than a bug: line 123 instructs the model *"Provide a structured clinical support analysis in Portuguese."* The output language is pinned to pt-BR at the prompt level, so **no amount of frontend i18n will localize AI-generated analysis text.** Whether that is correct is a product decision — but it should be a decision, not an accident.

Note the disclaimer at `service.go:272,316,354,384` is a legal/regulatory statement hardcoded four separate times. A single `const` would prevent the copies from drifting.

---

### Tier 3 — 12 sites intentionally not translated

| File:line | Text | Justification |
| --- | --- | --- |
| `frontend/src/modules/auth/Login.tsx:19` | `HealthCare` | Brand |
| `frontend/src/modules/patients/components/PatientsMetricsGrid.tsx:35` | `FHIR R4 Compliant` | Standard/spec name |
| `frontend/src/modules/patients/components/PatientsMetricsGrid.tsx:47` | `Cloud Healthcare API` | Product name |
| `frontend/src/shared/components/LanguageSwitcher.tsx:16-18` | `Português`, `English`, `Español` | Language autonyms — must not be translated |
| `frontend/src/modules/telemetry/components/TelemetryBedMonitor.tsx:93` | `BPM` | Unit |
| `frontend/src/modules/telemetry/components/TelemetryBedList.tsx:69` | `SpO₂` | Universal clinical symbol |
| `frontend/src/modules/patients/components/VitalSignValueDisplay.tsx:14` | `N/A` | Universal abbreviation |
| `frontend/src/modules/imaging/ImagingWorkspace.tsx:48` | `"MR"` | DICOM modality code |
| LOINC literals in the shared/vitals vocabulary — `shared/types/index.ts:105-111`, `patients/patient_schemas.ts:38-45`, `patients/components/vitalSignDisplay.ts:10-19` | `8867-4`, `8310-5`, … | LOINC codes. Correctly not translatable, but declared in 3 places — a Part B defect (B.4) |
| LOINC literals in `patients/components/ClinicalFormModal/clinicalFormConfigs.ts:34-37` | `58410-2`, `2345-7`, `24323-8`, `2093-3` | LOINC codes used as option `value`s. Already correct — each pairs with a `labelKey`, which is the pattern the rest of the app should copy |
| `frontend/src/modules/staff/components/StaffTable.tsx:21` | `console.log({ filteredStaff })` | Debug logging — **not** a translation issue, but should not ship |
| `frontend/src/shared/components/ErrorBoundary.tsx:28` | `console.error(...)` | Diagnostic logging |

**Explicitly ruled out — not mojibake.** The PowerShell console renders `í`, `ã`, `é`, `°` and `—` as replacement glyphs, so `staff_schemas.ts` *looks* corrupted on screen (`m??nimo`, `v??lida`) while being perfectly intact on disk. Strict-decoder validation of `PortalConditions.tsx`, `DicomToolControls.tsx`, `TelemetryBedMonitor.tsx`, and `staff_schemas.ts` (743 bytes) confirms **all four are strict-valid UTF-8**. There is no encoding corruption anywhere in the audited source. Do not "fix" these.

---

### Requested high-signal classes

| # | Class | Result | Count |
| --- | --- | --- | --- |
| 1 | Inline copy inside a ternary / `&&` in JSX | The text-node sweep returns **31 lines**, but the vast majority are TypeScript comparison operators (`currentPhase >= 6 && currentPhase < 12`) and a generic type signature, not JSX text. Triaged: the genuine user-visible text-node hits are `Dialog.tsx:47` and `TelemetryBedList.tsx:68` (Tier 1) plus `TelemetryBedMonitor.tsx:93` (Tier 3). **No ternary/`&&` violations.** | **31 raw / 0 in ternaries** |
| 2 | `placeholder` / `aria-label` / `title` / `alt` string props | Single hit across all of `frontend/src` | **1** (`ScheduleViewToggle.tsx:13`) |
| 3 | `document.title` / `<title>` / manifest | `document.title`: 0 (command exit 1). `index.html:16` + `manifest.json` | **2** |
| 4 | `toast.*` / `new Error` / `setError` literals | All toast calls route through `t()`. Backend `render.Error` is the real leak. | **0** frontend / **71** backend |
| 5 | Zod validation messages | 3 of 4 schema files use the factory pattern (verified 0 hits each); `staff_schemas.ts` hardcodes 4, spread across three distinct zod syntax forms so no single regex catches them | **4** |
| 6 | Duplicate literal sets across files | `N/I` ×5, `Carregando componente...` ×2, DICOM controls, telemetry copy | **~12** |
| 7 | Untranslated enum values rendered as labels | `StaffModal.tsx:114-117` (4 values), `AppointmentDetailsModal.tsx:64`, 5× portal status renders | **10 values / 4 files** |
| 8 | `ColumnDef` `header:` without `t()` | **All 25 headers use `t()`.** Fully compliant. | **0** |

### Brittle tests that lock the defects in

Fixing any Tier 1 item requires updating these in the same commit, or CI goes red:

| File:line | Asserts |
| --- | --- |
| `frontend/e2e/schedule_calendar.spec.ts:69,76` | `"Calendar view"` |
| `frontend/e2e/staff-form-validation.spec.ts:77` | `"NURSE"` |
| `frontend/e2e/telemetry_stats_staff.spec.ts:118` | `"NURSE"` |
| `frontend/e2e/imaging.spec.ts:23-25,37-39` | DICOM control labels |
| `frontend/src/modules/schedule/components/AppointmentDetailsModal.test.tsx:60` | `"Dr. André Silva — DOCTOR"` |
| `frontend/src/shared/components/ui/ToastContainer.test.tsx:29` | `"Sucesso"` |
| `frontend/src/modules/notifications/components/NotificationItem.test.tsx:21-22` | backend notification prose |

Conversely, tests that mock `t` as an identity function (`t: (key) => key`) are locale-agnostic and will survive a correct fix.

---

## Part B — Raw strings that should be closed enums

### The core observation

`backend/proto` contains **13 files, 224 `string` fields, and 0 `enum` declarations**:

```
rg -P "^enum\s" backend/proto -g "*.proto"     → 0 matches
```

32 of those string fields carry closed-domain names and are the enum-worthy set:

```
rg -P "^\s*string\s+(status|condition|role|priority|notification_type|type|modality|clinical_status|caller_role|sort_field|sort_direction|gender)\s*=" backend/proto
```

Meanwhile the Go layer *does* model three of them as named types. **The closed sets exist in-process and are lost at the wire boundary.**

### B.1 `backend/proto` — enum-worthy fields (32)

| File:line | Field | Proposed type |
| --- | --- | --- |
| `telemetry.proto:30,31` | `TelemetryBed.status`, `.condition` | `BedStatus`, `CardiacCondition` |
| `telemetry.proto:63,64` | `UpdateBedConditionRequest.status`, `.condition` | `BedStatus`, `CardiacCondition` |
| `telemetry.proto:26` | `BedPatient.gender` | `Gender` |
| `staff.proto:18,34` | `role` (create/update) | `StaffRole` |
| `auth.proto:21,28` | `role` (response / register) | `StaffRole` |
| `audit_logs.proto:15,33` | `caller_role` | `CallerRole` |
| `audit_logs.proto:16,34` | `method` | `AuditAction` |
| `encounter.proto:31,42` / `clinical.proto:45,102` | `status` | `EncounterStatus` |
| `medication.proto:33` / `clinical.proto:157` | `status` | `MedicationRequestStatus` |
| `imaging.proto:29,43` | `status` | `ImagingStudyStatus` |
| `imaging.proto:24,41` | `modality` | `DicomModality` |
| `diagnostic_report.proto:31` / `clinical.proto:183` | `status` | `DiagnosticReportStatus` |
| `condition.proto:17,32` / `clinical.proto:87,102` | `clinical_status` | `ConditionClinicalStatus` |
| `allergy.proto:16,32` / `clinical.proto:113,128` | `clinical_status` | `AllergyClinicalStatus` |
| `patients.proto:45,46` | `sort_field`, `sort_direction` | `PatientSortField`, `SortDirection` |

**Blast radius:** all 13 protos, all generated Go stubs, all generated TS stubs, and both wire clients. **This is the highest-risk item in Part B** and must not be bundled with a UI change.

**Mitigation:** add the enums and keep the `string` fields populated in a follow-up; only remove the string once both clients are regenerated. Because the current values are already closed sets, no data migration is implied — but existing FHIR resources already carry whatever was written, so the mapping must tolerate unknown values.

### B.2 `CardiacCondition` — display strings used as business-logic values (highest-risk semantic bug)

Exact observed literal set (`frontend/src/shared/types/index.ts:1-6`):

```ts
export const CardiacCondition = {
  Normal: "Normal",
  Bradycardia: "Bradicardia",
  Tachycardia: "Taquicardia",
  CardiacArrest: "Parada Cardíaca"
} as const
```

**Backend business logic branches on these Portuguese strings:**

`backend/internal/modules/telemetry/simulator.go:114`

```go
case "Parada Cardíaca":
```

This is a switch on a human-readable label. It is the clearest possible demonstration of why display text must not be a stored value.

- **Proposed:** `CardiacCondition = "normal" | "bradycardia" | "tachycardia" | "cardiac-arrest"`, with `t("telemetry.conditions.cardiacArrest")` for display.
- **Blast radius:** `shared/types/index.ts:1-8`; `telemetry/types.ts` (`BedPatient.condition`, `UpdateBedConditionPayload.condition`); `TelemetryBedList.tsx:59` (renders raw); `TelemetryBedMonitor.tsx` (5 comparisons); `telemetry/model.go:29` (`Condition string`); `telemetry/service.go:103`; `telemetry/simulator.go:82,114`; `telemetry.proto:31,64`.
- **Risk:** **High.** Values are persisted in FHIR and cross the proto boundary. Needs a coexistence window accepting both the old and new vocabularies.

### B.3 `StaffStatus` — same anti-pattern, ASCII-only (easy to miss)

`frontend/src/shared/types/index.ts:27-30`:

```ts
export const StaffStatus = {
  OnDuty: "Plantonista",
  OffDuty: "Fora de Escala"
} as const
```

These are pure ASCII, so a diacritic-based audit finds nothing. They are nonetheless pt-BR display text. `staff/api.ts:31` synthesizes them from `is_active`.

- **Proposed:** `type StaffStatus = "on-duty" | "off-duty"`.
- **Blast radius:** `shared/types/index.ts:27-32`; `staff/api.ts:31`; `StaffTable.tsx:89`; locale keys already exist at `staff:table.onDuty` / `staff:table.offDuty` (`*/staff.json:15-16`).
- **Risk:** **Low.** Purely client-side, not persisted, translations already authored.

### B.4 Enums duplicated across the frontend

Twelve const-object enums live in `frontend/src/shared/types/index.ts` (137 lines), but three of them are shadowed by local redefinitions:

| Canonical | Shadowed by | Problem |
| --- | --- | --- |
| `AppointmentStatus` (`schedule/types.ts:1`) | `AppointmentStatusLabel` (`schedule/types.ts:32-37`) | **Dead code.** The map is an identity (`scheduled: "scheduled"`) and is referenced nowhere but its own declaration |
| `ScheduleViewMode` (`schedule/components/ScheduleCalendar.tsx:12`) | `VIEW_MODES` array (`ScheduleViewToggle.tsx:8`) | Token list duplicated |
| `LoincCode` (`shared/types/index.ts:105-111`) | `vitalSignMetricDefinitions` (`patients/patient_schemas.ts:38-45`) and `vitalSignDisplayMetadata` (`patients/components/vitalSignDisplay.ts:10-19`) | LOINC codes declared in **3 places** |

**`LoincCode` is worse than duplicated — it is dead and stale.** It is exported at `shared/types/index.ts:105,111` and imported by nothing. Its `BloodPressure: "85354-9"` value is not one of the app's eight vital-sign metrics, and `frontend/src/modules/patients/components/vitalSignDisplay.test.tsx:30` *asserts* that `findVitalSignDisplay("85354-9")` returns `undefined`. The test suite documents the dead map as obsolete; the map was never removed.

`DicomModality` (`shared/types/index.ts:113-125`) has a related coverage gap: it declares 9 modalities, but the backend `displayNameMap` (`analytics/service.go:96-101`) only supplies pt-BR labels for 4 (`CT`, `MR`, `CR`, `US`). `DX`, `XA`, `PT`, `MG`, and `SR` have no label in any locale.

Two of these three also use untyped inline unions rather than the const-object pattern:

- `frontend/src/modules/schedule/types.ts:1` — `export type AppointmentStatus = "scheduled" | "confirmed" | "cancelled" | "finished"`
- `frontend/src/modules/portal/PortalAppointments.tsx:7-12` — `Record<string, string>` for the same domain, with no union at all, so a mistyped status key silently falls through to the unstyled branch

- **Proposed:** delete `AppointmentStatusLabel` and `LoincCode` outright; derive `VIEW_MODES` via `Object.values()`; make `vitalSignMetricDefinitions` the single LOINC source; type the portal map as `Record<AppointmentStatus, string>`; extend `DicomModality` coverage in the locale files.
- **Risk:** **Low.** No wire or persistence impact; both deletions are provably unreferenced.

### B.5 `shared/auth/types.ts` — the correct Go analogue, already written

`frontend/src/shared/auth/types.ts:1-29` defines `Action` and `Feature` as `as const` objects with derived union types, using FHIR resource names as values. This is the pattern `CardiacCondition` and `StaffStatus` should follow, and the Go equivalent already exists at `backend/internal/shared/role/role.go:5-13`. No change needed — cited as the reference implementation.

### B.6 Go named enums that already exist (positive baseline)

| File:line | Type |
| --- | --- |
| `backend/internal/shared/role/role.go:5-13` | `Role` |
| `backend/internal/modules/notifications/model.go:9-30` | `NotificationType` (9 values), `NotificationPriority` (4 values) |
| `backend/internal/modules/schedule/model.go:10-17` | `AppointmentStatus` (4 values) |

**The corresponding Go models drop the type.** `backend/internal/modules/telemetry/model.go:28-29` declares `Status string` / `Condition string` with no named type, even though `telemetry.proto:30,31,63,64` and `frontend/src/shared/types/index.ts:10,1` both model the same two closed sets. This is the clearest single example of the enum being lost in translation — it exists in the frontend types and in the Go enums, but not in the Go struct that actually holds the value.

**One consistency defect:** `backend/internal/shared/role/role.go:17-28` also accepts stray aliases `"RoleAdmin"`, `"RoleDoctor"`, etc. alongside the canonical values. This is migration debt, not a missing type — worth a cleanup ticket.

### B.7 Structural note — `clinical.proto` duplicates six per-resource protos

30 message names are declared twice, once in `backend/proto/clinical.proto` and once in the corresponding per-resource file:

`Condition`, `Encounter`, `Observation`, `MedicationRequest`, `DiagnosticReport`, `AllergyIntolerance` and all their `Create*`/`Get*` request/response pairs.

Any enum added to a resource message must be added in **both** files or the two definitions silently diverge. Resolve this before attempting the B.1 migration.

### Prioritization

**High value / low risk** — start here:

1. `StaffStatus` → `"on-duty" | "off-duty"` (B.3)
2. `staff.json` `table.roles.*` keys (B.4) — the translations may already be authored under another name; check first
3. Delete the dead `AppointmentStatusLabel` and `LoincCode` maps, derive `VIEW_MODES` (B.4)
4. `AppHeader` role key construction (B.3 in Part A)
5. `ScheduleViewToggle` `aria-label` + option labels (A.3)
6. The whole `portal` module (A.1) — copy `PortalAppointments.tsx` as the template; no new architecture required

**High value / high risk** — each needs its own branch, PR, and migration plan:

1. `CardiacCondition` → stable tokens (B.2) — backend switches on the literal
2. Backend `render.Error` → machine codes (B.1 in Part A) — wire-format change
3. `backend/proto` enums (B.1) — 224 fields, both clients regenerate
4. `clinical.proto` / per-resource duplication (B.7) — blocks B.1

**Low value:** Tier 3 items, the `console.log` at `StaffTable.tsx:21`.

---

## Recommendations

### Tooling

1. **Enable typed translation keys.** `frontend/src/shared/i18n/i18n.ts:9-22` registers resources with no type augmentation, so `t("roles.DOCTOR")` and `t("table.roles.doctor")` both compile. Adding i18next's `CustomTypeOptions.resources` augmentation makes every missing-key class in this report a **build error** rather than a silent fallback. This single change would have caught B.3 and B.4 automatically.
2. **Keep `interpolation.escapeValue: false`** (`i18n.ts` init block). It is correct here because the app renders through React, which escapes by default; note it in a comment so a future move to `dangerouslySetInnerHTML` is caught.
3. **Ban literal strings in JSX text nodes via ESLint.** `no-restricted-syntax` targeting JSX text containing a letter, and on `placeholder` / `aria-label` / `title` / `alt` attributes. The class-2 sweep returns exactly 1 hit today, so the rule is cheap and near-zero-false-positive.
4. **Require the schema-factory pattern.** Extend the existing convention from `patient_schemas.ts` / `schedule_schemas.ts` / `auth_schemas.ts` to `staff_schemas.ts` and forbid bare string literals as the second argument to any `z.*` validator.

### Process

5. **Every fix touches shared files** — `shared/i18n/locales/**`, `shared/components/**`, `app/routes.tsx`. Per `AGENTS.MD:19`, these follow the additive-anchored protocol: develop in a dedicated worktree, add entries in ordered/anchored position only, never reformat neighbouring lines.
6. **Add locale keys to all three locales in the same commit.** `pt-BR` / `en-US` / `es-ES` must move together; a key present in one and missing in the others produces a Portuguese fallback in production.
7. **Update the brittle tests in the same commit** (see the table above). A translation fix that leaves `e2e/schedule_calendar.spec.ts:69` asserting `"Calendar view"` is a broken fix.
8. **Extend the `portal` namespace rather than creating one.** It already exists in all three locales (`*/portal.json`) with 2 keys. The migration is to add keys to it, following the naming already established in `schedule.json:82-88` and `schedule.json:9-14`.

### Documentation

9. **The `kebab-case` filename** and the new `docs/research/` directory follow the repo's existing docs layout. This file is a point-in-time research note, not a living specification — it should be revised, not appended to, and it should not be treated as a source of truth for the enums, which belong in code.

---

## Appendix A — i18n configuration as found

| Property | Value | Location |
| --- | --- | --- |
| Libraries | `i18next ^26.3.0`, `react-i18next ^17.0.8` | `frontend/package.json:33-39` |
| Locales | `pt-BR` (default/fallback), `en-US`, `es-ES` | `frontend/src/shared/i18n/i18n.ts:9-22` |
| Namespaces | 13 per locale, 39 JSON files total | `frontend/src/shared/i18n/locales/en-US/index.ts:1-29` |
| Interpolation escaping | `escapeValue: false` | `i18n.ts` init block |
| Language detection | `LanguageDetector` + `initReactI18next` | `i18n.ts:9-22` |
| Module translator helper | `createModuleTranslator` | `i18n.ts:24-38` |
| Type safety | **None** — no `CustomTypeOptions` augmentation | — |
| Form / validation | `react-hook-form ^7.76.0`, `zod ^4.4.3` | `frontend/package.json:33-39` |

**Two subtleties in `createModuleTranslator` (`i18n.ts:24-38`) worth knowing before any fix:**

1. Lines 28-31 normalize both `namespace:key` and `namespace.key` to the same target, which is why the codebase mixes colon keys (`t("header:errorBoundary.title")` in `ErrorBoundary.tsx:44`) and dot keys (`t("table.roles.doctor")` in `StaffTable.tsx:26`). Both resolve. Pick one convention for new keys.
2. Line 35 tests `defaultValue` by **truthiness**, so an empty-string default is silently dropped and the raw key is returned instead. Not currently triggered, but it will bite any "hide this label" use case.

## Appendix B — Repository conventions consulted

| Convention | Location |
| --- | --- |
| Additive-anchored edits to shared files | `AGENTS.MD:19` |
| Table column definitions in a dedicated file | `AGENTS.MD:35` |
| Table column hooks using `useTranslation` | `AGENTS.MD:37` |
| Documented locale set | `CONTEXT.md:172` |
| Hexagonal architecture (`model.go` / `service.go` / `repository.go` / `register.go`) | `AGENTS.MD` |

## Appendix C — Reproduction commands

All commands run from the repository root and exclude `node_modules`, `dist`, `playwright-report`, `test-results`, `.git`, `backend/.gocache`, and `frontend/.vercel`.

> **Reproduction caveat — read before running on Windows PowerShell 5.1.** Two independent quoting failures corrupt any regex passed inline:
>
> 1. **Double quotes are stripped** when PowerShell hands an argument to a native executable, so `"[^"]*"` silently degrades to `[^*` and PCRE2 dies with `missing terminating ] for character class`.
> 2. **Non-ASCII literals in the pattern are re-encoded to the console codepage**, so a character class like `[áàâãéêíóôõúç]` arrives mangled and fails the same way.
>
> Both are avoided below by writing each pattern to an ASCII file with `Set-Content -Encoding ascii` and passing it via `rg -f`. Diacritics are written as `\x{...}` code-point escapes. Every count in this report was reproduced with the `-f` form.

```powershell
# Codebase size
rg --files -g '*.ts' -g '*.tsx' -g '*.go' -g '*.proto' -g '*.sql' `
  -g '!node_modules' -g '!dist' -g '!playwright-report' -g '!test-results' `
  -g '!.git' -g '!backend/.gocache' -g '!frontend/.vercel'

# Class 1 - JSX text nodes (frontend/src, excluding tests and locales) -> 31 raw, triaged
rg -n -P '>[^<{}]{0,120}?[A-Za-z][^<{}]{0,120}?<' frontend\src `
  -g '*.tsx' -g '!*.test.tsx' -g '!**/locales/**'

# Class 2 - user-visible string props -> 1 hit
Set-Content -LiteralPath "$env:TEMP\p2.txt" -Encoding ascii -Value `
  '\b(placeholder|aria-label|title|alt)\s*=\s*"[^"]*[A-Za-z][^"]*"'
rg -n -P -f "$env:TEMP\p2.txt" frontend\src -g '*.tsx' -g '!*.test.tsx'

# Class 3 - document title -> 0 matches, exit 1
rg -n 'document\.title' frontend\src frontend\e2e frontend\index.html

# Class 4 - backend error strings that reach the UI -> 71 across 16 files
#   \x{e1}\x{e0}\x{e2}\x{e3}\x{e9}\x{ea}\x{ed}\x{f3}\x{f4}\x{f5}\x{fa}\x{e7} = a a^ a~ a' e e^ i' o o^ o~ u' c- (pt-BR)
Set-Content -LiteralPath "$env:TEMP\p4.txt" -Encoding ascii -Value `
  'render\.Error\(\s*\w+,\s*http\.Status\w+,\s*"[^"]*[\x{e1}\x{e0}\x{e2}\x{e3}\x{e9}\x{ea}\x{ed}\x{f3}\x{f4}\x{f5}\x{fa}\x{e7}][^"]*"'
rg -n -P -f "$env:TEMP\p4.txt" backend\internal -g '*.go' -g '!*_test.go'

# Class 5 - zod messages bypassing the translator factory -> 4 hits
#   Zod messages appear in three syntax forms, so three patterns are unioned.
#   A: object-property form + positional with a preceding arg
Set-Content -LiteralPath "$env:TEMP\p5a.txt" -Encoding ascii -Value `
  'message:\s*"|\.\w+\(\s*[^,()]+,\s*"'
#   B: standalone message line (message of a multi-line .refine)
Set-Content -LiteralPath "$env:TEMP\p5b.txt" -Encoding ascii -Value '^\s*"[^"]+",?\s*$'
#   C: single string argument that contains a space (e.g. .email("..."))
Set-Content -LiteralPath "$env:TEMP\p5c.txt" -Encoding ascii -Value '\.\w+\(\s*"[^"]* [^"]*"'
foreach ($p in 'p5a','p5b','p5c') {
  rg -n -P -f "$env:TEMP\$p.txt" frontend\src\modules -g '*schemas*.ts' -g '!*.test.ts'
}   # -> staff_schemas.ts:5,7,19,21 only; the 3 factory schemas return 0

# Class 8 - table headers not routed through t() -> 0 matches, exit 1
rg -n -P 'header:\s*(?!.*\bt\()' frontend\src `
  -g '*Columns.tsx' -g 'use*Columns.tsx' -g '!*.test.tsx' -g '!*.test.ts'

# Part B - proto enums (expect zero) and string-field census
rg -P '^enum\s' backend\proto -g '*.proto'                     # 0 matches, exit 1
rg -c -P '^\s*(repeated\s+|optional\s+)?string\s+\w+\s*=' backend\proto -g '*.proto'   # 224 total
rg -n -P '^\s*string\s+(status|condition|role|priority|notification_type|type|modality|clinical_status|caller_role|sort_field|sort_direction|gender)\s*=' `
  backend\proto                                                       # 32 enum-worthy

# UTF-8 validation (rules out console mojibake)
#   Strict decode throws on any invalid byte sequence, unlike [Text.Encoding]::UTF8
#   which substitutes U+FFFD silently. Verified clean: PortalConditions.tsx,
#   DicomToolControls.tsx, TelemetryBedMonitor.tsx, staff_schemas.ts.
$utf8 = New-Object System.Text.UTF8Encoding($false, $true)
$utf8.GetString([System.IO.File]::ReadAllBytes('<path>'))
```

**Why the strict decoder matters:** `Get-Content` and `[Text.Encoding]::UTF8.GetString` both *replace* invalid bytes with U+FFFD rather than failing, so a corrupted file can masquerade as clean. Only the two-argument `UTF8Encoding($false, $true)` constructor throws, which is what makes this a real test.

---

*Generated as a research note. No source files were modified.*
