---
name: healthcare-commit
description: Publishes finished work following the Healthcare conventions — validates local integrity (backend build/tests/vet, frontend lint/tests/build), stages, commits with Conventional Commits v1.0.0, and pushes immediately to the active branch. Use when the user says "faça o commit", "commit e push", "publique", "finaliza o marco", or when a relevant feature/milestone is completed (AGENTS.MD automatic commit rule).
---

# Healthcare Commit & Push

Per AGENTS.MD, whenever a relevant feature is completed or a milestone is reached, the work must be committed and pushed automatically. Never wait for the user to ask explicitly.

## Step 1 — Inspect changes

```bash
git status
git diff --stat
git log --oneline -5
```

Review what is staged/unstaged and confirm there are no secrets (`.env`, credentials) being added.

## Step 2 — Validate local integrity

Run the checks that apply to the changed areas:

- Backend changed:
  ```bash
  cd backend
  go build ./...
  go test -v ./internal/...
  go vet ./...
  ```
- Frontend changed:
  ```bash
  cd frontend
  npm run lint
  npm run test
  npm run build
  ```
- E2E flows changed: `cd frontend && npm run test:e2e`

If anything fails, fix it before committing. Do not commit broken code.

## Step 3 — Stage

```bash
git add .
```

## Step 4 — Conventional Commit

Format: `<tipo>(<escopo>): <descrição curta e objetiva>`

Allowed types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`.

Scope = the module/area changed. Examples:

- `feat(patients): criar endpoint gRPC de histórico de pacientes`
- `fix(observation): corrigir parse de datas FHIR`
- `test(condition): adicionar cobertura de validação de ICD-10`
- `refactor(imaging): otimizar buffer de transmissão DICOM`

When the change spans frontend and backend, pick the dominant scope or omit it.

## Step 5 — Push

```bash
git push
```

Push to the active branch. If the push is rejected (remote moved), pull with rebase (`git pull --rebase`) and retry.

## Definition of Done

- [ ] Build/tests/vet clean on changed areas.
- [ ] No secrets staged.
- [ ] Commit message follows Conventional Commits.
- [ ] Changes pushed to the active branch.
