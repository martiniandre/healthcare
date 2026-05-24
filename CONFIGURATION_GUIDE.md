# Guia de Configurações Pendentes e Ambiente (Healthcare Ecosystem) 🏥🧪

Este documento consolida todas as informações, variáveis de ambiente, dependências de infraestrutura e etapas de autenticação que **ainda precisam ser configuradas** no projeto para o funcionamento completo e seguro em ambientes de Desenvolvimento (Local) e Produção.

---

## 🏛️ 1. Variáveis de Ambiente do Backend (`backend/.env`)

O arquivo `.env` do backend (`backend/.env`) possui definições cruciais de infraestrutura. Algumas variáveis já vêm populadas com valores padrão de desenvolvimento local, enquanto outras dependem diretamente da sua conta Google Cloud Platform (GCP).

### 🔴 Configurações do Google Cloud Platform (Obrigatório para Clínico & PACS)
Como os módulos de **Pacientes**, **Prontuário Clínico (Clinical)** e **PACS (Imaging)** operam 100% integrados à Google Cloud Healthcare API (padrão FHIR) e ao Google Cloud Storage (arquivos DICOM), as seguintes configurações na GCP precisam ser definidas:

*   **`GCP_PROJECT_ID`**: Atualmente está vazio no arquivo `.env`. Deve ser preenchido com o ID exato do seu projeto na GCP (ex: `GCP_PROJECT_ID=healthcare-portal-41221`).
*   **`GCP_LOCATION_ID`**: Região onde a infraestrutura GCP está implantada. O padrão é `us-central1` (usado na criação dos datasets e stores).
*   **`GCP_DATASET_ID`**: Nome do dataset clínico criado no console do GCP Healthcare (ex: `GCP_DATASET_ID=healthcare-dataset`).
*   **`GCP_FHIR_STORE_ID`**: Identificador do FHIR Store criado sob o dataset (ex: `GCP_FHIR_STORE_ID=fhir-store`).
*   **`GCP_DICOM_STORE_ID`**: Nome do Dicom Store (fallback padrão em código: `default-dicom`).
*   **`GCS_BUCKET_NAME`**: Nome do bucket do Google Cloud Storage utilizado para hospedar as imagens PACS do exame (fallback padrão em código: `default-bucket`).

> [!CAUTION]
> **Bloqueio no Boot:** Se o servidor backend for iniciado com `GCP_PROJECT_ID`, `GCP_DATASET_ID` ou `GCP_FHIR_STORE_ID` vazios, o validador de bootstrap de configuração irá falhar, impedindo que a aplicação suba.

### 🟡 Segurança & Autenticação
*   **`JWT_SECRET`**: Atualmente configurado com uma chave de teste de desenvolvimento. Para produção ou ambientes públicos de homologação, **gere uma chave criptografada aleatória de alta entropia (mínimo de 32 ou 64 caracteres)**.
*   **`SENTRY_DSN`**: Possui uma URL de teste apontando para o Sentry. Se deseja utilizar o monitoramento de erros em sua própria conta, configure o DSN do seu workspace do Sentry. Caso contrário, mantenha em branco para desabilitar o tracing local.

---

## 🔑 2. Autenticação Google Cloud (ADC - Application Default Credentials)

Os clientes Go da GCP (`google.DefaultClient` no Healthcare API e `storage.NewClient` para o GCS) utilizam o mecanismo padrão do Google para capturar credenciais.

### 💻 Em ambiente de Desenvolvimento Local:
Você precisará instalar o **Google Cloud SDK (gcloud CLI)** e autenticar sua máquina local para que a aplicação consiga interagir com as APIs da GCP de forma transparente:

1. Instale o [Google Cloud SDK](https://cloud.google.com/sdk/docs/install).
2. Execute o comando de login no seu terminal:
   ```bash
   gcloud auth application-default login
   ```
3. O comando abrirá uma janela do navegador para que você selecione a conta com acesso ao projeto configurado em `GCP_PROJECT_ID`.

### ☁️ Em ambiente de Produção / CI/CD:
Você deve criar uma **Service Account (Conta de Serviço)** dedicada na GCP, baixar a chave JSON correspondente e passá-la para o contêiner/servidor:

1. Crie uma Service Account com as seguintes permissões (IAM Roles):
   - **Administrador do Cloud Healthcare (ou roles específicas para Dataset e FHIR Store)**.
   - **Administrador de Objetos do Storage (Storage Object Admin)** (para leitura e escrita dos exames no bucket).
   - **Usuário do Vertex AI (Vertex AI User)** (para uso do módulo `exam_analyzer` assistido por IA).
2. Gere a chave privada JSON.
3. Configure a variável de ambiente do sistema para apontar para o caminho desse arquivo JSON:
   ```bash
   export GOOGLE_APPLICATION_CREDENTIALS="/caminho/para/sua/chave-privada.json"
   ```

---

## 💻 3. Configurações Pendentes no Frontend (Vite)

Atualmente, o endpoint de comunicação da interface React com a API do Backend REST Gateway está hardcoded no arquivo `frontend/src/shared/services/api.ts`:

```typescript
export const api = axios.create({
  baseURL: "http://localhost:8080/api",
  withCredentials: true,
})
```

### ⚡ Próximas Configurações Recomendadas no Frontend:
Para evitar o acoplamento de `localhost` e permitir a flexibilidade em servidores de homologação ou produção:

1. **Introduzir variáveis do Vite**: Criar um arquivo `.env` e `.env.production` no frontend.
2. **Definir variável base**:
   ```env
   VITE_API_BASE_URL=https://api.seu-dominio.com/api
   ```
3. **Refatorar o serviço de API**:
   ```typescript
   export const api = axios.create({
     baseURL: import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api",
     withCredentials: true,
   })
   ```

---

## 🐳 4. Porta de Conexão de Infraestrutura (Port Conflicts)

O arquivo `docker-compose.yml` da raiz mapeia as seguintes portas na máquina local:
*   **PostgreSQL**: `5432:5432`
*   **Redis**: `6379:6379`

> [!WARNING]
> Se você já tiver instâncias do PostgreSQL ou Redis instaladas nativamente rodando como serviço local em seu sistema operacional, o comando `docker compose up -d` irá falhar devido a conflitos de portas.
> **Solução:** Pare os serviços locais temporariamente (`pg_ctl` / `services.msc` no Windows) ou edite as portas esquerdas no `docker-compose.yml` e ajuste no `backend/.env` (`DB_URL` e `REDIS_URL`).

---

## 📋 Checklist de Prontidão de Configuração

- [ ] Instalar o Google Cloud SDK (gcloud CLI) na máquina de desenvolvimento.
- [ ] Executar `gcloud auth application-default login` para gerar as credenciais locais.
- [ ] Preencher `GCP_PROJECT_ID` no arquivo `backend/.env`.
- [ ] Criar o Dataset Healthcare na console do Google Cloud no local especificado.
- [ ] Criar o FHIR Store com o ID configurado no arquivo `.env`.
- [ ] Criar o bucket do Cloud Storage para PACS e atualizar a chave `GCS_BUCKET_NAME`.
- [ ] Gerar uma chave segura de produção para `JWT_SECRET` ao implantar em servidores.
- [ ] Refatorar `frontend/src/shared/services/api.ts` para ler dinamicamente o endpoint do backend via variáveis de ambiente do Vite.
