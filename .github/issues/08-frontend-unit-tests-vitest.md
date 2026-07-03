---
title: "Adicionar testes unitários frontend com Vitest"
labels: ["refactor", "frontend", "testing"]
---

## What to build

Implementar testes unitários no frontend com Vitest para hooks customizados, componentes reutilizáveis, e helpers de integração.

### Escopo

1. Configurar Vitest no projeto (vite.config.ts já existe)
2. Testar hooks: `useAuthInit`, `useDebounce`
3. Testar componentes compartilhados: `ui/` components (Button, Dialog, Select, etc.)
4. Testar helpers: `validators.ts`, `cn.ts`
5. Testar schemas: `patient_schemas.ts` (Zod validation)
6. Adicionar ao CI: `npm run test` no workflow ci-frontend.yml

### Acceptance criteria

- [ ] Vitest configurado e rodando com `npm run test`
- [ ] Hooks `useDebounce` e `useAuthInit` testados
- [ ] Componentes UI têm teste de renderização
- [ ] Schemas de validação testados (casos válidos e inválidos)
- [ ] CI executa `npm run test` no push/PR
- [ ] Mínimo 10 testes passando

### Blocked by

None — can start immediately
