---
title: "Quebrar módulo Clinical em subdomínios"
labels: ["refactor", "backend", "frontend"]
---

## What to build

Dividir o monolítico módulo `clinical` em módulos de domínio menores conforme decidido na sessão de domain-modeling.

### Motivação

Clinical atualmente agrupa 6 conceitos de domínio distintos (Encounter, Observation, Condition, AllergyIntolerance, MedicationRequest, DiagnosticReport) em um único módulo. Isso viola o princípio de single responsibility e dificulta navegação e manutenção.

### Proposta de subdomínios

1. **encounter** — CreateEncounter, GetEncounters
2. **observation** — CreateObservation, GetObservations
3. **condition** — CreateCondition, GetConditions
4. **allergy** — CreateAllergyIntolerance, GetAllergyIntolerances
5. **medication** — CreateMedicationRequest, GetMedicationRequests
6. **diagnostic_report** — CreateDiagnosticReport, GetDiagnosticReports

### Escopo

**Backend:**
- Criar 6 novos módulos em `internal/modules/`
- Dividir `proto/clinical.proto` em 6 proto files (ou manter um proto geral)
- Gerar stubs
- Cada módulo: model, repository (FHIR), service, grpc_handler, register
- Atualizar `permissions.go` com os novos endpoints
- Atualizar `main.go` para registrar todos

**Frontend:**
- Decidir se o frontend também separa ou mantém agrupado
- Atualizar imports e queries

### Acceptance criteria

- [ ] 6 módulos criados, cada um com model, repository, service, grpc_handler, register
- [ ] `clinical` original removido (ou mantido como compatibilidade)
- [ ] Todos os endpoints funcionando
- [ ] Permissions.go atualizado
- [ ] Testes para cada novo módulo
- [ ] Build, lint e testes passando

### Blocked by

Issue #6 (Padronizar builders FHIR) — todos os subdomínios usam builders FHIR
