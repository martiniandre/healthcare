---
name: frontend-module-gen
description: Creates a new frontend React (Vite + TypeScript) domain module in the Healthcare monorepo following FRONTEND_GUIDELINES.md (types.ts, api.ts, queries.ts, schemas, page master, components, route registration, Vitest unit tests and offline Playwright E2E). Use when the user says "crie o módulo frontend [Nome]", "novo módulo React", "frontend module", "gerar módulo frontend", or triggers the Antigravity frontend flow.
---

# Frontend Module Generator

Create a new domain module under `frontend/src/modules/{domain}/`, following the Domain-Driven scaffolding described in `FRONTEND_GUIDELINES.md`.

## Trigger

User request: `"Antigravity, crie o módulo frontend [Nome]"`.

## Step 0 — Clarification (MANDATORY)

Identify ambiguities before writing code (fields, screens, route path, roles, API endpoints). Ask the user and wait for the answer. Never assume silently.

## Step 0.5 — Worktree isolation (MANDATORY when other features may be in flight)

If the module is created inside a dedicated worktree, it is already claimed by the orchestrator. If the current working directory is the main repo, load `parallel-worktrees` Step 0 and activate a worktree first. All `routes.tsx`, `AppSidebar` and i18n edits must follow the shared-file additive-anchored protocol: add-only, sorted/anchored inserts, never reorder or rename existing entries.

## Step 1 — Module scaffolding

Create `frontend/src/modules/{domain}/` with this autonomous structure (mirror `frontend/src/modules/patients/`):

| File | Responsibility |
| --- | --- |
| `types.ts` | Strict DTO and domain types |
| `api.ts` | Isolated HTTP client (Axios) with typed input/output |
| `queries.ts` | TanStack Query hooks (queries + mutations + cache keys) |
| `{domain}_schemas.ts` | Zod schemas for forms |
| `{domain}_schemas.test.ts` | Vitest unit tests for the schemas |
| `{Domain}.tsx` | Master/orchestrating page component |
| `components/` | Domain-exclusive subcomponents |
| `components/modals/` | Reactive modals and dialogs |

## Step 2 — Conventions (MANDATORY)

- **TanStack Query:** all backend sync goes through TanStack Query hooks. Queries supporting filters/search must include active params in the `queryKey`. Mutations must call `queryClient.invalidateQueries` in `onSuccess` with the corresponding query keys.
- **Zod + React Hook Form:** every form is validated with a declarative Zod schema resolved via `zodResolver`. Form types inferred with `z.infer<typeof schema>`.
- **Visual identity:** premium, accessible UI — semantic HTML (`<main>`, `<section>`, `<aside>`, `role="img"`, `aria-label`), HSL tonal palettes (Slate/Emerald/Rose/Cyan), `duration-300` transitions, `animate-pulse` loading states, contextual cursors where the tool requires them.
- **Zero comments:** no explanatory comments in any file.
- **Descriptive variables:** no single-letter variables (`patients.map((patientItem) => ...)`), even in event handlers (`(event) => ...`).

## Step 3 — Route and navigation registration

- Register the page in `frontend/src/app/routes.tsx` using `lazy()` + `Suspense` (existing pattern) and a guarded `<Route path=... element={...} />`. Insert the `lazy()` import and the `<Route>` at their **path-ordered positions** — additive inserts only, so parallel worktrees merge cleanly.
- Add the menu entry in the `AppSidebar` (respecting role-based visibility when the module is role-specific), inserted after the alphabetical/role-ordered neighbour. Never reorder existing items.

## Step 4 — i18n

Add pt-BR / en-US / es-ES keys **only** (new `{domain}` key groups) when the module renders user-facing copy. Never rename or move existing keys that other modules or worktrees may rely on. Follow the existing i18n structure in `frontend/src/shared/i18n/`.

## Step 5 — Unit tests (MANDATORY)

Add Vitest unit tests for:
- Zod schemas (`*.test.ts` next to the schema file, covering valid + invalid inputs).
- Any reusable hook or integration helper in the module.

## Step 6 — E2E tests (MANDATORY for new/changed flows)

Create `frontend/e2e/{domain}.spec.ts`:

- **Titles strictly in English.**
- **100% offline/resilient:** intercept every external call and every auth call with native `page.route()`; reuse `mockAuthAPI` from `frontend/e2e/helpers.ts` for authentication. The suite must pass with zero backend running.
- Cover the happy path and at least one validation/error flow.

## Validation commands

```bash
cd frontend
npm run lint
npm run test
npm run build
npm run test:e2e
```

## Definition of Done

- [ ] `types.ts`, `api.ts`, `queries.ts`, schemas + tests, `{Domain}.tsx`, `components/` created.
- [ ] Route registered in `routes.tsx` and sidebar updated.
- [ ] i18n keys added for supported locales.
- [ ] Vitest unit tests pass (`npm run test`).
- [ ] Playwright E2E passes offline with mocked routes (`npm run test:e2e`).
- [ ] `npm run lint` and `npm run build` clean.
- [ ] Zero comments; descriptive variable names only.
