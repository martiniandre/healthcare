# Guia Passo a Passo: Configuração de Ambiente (Healthcare Ecosystem) 🏥🧪

Este documento fornece um passo a passo prático para configurar todas as pendências listadas no projeto, permitindo que a aplicação execute com integrações reais da Google Cloud Platform (GCP), segurança adequada e comunicação de rede flexível.

---

## 🛠️ Passo 1: Configuração da Infraestrutura na GCP (Console Google Cloud)

Os dados clínicos e imagens PACS do projeto residem na Google Cloud. Siga estas etapas para criar a estrutura necessária:

### 1.1 Criar o Projeto no GCP
1. Acesse o [Console do Google Cloud](https://console.cloud.google.com/).
2. No menu superior esquerdo, clique no seletor de projetos e selecione **Novo Projeto**.
3. Escolha um nome descritivo (ex: `Healthcare Portal`) e anote o **ID do Projeto** gerado automaticamente (ex: `healthcare-portal-45123`).
4. Atualize a variável no seu arquivo `backend/.env`:
   ```env
   GCP_PROJECT_ID=healthcare-portal-45123
   ```

### 1.2 Ativar as APIs Necessárias
No console do GCP, acesse o menu de pesquisa no topo e ative as seguintes APIs clicando em **Ativar**:
*   **Cloud Healthcare API** (para gerenciamento de dados FHIR e DICOM).
*   **Vertex AI API** (para análise assistida por IA no módulo `exam_analyzer`).
*   **Google Cloud Storage JSON API** (para armazenar os arquivos de imagem DICOM).

### 1.3 Criar o Dataset e o FHIR Store
1. No menu lateral do console GCP, navegue até **Healthcare** (ou pesquise por "Healthcare" na barra de busca).
2. Clique em **Criar Dataset**.
3. Configure os detalhes:
   * **ID do Dataset:** `healthcare-dataset` (ou o nome desejado).
   * **Região:** Escolha `us-central1` (Região Padrão configurada na aplicação em `GCP_LOCATION_ID`).
4. Clique em **Criar**.
5. Dentro do dataset recém-criado, clique em **Criar Armazenamento de Dados** (Data Store).
6. Configure os detalhes do Data Store:
   * **Tipo de dados:** Selecione **FHIR**.
   * **ID do Armazenamento de Dados:** `fhir-store` (ou o nome desejado).
   * **Versão do FHIR:** Escolha **R4** (versão de prontuários médicos compatível com o backend).
7. Clique em **Criar**.
8. Preencha as variáveis correspondentes no seu `backend/.env`:
   ```env
   GCP_DATASET_ID=healthcare-dataset
   GCP_FHIR_STORE_ID=fhir-store
   ```

### 1.4 Criar o Bucket no Google Cloud Storage (PACS)
1. No console GCP, navegue até **Cloud Storage** -> **Buckets**.
2. Clique em **Criar**.
3. Configure os detalhes do bucket:
   * **Nome do Bucket:** Escolha um nome globalmente único (ex: `healthcare-pacs-files-prod`).
   * **Região:** Escolha `us-central1`.
   * **Classe de Armazenamento:** Selecione **Standard**.
   * **Controle de Acesso:** Recomendado manter **Uniforme**.
4. Clique em **Criar**.
5. Preencha a variável no seu `backend/.env`:
   ```env
   GCS_BUCKET_NAME=healthcare-pacs-files-prod
   ```

---

## 🔑 Passo 2: Configuração de Autenticação do Google Cloud (Máquina Local)

Para rodar a aplicação localmente no seu computador e permitir que o código interaja com a nuvem de forma transparente, é necessário autenticar sua máquina.

### 2.1 Instalar o gcloud CLI (Windows)
1. Baixe o instalador do [Google Cloud SDK para Windows](https://dl.google.com/dl/cloudsdk/channels/rapid/GoogleCloudSDKInstaller.exe).
2. Siga o assistente de instalação padrão na sua máquina.
3. Certifique-se de marcar a opção para abrir o terminal de comandos do gcloud no final da instalação.

### 2.2 Autenticar o seu Terminal
1. Abra uma nova janela do PowerShell ou Terminal do Windows.
2. Execute o comando de inicialização rápida:
   ```powershell
   gcloud init
   ```
3. O terminal solicitará login. Faça login utilizando a sua conta Google associada ao console da GCP onde você criou o projeto.
4. Selecione o projeto criado na etapa 1 quando solicitado.

### 2.3 Gerar Credenciais Padrão da Aplicação (ADC)
Este é o passo crucial. Ele salva um arquivo JSON com permissões na sua máquina local, o qual é lido de forma automática pelas bibliotecas Go.
1. No seu terminal de comandos, execute:
   ```powershell
   gcloud auth application-default login
   ```
2. Uma janela do seu navegador se abrirá. Permita o acesso clicando em **Confirmar**.
3. O terminal exibirá uma mensagem confirmando que as credenciais foram gravadas no caminho `C:\Users\<Usuario>\AppData\Roaming\gcloud\application_default_credentials.json`.
4. **Pronto!** O backend Go agora tem acesso direto à GCP sem precisar de nenhuma chave JSON hardcoded no código.

---

## 💻 Passo 3: Configurando Variáveis de Ambiente no Frontend (Vite)

Para evitar que a URL do backend (`http://localhost:8080/api`) fique fixa no código, faremos com que ela leia de um arquivo `.env` do Vite.

### 3.1 Criar o Arquivo `.env` no Frontend
1. Na raiz da pasta `frontend/`, crie um novo arquivo chamado `.env`:
   * Caminho completo: `frontend/.env`
2. Adicione a seguinte linha dentro dele:
   ```env
   VITE_API_BASE_URL=http://localhost:8080/api
   ```

### 3.2 Atualizar o arquivo de configuração de API do Frontend
Substitua o hardcode em `frontend/src/shared/services/api.ts` para ler a nova variável.
1. Abra o arquivo [api.ts](file:///c:/Users/andre/Desktop/Projetos/healthcare/frontend/src/shared/services/api.ts).
2. Substitua o conteúdo pelo seguinte bloco reativo:
   ```typescript
   import axios from "axios"

   export const api = axios.create({
     baseURL: import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api",
     withCredentials: true,
   })
   ```

---

## 🐳 Passo 4: Solução de Conflitos de Portas no Docker (Postgres e Redis)

Se ao rodar `docker compose up -d` você receber erros de portas já utilizadas:

### 4.1 Identificar qual processo está ocupando as portas no Windows
1. Abra o PowerShell como Administrador.
2. Execute o comando para descobrir o processo usando a porta do Postgres (`5432`):
   ```powershell
   Get-Process -Id (Get-NetTCPConnection -LocalPort 5432).OwningProcess
   ```
3. Execute o mesmo para a porta do Redis (`6379`):
   ```powershell
   Get-Process -Id (Get-NetTCPConnection -LocalPort 6379).OwningProcess
   ```

### 4.2 Desativar Serviços Nativos do Windows
Se os processos identificados pertencerem aos serviços nativos de banco de dados instalados na sua máquina:
1. Pressione `Win + R`, digite `services.msc` e aperte Enter.
2. Procure por **PostgreSQL** ou **Redis** na lista de serviços.
3. Dê dois cliques no serviço, clique no botão **Parar** e altere o **Tipo de inicialização** para **Manual** para evitar conflito nos próximos boots.

---

## 🔒 Passo 5: Gerando um JWT Secret Forte de Produção

Nunca utilize chaves fracas em produção. Siga estes passos para gerar uma nova chave de 32 bytes segura:

### 5.1 Usando PowerShell no Windows
1. Abra o terminal do PowerShell.
2. Rode o seguinte script nativo de uma linha para obter uma string aleatória:
   ```powershell
   -join ((48..57) + (65..90) + (97..122) | Get-Random -Count 64 | % {[char]$_})
   ```
3. Copie a string gerada no terminal.
4. No arquivo `backend/.env`, substitua a linha do `JWT_SECRET`:
   ```env
   JWT_SECRET=SuaStringAltamenteAleatoriaGeradaDe64Caracteres
   ```
