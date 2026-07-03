---
title: "Correções de consistência no frontend"
labels: ["refactor", "frontend"]
---

## What to build

Corrigir inconsistências menores no frontend identificadas na auditoria.

### Itens

1. **`http` wrapper sem `.delete()`**: Adicionar método `delete()` em `shared/utils/http.ts`. Corrigir `exam_analyzer/api.ts` que importa `api` diretamente.
2. **Modal props inconsistentes**: Unificar para `onOpenChange: (open: boolean) => void` em todos os modais.
3. **`window.alert()` no Imaging**: Substituir por `toast.success()` em `ImagingWorkspace.tsx`.
4. **Variável `i` em loop**: `PatientsTable.tsx` usa `i` como index — renomear para `index`.
5. **Comentário em `useAuthInit.ts`**: Remover comentário em português.
6. **EmptyState compartilhado**: Migrar módulos para usar `shared/components/ui/EmptyState.tsx`.
7. **Maping snake_case do Staff**: Decidir se padroniza camelCase no backend ou mapeia no frontend. Unificar abordagem.
8. **Stats types.ts**: Criar `stats/types.ts` e mover types de `stats/api.ts`.
9. **Settings route**: Criar rota `/settings` em `routes.tsx` ou remover do sidebar.
10. **Telemetry `useUnlockRoomMutation`**: Adicionar `onSuccess` invalidation de queries.
11. **Stats query keys**: Adicionar `lists()` e `detail()` em `statsKeys`.

### Acceptance criteria

- [ ] `http.delete()` existe e exam_analyzer usa `http` em vez de `api` direto
- [ ] Todas as props de modal usam `onOpenChange`
- [ ] Imaging usa `toast.success()` em vez de `window.alert()`
- [ ] Nenhum `i` como variável de loop no código
- [ ] Nenhum comentário em português no código
- [ ] EmptyState compartilhado usado em pelo menos 2 módulos
- [ ] `stats/types.ts` criado
- [ ] Rota Settings resolvida (rota ou removida)
- [ ] `useUnlockRoomMutation` invalida queries no success
- [ ] Build e ESLint passando

### Blocked by

None — can start immediately
