---

### 📄 `agents.md`

```markdown
# Antigravity Agent Blueprint Specification (React + Vite Pipeline)

**REGRA ESTRITA:** Zero Comentários e Variáveis Descritivas. É terminantemente proibido adicionar comentários no código gerado (seja Go, TypeScript ou Protobuf). O código deve ser autoexplicativo através de boas nomenclaturas (use nomes de variáveis claros e completos, sendo proibido usar variáveis de uma única letra como `s`, `h`, `r`, mesmo em receivers do Go).

**REGRA DE PERSISTÊNCIA:** Antes de criar qualquer módulo, determinar onde o dado vive:
- **Google Cloud Healthcare API (FHIR):** Dados clínicos — `Patient`, `Observation`, `Encounter`, `Condition`, `DiagnosticReport`, `ImagingStudy`. Nesses módulos, `repository.go` bate no cliente `shared/healthcare.Client`, não no `pgxpool`.
- **PostgreSQL local:** Dados operacionais — `auth`, `staff`, `appointments`. Nesses módulos, `repository.go` usa `pgxpool` e requer migration SQL.

Este documento dita a execução técnica e a orquestração das automações do ecossistema para a stack modular Go + React (Vite).

---

## [AG-01] Backend Domain Generator (`go-module-gen`)

*   **Objetivo:** Criar um novo módulo de domínio independente no backend Go e plugá-lo no núcleo central (`app`).
*   **Gatilho:** `"Antigravity, crie o módulo backend [Nome] com os campos [Campos]."`

### Ações do Agente:
0.  **Clarificação (OBRIGATÓRIO):** Antes de iniciar qualquer implementação, identificar ambiguidades nos requisitos. Se houver dúvidas sobre campos, regras de negócio, integrações ou comportamentos esperados, **fazer as perguntas ao usuário e aguardar resposta antes de escrever código**. Nunca assumir silenciosamente.
1.  **Migrations:** Criar migration SQL **apenas para módulos operacionais** (auth, staff, appointments). Módulos clínicos (FHIR) não geram migrations — seus dados vivem na Healthcare API.
2.  **Scaffolding:** Criar a pasta `internal/modules/{domain}/`.
3.  **Protobuf:** Criar o arquivo `proto/{domain}.proto` e rodar a compilação (`protoc` ou `buf`) para gerar os stubs em Go e TypeScript.
4.  **Model & DTOs:** Gerar `model.go` para as entidades do banco de dados e utilizar as structs geradas pelo Protobuf como DTOs (adicionando validações customizadas se necessário).
5.  **Repository:** Gerar `repository.go` criando uma interface para persistência local (`pgxpool`) ou integração remota (`Google Cloud Healthcare API`).
6.  **Service:** Gerar `service.go` encapsulando as regras de negócio e validações clínicas, injetando a interface do repositório.
7.  **gRPC Handler:** Gerar `grpc_handler.go` implementando a interface gRPC gerada, atuando como controller e chamando o Service. Todos os erros devem ser convertidos via `apperrors.ToGRPCStatus(err)`.
8.  **Permissões (OBRIGATÓRIO):** Abrir `internal/app/interceptor/permissions.go` e adicionar cada novo endpoint na `methodPermissions` com as roles permitidas. **Endpoints não registrados são bloqueados por padrão.** A matriz de roles disponíveis é: `RoleAdmin`, `RoleDoctor`, `RoleNurse`, `RoleReception`, `RolePatient`.
9.  **Core Injection:** Criar `register.go` no módulo com uma função `Register(grpcServer *grpc.Server, dep Dependency)` que encapsula todo o wiring interno (repo → svc → handler → RegisterXxxServer). No `main.go`, adicionar apenas `modulename.Register(server.GRPCServer, dep)`.
10. **Testing:** Gerar os testes unitários (`tests/service_test.go`) com `MockRepository` no mesmo pacote do módulo, validando validações, defaults e falhas de repositório.

### Comando de Validação:
```bash
go test -v ./internal/modules/{domain}/...