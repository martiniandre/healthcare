# Healthcare — Fluxo de Skills Automatizados

Orquestração de automações do ecossistema Healthcare (stack Go + React/Vite). Este diretório contém as **skills de projeto** versionadas e o ponto central de convenções; a skill orquestradora `healthcare-autopilot` vive globalmente em `~/.agents/skills/`.

## Skills de projeto (`.opencode/skills/`)

| Skill | Gatilho | O que faz |
| --- | --- | --- |
| `go-module-gen` | "crie o módulo backend [Nome] com os campos [Campos]" | Gera módulo Go gRPC completo (proto, pb stubs, model, repository, service, grpc_handler, register, RBAC, migration opcional, testes unitários) e pluga no `main.go`. |
| `frontend-module-gen` | "crie o módulo frontend [Nome]" | Gera módulo React/Vite (types, api, queries, schemas Zod, página, componentes, rota, i18n, testes Vitest + E2E Playwright offline). |
| `migration-gen` | "crie a migration", "nova migration" | Cria par `{NNN}_{snake_case}.up.sql`/`.down.sql` em `backend/migrations/` (golang-migrate), numerado após o maior existente, com DDL reversível. |
| `issue-to-feature` | "implemente a issue", "issue #N" | Lê a issue com `gh`, cria branch `feature/<kebab>`, executa `go-module-gen`/`frontend-module-gen`, registra RBAC e abre PR referenciando a issue. |
| `healthcare-commit` | "commit e push", "publique", "finaliza" | Valida integridade local (build/testes/vet/lint), staged com `git add .`, commit Conventional Commits e push automático. |

## Skill orquestradora (global)

`healthcare-autopilot` (`~/.agents/skills/healthcare-autopilot/`) roteia os pedidos para o fluxo certo e encadeia pipelines de ponta a ponta:

- **Backend module:** `go-module-gen` → validação → `healthcare-commit`
- **Frontend module:** `frontend-module-gen` → validação → `healthcare-commit`
- **Full feature:** backend + testes unitários (`go-module-gen`) → frontend + testes Vitest + E2E (`frontend-module-gen`) → validações completas → commit+push → PR
- **Migration:** `migration-gen` → validação → `healthcare-commit`
- **Issue to Feature:** `issue-to-feature` → validações completas → commit+push → PR
- **Changelog:** gera `CHANGELOG.md` (`node scripts/changelog.mjs`) → commit
- **Finalize:** `healthcare-commit`
- **Review:** validações + PR

## Scripts de automação (`scripts/`)

| Script | O que faz |
| --- | --- |
| `conventions-check.mjs` | Gate de convenções: erro fatal para variáveis de 1 letra e aviso para comentários (com exceções: swagger/doc blocks/diretivas/receivers de teste). Usa `--files <paths>` para checar só arquivos de um push e `--strict` para tratar avisos como erro. |
| `changelog.mjs` | Gera `CHANGELOG.md` na raiz a partir do histórico de Conventional Commits (`git log`), agrupado por tipo. |

## Plugin Skills & Flows (global)

`~/.config/opencode/plugins/skills-flows.ts` — plugin global que gerencia skills e flows de qualquer projeto. Tools expostas ao agente:

| Tool | O que faz |
| --- | --- |
| `list_skills` | Lista skills de projeto + globais com descrição e escopo. |
| `create_skill` | Cria `SKILL.md` com frontmatter válido (projeto ou global), sem sobrescrever. Aceita template pronto (`go-module`, `frontend-module`, `commit`, `review`, `migration`, `issue-feature`). |
| `validate_skills` | Varre todas as skills e reporta frontmatter ausente ou nome fora do padrão. |
| `list_flows` | Lista os flows do registry do projeto atual. |
| `run_flow` | Executa as validações de um flow via shell e orienta os próximos passos (carregar skill, commitar). |
| `chat.message` (hook) | Detecta frases de trigger nos flows do registry e injeta uma sugestão sintética apontando o `run_flow` adequado. |

Os flows ficam no registry do projeto em `.opencode/flows/registry.json` (neste repo: `backend-module`, `frontend-module`, `full-feature`, `migration`, `changelog`, `issue-feature`, `finalize`, `review`). O plugin resolve o registry do repo atual automaticamente (`<worktree>/.opencode/flows/registry.json`).

## Automação Git (`.githooks/` e CI)

- `.githooks/commit-msg` — valida Conventional Commits v1.0.0 a cada commit.
- `.githooks/pre-push` — check antes do push: `[1/4]` go vet, `[2/4]` RBAC gate (todo RPC de proto compilado registrado em `policy.go`), `[3/4]` convenções nos arquivos do push (`node scripts/conventions-check.mjs --files`), `[4/4]` npm lint.
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
