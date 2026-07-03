---
title: "Desacoplar módulo staff do módulo auth"
labels: ["refactor", "backend"]
---

## What to build

Remover a dependência direta do módulo `staff` no módulo `auth`.

### Problema

`internal/modules/staff/model.go` importa `auth.Role` como tipo do campo `Employee.Role`. `service.go` também importa `auth.ParseRole()` e `auth.ErrInvalidRole`. Isso cria acoplamento direto entre módulos — staff não pode existir sem auth.

### Solução

Extrair `Role` e funções auxiliares para um pacote compartilhado (ex: `internal/shared/role`):

1. Mover `Role` type e constantes para `internal/shared/role/role.go`
2. Mover `ParseRole()` para `internal/shared/role/`
3. Atualizar `auth` para reexportar ou consumir de `shared/role`
4. Atualizar `staff` para consumir de `shared/role`
5. Atualizar todos os outros módulos que importam `auth.Role` (http_handlers)

### Acceptance criteria

- [ ] `staff` não importa mais `auth`
- [ ] `auth` continua exportando `Role` (pode ser reexport ou shared)
- [ ] Todos os HTTP handlers que usam `auth.Role` para permissões atualizados
- [ ] Build e testes passando

### Blocked by

None — can start immediately
