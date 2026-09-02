---
name: issue-to-feature
description: Turns a detailed GitHub issue into a complete feature with a single command — reads the issue via gh (title, body, labels, comments, acceptance criteria), extracts the domain and fields, creates a feature branch, scaffolds backend (go-module-gen) and frontend (frontend-module-gen), registers RBAC, validates, commits, pushes and opens a PR that closes the issue. Asks for clarification only when the issue is genuinely ambiguous. Use when the user says "implemente a issue", "feature from issue", "feature da issue", "issue-to-feature", "issue #N".
---

# Issue → Feature (single command)

Deliver a complete feature from a detailed GitHub issue in one command, following the AGENTS.MD pipeline.

## Step 0 — Read the issue in detail (MANDATORY)

```bash
gh issue view <N> --json number,title,body,labels,assignees,milestone,state,url --jq '{number,title,body,labels:[.labels[].name],assignees:[.assignees[].login],milestone,state,url}'
```

Comments usually carry extra requirements — read them too:

```bash
gh issue view <N> --json comments --jq '.comments[].body'
```

## Step 1 — Extract the feature blueprint

From title/body/comments extract:

| Item | Source |
| --- | --- |
| Domain/module name (kebab-case, ex: `appointments`) | title, `[Domain]` prefix, or section headers in the body |
| Fields + types | "campos", "fields", lists or tables in the body |
| Business rules | bullets starting with "deve", "must", "regra", "não pode" |
| Scope (backend, frontend, or both) | labels `backend-only`/`frontend-only`; otherwise default **backend + frontend** |
| Acceptance criteria | "acceptance criteria", "definition of done", checklist items |

## Step 2 — Ambiguity gate (ask, do not guess)

Proceed automatically, but STOP and ask the user (via the `question` tool) only when genuinely blocking:

- No domain name can be derived from the issue.
- Scope is unclear AND the labels do not say `backend-only`/`frontend-only`.
- Required fields have no types and no example payload.
- The issue references an external integration that does not exist in the repo.

Otherwise use sensible defaults:
- Scope default = backend + frontend.
- Field defaults: `id` (UUID PK), `status` (enum), `createdAt`/`updatedAt` (timestamps).

## Step 3 — Activate an isolated worktree (MANDATORY)

Load the `parallel-worktrees` skill and follow **Step 0 — Activate** before writing any code:

1. `git fetch origin main`
2. `node scripts/worktrees.mjs list` — verify no other feature claimed this slug or the required identifiers.
3. Create the worktree from `origin/main`, never in the main repo directory:
   ```bash
   git worktree add ../healthcare-worktrees/<kebab-slug-from-issue-title> -b feat/<kebab-slug-from-issue-title> origin/main
   ```
4. Claim it in the registry:
   ```bash
   node scripts/worktrees.mjs claim --slug <kebab-slug-from-issue-title> --type feat --note "<issue #N: <title>>"
   ```
5. All remaining steps (backend, frontend, validation, commit, push, PR) run **inside the worktree** (`cd ../healthcare-worktrees/<slug>`), never in the main repo directory. If another feature already claimed a conflicting slug or identifier, STOP and choose another slug before proceeding.

When a migration number is needed, reserve it via `node scripts/worktrees.mjs reserve-migration --slug <slug> --name <snake_case>` instead of computing "local max + 1".

## Step 4 — Backend (when in scope)

Load and execute the `go-module-gen` skill:

- proto + compiled pb stubs, `model.go`, `repository.go`, `service.go`, `grpc_handler.go`, `register.go`.
- Register every new RPC in `backend/internal/app/policy/policy.go` (unregistered endpoints are blocked).
- Numbered migration in `backend/migrations/` only for operational modules (FHIR modules need none).
- Unit tests in `tests/service_test.go` with a `MockRepository` in the same package.

## Step 5 — Frontend (when in scope)

Load and execute the `frontend-module-gen` skill:

- types, api, queries, Zod schemas, page, components, route, i18n.
- Vitest unit tests for reusable logic/hooks/helpers.
- Playwright E2E offline with `page.route()` + `helpers.ts` (test titles strictly in English).

## Step 6 — Full validation

```bash
cd backend; go build ./...; go vet ./...; go test -v ./internal/...
cd frontend; npm run lint; npm run test; npm run build; npm run test:e2e
```

RBAC gate must be green: `go test -run TestEveryCompiledGRPCRPCIsRegisteredInPolicy ./internal/app/policy/...`

## Step 7 — Commit, push, open PR

1. Commit with Conventional Commits: `feat(<domain>): implement #<N>` (separate commits per layer when sensible).
2. `git push -u origin <branch>`.
3. Open the PR closing the issue:

```bash
gh pr create --base main --title "<subject>" --body "$(printf 'Closes #%s\n\n## What\n\n## Validation' '<N>')"
```

Use the official template from `.github/pull_request_template.md` if present.

## Definition of Done

- [ ] Issue read in detail (title, body, labels, comments) and blueprint extracted.
- [ ] Ambiguities asked; non-blocking defaults applied.
- [ ] Worktree activated at `../healthcare-worktrees/<slug>` and claimed in the registry; nothing committed on `main`.
- [ ] Backend/frontend scaffolded with mandatory tests.
- [ ] All endpoints registered in `policy.go`; RBAC gate green.
- [ ] Full validation matrix green (BE build/vet/test, FE lint/test/build, E2E).
- [ ] PR opened referencing and closing the issue (`Closes #N`).
