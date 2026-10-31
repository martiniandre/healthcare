# Trunk-Based Development com Feature Flags

Este documento define o padrão do projeto Healthcare para integrar código continuamente em `main` usando **feature flags** para esconder trabalho incompleto. Complementa a estratégia de branching em [GIT_WORKFLOW.md](../GIT_WORKFLOW.md) e a decisão registrada em [ADR 0004](adr/0004-trunk-based-development-with-feature-flags.md).

---

## 1. Princípios

- **`main` é sempre deployável.** Trabalho incompleto entra em `main`, mas invisível, atrás de uma flag desligada.
- **Flag é mecanismo de release, não de segurança.** Autorização (`internal/app/policy/policy.go`), auditoria, criptografia e consentimento nunca dependem de flag.
- **Padrão é OFF (fail-closed).** Flag ausente, não cadastrada ou com erro de leitura resolve para `false`, mantendo o comportamento anterior.
- **Compatibilidade de dados.** Com a flag ON, os dados gravados devem continuar válidos quando a flag voltar a OFF. Dados clínicos persistem somente no FHIR (GCP Healthcare API) e não podem ganhar um formato incompatível sem migração planejada.
- **Uma flag, uma intenção.** Cada flag controla um comportamento. Não existem flags "guarda-chuva" que ligam vários comportamentos juntos.

---

## 2. Tipos de Flag

| Tipo | Propósito | Vida útil esperada | Exemplo |
|------|-----------|--------------------|---------|
| `release` | Esconder feature em desenvolvimento ou liberar gradualmente | Dias a poucas semanas. Removida após 100% de rollout | `ff_clinical_timeline_v2` |
| `ops` | Kill switch de integração externa ou capacidade pesada (ex.: IA de imagem) | Longa, com dono e motivo documentados | `ff_imaging_llm_analysis` |
| `experiment` | Comparação A/B de comportamento | Curta, com métrica de decisão definida | `ff_schedule_layout_b` |

Flags **não** são usadas para conceder permissões a usuários. Permissões são sempre RBAC.

---

## 3. Nomenclatura

- **Chave da flag:** `ff_<dominio>_<feature>`, em `snake_case`, descrevendo a feature e não o ticket. Ex.: `ff_patients_birth_date_localized`. Nunca usar número de issue ou nome de pessoa.
- **Variável de ambiente (backend):** a mesma chave em `UPPER_SNAKE_CASE`, ex.: `FF_PATIENTS_BIRTH_DATE_LOCALIZED=true`.
- **Campo Go:** o nome em `PascalCase` dentro da struct de flags, ex.: `FeatureFlags.PatientsBirthDateLocalized`.
- **Domínios:** os mesmos módulos do backend e as áreas do frontend (`auth`, `patients`, `clinical`, `imaging`, `staff`, `schedule`, `telemetry`, `llm`).

---

## 4. Onde as Flags Vivem

### Backend (Go)

- Flags são declaradas em uma struct dentro de `internal/shared/config` e carregadas de variáveis de ambiente em `config.Load()`, com default `false`.
- A struct é injetada nas dependências dos módulos (`Dependency` do `register.go`). **Não** ler `os.Getenv` nem variáveis globais dentro de services, repositórios ou handlers.

### Frontend (React)

- O frontend recebe as flags do backend por um endpoint autenticado (`GET /api/v1/feature-flags`), para que backend e frontend nunca divirjam sobre o estado de uma feature.
- O acesso é feito por um hook único, `useFeatureFlag(key)`. Enquanto as flags carregam, o valor é `false`.

### Registro

- Toda flag existente deve constar na tabela da seção 9. A inclusão de uma flag e sua remoção acontecem no mesmo PR que a cria ou a remove.

---

## 5. Padrões de Código

### Go

Verificar a flag em **um único ponto de decisão**, na borda do fluxo (handler ou entrada do service). Repositórios não conhecem flags.

```go
func (service *TimelineService) BuildTimeline(ctx context.Context, patientID string) (*Timeline, error) {
    if service.flags.ClinicalTimelineV2 {
        return service.buildTimelineV2(ctx, patientID)
    }
    return service.buildTimelineLegacy(ctx, patientID)
}
```

Regras:
- Os dois caminhos ficam em métodos com nomes explícitos (`...V2`, `...Legacy`). Nunca `if` aninhado dentro de outro `if`.
- Novos endpoints gRPC protegidos por flag continuam registrados em `policy.go`. Endpoint bloqueado pela flag não pode ser "liberado" pela ausência de política.
- Erros do caminho novo passam por `apperrors.ToGRPCStatus(err)` como qualquer outro.

### Frontend

Gate na **rota ou no ponto de composição**, não espalhado em componentes filhos.

```tsx
const isTimelineV2Enabled = useFeatureFlag('ff_clinical_timeline_v2')

return isTimelineV2Enabled ? <ClinicalTimelineV2 /> : <ClinicalTimeline />
```

Regras:
- Nunca `if` de flag dentro de um componente que já recebe os dados prontos.
- O caminho antigo continua funcionando e testado enquanto a flag existir.

---

## 6. Ciclo de Vida de uma Mudança

1. **Branch curta** a partir de `main`, nomeada conforme `GIT_WORKFLOW.md` (`feat/`, `fix/`, `refactor/`). Duração máxima de **2 dias úteis**.
2. **Criar a flag** no mesmo PR inicial, com default `false`, e cadastrá-la na seção 9.
3. **Implementar atrás da flag.** Adicionar testes para os dois valores.
4. **Merge em `main` com a flag OFF** (squash merge, conforme `GIT_WORKFLOW.md`). A feature fica invisível em produção.
5. **Ativar por ambiente**, na ordem `dev` → `staging` → `produção`, com a ativação registrada no PR de rollout ou no changelog.
6. **Remover a flag** em até **30 dias** após 100% de rollout: apagar a flag, o caminho antigo e a entrada na tabela. A remoção é um PR `refactor` ou `chore`.

Flag com mais de 30 dias após 100% de rollout é tratada como dívida técnica e revisada na sprint seguinte.

---

## 7. Testes

- **Go:** testes table-driven que executam os dois valores da flag com o mesmo `MockRepository`. Validar: caminho ON, caminho OFF, default quando a flag não está configurada e falha do repositório em cada caminho.
- **Frontend:** testes com um provider de flags injetando `true` e `false`, validando a renderização dos dois caminhos.
- **CI (meta):** a pipeline deve executar a suíte com todas as flags em OFF e em ON. Hoje `.github/workflows/ci-backend.yml` e `ci-frontend.yml` rodam apenas com os defaults; a matriz ainda precisa ser implementada.

---

## 8. Segurança e Compliance (HIPAA/LGPD)

- Uma flag nunca desliga RBAC, auditoria, criptografia em trânsito ou repouso, nem registro de consentimento.
- Toda alteração de flag em produção deve ser rastreável: quem alterou, quando e por quê. Registrar no PR de rollout ou no changelog.
- Flags que alteram fluxos clínicos (`clinical`, `imaging`) exigem revisão de um responsável clínico no PR de ativação, além da aprovação técnica.
- O valor de uma flag nunca carrega segredo. Credenciais continuam fora do repositório, em variáveis de ambiente de deploy.

---

## 9. Registro de Flags

| Chave | Tipo | Dono | Criada em | Status | Remover até | Descrição |
|-------|------|------|-----------|--------|-------------|-----------|
| _nenhuma flag cadastrada ainda_ | | | | | | |

---

## 10. Anti-padrões

- Flag sem default explícito ou com default `true`.
- Flag dentro de flag (`if a { if b { ... } }`).
- Flag usada como controle de acesso de usuário.
- Flag lida em runtime via `os.Getenv` ou por um componente que não é o ponto de decisão.
- Código do caminho antigo removido antes da flag ser removida, ou flag removida com o caminho antigo ainda vivo.
- Flag de `release` sem data de remoção.
- Branch que vive mais que 2 dias úteis sem merge, mesmo com flag.
