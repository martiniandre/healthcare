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
| `GCP_PROJECT_ID` | Identificador do projeto Google Cloud Platform | `my-healthcare-project` |
| `GCP_LOCATION_ID` | Região onde a Healthcare API está provisionada | `us-central1` |
| `GCP_DATASET_ID` | Nome do Dataset Clínico na GCP | `hospital-dataset` |
| `GCP_FHIR_STORE_ID` | Identificador do FHIR Store ativo | `main-fhir-store` |

---

## 🐳 Inicialização & Orquestração Local

### 1. Inicializar Serviços de Infraestrutura (Postgres & Redis)
```bash
docker compose up -d
```

### 2. Executar Migrações do Banco de Dados
O backend executa migrações automáticas de banco de dados (`migrations.Run`) no bootstrapping do `main.go`. Caso necessite de ferramentas manuais, aplique os comandos contra a pasta `/migrations`.

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
