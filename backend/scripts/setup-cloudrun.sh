#!/usr/bin/env bash
set -euo pipefail

# =============================================================================
# Healthcare API - Google Cloud Run Setup Script
# Execute uma vez para provisionar toda a infraestrutura GCP necessária.
# =============================================================================

# --- Configuração (edite estas variáveis antes de rodar) ---
PROJECT_ID="${GCP_PROJECT_ID:-}"
REGION="us-central1"
SERVICE_NAME="healthcare-api"
REGISTRY_REPO="healthcare"
SQL_INSTANCE_NAME="healthcare-db"
SQL_DB_NAME="healthcare_db"
SQL_USER="healthcare_user"
REDIS_INSTANCE_NAME="healthcare-redis"
SA_NAME="healthcare-cloudrun"
# -----------------------------------------------------------

if [[ -z "$PROJECT_ID" ]]; then
  echo "❌  Defina a variável GCP_PROJECT_ID antes de rodar:"
  echo "    export GCP_PROJECT_ID=seu-projeto-id"
  exit 1
fi

echo "🚀  Iniciando setup para projeto: $PROJECT_ID"
gcloud config set project "$PROJECT_ID"

echo ""
echo "📦  [1/7] Habilitando APIs necessárias..."
gcloud services enable \
  run.googleapis.com \
  sqladmin.googleapis.com \
  redis.googleapis.com \
  artifactregistry.googleapis.com \
  cloudresourcemanager.googleapis.com \
  --quiet

echo ""
echo "🗄️   [2/7] Criando Artifact Registry repository..."
gcloud artifacts repositories create "$REGISTRY_REPO" \
  --repository-format=docker \
  --location="$REGION" \
  --description="Healthcare API Docker images" \
  --quiet 2>/dev/null || echo "   → Repository já existe, continuando."

echo ""
echo "🐘  [3/7] Criando instância Cloud SQL PostgreSQL (db-f1-micro)..."
echo "   (pode levar ~5 minutos)"
gcloud sql instances create "$SQL_INSTANCE_NAME" \
  --database-version=POSTGRES_16 \
  --tier=db-f1-micro \
  --region="$REGION" \
  --storage-type=SSD \
  --storage-size=10GB \
  --no-storage-auto-increase \
  --quiet 2>/dev/null || echo "   → Instância já existe, continuando."

echo ""
echo "   Criando banco de dados e usuário..."
SQL_PASSWORD=$(openssl rand -base64 24 | tr -d "=+/")

gcloud sql databases create "$SQL_DB_NAME" \
  --instance="$SQL_INSTANCE_NAME" \
  --quiet 2>/dev/null || echo "   → Banco já existe."

gcloud sql users create "$SQL_USER" \
  --instance="$SQL_INSTANCE_NAME" \
  --password="$SQL_PASSWORD" \
  --quiet 2>/dev/null || echo "   → Usuário já existe."

SQL_CONNECTION_NAME=$(gcloud sql instances describe "$SQL_INSTANCE_NAME" \
  --format="value(connectionName)")
DB_URL="postgres://${SQL_USER}:${SQL_PASSWORD}@/cloudsql/${SQL_CONNECTION_NAME}/${SQL_DB_NAME}"

echo ""
echo "⚡  [4/7] Criando instância Memorystore Redis (1GB Basic)..."
echo "   (pode levar ~5 minutos)"
gcloud redis instances create "$REDIS_INSTANCE_NAME" \
  --size=1 \
  --region="$REGION" \
  --tier=basic \
  --quiet 2>/dev/null || echo "   → Instância Redis já existe, continuando."

REDIS_HOST=$(gcloud redis instances describe "$REDIS_INSTANCE_NAME" \
  --region="$REGION" \
  --format="value(host)" 2>/dev/null || echo "PENDENTE")
REDIS_URL="${REDIS_HOST}:6379"

echo ""
echo "👤  [5/7] Criando Service Account para Cloud Run..."
gcloud iam service-accounts create "$SA_NAME" \
  --display-name="Healthcare Cloud Run Service Account" \
  --quiet 2>/dev/null || echo "   → Service Account já existe."

SA_EMAIL="${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

echo "   Concedendo permissões mínimas..."
for ROLE in \
  "roles/cloudsql.client" \
  "roles/redis.viewer" \
  "roles/healthcare.fhirResourceEditor" \
  "roles/aiplatform.user" \
  "roles/storage.objectAdmin"; do
  gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:${SA_EMAIL}" \
    --role="$ROLE" \
    --quiet > /dev/null
done

echo ""
echo "🔑  [6/7] Gerando chave JSON da Service Account..."
SA_KEY_FILE="sa-key-${SA_NAME}.json"
gcloud iam service-accounts keys create "$SA_KEY_FILE" \
  --iam-account="$SA_EMAIL" \
  --quiet

echo ""
echo "🔑  [7/7] Configurando Cloud Run para usar Cloud SQL via socket..."
gcloud run services update "$SERVICE_NAME" \
  --region="$REGION" \
  --add-cloudsql-instances="$SQL_CONNECTION_NAME" \
  --quiet 2>/dev/null || echo "   → Serviço não existe ainda, será configurado no primeiro deploy."

echo ""
echo "========================================================================="
echo "✅  Setup concluído!"
echo "========================================================================="
echo ""
echo "📋  ADICIONE ESTES SECRETS NO GITHUB (Settings → Secrets → Actions):"
echo ""
echo "   GCP_SA_KEY          → conteúdo do arquivo: $SA_KEY_FILE"
echo "   GCP_PROJECT_ID      → $PROJECT_ID"
echo "   GCP_LOCATION_ID     → $REGION"
echo "   GCP_DATASET_ID      → (seu dataset FHIR)"
echo "   GCP_FHIR_STORE_ID   → (seu FHIR store)"
echo "   GCP_VERTEX_MODEL    → gemini-2.0-flash-001"
echo "   GCS_BUCKET_NAME     → (seu bucket GCS)"
echo "   CLOUD_RUN_DB_URL    → $DB_URL"
echo "   CLOUD_RUN_REDIS_URL → $REDIS_URL"
echo "   CLOUD_RUN_JWT_SECRET → (gere com: openssl rand -base64 32)"
echo "   SENTRY_DSN          → (opcional)"
echo ""
echo "⚠️   IMPORTANTE: O arquivo '$SA_KEY_FILE' contém credenciais sensíveis."
echo "    Adicione-o ao .gitignore e delete após copiar o conteúdo."
echo ""
echo "🚀  Após adicionar os secrets, faça push para 'main' para iniciar o deploy."
echo "========================================================================="
