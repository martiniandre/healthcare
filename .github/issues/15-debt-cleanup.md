---
title: "Cleanup final: dívida técnica remanescente"
labels: ["refactor", "backend", "frontend"]
---

## What to build

Resolver os itens de dívida técnica menores que sobraram das outras issues.

### Itens

1. **CI/CD pre-push hook**: O hook atual é pesado (go vet + go test + eslint + build). Considerar mover validações para CI e deixar hook só com `go vet` + linter rápido.
2. **Pacote `imaging/` com estrutura extra**: `hooks/`, `utils/`, `constants.ts` fogem da convenção. Avaliar se devem ser padronizados ou documentados como exceção permitida.
3. **`patients/grpc_handler.go` ignora params**: ListPatients recebe `(search, sortField, sortDirection, page, limit)` mas passa hardcoded `("", "", "", 1, 100)`. Implementar os parâmetros reais.
4. **Clinical service usa `ErrEncounterNotFound` para validação**: `CreateEncounter` retorna `ErrEncounterNotFound` quando `patientID` está vazio — deveria ser erro de validação, não "not found".
5. **FHIR error messages vazam detalhes**: `client.go` retorna `healthcare api error %d: %s` com body cru da API. Sanitizar mensagens de erro.
6. **Patients HTTP handler em português**: Mensagens de erro em português vs gRPC em inglês. Padronizar idioma.
7. **Stats defaults hardcoded**: Quando FHIR retorna vazio, stats usa defaults hardcoded. Documentar ou remover.
8. **`.env.production`**: Verificar se deve ser adicionado ao `.gitignore`.
9. **PageViewLogger acoplamento**: `routes.tsx` importa `auditLogsApi` diretamente. Extrair para shared hook.

### Acceptance criteria

- [ ] Cada item resolvido ou documentado como decisão consciente
- [ ] `.env.production` no `.gitignore`
- [ ] `ListPatients` usa parâmetros recebidos
- [ ] Mensagens de erro padronizadas (inglês)
- [ ] Build, lint e testes passando

### Blocked by

Todas as issues anteriores (é o cleanup final)
