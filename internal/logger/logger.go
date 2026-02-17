package logger

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LogConfig holds logging configuration
type LogConfig struct {
	Level string `mapstructure:"level"` // debug, info, warn, error
}

// InitLogger initializes the global slog logger based on config
func InitLogger(cfg LogConfig) {
	level := parseLevel(cfg.Level)

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Info("Logger initialized", "level", cfg.Level)
}

func parseLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// GinLoggerMiddleware returns a Gin middleware that logs HTTP requests using slog
func GinLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		attrs := []any{
			"status", status,
			"method", method,
			"path", path,
			"latency_ms", latency.Milliseconds(),
			"ip", clientIP,
		}

		if query != "" {
			attrs = append(attrs, "query", query)
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
		}

		switch {
		case status >= 500:
			slog.Error("HTTP Request", attrs...)
		case status >= 400:
			slog.Warn("HTTP Request", attrs...)
		default:
			slog.Info("HTTP Request", attrs...)
		}
	}
}
