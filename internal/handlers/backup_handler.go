package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/akatakan/nobetgo/internal/services"
	"github.com/akatakan/nobetgo/util"
	"github.com/gin-gonic/gin"
)

type BackupHandler struct {
	service *services.BackupService
}

func NewBackupHandler(service *services.BackupService) *BackupHandler {
	return &BackupHandler{service: service}
}

// ExportBackup handles the request to trigger and download a DB backup.
func (h *BackupHandler) ExportBackup(c *gin.Context) {
	filePath, err := h.service.CreateBackup()
	if err != nil {
		util.InternalError(c, "Yedek oluşturulamadı", err)
		return
	}

	// Ensure cleanup of the temp file after download
	defer os.Remove(filePath)

	fileName := fmt.Sprintf("nobetgo_backup_%s.sql", time.Now().Format("2006-01-02"))

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}
