# Healthcare — Fluxo de Skills Automatizados

Orquestração de automações do ecossistema Healthcare (stack Go + React/Vite). Este diretório contém as **skills de projeto** versionadas e o ponto central de convenções; a skill orquestradora `healthcare-autopilot` vive globalmente em `~/.agents/skills/`.

## Skills de projeto (`.opencode/skills/`)

| Skill | Gatilho | O que faz |
| --- | --- | --- |
| `go-module-gen` | "crie o módulo backend [Nome] com os campos [Campos]" | Gera módulo Go gRPC completo (proto, pb stubs, model, repository, service, grpc_handler, register, RBAC, migration opcional, testes unitários) e pluga no `main.go`. |
| `frontend-module-gen` | "crie o módulo frontend [Nome]" | Gera módulo React/Vite (types, api, queries, schemas Zod, página, componentes, rota, i18n, testes Vitest + E2E Playwright offline). |
| `healthcare-commit` | "commit e push", "publique", "finaliza" | Valida integridade local (build/testes/vet/lint), staged com `git add .`, commit Conventional Commits e push automático. |

## Skill orquestradora (global)

`healthcare-autopilot` (`~/.agents/skills/healthcare-autopilot/`) roteia os pedidos para o fluxo certo e encadeia pipelines de ponta a ponta:

- **Backend module:** `go-module-gen` → validação → `healthcare-commit`
- **Frontend module:** `frontend-module-gen` → validação → `healthcare-commit`
- **Full feature:** módulo + commit + push + abertura de PR (`gh pr create`)
- **Finalize:** `healthcare-commit`
- **Review:** validações + PR

## Automação Git (`.githooks/` e CI)

- `.githooks/commit-msg` — valida Conventional Commits v1.0.0 a cada commit.
- `.githooks/pre-push` — já existente em `.git/hooks`; agora versionado como referência (go vet + npm lint).
- `.github/workflows/pr-conventions.yml` — valida título do PR e todos os commits na CI.

### Instalar hooks (uma vez por clone)

```powershell
powershell -ExecutionPolicy Bypass -File .githooks/setup.ps1
```

```bash
chmod +x .githooks/* && ./.githooks/setup.sh
```

## Convenções obrigatórias

- **Zero comentários** no código gerado; variáveis sempre descritivas (nunca `x`, `p`, `e`).
- **Persistência:** clínico → GCP Healthcare API (FHIR); operacional → PostgreSQL local com migration.
- **RBAC:** todo endpoint novo deve ser registrado em `backend/internal/app/policy/policy.go` — sem registro, é bloqueado.
- **Testes:** service/handler no Go com MockRepository; schemas/hooks no React com Vitest; fluxos novos com Playwright E2E offline (`page.route()` + `helpers.ts`).
- **Commits:** Conventional Commits v1.0.0, com push automático após marco concluído.
