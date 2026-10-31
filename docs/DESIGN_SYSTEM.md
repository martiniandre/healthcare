# Design System

Catálogo vivo no app: rota `/design-system` (qualquer perfil de equipe: admin, médico, enfermagem e recepção). Mostra cores com contraste calculado em tempo de execução, tipografia, componentes com variantes/estados e diretrizes. Código em `frontend/src/modules/design_system`.

## Tokens

Definidos em `frontend/src/index.css` como variáveis `--color-*` (tema claro em `@theme`, escuro na classe `.dark`).

| Token | Uso |
| --- | --- |
| `background` / `foreground` | Fundo e texto da página |
| `card`, `popover` (+ `-foreground`) | Superfícies elevadas |
| `primary`, `secondary`, `accent`, `muted` (+ `-foreground`) | Ações e áreas de apoio |
| `destructive`, `danger` | Erro e ações destrutivas |
| `success`, `warning` | Estados positivo e de atenção |
| `border`, `input`, `ring` | Bordas, campos e foco |

Regras: usar sempre tokens semânticos (`bg-muted`, `text-foreground`), nunca `gray-*` ou cores hexadecimais soltas. Texto exige contraste mínimo 4,5:1 (AA); bordas e foco, 3:1.

Tipografia: Inter Variable. Use `tabular-nums` em sinais vitais, tabelas e contadores.

## Adicionar um componente shadcn

O `frontend/components.json` já aponta os aliases (`@` = `src`): `ui` em `@/shared/components/ui`, `utils` em `@/shared/utils/cn`.

```bash
cd frontend
npx shadcn@latest add <nome>
```

Depois de adicionar:

1. Confirme que o arquivo ficou em `src/shared/components/ui/` com nome em kebab-case minúsculo (`dropdown-menu.tsx`).
2. Troque textos fixos por i18n (`useTranslation`) e adicione as chaves nos 3 idiomas.
3. Garanta alvos de toque de 44px no mobile e estado `aria-invalid`/`disabled` quando houver.
4. Adicione um exemplo em `modules/design_system/components` e, se tiver lógica, um teste `*.test.tsx`.

## Regras de código e estilo

- Sem comentários no código; sem variáveis de uma letra.
- Todo texto visível passa pelo i18n (`pt-BR`, `en-US`, `es-ES`); a regra ESLint `no-hardcoded-user-text` bloqueia texto fixo. Valide com `node scripts/check-locale-parity.mjs`.
- Arquivos de componentes shadcn em kebab-case minúsculo; componentes de módulo em PascalCase.
- Status nunca apenas por cor: use `StatusBadge` (ícone + texto).
- Imports pelo alias `@/`.

## Verificação

```bash
cd frontend
npx tsc -b
npx eslint src
npx vitest run --maxWorkers=2
node scripts/check-locale-parity.mjs
```
