package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/akatakan/nobetgo/config"
)

type BackupService struct {
	cfg config.DatabaseConfig
}

func NewBackupService(cfg config.DatabaseConfig) *BackupService {
	return &BackupService{cfg: cfg}
}

// CreateBackup generates a .sql backup file using pg_dump.
// It returns the path to the temporary backup file.
func (s *BackupService) CreateBackup() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	tempDir := os.TempDir()
	fileName := fmt.Sprintf("nobetgo_backup_%s.sql", timestamp)
	filePath := filepath.Join(tempDir, fileName)

	// Prepare pg_dump command
	cmd := exec.Command("pg_dump",
		"-h", s.cfg.Host,
		"-p", s.cfg.Port,
		"-U", s.cfg.User,
		"-d", s.cfg.DBName,
		"-f", filePath,
		"--no-owner",
		"--no-privileges",
	)
	// Pass PGPASSWORD via subprocess environment (thread-safe)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+s.cfg.Password)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pg_dump failed: %v, output: %s", err, string(output))
	}

	return filePath, nil
}
