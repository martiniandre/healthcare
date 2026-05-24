# Healthcare Architecture Notes

Este documento complementa `AGENTS.MD`. As regras operacionais obrigatórias para agentes ficam em `AGENTS.MD`; este arquivo registra decisões arquiteturais de alto nível para orientar PRs e revisões.

## Arquitetura

- Backend em Go com módulos por domínio em `backend/internal/modules`.
- Frontend em React + Vite com módulos por domínio em `frontend/src/modules`.
- Componentes reutilizáveis ficam em `frontend/src/shared/components`.
- Utilitários compartilhados ficam em `frontend/src/shared/utils`.
- Hooks de dados devem ficar próximos do domínio que consome a API.

## Arquitetura Hexagonal e Domínios

- **Fronteiras de Domínio:** Módulos do Core Clínico (`Patients`, `Clinical`, `Imaging`) NÃO devem utilizar PostgreSQL, operando via Google Cloud Healthcare API (FHIR) para compliance HIPAA/LGPD. Módulos Operacionais (`Auth`, `Staff`, `Telemetry`) utilizam PostgreSQL local para controle e performance.
- **Portas e Adaptadores:** Todo módulo backend deve seguir a arquitetura hexagonal. `model.go` para entidades puras; `repository.go` define interfaces (Portas) e implementa o acesso (Adaptadores); `service.go` orquestra a lógica de negócio desconhecendo a camada de transporte; `grpc_handler.go`/`http_handler.go` atuam como Portas de Entrada.
- **Injeção e Composição:** `cmd/api/main.go` atua puramente como *Compositor Root*. A inicialização de dependências, bancos de dados e _wiring_ de serviços ocorre aqui, sendo os serviços registrados através de funções `Register(server, dep)`.
- **Segurança (AOP):** Autenticação, autorização (RBAC), rate limit e tracing devem viver em gRPC Interceptors / Middlewares HTTP, separando requisitos não-funcionais da lógica de negócio.

## Persistência

- Dados clínicos ficam na Google Cloud Healthcare API usando recursos FHIR.
- `Patient`, `Observation`, `Encounter`, `Condition`, `DiagnosticReport` e `ImagingStudy` não devem ser persistidos como dados clínicos em tabelas SQL locais.
- Dados operacionais locais ficam no PostgreSQL.
- Atualmente os dados operacionais planejados são autenticação e equipe.
- O projeto não terá módulo de agendamentos.

## Backend

- Novos módulos devem expor `register.go`, `repository.go`, `service.go`, `grpc_handler.go`, `model.go` quando aplicável e testes unitários.
- Endpoints gRPC novos devem ser registrados em `backend/internal/app/interceptor/permissions.go`.
- Tipos de contrato HTTP ou gRPC auxiliares devem ficar no módulo de domínio, não em `cmd/api/main.go`.
- `cmd/api/main.go` deve permanecer focado em bootstrap, wiring e handlers mínimos.

## Frontend

- Cada domínio deve ter seus próprios tipos, API client e hooks de query/mutation.
- Evitar um client global gigante para todas as features.
- Componentes de página devem orquestrar dados e navegação.
- Componentes filhos devem receber dados e callbacks por props.
- Fluxos novos ou alterados exigem testes E2E em `frontend/e2e`.
- Chamadas externas e autenticação nos E2E devem ser interceptadas com `page.route()`.

## Testes

- Backend: `go test -v ./...`.
- Frontend: `npm run lint`, `npm run build` e `npm run test:e2e`.
- Testes E2E devem ter títulos em inglês.
