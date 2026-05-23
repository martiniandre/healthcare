# Diretrizes do Projeto - Healthcare Backend

Este documento estabelece as regras e padrões arquiteturais adotados para o projeto Healthcare, garantindo manutenibilidade, legibilidade, segurança e performance.

## 1. Arquitetura em Módulos (Domain-Driven)

O backend é organizado em módulos independentes dentro de `internal/modules/`. Cada módulo (ex: `auth`, `clinical`, `patients`) deve ter:
- `service.go`: Contém as regras de negócio e validações. Não deve lidar com lógica HTTP.
- `repository.go`: Lida com a persistência de dados.
- `grpc_handler.go`: Exclusivo para os endpoints gRPC.
- `register.go`: Ponto de entrada do módulo, responsável por inicializar dependências (Service/Repo) e retornar as interfaces injetadas.

A camada HTTP (`internal/api/`) é responsável pelo roteamento (`router.go`) e tratamento REST, delegando toda a lógica de negócio para as interfaces de `Service` fornecidas pelos módulos.

## 2. Injeção de Dependências em `main.go`

O `cmd/api/main.go` é estritamente um ponto de montagem (bootstrap).
- **Sem Lógica de Roteamento ou Regra de Negócio**: Nenhuma lógica de negócio, extração de cookies, ou validação direta deve residir no `main.go`.
- **Registro Padrão**: Todo módulo deve expor uma função de registro (`auth.Register(grpcServer, db)`, etc.) que será consumida no `main.go`.

## 3. Segurança e Middlewares

- **CORS e HTTP Headers**: O tratamento de CORS é isolado no pacote `internal/api/middleware`. O roteador utiliza middlewares nativos (`mux`) para aplicar regras de CORS, baseadas em variáveis de ambiente (`secureCookies`).
- **Autenticação**: O middleware de autenticação (`validateHTTPAuth`) injeta de forma segura o `UserID` e a `Role` no contexto utilizando chaves tipadas (`ctxkeys.ContextKey`).
- **Validação JWT / CSRF**: O backend usa verificação JWT em cookies HTTPOnly, com validação paralela de tokens anti-CSRF para requisições mutáveis (`POST`, `PUT`, `DELETE`).

## 4. Gerenciamento de Contexto (`context.Context`)

Sempre passe o `context.Context` como o **primeiro parâmetro** para funções que realizam I/O, banco de dados ou RPCs.

### Uso de Chaves Tipadas (Type-Safe Context Keys)
- Evite o uso de `string` pura ao salvar ou ler valores do contexto.
- Utilize chaves fortemente tipadas localizadas em `internal/shared/ctxkeys` (ex: `ctxkeys.UserIDKey`, `ctxkeys.RoleKey`) para evitar colisões entre pacotes.
- `audit_trail`, `logging` e outros interceptors gRPC sempre leem as permissões a partir destas chaves padronizadas.

## 5. Tratamento de Erros e Performance

- **Mitigação de Timing Attacks**: Funções sensíveis (ex. `auth.Login`) devem realizar rotinas de validação pseudo-constantes para evitar ataques de tempo em contas não existentes.
- **Race Conditions**: Cache/Rate Limiters (como no Redis) devem utilizar processamento atômico (ex. `redis.Pipeline`) para agrupar operações (`Incr`, `Expire`) e evitar condições de corrida em alta concorrência.
- **Filas em Background**: Processos custosos (ex. processamento DICOM) devem usar Workers com goroutines para leitura em blocos e processamento em background (ver `imaging.Worker`).

## 6. Boas Práticas de Código Idiomático (Go)

- Os manipuladores gRPC devem sempre retornar os erros empacotados em códigos gRPC (usando conversão como `apperrors.ToGRPCStatus(err)`).
- Retorne explicitamente interfaces (`auth.Service`, `clinical.Service`) para facilitar injeção de dependências e uso de mocks em testes unitários.
- Use `log/slog` de maneira estruturada (`slog.Error("Mensagem", "chave", valor)`) em invés da biblioteca global padrão de log.
