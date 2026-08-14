---
name: issue-to-feature
description: Turns a GitHub issue into a complete feature — reads the issue with gh, creates a branch, scaffolds the backend (go-module-gen) and/or frontend (frontend-module-gen) modules, registers the RBAC endpoints, and opens a PR referencing the issue. Use when the user says "implemente a issue", "issue-to-feature", "issue #N", "crie a branch e implemente a issue".
---

# Issue → Feature

Deliver a GitHub issue as a shippable feature following the AGENTS.MD pipeline.

## Step 0 — Read the issue (MANDATORY)

```bash
gh issue view <N> --json number,title,body,labels,assignees --jq '{number,title,body,labels:[.labels[].name]}'
```

Extract from the title/body:

- The domain name (module).
- The fields and business rules.
- Whether it needs backend, frontend, or both.
- Any acceptance criteria listed.

If the issue is ambiguous, ask the user before writing code.

## Step 1 — Create the branch

```bash
git checkout -b feature/<kebab-case-slug-from-issue-title>
```

Never implement directly on `main`.

## Step 2 — Scope the work

- **Backend needed** → load and execute the `go-module-gen` skill (proto, pb stubs, model, repository, service, handler, register, RBAC, migration when operational, unit tests).
- **Frontend needed** → load and execute the `frontend-module-gen` skill (types, api, queries, schemas, page, components, route, i18n, Vitest tests, E2E tests offline).
- **Both** → run the `full-feature` flow.

## Step 3 — Register security

Open `backend/internal/app/policy/policy.go` and add every new endpoint to `grpcMethodRoles`/`httpRouteRoles`. The RBAC gate test (`go test -run TestEveryCompiledGRPCRPCIsRegisteredInPolicy ./internal/app/policy/...`) fails on unregistered proto RPCs.

## Step 4 — Validate

```bash
cd backend; go build ./...; go vet ./...; go test -v ./internal/...
cd frontend; npm run lint; npm run test; npm run build; npm run test:e2e
```

## Step 5 — Commit, push, open PR

1. Commit with Conventional Commits (e.g. `feat(<domain>): implement #<N>`).
2. `git push -u origin <branch>`.
3. Open the PR referencing the issue:

```bash
gh pr create --base main --title "<commit subject>" --body "Closes #<N> ..."
```

Use the official PR template from `.github/pull_request_template.md` if present.

## Definition of Done

- [ ] Issue read and requirements extracted.
- [ ] Work done on a feature branch.
- [ ] Backend/frontend scaffolds complete with tests.
- [ ] All endpoints registered in policy.go and RBAC gate green.
- [ ] Build/vet/lint/tests passing.
- [ ] PR opened and referencing the issue (`Closes #N`).
