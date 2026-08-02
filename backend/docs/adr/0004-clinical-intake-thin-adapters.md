# Clinical intake — thin transport adapters, domain rules behind the module interface

For the pilot modules (encounter, diagnostic_report, observation), the module's Service interface owns required-field validation, defaults, patch-merge on update, and typed errors via `apperrors`; the gRPC and HTTP adapters translate wire shapes only. The HTTP adapter renders errors from `apperrors.HTTPCode` (killing silent 500s), and the shared FHIR client distinguishes 404 so the module can return a real not-found.

Reason: both transports previously carried validation, defaults, and error mapping independently — producing real divergence (encounter status, hardcoded report code, HTTP always 500) and a validation surface three layers wide (handler, handler, service).

Deferred: RBAC policy (permissions.go map vs per-route RequireRoles) is only aligned in the pilot, not unified; full FHIR read-side parsing consolidation is a separate deepening. These are not rejected — sequenced after the pilot.

Consequences: `CreateEncounter`/`UpdateEncounter` take dedicated input types; the frontend's practitioner field becomes required and a report-type (LOINC) selector is added.
