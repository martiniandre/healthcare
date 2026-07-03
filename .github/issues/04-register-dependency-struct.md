---
title: "Unificar Register.go com Dependency struct em todos os módulos"
labels: ["refactor", "backend"]
---

## What to build

Padronizar o wiring de todos os módulos para seguir o padrão documentado: `func Register(grpcServer *grpc.Server, dep Dependency)`.

### Problemas

Atualmente existem **3 variantes** do `Register()`:

| Módulo | Assinatura Atual | Problema |
|--------|-----------------|----------|
| auth, staff, patients, clinical, imaging, telemetry, audit_logs | `Register(grpcServer, rawDep1, rawDep2...)` | Parâmetros crus em vez de struct |
| exam_analyzer | `Register(databasePool, projectID, locationID, vertexModel) (Repository, Service, *Worker)` | Sem gRPC server, retorna tripleto |
| stats | `Register(databasePool, fhirClient) (Service, *HTTPHandler)` | Sem gRPC server, retorna handler direto |
| health | `Register(grpcServer, pool, redis)` | Sem service, grpc_handler injeta conexões direto |

### Escopo

1. Criar `Dependency` struct em cada módulo (ou uma shared) com os campos necessários
2. Unificar todas as assinaturas para `Register(grpcServer *grpc.Server, dep Dependency)`
3. `exam_analyzer`: Refatorar service para depender de Repository (não de strings)
4. `stats`: Mover a criação do HTTPHandler para fora do Register
5. `health`: Decidir se vira módulo hexagonal ou sai de `modules/`
6. Limpar `main.go`: remover criação duplicada de repositório imaging (linha 93)

### Acceptance criteria

- [ ] Todos os módulos com gRPC usam `Register(grpcServer, dep)` com Dependency struct
- [ ] `exam_analyzer` service aceita Repository interface
- [ ] `health` module segue padrão hexagonal ou é movido para `internal/app/health/`
- [ ] `main.go` não duplica criação de repositórios
- [ ] Build e testes passando

### Blocked by

None — can start immediately
