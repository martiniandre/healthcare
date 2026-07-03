---
title: "Padronizar módulo auth do frontend (types.ts, api.ts, queries.ts)"
labels: ["refactor", "frontend"]
---

## What to build

Fazer o módulo `auth` do frontend seguir o mesmo padrão dos demais módulos.

### Problema

`frontend/src/modules/auth/` é o único módulo que NÃO tem `types.ts`, `api.ts` e `queries.ts`. Em vez disso:
- Lógica de API está em `shared/services/auth_api.ts` (fora do módulo)
- Tipos estão inline em `auth_schemas.ts`
- `Login.tsx` chama `authApi` diretamente sem camada de TanStack Query

### Acceptance criteria

- [ ] `auth/types.ts` criado com tipos de domínio
- [ ] `auth/api.ts` criado com chamadas HTTP (login, register, logout)
- [ ] `auth/queries.ts` criado com hooks TanStack Query
- [ ] `Login.tsx` usa hooks de `queries.ts` em vez de chamar `auth_api` direto
- [ ] `shared/services/auth_api.ts` removido ou redirecionado para módulo
- [ ] Testes atualizados

### Blocked by

None — can start immediately
