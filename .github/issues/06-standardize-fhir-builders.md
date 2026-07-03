---
title: "Padronizar builders de recursos FHIR para structs tipadas"
labels: ["refactor", "backend"]
---

## What to build

Unificar os 8 builders de recursos FHIR para usar structs tipadas (Style A) em vez de `map[string]interface{}` (Style B).

### Problema

Metade dos builders usa structs Go tipadas com JSON tags (Patient, Encounter, Observation, ImagingStudy), metade usa `map[string]interface{}` (Condition, DiagnosticReport, MedicationRequest, AllergyIntolerance). O estilo `map[string]interface{}` não tem segurança de tipo, não tem autocomplete, e é propenso a typos silenciosos em chaves.

### Escopo

1. Migrar Condition, DiagnosticReport, MedicationRequest, AllergyIntolerance para structs tipadas
2. Adicionar `NewImagingStudyResource()` builder function (atualmente construído manualmente em imaging/worker.go)
3. Corrigir `PatientResource`: `FullName` deve ser dividido em `Given` + `Family` (hoje vai tudo em `Family`)
4. Remover campos `json:"-"` mortos em `PatientResource` (`BloodType`, `Allergies`)
5. Adicionar validação de campos obrigatórios nos builders (ex: patient reference, code)

### Acceptance criteria

- [ ] Todos os 8 builders usam structs tipadas
- [ ] `ImagingStudyResource` tem builder function
- [ ] `PatientResource` usa `Name[0].Given` + `Name[0].Family` corretamente
- [ ] Campos mortos `json:"-"` removidos
- [ ] Build e testes passando

### Blocked by

None — can start immediately
