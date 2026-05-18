## 📝 Descrição
<!-- Explique resumidamente o objetivo deste Pull Request e o problema que ele resolve -->

## 🛠️ O que foi feito?
- [ ] Implementação de lógica de negócio / UI
- [ ] Criação/Modificação de Schemas ou DTOs Protobuf
- [ ] Criação de Migrations SQL (apenas dados operacionais - PostgreSQL)
- [ ] Criação de testes unitários ou de integração

## 🔒 Segurança & Healthcare Compliance (HIPAA/LGPD)
- [ ] **Persistência Correta:** Dados clínicos salvos estritamente no FHIR Store (GCP Healthcare API); Dados operacionais salvos no Postgres local.
- [ ] **RBAC:** O novo endpoint foi devidamente registrado no interceptor de permissões (`internal/app/interceptor/permissions.go`)?
- [ ] **Sem Segredos:** Verificado que nenhuma chave `.env` ou credenciais foram adicionadas acidentalmente.
- [ ] **Zero Comentários (Regra Estrita):** O código gerado está 100% livre de comentários redundantes e documentação em linha.
- [ ] **Variáveis Descritivas:** Todas as variáveis possuem nomes claros e significativos (sem letras ou siglas únicas).

## 🧪 Como testar?
<!-- Descreva detalhadamente como o revisor pode validar a sua alteração -->
1. Comando para rodar os testes: `go test -v ./internal/modules/{domain}/...` ou `npm run test`
2. Fluxo manual verificado: [Passo a Passo]
