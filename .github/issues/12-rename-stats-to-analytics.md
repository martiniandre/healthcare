---
title: "Renomear módulo Stats para Analytics"
labels: ["refactor", "backend", "frontend"]
---

## What to build

Renomear o módulo `stats` para `analytics` no backend e frontend, conforme definido no CONTEXT.md.

### Escopo

**Backend:**
- Renomear `internal/modules/stats/` → `internal/modules/analytics/`
- Atualizar package name de `stats` para `analytics`
- Atualizar imports em `main.go`
- Atualizar rotas HTTP
- Atualizar protobuf se aplicável

**Frontend:**
- Renomear `frontend/src/modules/stats/` → `frontend/src/modules/analytics/`
- Atualizar imports em `routes.tsx`
- Atualizar sidebar
- Renomear arquivos de tradução `stats.json` → `analytics.json`
- Atualizar namespace i18n

**Banco de dados:**
- N/A (stats não tem tabela própria — usa FHIR + queries SQL em tabelas existentes)

### Acceptance criteria

- [ ] Backend package renomeado sem quebrar imports
- [ ] Frontend módulo renomeado com todas as referências atualizadas
- [ ] Rotas HTTP `/api/stats` → `/api/analytics` (ou alias de compatibilidade)
- [ ] i18n namespace atualizado
- [ ] Build, lint e testes passando

### Blocked by

Issue #5 (Desacoplamento staff → auth) — stats http_handler importa auth
