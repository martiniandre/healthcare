---
title: "Padronizar tratamento de erros em todos os módulos (apperrors)"
labels: ["refactor", "backend"]
---

## What to build

Unificar o tratamento de erros para usar `apperrors` consistentemente em TODOS os módulos do backend.

### Problemas

1. **Clinical**: `mapClinicalError` default case retorna `apperrors.ErrInternalServer.ToGRPC()` em vez de `apperrors.ToGRPCStatus(err)`. Isso engole a mensagem de erro original (network, parse, etc.).

2. **exam_analyzer**: Zero uso de `apperrors`. Retorna `fmt.Errorf` e `errors.New` crus diretamente.

3. **stats**: Zero uso de `apperrors`. HTTP handler usa `slog.Error` + `render.Error(http.StatusInternalServerError)` sem mapping de códigos.

4. **health**: Usa `status.Error(codes.Unimplemented)` — o único módulo que bypassa `apperrors` completamente.

### Acceptance criteria

- [ ] Clinical `mapClinicalError` fallback usa `apperrors.ToGRPCStatus(err)`
- [ ] exam_analyzer service e handler usam `apperrors` para todos os erros de domínio
- [ ] stats service e handler usam `apperrors` para todos os erros
- [ ] health module usa `apperrors` em vez de raw `status.Error`
- [ ] Testes atualizados para verificar erros corretos

### Blocked by

None — can start immediately
