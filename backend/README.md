# HealthCare Core Backend Engine

Motor de alta performance baseado em Go, gRPC e Google Cloud Healthcare API (FHIR/DICOM) estruturado de forma modular para fornecer suporte a prontuários eletrônicos robustos, gestão operacional de equipes médicas e processamento assíncrono de exames PACS de alta fidelidade.

---

## 🏛️ Filosofia Arquitetural & Regras de Persistência

O sistema adota uma divisão estrita sobre onde os dados residem, garantindo compliance com normas HIPAA/LGPD e isolamento completo de dados sensíveis de pacientes (PHI):

### 1. Nuvem Clínica (Google Cloud Healthcare API - FHIR Store)
Nenhum dado clínico evolutivo de paciente é armazenado localmente em banco de dados relacional. Toda a evolução médica reside na nuvem segura da Google Cloud via recursos padronizados FHIR:
* **Recursos Gerenciados**: `Patient`, `Observation` (sinais vitais), `Encounter` (consultas), `Condition` (diagnósticos/patologias), `DiagnosticReport` (laudos médicos).
* **Conectividade**: Realizada através do cliente centralizado `shared/healthcare.Client`, autenticado nativamente via Google Cloud SDK.

### 2. Banco Operacional Local (PostgreSQL)
Persiste apenas dados que coordenam o sistema, acessos operacionais e o catálogo de rastreabilidade de arquivos pesados DICOM:
* **Tabelas Locais**: Controle de acesso (`auth`), credenciamento de colaboradores (`staff`), registros estruturados de uploads PACS/DICOM concluídos (`imaging_studies`).
* **Conectividade**: Gerenciamento pool concorrente seguro com `pgxpool`.

---

## 🔒 Camadas de Segurança, Compliance & Hardening

Toda a troca de dados gRPC é submetida a uma cadeia contínua de segurança aplicada a requisições Unary e streams bidirecionais contínuos:

1. **Autenticação JWT (HttpOnly & Secure Cookies)**:
   Decodificação automática de claims em cabeçalhos de metadados com transporte seguro, mitigando vulnerabilidades de Cross-Site Scripting (XSS).
2. **Defesa Ativa contra Ataques de Assinatura JWT**:
   Validação explícita de `SigningMethodHMAC` que impossibilita exploits usando algoritmos `none` ou chaves assimétricas fracas.
3. **Double-Submit CSRF Protection**:
   Exigência de envio casado do cabeçalho `x-csrf-token` idêntico ao cookie encriptado `csrf_token` gerado em tempo de login, blindando o sistema contra Cross-Site Request Forgery em mutações. Métodos com prefixo `Get`, `List`, `Search` e `Health` são isentos.
4. **Matriz RBAC Declarativa Unificada**:
   Toda permissão de rota é declarada em `internal/app/interceptor/permissions.go` mapeando roles (`RoleAdmin`, `RoleDoctor`, `RoleNurse`, `RoleReception`, `RolePatient`). Métodos não mapeados são bloqueados por padrão.
5. **Logs de Auditoria de Prontuário (Audit Trail)**:
   Emissão assíncrona de registros contendo `caller_user_id`, `caller_role`, `method` e dados mutados para histórico contábil de acessos a prontuários médicos.

---

## 📁 Organização de Pastas do Backend

```text
backend/
├── cmd/
│   └── api/                    # Ponto de entrada (main.go - Bootstrapping)
├── migrations/                 # Migrações puras do banco de dados (SQL)
├── proto/                      # Contratos e definições do Protocol Buffers (.proto)
├── internal/
│   ├── app/                    # Kernel central gRPC e Registro de Interceptors
│   │   ├── interceptor/        # Cadeias de Segurança (Auth, CSRF, RateLimit, Timeout)
│   │   └── server.go           # Instanciação unificada do gRPC Engine
│   ├── shared/                 # Componentes transversais agnósticos a domínio
│   │   ├── apperrors/          # Catálogo centralizado de erros e enums gRPC
│   │   ├── cache/              # Instanciação da conexão do Redis
│   │   ├── database/           # Pool e gerenciamento pgxpool
│   │   ├── fhir/               # Schemas de conversão FHIR
│   │   ├── healthcare/         # Integração direta com Google Cloud API
│   │   └── logger/             # Estruturação e inicialização do slog
│   └── modules/                # Módulos encapsulados baseados em domínio
│       ├── auth/               # Autenticação, registro e geração de credenciais
│       ├── clinical/           # Gerenciador de prontuários eletrônicos (FHIR)
│       ├── health/             # Monitoramento de status ativo da infraestrutura
│       ├── imaging/            # Ingestão DICOM PACS, streaming bidirecional e workers
│       ├── patients/           # Ficha cadastral e rastreamento de pacientes (FHIR)
│       └── staff/              # Gerenciamento de equipe clínica (PostgreSQL local)
```

---

## ⚙️ Variáveis de Ambiente (`.env`)

Copie o arquivo `.env.example` para `.env` e configure conforme as credenciais do seu ecossistema local e da nuvem GCP:

| Variável | Descrição | Exemplo |
| :--- | :--- | :--- |
| `APP_PORT` | Porta de escuta do servidor gRPC | `50051` |
| `APP_ENV` | Ambiente de execução | `development` |
| `DB_URL` | String de conexão do PostgreSQL operacional | `postgres://healthcare_user:healthcare_password@localhost:5432/healthcare_db?sslmode=disable` |
| `REDIS_URL` | Endereço do Redis para cache e controle de tráfego | `localhost:6379` |
| `SENTRY_DSN` | Chave de monitoria e tracking de exceções Sentry | `https://examplePublicKey@o0.ingest.sentry.io/0` |
| `JWT_SECRET` | Chave simétrica secreta para criptografia de tokens | `32-character-secret-key-goes-here` |
| `JWT_SECRET_PREVIOUS` | (Opcional) Chaves antigas aceitas apenas na validação, separadas por vírgula; usada no rodízio de `JWT_SECRET` | `old-secret-1,old-secret-2` |
| `GCP_PROJECT_ID` | Identificador do projeto Google Cloud Platform | `my-healthcare-project` |
| `GCP_LOCATION_ID` | Região onde a Healthcare API está provisionada | `us-central1` |
| `GCP_DATASET_ID` | Nome do Dataset Clínico na GCP | `hospital-dataset` |
| `GCP_FHIR_STORE_ID` | Identificador do FHIR Store ativo | `main-fhir-store` |

---

## 🔑 Sessão e Rodízio de `JWT_SECRET`

O access token (JWT, cookie `token`) dura 15 minutos. O refresh token opaco (cookie `refresh_token`, `HttpOnly`, `Path=/api/v1/auth`) é rotativo, guardado apenas como hash SHA-256 na tabela `refresh_tokens`, e a sessão absoluta dura 12 horas. O reuso de um refresh token já rotado revoga toda a família (sessão) e o logout revoga apenas a família atual. `POST /api/v1/auth/refresh` exige o cabeçalho `X-CSRF-Token` igual ao cookie `csrf_token` e revalida `is_active` e a role atual do usuário.

Todo JWT carrega o cabeçalho `kid` (8 primeiros caracteres hex do SHA-256 da chave). A assinatura usa sempre `JWT_SECRET`; a validação aceita `JWT_SECRET` e `JWT_SECRET_PREVIOUS`. Tokens sem `kid` (emitidos antes do rodízio) são validados tentando todas as chaves.

Procedimento de rotação sem derrubar sessões:

1. Gere a nova chave (mínimo de 32 caracteres; em produção todas as chaves são validadas) e defina `JWT_SECRET=<nova>` e `JWT_SECRET_PREVIOUS=<antiga>` (várias antigas separadas por vírgula).
2. Faça o deploy de todas as réplicas. Access tokens antigos continuam válidos e novos tokens saem assinados pela chave nova.
3. Após 15 minutos (vida do access token), remova `JWT_SECRET_PREVIOUS` e faça novo deploy. Os refresh tokens não dependem da chave JWT e seguem válidos.

---

## 📮 Outbox de Eventos (`internal/shared/outbox`)

Os eventos de dominio (`telemetry.alert`, `exam.complete`, `encounter.created`, `patient.created`, `report.ready`, `system.notification`, `appointment.*`) passam por uma outbox transacional no PostgreSQL (tabela `outbox_events`, migration 025). O `Publish()` do bus grava o evento de forma duravel e um dispatcher dentro do proprio processo da API (polling de 1s, lote de 50, `FOR UPDATE SKIP LOCKED`, sem broker externo) entrega aos handlers registrados com `Subscribe()`.

- **Entrega at-least-once:** falha do handler agenda nova tentativa com backoff exponencial e jitter (2s, 4s, ... ate 5min); apos 8 tentativas o evento vai para dead letter (`dead_at`). Um claim abandonado por queda do processo volta a ser elegivel quando o lease de 60s expira.
- **Idempotencia:** cada evento tem um `event_id` (`eventbus.Event.ID` e `eventbus.EventIDFromContext(ctx)`). Notificacoes derivam o id da notificacao do `event_id`, entao reentrega nao duplica.
- **Eventos sem assinante** nao sao gravados. Eventos processados sao removidos apos 24h; dead letters ficam para inspecao (`SELECT name, attempts, last_error FROM outbox_events WHERE dead_at IS NOT NULL`).
- **Dentro de uma transacao:** use `outbox.Bus.PublishWith(ctx, tx, event)` para gravar o evento na mesma transacao do dado. Os publishers atuais ainda usam `Publish()` apos o commit (janela minima entre commit e gravacao).
- **Dados sensiveis (PHI):** o payload e exatamente o `Data` que os eventos ja carregavam; nao inclua nada alem disso. Os logs de retry e dead letter trazem apenas nome do evento, tentativa e erro, nunca o payload. `last_error` guarda a mensagem de erro truncada em 500 caracteres.
- **Observabilidade:** `Dispatcher.Stats(ctx)` retorna pendentes, dead letters e a idade do evento pendente mais antigo.
- O bus em memoria (`eventbus.New()`) continua disponivel para testes unitarios.

## Seguranca de prescricao (`internal/modules/medication`)

`POST /api/v1/patients/{patientFhirId}/medications/check` (Admin, Doctor) e o `POST /api/v1/encounters/{encounterFhirId}/medications` cruzam o medicamento pedido com as alergias ativas e os `MedicationRequest` ativos do paciente.

- **Alergia ativa (por principio ativo ou classe) bloqueia** a prescricao com `409 {"code":"MEDICATION_ALLERGY_BLOCKED","alerts":[...]}`, a menos que `overrideJustification` (minimo 10 caracteres) seja enviada por um DOCTOR ou ADMIN. A sobreposicao grava uma linha em `medication_safety_overrides` e um registro em `audit_logs` (`action = medication_safety_override`, so ids e codigo do medicamento, sem PHI livre). **Interacao medicamentosa so avisa** (`blocking: false`) e volta em `warnings`.
- **Catalogo proprio versionado** (migration 026: `medication_classes`, `medication_catalog`, `medication_interactions`). O servico depende apenas da interface `SafetyCatalog` (`LoadSnapshot`), entao a fonte pode ser trocada (RxNorm, ANVISA, terceiros) sem mexer nas regras. Mensagens sao chaves i18n (`messageKey`) traduzidas pelo frontend.
- **AVISO CLINICO:** a migration 027 carrega apenas um catalogo de AMOSTRA (`catalog_version = sample-1`: penicilinas, cefalosporinas, AINEs, dipirona, varfarina e tres interacoes). Ele NAO foi validado clinicamente. **Um medico ou farmaceutico DEVE revisar e substituir esse conteudo antes de qualquer uso em producao.**
- A verificacao falha fechada: se catalogo, alergias ou medicamentos ativos nao puderem ser lidos, a prescricao retorna 503 em vez de prosseguir sem checagem. O extrato de medicamentos ativos segue o limite de paginacao do cliente FHIR (`maxPaginationPages`).

## 🐳 Inicialização & Orquestração Local

### 1. Inicializar Serviços de Infraestrutura (Postgres & Redis)
```bash
docker compose up -d
```

### 2. Executar Migrações do Banco de Dados
Em desenvolvimento/teste a API migra o banco no startup (`RUN_MIGRATIONS_ON_START` default `true`). Em produção (`APP_ENV=production`) o default é `false`: a API apenas valida que a versão do schema é >= a última migration disponível e aborta com erro claro caso contrário (ou se o schema estiver `dirty`).

As migrations rodam como etapa separada via `cmd/migrate` (lê `DB_URL`; `MIGRATIONS_PATH` opcional, default `file://migrations`):
```bash
DB_URL=postgres://... go run ./cmd/migrate
```
A imagem Docker contém `/app/migrate` e `/app/migrations`.

**Deploy (Fly.io):** o app `healthcare-api-small-dust-5285` roda no Fly.io. O workflow `.github/workflows/deploy-backend.yml` dispara o deploy a cada push na `main` que altere `backend/`, `Dockerfile`, `deploy/`, `fly.toml` ou o próprio workflow, e também manualmente (**Actions → Deploy Backend to Fly.io → Run workflow**).

**Migrations antes do rollout.** O `fly.toml` define `release_command = '/app/migrate'`. O Fly executa a migração em uma máquina temporária com a imagem nova e os secrets do app, e só troca as réplicas se ela terminar sem erro. Se a migração falhar, a versão atual continua no ar. Como `APP_ENV=production`, a API não migra no startup: ela só verifica se o schema está na versão esperada.

**Configuração única:**

1. Gerar o token de deploy e cadastrá-lo como secret do GitHub (`FLY_API_TOKEN`):

```bash
fly tokens create deploy -a healthcare-api-small-dust-5285
gh secret set FLY_API_TOKEN --repo martiniandre/healthcare
```

2. Conferir os secrets do app no Fly (os nomes aparecem, os valores não):

```bash
fly secrets list -a healthcare-api-small-dust-5285
```

Os necessários são `DB_URL`, `JWT_SECRET`, `REDIS_URL`, `GCP_PROJECT_ID`, `GCP_DATASET_ID`, `GCP_FHIR_STORE_ID`, `GCP_KEY_B64` (JSON da service account em base64) e `FRONTEND_URL`. Opcionais: `JWT_SECRET_PREVIOUS`, `GCP_LOCATION_ID`, `GCP_VERTEX_MODEL`, `GCS_BUCKET_NAME`, `SENTRY_DSN`. Para cadastrar um valor sem deixá-lo no histórico do shell:

```bash
fly secrets set NOME_DO_SECRET -a healthcare-api-small-dust-5285
```

**Antes do primeiro deploy com a migration `028`:** rode a consulta de verificação descrita no cabeçalho de `backend/migrations/032_imaging_study_identity.up.sql`.

**Verificação do app:** `fly apps list` mostrava o app como `suspended`. Confirme que ele está ativo (ou reative-o pelo painel do Fly) antes de disparar o deploy.

### 3. Compilar e Rodar o Servidor Go
```bash
go run cmd/api/main.go
```

---

## 🛠️ Contratos gRPC Disponíveis

### `AuthService`
* `Login` - Realiza a autenticação de credenciais, injetando cookies `token` e `csrf_token` via headers.
* `Register` - Registra novos colaboradores ou administradores.
* `Logout` - Invalida e limpa cookies de credenciais ativas.

### `StaffService`
* `CreateEmployee` - Cria novo registro de profissional de saúde.
* `GetEmployee` - Resgata perfil operacional de um colaborador.
* `ListEmployees` - Listagem paginada da equipe médica.
* `DeactivateEmployee` - Desativação administrativa de permissões.

### `PatientService`
* `CreatePatient` - Registra uma ficha no FHIR Store na Google Cloud.
* `GetPatient` - Retorna perfil cadastral do paciente via FHIR ID.
* `GetPatientByDocument` - Resgata perfil de paciente cruzando documento CPF/RG.
* `ListPatients` - Retorna lista paginada de todos os registros clínicos.

### `ClinicalService`
* `CreateEncounter` - Abre uma consulta clínica de atendimento.
* `CreateObservation` - Registra medições e sinais vitais evolutivos.
* `CreateCondition` - Anexa diagnósticos com CID-10 e status de atividade.
* `ListPatientEncounters` - Recupera histórico de atendimentos.
* `ListPatientObservations` - Recupera métricas históricas do paciente.
* `ListPatientConditions` - Recupera lista de patologias ativas.

### `ImagingService`
* `UploadPACS` - Canal stream bidirecional de ingestão de arquivos DICOM, validando magic bytes `DICM` no offset 128 com retorno contínuo de progressos em bytes.
* `GetPresignedURL` - Gera link assinado temporário e seguro do GCS para visualização imediata do exame no frontend.
* `ListPatientStudies` - Lista os estudos DICOM persistidos associados a um paciente.

### `Health`
* `Check` - Serviço padrão de saúde (`grpc.health.v1`) que valida de forma integrada conexões do Postgres e do Redis, fornecendo relatórios automatizados de integridade.

---

## 🧪 Rodando Testes Unitários

Todos os testes unitários são desenvolvidos usando o framework padrão Go acoplado a asserções avançadas do `testify`:

```bash
go test -v ./...
```
