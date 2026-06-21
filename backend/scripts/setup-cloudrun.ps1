# =============================================================================
# Healthcare API - Google Cloud Run Setup Script (PowerShell)
# Execute uma vez para provisionar toda a infraestrutura GCP necessária.
# Pré-requisito: gcloud CLI instalado e autenticado (gcloud auth login)
# =============================================================================

param(
    [Parameter(Mandatory=$true)]
    [string]$ProjectId
)

$Region            = "us-central1"
$ServiceName       = "healthcare-api"
$RegistryRepo      = "healthcare"
$SqlInstanceName   = "healthcare-db"
$SqlDbName         = "healthcare_db"
$SqlUser           = "healthcare_user"
$RedisInstanceName = "healthcare-redis"
$SaName            = "healthcare-cloudrun"

function Write-Step($step, $msg) {
    Write-Host ""
    Write-Host "[$step] $msg" -ForegroundColor Cyan
}

function Invoke-Gcloud($args) {
    $result = & gcloud @args 2>&1
    if ($LASTEXITCODE -ne 0 -and $result -notmatch "already exists") {
        Write-Host "   AVISO: $result" -ForegroundColor Yellow
    }
    return $result
}

Write-Host ""
Write-Host "=========================================================================" -ForegroundColor Green
Write-Host "  Healthcare API - Cloud Run Setup" -ForegroundColor Green
Write-Host "  Projeto: $ProjectId" -ForegroundColor Green
Write-Host "=========================================================================" -ForegroundColor Green

gcloud config set project $ProjectId --quiet

Write-Step "1/7" "Habilitando APIs necessarias..."
gcloud services enable `
    run.googleapis.com `
    sqladmin.googleapis.com `
    redis.googleapis.com `
    artifactregistry.googleapis.com `
    cloudresourcemanager.googleapis.com `
    --quiet

Write-Step "2/7" "Criando Artifact Registry repository..."
$registryResult = gcloud artifacts repositories create $RegistryRepo `
    --repository-format=docker `
    --location=$Region `
    --description="Healthcare API Docker images" `
    --quiet 2>&1
if ($registryResult -match "already exists") {
    Write-Host "   -> Repository ja existe, continuando." -ForegroundColor Yellow
} else {
    Write-Host "   -> Repository criado com sucesso." -ForegroundColor Green
}

Write-Step "3/7" "Criando instancia Cloud SQL PostgreSQL db-f1-micro..."
Write-Host "   (pode levar ~5 minutos)" -ForegroundColor Gray
$sqlResult = gcloud sql instances create $SqlInstanceName `
    --database-version=POSTGRES_16 `
    --tier=db-f1-micro `
    --region=$Region `
    --storage-type=SSD `
    --storage-size=10GB `
    --no-storage-auto-increase `
    --quiet 2>&1
if ($sqlResult -match "already exists") {
    Write-Host "   -> Instancia Cloud SQL ja existe, continuando." -ForegroundColor Yellow
} else {
    Write-Host "   -> Instancia Cloud SQL criada." -ForegroundColor Green
}

Write-Host "   Criando banco de dados e usuario..."
$chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
$SqlPassword = -join ((1..24) | ForEach-Object { $chars[(Get-Random -Maximum $chars.Length)] })

gcloud sql databases create $SqlDbName `
    --instance=$SqlInstanceName `
    --quiet 2>&1 | Out-Null

gcloud sql users create $SqlUser `
    --instance=$SqlInstanceName `
    --password=$SqlPassword `
    --quiet 2>&1 | Out-Null

$SqlConnectionName = gcloud sql instances describe $SqlInstanceName `
    --format="value(connectionName)" 2>&1
$DbUrl = "postgres://${SqlUser}:${SqlPassword}@/cloudsql/${SqlConnectionName}/${SqlDbName}"

Write-Step "4/7" "Criando instancia Memorystore Redis 1GB Basic..."
Write-Host "   (pode levar ~5 minutos)" -ForegroundColor Gray
$redisResult = gcloud redis instances create $RedisInstanceName `
    --size=1 `
    --region=$Region `
    --tier=basic `
    --quiet 2>&1
if ($redisResult -match "already exists") {
    Write-Host "   -> Instancia Redis ja existe, continuando." -ForegroundColor Yellow
} else {
    Write-Host "   -> Instancia Redis criada." -ForegroundColor Green
}

$RedisHost = gcloud redis instances describe $RedisInstanceName `
    --region=$Region `
    --format="value(host)" 2>&1
$RedisUrl = "${RedisHost}:6379"

Write-Step "5/7" "Criando Service Account para Cloud Run..."
gcloud iam service-accounts create $SaName `
    --display-name="Healthcare Cloud Run Service Account" `
    --quiet 2>&1 | Out-Null

$SaEmail = "${SaName}@${ProjectId}.iam.gserviceaccount.com"

Write-Host "   Concedendo permissoes minimas..."
$roles = @(
    "roles/cloudsql.client",
    "roles/redis.viewer",
    "roles/healthcare.fhirResourceEditor",
    "roles/aiplatform.user",
    "roles/storage.objectAdmin",
    "roles/run.invoker"
)
foreach ($role in $roles) {
    gcloud projects add-iam-policy-binding $ProjectId `
        --member="serviceAccount:${SaEmail}" `
        --role=$role `
        --quiet 2>&1 | Out-Null
    Write-Host "   -> $role" -ForegroundColor Gray
}

Write-Step "6/7" "Gerando chave JSON da Service Account..."
$SaKeyFile = "sa-key-${SaName}.json"
gcloud iam service-accounts keys create $SaKeyFile `
    --iam-account=$SaEmail `
    --quiet

Write-Step "7/7" "Configurando Cloud SQL no servico Cloud Run (proximo deploy)..."
Write-Host "   -> Sera aplicado automaticamente no primeiro deploy via GitHub Actions." -ForegroundColor Gray

$JwtHint = -join ((1..32) | ForEach-Object { $chars[(Get-Random -Maximum $chars.Length)] })

Write-Host ""
Write-Host "=========================================================================" -ForegroundColor Green
Write-Host "  Setup concluido!" -ForegroundColor Green
Write-Host "=========================================================================" -ForegroundColor Green
Write-Host ""
Write-Host "ADICIONE ESTES SECRETS NO GITHUB:" -ForegroundColor Yellow
Write-Host "(Settings -> Secrets and variables -> Actions -> New repository secret)" -ForegroundColor Gray
Write-Host ""

$secrets = [ordered]@{
    "GCP_SA_KEY"           = "Conteudo do arquivo: $SaKeyFile"
    "GCP_PROJECT_ID"       = $ProjectId
    "GCP_LOCATION_ID"      = $Region
    "GCP_DATASET_ID"       = "<seu-dataset-fhir>"
    "GCP_FHIR_STORE_ID"    = "<seu-fhir-store-id>"
    "GCP_VERTEX_MODEL"     = "gemini-2.0-flash-001"
    "GCS_BUCKET_NAME"      = "<seu-bucket-gcs>"
    "CLOUD_RUN_DB_URL"     = $DbUrl
    "CLOUD_RUN_REDIS_URL"  = $RedisUrl
    "CLOUD_RUN_JWT_SECRET" = $JwtHint
    "SENTRY_DSN"           = "(opcional)"
}

foreach ($key in $secrets.Keys) {
    Write-Host ("   {0,-25} -> {1}" -f $key, $secrets[$key]) -ForegroundColor White
}

Write-Host ""
Write-Host "AVISO: O arquivo '$SaKeyFile' contem credenciais sensiveis." -ForegroundColor Red
Write-Host "Delete-o apos copiar o conteudo para o GitHub Secret GCP_SA_KEY." -ForegroundColor Red
Write-Host ""
Write-Host "Apos adicionar os secrets, faca push para 'main' e o deploy" -ForegroundColor Green
Write-Host "acontecera automaticamente via GitHub Actions." -ForegroundColor Green
Write-Host "=========================================================================" -ForegroundColor Green
