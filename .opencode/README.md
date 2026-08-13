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
- **Full feature:** backend + testes unitários (`go-module-gen`) → frontend + testes Vitest + E2E (`frontend-module-gen`) → validações completas → commit+push → PR
- **Finalize:** `healthcare-commit`
- **Review:** validações + PR

## Plugin Skills & Flows (global)

`~/.config/opencode/plugins/skills-flows.ts` — plugin global que gerencia skills e flows de qualquer projeto. Tools expostas ao agente:

| Tool | O que faz |
| --- | --- |
| `list_skills` | Lista skills de projeto + globais com descrição e escopo. |
| `create_skill` | Cria `SKILL.md` com frontmatter válido (projeto ou global), sem sobrescrever. |
| `validate_skills` | Varre todas as skills e reporta frontmatter ausente ou nome fora do padrão. |
| `list_flows` | Lista os flows do registry do projeto atual. |
| `run_flow` | Executa as validações de um flow via shell e orienta os próximos passos (carregar skill, commitar). |

Os flows ficam no registry do projeto em `.opencode/flows/registry.json` (neste repo: `backend-module`, `frontend-module`, `full-feature`, `finalize`, `review`). O plugin resolve o registry do repo atual automaticamente (`<worktree>/.opencode/flows/registry.json`).

## Automação Git (`.githooks/` e CI)

- `.githooks/commit-msg` — valida Conventional Commits v1.0.0 a cada commit.
- `.githooks/pre-push` — check rápido (go vet + npm lint) antes do push.
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
