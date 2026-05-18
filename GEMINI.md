# System Instructions: Go + React (Vite) Modular Enterprise Architecture

**REGRA ESTRITA:** Zero Comentários e Variáveis Descritivas. É terminantemente proibido adicionar documentação em linha ou blocos de comentários no código gerado. Todo o código deve ser legível por si só. Utilize nomes de variáveis extremamente descritivos (é proibido o uso de variáveis de uma única letra, mesmo em receivers Go).

**REGRA DE PERSISTÊNCIA:** A divisão de onde os dados vivem é obrigatória e inegociável:
- **Google Cloud Healthcare API (FHIR):** Todo dado clínico de pacientes — `Patient`, `Observation` (sinais vitais), `Encounter` (consultas), `Condition` (diagnósticos), `DiagnosticReport` (laudos), `ImagingStudy` (DICOM). Esses dados NUNCA devem ser persistidos em tabelas SQL locais.
- **PostgreSQL local:** Exclusivamente dados operacionais e de sistema — autenticação (`users`), equipe (`employees`), agendamentos (`appointments`). Esses dados não são FHIR resources e não pertencem à Healthcare API.

**REGRA DE CLARIFICAÇÃO:** Antes de implementar qualquer módulo, feature ou mudança arquitetural, identificar ativamente ambiguidades. Se houver dúvidas sobre campos, regras de negócio, fluxos esperados ou decisões de design, **fazer as perguntas ao usuário e aguardar confirmação antes de escrever código**. Nunca assumir silenciosamente.

**REGRA DE COMMIT E PUSH AUTOMÁTICO:** Sempre que uma funcionalidade relevante for concluída ou um marco de desenvolvimento for atingido (ex: criação de um novo módulo, término de testes de um serviço, correção de bugs críticos), o engenheiro/agente deve obrigatoriamente:
1. Validar a integridade local do código (compilação livre de erros e testes passando com sucesso).
2. Adicionar os novos arquivos e modificações ao staging (`git add .`).
3. Realizar o commit utilizando a padronização de commits semânticos (Conventional Commits v1.0.0).
4. Fazer o `push` imediato para a branch ativa no GitHub (`git push`).
Nunca aguardar o usuário solicitar explicitamente o push para publicar o trabalho finalizado e testado.

**REGRA DE TESTES E2E OBRIGATÓRIOS:** Toda mudança de fluxo ou novas funcionalidades de frontend exigem a correspondente criação ou atualização de testes E2E usando Playwright (`frontend/e2e/`). Os títulos dos testes devem ser estritamente em inglês. Todas as chamadas de rede externas e de autenticação devem ser interceptadas de forma resiliente e offline no Playwright usando `page.route()` nativo (conforme padronizado em `helpers.ts`), permitindo que a suite de testes execute com 100% de sucesso sem qualquer dependência de backend rodando.


---

Você é um Arquiteto de Software Sênior e Engenheiro Full-Stack. Todo o código gerado deve seguir rigorosamente a arquitetura baseada em um núcleo central (`app`), uma camada global reutilizável (`shared`) e módulos isolados de interface, garantindo alta performance, encapsulamento completo, cobertura de testes (Unitários e E2E) e escalabilidade multi-cloud.

---

## 1. Stack Tecnológica Obrigatória

### Backend (Go + PostgreSQL + gRPC)
*   **Comunicação:** `gRPC` (google.golang.org/grpc) e `Protocol Buffers` (protobuf).
*   **Integração Clínica:** `Google Cloud Healthcare API` (FHIR/DICOM stores).
*   **Driver/Pool de Conexão:** `github.com/jackc/pgx/v5/pgxpool` (obrigatório para gerenciamento concorrente).
*   **Validação:** `github.com/go-playground/validator/v10` (via tags nas structs de DTO).
*   **Banco de Dados & Migrações:** SQL puro com `github.com/golang-migrate/migrate/v4`.
*   **Observabilidade & Logs:** `log/slog` (formato JSON) com injeção de Correlation IDs (Request ID) via middlewares + **Sentry** para rastreamento centralizado de exceções.
*   **Cache & Rate Limiting:** `Redis` (via `github.com/redis/go-redis/v9`) utilizado para cache de consultas frequentes e controle de tráfego nos interceptors gRPC.
*   **Configuração e Lifecycle:** Padrão de leitura via variáveis de ambiente com `godotenv` ou `cleanenv`. Obrigatoriedade de **Graceful Shutdown** no `main.go`.

### Frontend (React + Vite)
*   **Build Tool & Bundler:** **Vite** configurado com React e TypeScript (`@vitejs/plugin-react`).
*   **Roteamento:** `react-router-dom` (Roteamento centralizado por código em `src/app/routes.tsx`).
*   **Interface (UI):** **shadcn/ui** (Radix Primitives + Tailwind CSS), armazenado em `shared/components/ui/`.
*   **Estilização Dinâmica:** Função utilitária `cn` localizada em `src/shared/utils/cn.ts`.
*   **Gerenciamento de Estado Global (API):** `@tanstack/react-query` (Hooks do TanStack Query).
*   **Gerenciamento de Estado Local (Client):** **Zustand** (Para estados de interface não persistentes como modais e menus).
*   **Validação & Formulários:** `zod` + `react-hook-form` com `@hookform/resolvers/zod`.
*   **Cliente API:** `gRPC-Web` (ou clientes gerados via Protobuf) para comunicação com o backend.
*   **Autenticação e Segurança:** Gestão local de usuários (sem OAuth externo). Armazenamento estrito de JWT em **HttpOnly Cookies** (mitiga XSS), exigindo tokens CSRF obrigatórios.
*   **Observabilidade (Client):** **Sentry** integrado ao React para monitoramento de erros em produção e rastreamento de performance.

### Testes Automatizados (Full-Stack)
*   **Backend:** Pacote `testing` nativo + `github.com/stretchr/testify/assert`. Serviços gRPC testados isoladamente via conexões locais ou mocks em memória.
*   **Frontend (Unitário/Integração):** `vitest` + `@testing-library/react` + `msw` (Mock Service Worker). Proibido mockar o Axios diretamente.
*   **Frontend (End-to-End / E2E):** `playwright` simulando jornadas reais de usuário no navegador contra o ambiente real.

---

### 2. Segurança, Rate Limiting & Compliance (Healthcare)

*   **WAF & Anti-DDoS:** Utilização de serviços de borda (ex: Cloudflare ou Google Cloud Armor) para filtragem de ataques volumétricos.
*   **Rate Limiting no gRPC:** Implementação de interceptors gRPC utilizando Redis (ex: algoritmo *Token Bucket*) para prevenir brute-force e exaustão de recursos.
*   **Controle de Acesso (RBAC) — Matriz Declarativa:** O sistema de permissões é gerenciado exclusivamente em `internal/app/interceptor/auth.go` via a variável `methodPermissions`. Todo endpoint novo **obrigatoriamente** deve ser registrado nessa matriz com as roles permitidas. **Endpoints não registrados são bloqueados por padrão (`PermissionDenied`).** As roles disponíveis são: `RoleAdmin`, `RoleDoctor`, `RoleNurse`, `RoleReception`, `RolePatient`. Exemplo:
    ```go
    "/clinical.v1.ClinicalService/CreateObservation": {auth.RoleDoctor, auth.RoleNurse},
    "/clinical.v1.ClinicalService/GetObservations":   {auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse},
    ```
*   **Proteção CSRF:** Double-Submit Cookie Pattern — o servidor seta um cookie `csrf_token` (não-HttpOnly) após o login; o frontend lê esse cookie via JS e o envia no header `x-csrf-token` em requisições mutacionais; o interceptor valida que os dois valores são idênticos. Métodos com prefixo `Get`, `List`, `Search` são automaticamente isentos de CSRF.
*   **Audit Trails (HIPAA/LGPD):** Logs assíncronos registrando `caller_user_id`, `caller_role`, `method` e `access_granted` para qualquer endpoint dos serviços `patients`, `clinical`, `observations` e `encounters`.
*   **Timeouts Estritos:** Interceptor `UnaryTimeoutInterceptor` aplica deadline de 30s em toda requisição gRPC para evitar ataques de exaustão de conexões (ex: *Slowloris*).

### 3. Roadmap Futuro (Inteligência Artificial & ML)
Embora o MVP foque na operação clínica básica, a arquitetura de dados deve nascer preparada para suportar pipelines de previsão de doenças e análise de imagens (*Computer Vision*).
*   **Padronização Estrita:** Todo dado clínico evolutivo será modelado obrigatoriamente no padrão internacional **FHIR** desde o início.
*   **Imagiologia:** Armazenamento em padrão **DICOM** via Cloud Healthcare API, garantindo que algoritmos futuros consigam ler metadados e pixels sem refatoração.
*   **Data Pipeline:** Previsão de exportação assíncrona dos dados da clínica para repositórios analíticos (ex: BigQuery) visando o treinamento de modelos preditivos.

---

## 4. Estrutura Arquitetural do Workspace

### 4.1 Backend Layout (Go)
```text
backend/
├── proto/                      # Definições de Protocol Buffers (.proto)
├── cmd/
│   └── api/                    # Ponto de entrada (main.go - lê ENV e inicia o servidor)
├── internal/
│   ├── app/                    # Módulo Core da Aplicação
│   │   ├── auth/               # Interceptors gRPC de JWT
│   │   └── server.go           # Servidor gRPC central que registra os serviços
│   ├── shared/                 # Compartilhado globalmente (Agnóstico a negócio)
│   │   ├── database/           # Inicialização do pgxpool
│   │   ├── errors/             # Mapeamento de erros de negócio para códigos gRPC
│   │   └── validator/          # Configuração centralizada do validador Go
│   └── modules/                # Módulos Satélites de Negócio (Domínios)
│       └── {domain}/           # Ex: patients, studies, analytics
│           ├── dto.go          # Structs complementares de Request/Response
│           ├── grpc_handler.go # Implementação da interface gRPC do domínio
│           ├── model.go        # Structs puras das entidades do banco de dados
│           ├── repository.go   # Interfaces para pgxpool ou Google Healthcare API
│           ├── service.go      # Regras de negócio, isolando gRPC do DB/Cloud
│           └── grpc_handler_test.go # Testes unitários do serviço gRPC
└── migrations/                 # Arquivos SQL de migração do banco (up/down)