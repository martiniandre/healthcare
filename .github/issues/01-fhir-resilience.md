---
title: "Resiliência FHIR: paginação, retry, DeleteResource, zero-value bug"
labels: ["refactor", "backend"]
---

## Parent

CONTEXT.md criado como baseline do domínio.

## What to build

Corrigir a camada de integração com GCP Healthcare API para ser resiliente a falhas e retornar dados completos.

### Problemas

1. **Sem paginação**: `SearchResources` retorna apenas uma página. FHIR bundles podem ter links `next` que não são seguidos. `_count=100` nos callers (stats, patients) retorna dados parciais silenciosamente.

2. **Sem retry/backoff**: GCP Healthcare API retorna 429 (rate limit) e 5xx sob carga. Zero tratamento de retry.

3. **Sem `DeleteResource`**: FHIR suporta DELETE, mas o cliente não expõe essa operação.

4. **Imaging worker usa `*healthcare.Client` concreto**: `internal/modules/imaging/worker.go` usa tipo concreto em vez da interface `healthcare.FHIRClient`.

5. **NewObservationResource descarta valor 0**: `shared/fhir/observation_resource.go` — `if valueQuantity != 0` impede medições clinicamente válidas com valor zero.

6. **Sem timeout configurado**: Cliente HTTP subjacente não tem timeout. Context passado pelo caller é o único mecanismo.

### Acceptance criteria

- [ ] `SearchResources` segue links `next` em bundles FHIR ou usa `_getpagesoffset`
- [ ] `FHIRClient` interface inclui `DeleteResource`
- [ ] Cliente FHIR implementa retry com backoff exponencial para 429/5xx
- [ ] Timeout configurável no cliente HTTP subjacente
- [ ] `NewObservationResource` aceita valor 0 como medição válida
- [ ] Imaging worker usa `FHIRClient` interface
- [ ] Testes unitários cobrindo paginação, retry e zero-value

### Blocked by

None — can start immediately
