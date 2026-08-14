# Changelog

Todas as mudanças relevantes por tipo, geradas automaticamente a partir dos commits Conventional Commits.

## Features

- **automation:** expand full-feature flow with backend, frontend and full test pipeline (`ea41578`)
- **automation:** add flows registry and document skills-flows plugin (`d55273b`)
- **fhir:** add offline FHIR base URL mode for local development and tests (`c751dd4`)
- **ui:** harden form validations, fix API payload/endpoint bugs, and improve loading/error states (`b26c30b`)
- **api:** surface backend error payload in axios interceptor (`9b3fa1f`)
- **notifications:** support report_ready notification type (`e864ae3`)
- **portal:** patient portal with appointments and report version badge (`3484211`)
- **audit_logs:** render resource and payload diff in audit table (`1e1bad7`)
- **schedule:** add frontend scheduling module with routes, navigation and i18n (`3acaaa0`)
- **patients:** encounter detail workspace with clinical modals and report version history (`324ec05`)
- **auth:** support appointment feature and extract useAbility hook (`8a7a031`)
- **schedule:** add authenticated patient my-appointments endpoint (`9ff6a24`)
- **diagnostic_report:** surface report version in http responses (`d599519`)
- **portal:** expose FHIR meta.versionId as report version (`d97beab`)
- **audit_logs:** add rich resource audit with deterministic payload diff (`11f0ff5`)
- **schedule:** add appointment booking backend with idempotency and conflict prevention (`f160fbf`)
- **diagnostic_report:** add version history and report update endpoint (`9d24d05`)
- **diagnostic_report:** publish report.ready event on report creation (`560408b`)
- **notifications:** add report_ready type routed to patients (`02aadcb`)
- standardize data structures and add validation across all modules (`977e8f6`)
- publish domain events from exam_analyzer (exam.complete) and auth (login) (`26c98d9`)
- add in-app notification system with event bus, SSE streaming, and role-based routing (`5283b77`)
- add Patient Portal and Clinical Dashboard (`75e44a9`)
- docker-compose, API versioning, health endpoint, request ID propagation (`bdc3acb`)
- swagger docs, exam_analyzer cleanup, production infra, OTEL (`fe7b16b`)
- complete remaining 3 issues - rename stats to analytics, split clinical into 6 subdomains, final cleanup (`f039f73`)
- implement issues #1, #4-#11, #14 in parallel multi-agent groups (`8c4e2c5`)
- **telemetry:** add background vital signs simulator and fix Bed JSON serialization (`f1624c3`)
- implement database-backed audit and action log system (`d1fe899`)
- **frontend:** integrate stats module with GET /api/stats query and dynamic charts (`99a7fe4`)
- **backend/stats:** implement statistics module using hexagonal architecture (`8c074fa`)
- **validation:** implement comprehensive frontend and backend input validations and fix E2E tests (`b005ba7`)
- configure docker-in-docker with nested frontend and backend services (`ec2363e`)
- add docker-in-docker templates and documentation (`e8563e6`)
- **frontend:** implement medications module and exam analyzer e2e tests (`ee2a2ba`)
- **frontend:** replace alert calls with a sleek custom toast notification system (`2b68e0c`)
- **frontend:** integrate AllergyIntolerance clinical module and add FRONTEND_GUIDELINES (`df371c9`)
- **imaging:** extract informative canvas overlay labels and fallbacks to constants (`6890efb`)
- **clinical:** expose allergy and medication request endpoints in rest api (`75643b7`)
- **frontend:** integrate all remaining mock clinical features and extract response DTOs (`a3a4818`)
- **frontend:** set /api as default baseURL in Axios and clean up routes (`ee7e483`)
- **frontend:** optimize pacs viewer performance and enhance accessibility (`a191b01`)
- **clinical:** integrate clinical data models and REST endpoints with front-end Axios client and resilient Playwright E2E mocks (`f7f7e31`)
- **patients:** integrate frontend client with secure backend REST endpoints and E2E mocks (`00a6700`)
- **e2e:** implement robust full-coverage Playwright test suite for all modules (`24d9c9f`)
- **auth:** integrate Axios request factory, real HTTP auth gateway, and Playwright tests (`e998229`)
- **frontend:** integrate PWA manifest, service worker caching, and connectivity UI warnings (`a9cda93`)
- **frontend:** add github pages deployment pipeline (`d6f4779`)

## Bug fixes

- **i18n:** resolve module namespaces with ':' separator in createModuleTranslator (`380e167`)
- **medication:** validate requests via apperrors and make requester optional (`acb088b`)
- parse all FHIR fields in observation, diagnostic report and medication bundle listing (`5c66fb8`)
- align FE ability permissions with BE endpoints; hide restricted tabs from PatientDetails (`bb91a1e`)
- parse all fields in encounter, condition and allergy FHIR bundle parsers; add encounter_id to condition creation (`ef9f7ee`)
- **exam_analyzer:** route gRPC handler through service layer instead of direct repo access (`f6e9978`)
- **patients:** remove duplicate UUID, use deterministic ID from FHIR resource, extract meta.lastUpdated, fix identifier search (`56dbaa3`)
- **tests:** update mock signatures for ListEmployees and ListPatients to match updated interfaces (`20900fe`)
- **backend:** standardize apperrors across clinical, health, exam_analyzer, and stats modules (`1191053`)
- **staff:** fix ListEmployees test call to match updated signature (`bf76eb4`)
- **audit-logs:** sync model types and handle access_granted status on frontend (`746d910`)
- **frontend:** handle paginated response and add error state in audit logs page (`0c08176`)
- **frontend:** resolve lint error and unused declarations in PatientDetails (`2c0a72f`)
- **frontend:** export base schemas to fix unused variable lint error (`f1b2d3d`)
- **frontend:** remove default bed selection and fix E2E test assertions (`17b4ec8`)
- **lint:** resolve E2E helpers any type and staff schema regex escape ESLint errors (`6ee9c89`)
- **csrf:** resolve cross-origin token issue with vite dev proxy (`54de34f`)
- **frontend:** manually inject CSRF token from cookie to avoid cross-origin Axios constraints (`1814fe0`)
- **frontend:** remove unsafe any cast in Staff.tsx to satisfy ESLint constraints (`51b9182`)
- **frontend:** handle undefined staff fields safely to prevent crashes (`4a6bcd1`)
- **frontend:** configure axios to read and send csrf token (`7c182df`)
- **security/staff:** enforce patch csrf, protect logout, align rbac, fix null crm scan (`cc418b3`)
- **frontend:** fix any to unknown in e2e tests (`c4857de`)
- **backend:** resolve Go compilation errors and interface type mismatches (`4f79bcd`)
- **frontend:** resolve react-hooks exhaustive-deps and ref-render linting errors (`fd3d675`)

## Documentation

- **plan:** mark i18n accessor phase done (`da8e70e`)
- **adr:** record clinical intake deepening and Encounter lifecycle decisions (`bdfe0a9`)
- update CONTEXT.md with infra section and OTEL term (`cb05a21`)
- update CONTEXT.md with final architecture state after all 15 refactoring issues (`8816272`)
- **plan:** add 15 refactoring issues as markdown files (`0a56b8b`)
- **context:** add ubiquitous language glossary (CONTEXT.md) (`2d3e634`)
- **config:** add detailed step-by-step setup guides and refactor frontend dynamic api client (`1c96c90`)
- **config:** add project configuration checklist and guide (`8327715`)
- add hexagonal architecture and domain boundaries rules to agents and gemini docs (`b8f3393`)
- **readme:** add root README with comprehensive monorepo documentation (`648c82e`)

## Refactors

- **i18n:** single module-bound translation accessor (`beacb57`)
- **security:** unify RBAC policy across gRPC and HTTP (`a9d6f57`)
- **imaging:** move validation behind service interface (ADR-0004 rollout) (`04dff9d`)
- **telemetry:** move validation behind service interface (ADR-0004 rollout) (`3302a1b`)
- **staff:** move validation behind service interface (ADR-0004 rollout) (`80b1404`)
- **patients:** move validation behind service interface (ADR-0004 rollout) (`8b51771`)
- **medication:** move validation behind service interface (ADR-0004 rollout) (`4f777b0`)
- **condition:** move validation behind service interface (ADR-0004 rollout) (`eca70b5`)
- **allergy:** move validation behind service interface (ADR-0004 rollout) (`ce8c84d`)
- **fhir:** consolidate FHIR read-side parsing into shared deep module (`c8e47fe`)
- **observation:** default observed time in service input with thin transports (`b7dac34`)
- **diagnostic_report:** require report code via service input and thin transports (`f8bc482`)
- **encounter:** move domain rules behind service inputs with thin transports (`9c770d8`)
- **shared:** type FHIR 404s and render errors from apperrors HTTP codes (`de2154b`)
- transport per module - notifications HTTP-only, exam_analyzer gains gRPC (`2369ff2`)
- **frontend:** delete unused Navbar component (`84f4cf1`)
- **frontend:** consolidate redundant useState hooks in multiple components (`e6ef840`)
- **frontend:** replace free-text string types with const objects and union types (`fa3d5d8`)
- **frontend:** align pattern violations and enable route lazy loading (`71b3596`)
- fix architecture, security and resilience issues (`4eb4ef7`)
- fix architecture, security and resilience issues (`360f780`)
- **api:** modular router with RouteRegisterer, declarative auth middleware, render helpers, and structured observability (`aa8d0ef`)
- **api:** extract routes, decouple main.go, and implement auth/cors middlewares (`33f78d8`)

## Tests

- **integration:** add full-stack E2E HTTP suite for all modules (`a663abd`)
- **integration:** add end-to-end HTTP suite with in-memory FHIR client (`37d2a7b`)
- **e2e:** stabilize Playwright suite and add portal, report versioning and report.ready coverage (`8a5198c`)
- **clinical:** add backend service unit tests for GetObservationsByPatient (`2cdb928`)

## Style

- enhance ErrorBoundary reset button aesthetics (`19e2ad7`)
- **frontend:** refine PWA manifest and service worker to use relative paths (`3097395`)
- **frontend:** make staff and stats modules fully responsive (`3b52a28`)

## Chores

- **hooks:** version pre-push check and document hooks setup (`314a52c`)
- **repo:** version project skills, git hooks and PR conventions workflow (`8c3eec9`)
- **migrations:** renumber diagnostic_report_versions to 012 (`8a04d6b`)
- **docker:** optimize Dockerfile for production and remove .env copy (`ef9b847`)
- temporarily disable Redis connection (`c59cb1f`)
- **instructions:** enforce mandatory unit testing for all business logic and API endpoints in backend and frontend (`4ccf530`)
- **git:** add strict automated commit and push rule to agents.md and gemini.md (`ff28dde`)
- **git:** add CODEOWNERS for automated code review routing (`b6943e8`)
- **git:** add pull request template for automated PR checklists (`fec4a69`)
- **git:** add CI pipeline and comprehensive git workflow specification (`87d9b7c`)
- initial commit with monorepo structure and .gitignore (`b735636`)

## CI

- **frontend:** configure CI and CD workflows to trigger exclusively on frontend changes (`0d6e07f`)

## Build

- **deps:** resolve and download all required Go module dependencies (`4c01dd6`)

## Outros

- merge branch 'feat/clinical-integration' into main (`4d12d6f`)

Últimos 127 commits incluídos.
