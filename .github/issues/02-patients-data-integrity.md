---
title: "Corrigir integridade de dados do módulo Patients"
labels: ["bug", "backend"]
---

## What to build

Corrigir 4 bugs de integridade de dados no módulo patients.

### Problemas

1. **UUID duplicado**: `service.go:44` gera `patient.ID = uuid.New()`, depois `repository.go:54` sobrescreve com outro UUID. Duas atribuições para o mesmo campo.

2. **UUID novo em cada leitura**: `parsePatientFromFHIR` gera `ID: uuid.New()` em toda fetch. O mesmo paciente FHIR ganha um UUID local diferente cada vez que é lido — quebra referência com dados locais que referenciam Patient ID.

3. **CreatedAt/UpdatedAt falsos**: Ambos setados para `time.Now()` em cada leitura, nunca refletem o `meta.lastUpdated` do FHIR ou um valor persistido localmente.

4. **GetPatientByDocumentID formato errado**: FHIR busca `identifier` como `system|value`, mas o código passa só o `value`. A busca pode não encontrar registros.

### Acceptance criteria

- [ ] Apenas UM local define `patient.ID` (service ou repository, não ambos)
- [ ] `parsePatientFromFHIR` preserva o Patient ID entre leituras
- [ ] `CreatedAt`/`UpdatedAt` vêm do FHIR `meta.lastUpdated` ou são omitidos
- [ ] `GetPatientByDocumentID` usa formato `system|value` do FHIR
- [ ] Testes atualizados para validar IDs consistentes

### Blocked by

None — can start immediately
