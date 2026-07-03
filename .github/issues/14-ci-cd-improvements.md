---
title: "Melhorar pipeline CI/CD: Go tests + Vitest"
labels: ["refactor", "devops"]
---

## What to build

Adicionar execução de testes no pipeline CI/CD do GitHub Actions.

### Problemas

- CI frontend atual (`ci-frontend.yml`) só roda `npm run lint` + `npm run build` — **não roda testes**
- Backend não tem workflow CI próprio — `go vet` e `go test` só rodam no pre-push hook local
- Vitest não está configurado no projeto (Issue #08 cria os testes)

### Escopo

1. Adicionar `npm run test` (Vitest) ao `ci-frontend.yml`
2. Criar `ci-backend.yml` com `go vet` + `go test ./...`
3. Ou unificar em um workflow `ci.yml` que roda ambos

### Acceptance criteria

- [ ] Frontend CI roda `npm run test`
- [ ] Backend CI roda `go vet ./...` + `go test -v ./internal/...`
- [ ] Workflows falham se testes quebram
- [ ] Cache de módulos Go e node_modules para performance

### Blocked by

Issue #08 (Testes frontend Vitest) — precisa existirem testes para rodar
