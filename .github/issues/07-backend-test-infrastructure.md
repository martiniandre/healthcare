---
title: "Padronizar infraestrutura de testes do backend"
labels: ["refactor", "backend", "testing"]
---

## What to build

Unificar os padrões de teste em todos os módulos do backend.

### Problemas

3 padrões diferentes de teste:

| Módulo | Localização | Mocks |
|--------|------------|-------|
| auth, staff, patients | `tests/` subdir | `mocks/` dir |
| clinical, exam_analyzer, stats | Module root | Inline function mock |
| audit_logs | `tests/` subdir | Sem mocks dir (inline) |
| health | `tests/` subdir | N/A |

### Acceptance criteria

- [ ] Todos os módulos com testes usam `tests/` subdirectory
- [ ] Todos os módulos com repositório têm `mocks/` com mock manual da interface
- [ ] `clinical`, `exam_analyzer`, `audit_logs` migrados para `mocks/` + `tests/`
- [ ] `stats` tem testes (atualmente sem mock/repository testável)
- [ ] CI executa `go test -v ./internal/modules/...`
- [ ] Cobertura mínima: service + handler para cada módulo

### Blocked by

None — can start immediately
