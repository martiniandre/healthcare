# Diretrizes de Engenharia e Padrões de Front-End

Este manual define as diretrizes oficiais de arquitetura, qualidade de código, padrões de estado assíncrono e decisões de design visual para a aplicação front-end baseada em React (Vite + TypeScript) do ecossistema de saúde.

---

## 🏛️ 1. Arquitetura Modular (`Domain-Driven Scaffolding`)

O front-end é estruturado seguindo conceitos de encapsulamento por domínio, organizando as pastas em módulos independentes localizados em `src/modules/`. Recursos globais de utilidades, componentes reutilizáveis e layouts compartilhados residem na pasta `src/shared/`.

### Estrutura de Diretórios de Domínio
Cada módulo de negócio (ex: `patients`, `imaging`, `telemetry`) deve conter sua estrutura autônoma:
*   `types.ts`: Definições estritas de DTOs e tipos de dados do domínio.
*   `api.ts`: Isolamento do cliente HTTP (Axios) com tipagem de entrada e saída.
*   `queries.ts`: Hooks do TanStack Query (consultas e mutações reativas).
*   `components/`: Subcomponentes especializados exclusivos do domínio.
*   `modals/`: Janelas modais e diálogos reativos do domínio.
*   `{Domain}.tsx`: Componente de página mestre/orquestradora do módulo.

---

## 📡 2. Gerenciamento de Estado e Sincronização Assíncrona

*   **TanStack Query (React Query)**: Toda a sincronização de dados com o backend Go (REST Gateway) é delegada ao TanStack Query.
*   **Chaves de Cache Dinâmicas**: Consultas que suportam filtros ou buscas (ex: busca de pacientes) devem incluir os parâmetros ativos em suas chaves de query (ex: `queryKey: [...patientQueryKeys.lists(), searchQueryValue]`) para forçar o carregamento transparente de estados reativos.
*   **Mutações e Invalidação**: Operações de gravação (`useMutation`) devem invocar `queryClient.invalidateQueries` no callback `onSuccess` utilizando as chaves de query correspondentes ao recurso atualizado para consistência de dados em tempo real.

---

## 📝 3. Validação de Formulários com Zod

*   **Zod Schemas**: Todos os formulários devem ser validados de forma estrita no lado do cliente utilizando schemas declarativos do Zod (ex: `newAllergySchema` em `patient_schemas.ts`).
*   **React Hook Form**: Integração nativa com React Hook Form utilizando o resolver do Zod (`zodResolver`).
*   **Segurança de Tipos**: Tipos de dados de formulário devem ser inferidos dinamicamente dos schemas declarados usando `z.infer<typeof schema>`.

---

## 🎨 4. Identidade Visual Premium e Acessibilidade

O design da interface de usuário deve se destacar à primeira vista, aplicando princípios de design de alta fidelidade e acessibilidade web:
*   **Acessibilidade Semântica (WCAG)**: Utilização de tags HTML5 estruturadas (`<main>`, `<section>`, `<aside>`) e atributos ARIA no Canvas/Viewports (`role="img"`, `aria-label`).
*   **Cores Harmoniosas**: Utilização de paletas modernas e elegantes baseadas em HSL (como tons suaves de Slate, Emerald, Rose e Cyan) em detrimento de cores puras primárias.
*   **Cursores Contextuais**: Alterar dinamicamente a propriedade CSS do cursor de acordo com a ferramenta ativa (ex: `cursor-zoom-in` para lupas, `cursor-crosshair` para réguas e miras, `cursor-ew-resize` para gradientes de brilho e contraste).
*   **Feedback de Estado**: Transições suaves (`duration-300`), animações de pulso suave (`animate-pulse`) e estados de foco/hover responsivos.

---

## 🚨 5. Regras Estritas de Qualidade de Código (AGENTS.MD)

Como determinado na especificação técnica central do projeto, as seguintes regras são de observância obrigatória e sem qualquer exceção:

### 🚫 A. Zero Comentários
É terminantemente proibido adicionar qualquer comentário explicativo no código fonte final. O código gerado deve ser autoexplicativo por si só.
*   *Incorreto*: `// Atualiza o ID do leito selecionado`
*   *Incorreto*: `/* Callback de sucesso */`
*   *Correto*: Ausência completa de linhas de comentário no arquivo.

### 📝 B. Variáveis Altamente Descritivas
Proibido utilizar variáveis ou abreviações de uma única letra (como `x`, `y`, `d`, `e`), mesmo em loops, iteradores de formulário ou manipuladores de eventos.
*   *Incorreto*: `const d = new Date()`
*   *Incorreto*: `patients.map((p) => ...)`
*   *Incorreto*: `onChange={(e) => ...)`
*   *Correto*: `const currentDate = new Date()`
*   *Correto*: `patients.map((patientItem) => ...)`
*   *Correto*: `onChange={(event) => ...)`
