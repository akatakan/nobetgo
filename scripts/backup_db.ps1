# NobetGo Database Backup Script
# Requires pg_dump to be in PATH

$DB_NAME = "nobetgo"
$DB_USER = "postgres"
$DB_PASS = "postgres"
$DB_HOST = "localhost"
$DB_PORT = "5432"

$TIMESTAMP = Get-Date -Format "yyyyMMdd_HHmmss"
$BACKUP_DIR = "./backups"
$FILE_NAME = "nobetgo_manual_backup_$TIMESTAMP.sql"
$FILE_PATH = Join-Path $BACKUP_DIR $FILE_NAME

if (-not (Test-Path $BACKUP_DIR)) {
    New-Item -ItemType Directory -Path $BACKUP_DIR
}

$env:PGPASSWORD = $DB_PASS
Write-Host "Creating backup: $FILE_PATH ..." -ForegroundColor Cyan

& pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $FILE_PATH --no-owner --no-privileges

if ($LASTEXITCODE -eq 0) {
    Write-Host "Backup successful!" -ForegroundColor Green
} else {
    Write-Host "Backup failed!" -ForegroundColor Red
}

$env:PGPASSWORD = $null
