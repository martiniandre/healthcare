---
title: "Extrair componentes do AuditLogs + adicionar i18n"
labels: ["refactor", "frontend"]
---

## What to build

Refatorar `AuditLogs.tsx` (296 linhas) para seguir o padrão de componentes dos demais módulos.

### Problema

`AuditLogs.tsx` renderiza tudo inline: filtros, tabela, skeleton loading, linhas expandidas. É o único módulo sem diretório `components/`. Além disso, não há arquivos de tradução `auditLogs.json` em nenhum dos 3 idiomas (pt-BR, en-US, es-ES).

### Acceptance criteria

- [ ] `audit_logs/components/` criado com: `AuditLogsFilters.tsx`, `AuditLogsTable.tsx`, `AuditLogsLoadingState.tsx`
- [ ] `AuditLogs.tsx` reduzido para < 80 linhas (apenas composição)
- [ ] `shared/i18n/locales/pt-BR/auditLogs.json` criado
- [ ] `shared/i18n/locales/en-US/auditLogs.json` criado
- [ ] `shared/i18n/locales/es-ES/auditLogs.json` criado
- [ ] Namespace `auditLogs` registrado no i18n config
- [ ] Componentes usam `useTranslation("auditLogs")`

### Blocked by

None — can start immediately
