# Plano de Ação — Novas Features e Melhorias

Este plano cobre todas as features sugeridas na revisão do projeto. Ele segue as convenções do `AGENTS.MD` (módulos backend com `register.go`, permissões em `policy.go`, migrations só para módulos operacionais, testes em `tests/`) e do `FRONTEND_GUIDELINES.md`.

Tamanhos são estimativas grosseiras de esforço de uma pessoa: **P** (até 2 dias), **M** (3 a 7 dias), **G** (mais de 1 semana). Valores sujeitos a revisão após o refinamento.

> Regra do projeto: antes de implementar qualquer módulo backend, resolver as **perguntas em aberto** de cada item com o responsável de produto. Nada deve ser assumido em silêncio.

---

## Visão geral e ordem sugerida

| Fase | Foco | Itens |
|------|------|-------|
| 0 | Fundação e pendências | F0.1 a F0.4 |
| 1 | Segurança e base de plataforma | S1 a S5, P1 a P3 |
| 2 | Segurança do paciente | C1, C2 |
| 3 | Operação e relacionamento | C3, C4, C8 |
| 4 | Inteligência e dados em tempo real | C5, C6 |
| 5 | Autonomia do paciente | C7 |
| 6 | Maturidade técnica | Q1 a Q4, P4, P5 |

Dependências principais: S2 (refresh token) antes de S1 (2FA) e C7; P1 (outbox) antes de C3 (lembretes) e C6 (telemetria real); C1 (interações) reaproveita os cabeçalhos de alergia já entregues no frontend.

---

## Fase 0 — Fundação e pendências da migração de UI

| ID | Item | Tamanho |
|----|------|---------|
| F0.1 | Revisão visual no navegador das telas migradas (light/dark, es-ES, mobile 375 px, portal) e correção de desvios | M |
| F0.2 | Rodar a suíte completa de testes em CI com poucos workers e corrigir testes sensíveis a carga (timeout global já elevado para 15 s) | P |
| F0.3 | Mover as chaves `retry` por módulo para `common.retry` e unificar duplicações restantes (`VitalSignValueDisplay` vs `PortalVitalValue`) | P |
| F0.4 | Estender o `ci-frontend.yml` com `tsc -b`, verificação de paridade de idiomas, build e Playwright | P |

Critério de aceite: nenhuma tela com contraste abaixo de AA nos dois temas; CI frontend verde com todos os checks.

---

## Fase 1 — Segurança e plataforma

### S1. Autenticação de dois fatores (TOTP)
- **Objetivo:** exigir segundo fator para Admin e Doctor; opcional para os demais.
- **Backend:** migration com `user_totp` (segredo cifrado, estado, códigos de recuperação com hash); endpoints de enroll, verify e disable; etapa extra no login; bloqueio temporário após falhas; entradas novas em `policy.go`; auditoria de eventos.
- **Frontend:** fluxo de configuração (QR code + códigos de recuperação), tela de desafio no login, aviso de expiração.
- **Perguntas em aberto:** 2FA obrigatório para quais perfis? Haverá recuperação por administrador? Lembrar dispositivo por quantos dias?
- **Aceite:** login de Admin sem segundo fator é bloqueado; códigos de recuperação são de uso único; todos os eventos aparecem nos audit logs.
- **Tamanho:** G. **Depende de:** S2.

### S2. Refresh token e sessão curta
- **Backend:** access token de 15 min, refresh token rotativo com `jti` persistido (tabela `refresh_tokens`), detecção de reuso (revoga a família), `is_active` e role revalidados no refresh.
- **Frontend:** o interceptor de 401 tenta o refresh uma vez antes de deslogar; fila de requisições durante o refresh.
- **Perguntas em aberto:** duração da sessão absoluta? Sessões simultâneas permitidas?
- **Aceite:** desativar um usuário derruba o acesso em até 15 min; reuso de refresh token revoga a sessão.
- **Tamanho:** M.

### S3. Alerta de "quebra de vidro" (acesso sem vínculo)
- **Backend:** modelo de vínculo staff↔paciente (via encontros/agenda); ao acessar paciente sem vínculo, exigir justificativa, registrar no audit log com severidade alta e notificar o Admin (`audit_alert`).
- **Frontend:** diálogo de justificativa antes de abrir o prontuário; selo "acesso excepcional" no cabeçalho.
- **Perguntas em aberto:** o que define "vínculo" e por quanto tempo vale? Quais perfis estão sujeitos?
- **Aceite:** acesso sem vínculo só ocorre com justificativa; relatório de acessos excepcionais disponível para Admin.
- **Tamanho:** M.

### S4. Rodízio de `JWT_SECRET`
- **Backend:** aceitar lista de chaves (`kid` no cabeçalho), assinar com a atual e validar com atuais e anteriores; documentar o procedimento de rotação.
- **Aceite:** trocar a chave sem derrubar sessões válidas.
- **Tamanho:** P.

### S5. Matriz de permissões única
- **Backend:** definir capacidades (`patient:read`, `report:sign`…) uma vez e derivar gRPC e HTTP; manter `rbac_gate_test` como rede de proteção.
- **Frontend:** consumir a mesma lista para menu e rotas (já existe `route_access.ts`; alinhar por geração ou teste de contrato).
- **Aceite:** mudar uma permissão exige editar um único lugar; testes de contrato garantem alinhamento.
- **Tamanho:** M.

### P1. Outbox de eventos
- **Backend:** tabela `outbox_events` gravada na mesma transação do dado; dispatcher com retry e backoff; idempotência por `event_id`; métricas de atraso.
- **Aceite:** nenhuma notificação se perde em queda do processo; reentrega é idempotente.
- **Tamanho:** M.

### P2. Rate limit compartilhado no Redis
- **Backend:** janela deslizante no Redis (script Lua), fallback local se o Redis cair, chaves por IP e por usuário; reaproveitar `TRUSTED_PROXIES`.
- **Aceite:** o limite vale entre múltiplas instâncias; teste de carga confirma o comportamento.
- **Tamanho:** P.

### P3. Job de migração separado
- **Deploy:** `deploy.yml` executa as migrations antes de subir as réplicas; a aplicação só valida a versão no startup.
- **Aceite:** falha de migration não derruba as réplicas em produção.
- **Tamanho:** P.

---

## Fase 2 — Segurança do paciente

### C1. Verificação de interações e alergias na prescrição
- **Backend (módulo `medication`):** serviço de verificação que cruza o medicamento solicitado com alergias ativas e medicamentos ativos do paciente; resposta com lista de alertas (severidade, motivo, fonte); regra de bloqueio ou aviso configurável; registro de sobreposição ("prescrevi mesmo assim") com justificativa e auditoria.
- **Dados:** decidir a fonte de interações (base própria, RxNorm/ANVISA ou serviço de terceiros).
- **Frontend:** painel de alertas no formulário de prescrição, com `StatusBadge` crítico/atenção e confirmação explícita para sobrepor.
- **Perguntas em aberto:** o alerta bloqueia ou apenas avisa? Qual base de interações? Quem pode sobrepor? Cobre alergia cruzada por classe?
- **Aceite:** prescrever medicamento ao qual o paciente é alérgico exige justificativa registrada; testes unitários cobrem os cenários de interação.
- **Tamanho:** G.

### C2. Timeline clínica no prontuário
- **Backend:** endpoint de leitura que mescla Encounter, Observation, Condition, MedicationRequest, DiagnosticReport, ImagingStudy e AllergyIntolerance, ordenado por tempo registrado (ADR 0001) e com sucesso parcial explícito (ADR 0002); paginação por cursor; entradas em `policy.go`.
- **Frontend:** nova aba com filtros por tipo e período, agrupamento por dia, estados completos e `StatusBadge`.
- **Perguntas em aberto:** janela padrão? Incluir agendamentos? (O `CONTEXT.md` diz que não.)
- **Aceite:** ordem cronológica correta, falha de um recurso não esconde os demais.
- **Tamanho:** M.

---

## Fase 3 — Operação e relacionamento

### C3. Lembretes de consulta, remarcação e lista de espera
- **Backend (`schedule`, `notifications`):** agendador de lembretes (24 h e 2 h antes) via outbox; canais e-mail, SMS e WhatsApp atrás de uma interface; preferências do paciente; lista de espera com preenchimento automático de vagas canceladas; ligação agendamento → encontro com um clique.
- **Frontend:** preferências de contato no portal; fila de espera na agenda.
- **Perguntas em aberto:** quais provedores de envio? Custos e consentimento de contato? Janela de silêncio?
- **Aceite:** lembrete enviado uma única vez por consulta; cancelar remove lembretes pendentes.
- **Tamanho:** G. **Depende de:** P1.

### C4. Assinatura e fechamento de laudo
- **Backend (`diagnostic_report`):** estado "assinado" imutável, assinatura digital (ICP-Brasil ou assinatura simples), carimbo de tempo, hash do conteúdo, retificação por nova versão vinculada.
- **Frontend:** fluxo de assinar, selo de laudo assinado, verificação de integridade.
- **Perguntas em aberto:** exigência legal de ICP-Brasil? Quais perfis assinam?
- **Aceite:** laudo assinado não pode ser alterado; qualquer correção gera nova versão.
- **Tamanho:** G.

### C8. Direitos do titular (LGPD) e consentimento
- **Backend:** registro de consentimentos com base legal; exportação dos dados do paciente (FHIR bundle + dados locais); anonimização sob pedido; trilha de auditoria.
- **Frontend:** área no portal para consentimentos, exportação e solicitações.
- **Perguntas em aberto:** prazos legais e responsável (DPO)? Retenção mínima de prontuário versus direito ao esquecimento?
- **Aceite:** exportação completa em formato aberto; pedidos rastreáveis até a conclusão.
- **Tamanho:** G.

---

## Fase 4 — Inteligência e dados em tempo real

### C5. Resumo do paciente por IA (rascunho)
- **Backend (reaproveitar o padrão do `exam_analyzer`):** geração assíncrona via Vertex AI, com consentimento, anonimização e registro de auditoria; resposta sempre marcada como rascunho com fontes citadas.
- **Frontend:** cartão "Resumo (rascunho)" com aviso, citações clicáveis e botão de descartar.
- **Perguntas em aberto:** quais dados entram no contexto? Retenção do texto gerado? Avaliação clínica de segurança?
- **Aceite:** nenhum resumo é exibido sem o aviso de rascunho; toda geração fica auditada.
- **Tamanho:** M. **Depende de:** C2 (fonte de dados consolidada).

### C6. Telemetria com dados reais
- **Backend:** ingestão de leituras (MQTT ou HL7 v2 via gateway), tabela de séries temporais, limiares por paciente, alertas via outbox; manter o simulador apenas em desenvolvimento.
- **Frontend:** histórico além das últimas 30 amostras, limiares editáveis e ack de alarmes.
- **Perguntas em aberto:** quais dispositivos e protocolos? Volume por segundo? Retenção das séries?
- **Aceite:** leituras reais substituem o simulador em produção; alarme crítico gera notificação auditável.
- **Tamanho:** G. **Depende de:** P1.

---

## Fase 5 — Autonomia do paciente

### C7. Portal com escrita
- **Backend:** endpoints para o paciente agendar, remarcar e cancelar (regras de antecedência), baixar laudos e exames em PDF, atualizar contato; sempre limitados ao próprio vínculo (`patient_user_links`).
- **Frontend:** fluxos mobile-first no portal, confirmação por e-mail, estados completos.
- **Perguntas em aberto:** limites de antecedência para cancelamento? Quais laudos ficam liberados ao paciente e quando?
- **Aceite:** paciente nunca acessa dados de outro; todas as ações geram audit log.
- **Tamanho:** G. **Depende de:** S2 (e S1 opcional).

---

## Fase 6 — Maturidade técnica

| ID | Item | Detalhe | Tamanho |
|----|------|---------|---------|
| Q1 | CI backend reforçado | `go-version-file`, `golangci-lint`, `govulncheck`, `go test -race`, verificação de proto/swagger gerados, testes de integração | P |
| Q2 | Tipos do frontend gerados | `openapi-typescript` a partir do swagger e testes de contrato | M |
| Q3 | Acessibilidade | auditoria com axe/Playwright nas telas principais e correções | M |
| Q4 | Catálogo de componentes | evoluir `/design-system` (Storybook opcional) e cobrir novas variantes | P |
| P4 | Observabilidade de negócio | métricas (agendamentos, laudos, falhas FHIR), dashboards e alertas | M |
| P5 | Cache de leituras do FHIR | cache com invalidação por escrita para dados quentes | M |

---

## Riscos transversais

- **Dados clínicos e privacidade:** toda feature nova passa por revisão de LGPD/HIPAA, auditoria e controle de acesso antes de ir a produção.
- **Dependências externas:** provedores de mensagem, base de interações e ICP-Brasil exigem contratos e prazos; iniciar cedo.
- **Migrations:** apenas módulos operacionais geram migrations; módulos FHIR guardam dados na Healthcare API.
- **Compatibilidade:** mudanças de sessão (S2) invalidam tokens antigos; planejar comunicação e janela de deploy.

## Definição de pronto (para toda feature)

1. Perguntas em aberto respondidas e registradas.
2. Proto e stubs gerados, permissões em `policy.go`, `Register` e testes de serviço com `MockRepository`.
3. Frontend com estados de carregamento, vazio e erro, textos em 3 idiomas, componentes shadcn e `StatusBadge` para status.
4. `go test`, `tsc -b`, lint, testes unitários e paridade de idiomas verdes em CI.
5. Entrada no `CHANGELOG.md` e atualização do `CONTEXT.md` quando surgir um termo novo de domínio.
