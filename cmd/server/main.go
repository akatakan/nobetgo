package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akatakan/nobetgo/config"
	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/database"
	"github.com/akatakan/nobetgo/internal/handlers"
	"github.com/akatakan/nobetgo/internal/logger"
	"github.com/akatakan/nobetgo/internal/middleware"
	"github.com/akatakan/nobetgo/internal/repositories"
	"github.com/akatakan/nobetgo/internal/services"
	"github.com/akatakan/nobetgo/util"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		slog.Error("Cannot load config", "error", err)
		os.Exit(1)
	}

	// Initialize structured logger
	logger.InitLogger(logger.LogConfig{Level: cfg.Log.Level})

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		slog.Error("Cannot connect to database", "error", err)
		os.Exit(1)
	}

	// Auto Migration
	err = db.AutoMigrate(
		&core.Department{},
		&core.Title{},
		&core.Employee{},
		&core.ShiftType{},
		&core.RotationPlan{},
		&core.Schedule{},
		&core.TimeEntry{},
		&core.LeaveType{},
		&core.Leave{},
		&core.LeaveBalance{},
		&core.OvertimeRule{},
		&core.PublicHoliday{},
		&core.AuditLog{},
		&core.Notification{},
		&core.PasswordResetToken{},
	)
	if err != nil {
		slog.Error("Cannot migrate database", "error", err)
		os.Exit(1)
	}
	// Partial unique index: enforce uniqueness only on non-empty emails
	db.Exec("DROP INDEX IF EXISTS idx_employees_email")
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_employees_email_unique ON employees (email) WHERE email != '' AND deleted_at IS NULL")
	slog.Info("Database migration completed")

	// Seed initial admin if no employees exist
	var employeesCount int64
	db.Model(&core.Employee{}).Count(&employeesCount)
	if employeesCount == 0 {
		// 1. Ensure a default Department exists
		var dept core.Department
		if err := db.First(&dept).Error; err != nil {
			dept = core.Department{Name: "Yönetim", Floor: 0}
			db.Create(&dept)
		}

		// 2. Ensure a default Title exists
		var title core.Title
		if err := db.First(&title).Error; err != nil {
			title = core.Title{Name: "Yönetici"}
			db.Create(&title)
		}

		// 3. Create Admin
		hashedPassword, hashErr := util.HashPassword("admin")
		if hashErr != nil {
			slog.Error("Cannot hash admin password", "error", hashErr)
			os.Exit(1)
		}
		admin := core.Employee{
			FirstName:    "Sistem",
			LastName:     "Yöneticisi",
			Username:     "admin",
			Email:        "admin@nobetgo.com",
			PasswordHash: hashedPassword,
			Role:         "admin",
			IsActive:     true,
			DepartmentID: dept.ID,
			TitleID:      title.ID,
		}
		db.Create(&admin)
		slog.Info("Initial admin user created: admin / admin")

		// 4. Ensure a default Overtime Rule exists
		var rule core.OvertimeRule
		if err := db.First(&rule).Error; err != nil {
			rule = core.OvertimeRule{
				Name:               "Standart Mesai Kuralları",
				WeeklyHourLimit:    45,
				DailyHourLimit:     11,
				OvertimeMultiplier: 1.5,
				WeekendMultiplier:  2.0,
				HolidayMultiplier:  2.5,
				NightShiftExtra:    0.1,
				IsActive:           true,
			}
			db.Create(&rule)
			slog.Info("Default overtime rule created")
		}
	}

	// ===== Initialize Layers =====

	// Department
	departmentRepo := repositories.NewDepartmentRepository(db)
	departmentService := services.NewDepartmentService(departmentRepo)
	departmentHandler := handlers.NewDepartmentHandler(departmentService)

	// Title
	titleRepo := repositories.NewTitleRepository(db)
	titleService := services.NewTitleService(titleRepo)
	titleHandler := handlers.NewTitleHandler(titleService)

	// Employee
	employeeRepo := repositories.NewEmployeeRepository(db)
	employeeService := services.NewEmployeeService(employeeRepo, departmentRepo, titleRepo)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

	// ShiftType
	shiftTypeRepo := repositories.NewShiftTypeRepository(db)
	shiftTypeService := services.NewShiftTypeService(shiftTypeRepo)
	shiftTypeHandler := handlers.NewShiftTypeHandler(shiftTypeService)

	// Leave (İzin)
	leaveRepo := repositories.NewLeaveRepository(db)
	leaveService := services.NewLeaveService(leaveRepo)
	leaveHandler := handlers.NewLeaveHandler(leaveService)

	// Schedule
	scheduleRepo := repositories.NewScheduleRepository(db)
	schedulerService := services.NewSchedulerService(scheduleRepo, employeeRepo, shiftTypeRepo, leaveRepo)
	scheduleHandler := handlers.NewScheduleHandler(schedulerService)

	// TimeEntry (Puantaj)
	timeEntryRepo := repositories.NewTimeEntryRepository(db)
	timekeepingService := services.NewTimekeepingService(timeEntryRepo, scheduleRepo)
	timeEntryHandler := handlers.NewTimeEntryHandler(timekeepingService)

	// Overtime (Mesai)
	overtimeRuleRepo := repositories.NewOvertimeRuleRepository(db)
	overtimeService := services.NewOvertimeService(overtimeRuleRepo, timeEntryRepo)
	overtimeHandler := handlers.NewOvertimeHandler(overtimeService)

	// Approval (Onay)
	auditLogRepo := repositories.NewAuditLogRepository(db)
	approvalService := services.NewApprovalService(auditLogRepo, timeEntryRepo, leaveRepo)
	approvalHandler := handlers.NewApprovalHandler(approvalService)

	// Reporting (Raporlama)
	reportingService := services.NewReportingService(timeEntryRepo, leaveRepo, overtimeRuleRepo)
	reportingHandler := handlers.NewReportingHandler(reportingService)

	// Notification (Bildirimler)
	notificationRepo := repositories.NewNotificationRepository(db)
	notificationService := services.NewNotificationService(notificationRepo)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Auth
	resetTokenRepo := repositories.NewPasswordResetTokenRepository(db)
	authService := services.NewAuthService(employeeRepo, resetTokenRepo, cfg.Server.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)
	loginRateLimiter := middleware.NewIPRateLimiter(1.0/60.0, 5) // 5 attempts per minute

	// Backup
	backupService := services.NewBackupService(cfg.Database)
	backupHandler := handlers.NewBackupHandler(backupService)

	// ===== Router =====

	router := gin.New()

	// Request logging middleware
	router.Use(logger.GinLoggerMiddleware())
	router.Use(gin.Recovery())

	// CORS Middleware — restrict to known origins
	allowedOrigins := map[string]bool{
		"http://localhost:5173": true,
		"http://localhost:3000": true,
	}
	router.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"mode":    cfg.Server.Mode,
		})
	})

	api := router.Group("/api/v1")
	{
		// Public Auth
		auth := api.Group("/auth")
		{
			auth.POST("/login", middleware.RateLimit(loginRateLimiter), authHandler.Login)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.POST("/change-password", middleware.AuthMiddleware(cfg.Server.JWTSecret), authHandler.ChangePassword)
		}

		// Protected Routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.Server.JWTSecret))
		{
			// Department
			departments := protected.Group("/departments")
			{
				departments.POST("", middleware.RoleMiddleware("admin"), departmentHandler.CreateDepartment)
				departments.GET("", departmentHandler.GetAllDepartments)
				departments.GET("/:id", departmentHandler.GetDepartment)
				departments.PUT("/:id", middleware.RoleMiddleware("admin"), departmentHandler.UpdateDepartment)
				departments.DELETE("/:id", middleware.RoleMiddleware("admin"), departmentHandler.DeleteDepartment)
			}

			// Title
			titles := protected.Group("/titles")
			{
				titles.POST("", middleware.RoleMiddleware("admin"), titleHandler.CreateTitle)
				titles.GET("", titleHandler.GetAllTitles)
				titles.GET("/:id", titleHandler.GetTitle)
				titles.PUT("/:id", middleware.RoleMiddleware("admin"), titleHandler.UpdateTitle)
				titles.DELETE("/:id", middleware.RoleMiddleware("admin"), titleHandler.DeleteTitle)
			}

			// Employee
			employees := protected.Group("/employees")
			{
				employees.POST("", middleware.RoleMiddleware("admin"), employeeHandler.CreateEmployee)
				employees.GET("", employeeHandler.GetAllEmployees)
				employees.POST("/import", middleware.RoleMiddleware("admin"), employeeHandler.ImportEmployees)
				employees.GET("/:id", employeeHandler.GetEmployee)
				employees.PUT("/:id", middleware.RoleMiddleware("admin"), employeeHandler.UpdateEmployee)
				employees.DELETE("/:id", middleware.RoleMiddleware("admin"), employeeHandler.DeleteEmployee)
			}

			// ShiftType (Çalışma Tipleri)
			shiftTypes := protected.Group("/shift-types")
			{
				shiftTypes.POST("", middleware.RoleMiddleware("admin"), shiftTypeHandler.CreateShiftType)
				shiftTypes.GET("", shiftTypeHandler.GetAllShiftTypes)
				shiftTypes.GET("/:id", shiftTypeHandler.GetShiftType)
				shiftTypes.PUT("/:id", middleware.RoleMiddleware("admin"), shiftTypeHandler.UpdateShiftType)
				shiftTypes.DELETE("/:id", middleware.RoleMiddleware("admin"), shiftTypeHandler.DeleteShiftType)
			}

			// Schedule (Nöbet Çizelgesi)
			schedules := protected.Group("/schedules")
			{
				schedules.POST("/generate", middleware.RoleMiddleware("admin"), scheduleHandler.GenerateSchedule)
				schedules.POST("", middleware.RoleMiddleware("admin"), scheduleHandler.CreateSingleSchedule)
				schedules.PUT("/:id", middleware.RoleMiddleware("admin"), scheduleHandler.UpdateSchedule)
				schedules.GET("", scheduleHandler.GetSchedule)
				schedules.DELETE("/clear", middleware.RoleMiddleware("admin"), scheduleHandler.ClearSchedule)
				schedules.DELETE("/:id", middleware.RoleMiddleware("admin"), scheduleHandler.DeleteSchedule)
			}

			// TimeEntry (Puantaj — Giriş/Çıkış)
			timeEntries := protected.Group("/time-entries")
			{
				timeEntries.POST("/clock-in", timeEntryHandler.ClockIn)
				timeEntries.POST("/clock-out", timeEntryHandler.ClockOut)
				timeEntries.POST("", middleware.RoleMiddleware("admin"), timeEntryHandler.CreateTimeEntry)
				timeEntries.GET("", timeEntryHandler.ListTimeEntries)
				timeEntries.GET("/:id", timeEntryHandler.GetTimeEntry)
				timeEntries.PUT("/:id", middleware.RoleMiddleware("admin"), timeEntryHandler.UpdateTimeEntry)
				timeEntries.DELETE("/:id", middleware.RoleMiddleware("admin"), timeEntryHandler.DeleteTimeEntry)
			}

			// Leave (İzin Yönetimi)
			leaves := protected.Group("/leaves")
			{
				leaves.POST("", leaveHandler.RequestLeave)
				leaves.GET("", leaveHandler.ListLeaves)
				leaves.GET("/balance", leaveHandler.GetLeaveBalance)
				leaves.GET("/:id", leaveHandler.GetLeave)
				leaves.POST("/:id/approve", middleware.RoleMiddleware("admin"), leaveHandler.ApproveLeave)
				leaves.POST("/:id/reject", middleware.RoleMiddleware("admin"), leaveHandler.RejectLeave)
			}

			// LeaveType (İzin Türleri)
			leaveTypes := protected.Group("/leave-types")
			{
				leaveTypes.POST("", middleware.RoleMiddleware("admin"), leaveHandler.CreateLeaveType)
				leaveTypes.GET("", leaveHandler.GetAllLeaveTypes)
				leaveTypes.GET("/:id", leaveHandler.GetLeaveType)
				leaveTypes.PUT("/:id", middleware.RoleMiddleware("admin"), leaveHandler.UpdateLeaveType)
				leaveTypes.DELETE("/:id", middleware.RoleMiddleware("admin"), leaveHandler.DeleteLeaveType)
			}

			// Overtime (Mesai Hesaplama)
			overtime := protected.Group("/overtime")
			{
				overtime.GET("/calculate", overtimeHandler.CalculateOvertime)
				overtime.GET("/summary", overtimeHandler.GetDepartmentSummary)
			}

			// OvertimeRule (Mesai Kuralları)
			overtimeRules := protected.Group("/overtime-rules")
			{
				overtimeRules.POST("", middleware.RoleMiddleware("admin"), overtimeHandler.CreateRule)
				overtimeRules.GET("", overtimeHandler.GetAllRules)
				overtimeRules.GET("/:id", overtimeHandler.GetRule)
				overtimeRules.PUT("/:id", middleware.RoleMiddleware("admin"), overtimeHandler.UpdateRule)
				overtimeRules.DELETE("/:id", middleware.RoleMiddleware("admin"), overtimeHandler.DeleteRule)
			}

			// Public Holiday (Resmi Tatiller)
			holidays := protected.Group("/public-holidays")
			{
				holidays.POST("", middleware.RoleMiddleware("admin"), overtimeHandler.CreateHoliday)
				holidays.GET("", overtimeHandler.GetHolidays)
				holidays.PUT("/:id", middleware.RoleMiddleware("admin"), overtimeHandler.UpdateHoliday)
				holidays.DELETE("/:id", middleware.RoleMiddleware("admin"), overtimeHandler.DeleteHoliday)
			}

			// Approval (Onay Mekanizması)
			approvals := protected.Group("/approvals")
			{
				approvals.GET("/pending", middleware.RoleMiddleware("admin"), approvalHandler.GetPendingApprovals)
				approvals.POST("/time-entry/:id/approve", middleware.RoleMiddleware("admin"), approvalHandler.ApproveTimeEntry)
				approvals.POST("/time-entry/:id/reject", middleware.RoleMiddleware("admin"), approvalHandler.RejectTimeEntry)
			}

			// AuditLog (Denetim İzi)
			auditLogs := protected.Group("/audit-logs")
			{
				auditLogs.GET("", middleware.RoleMiddleware("admin"), approvalHandler.GetAuditLogs)
			}

			// Reports (Raporlama)
			reports := protected.Group("/reports")
			{
				reports.GET("/work-hours", reportingHandler.GetWorkHoursReport)
				reports.GET("/absences", reportingHandler.GetAbsenceReport)
				reports.GET("/employee-summary", reportingHandler.GetEmployeeSummary)
				reports.GET("/trends", reportingHandler.GetTrendAnalysis)
				reports.GET("/cost-analysis", reportingHandler.GetCostAnalysis)
			}

			// Notifications (Bildirimler)
			notifications := protected.Group("/notifications")
			{
				notifications.GET("/unread", notificationHandler.GetUnread)
				notifications.POST("/:id/read", notificationHandler.MarkAsRead)
				notifications.POST("/read-all", notificationHandler.MarkAllAsRead)
			}

			// Backups
			backups := protected.Group("/backups")
			{
				backups.GET("/export", middleware.RoleMiddleware("admin"), backupHandler.ExportBackup)
			}
		}
	}

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling below
	go func() {
		slog.Info("Server starting", "address", addr, "mode", cfg.Server.Mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so no need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown:", "error", err)
	}

	// Properly close database connection if possible
	// (Assuming the database package provides a way to get the underlying sql.DB or close directly)
	if sqlDB, err := db.DB(); err == nil {
		slog.Info("Closing database connection...")
		sqlDB.Close()
	}

	slog.Info("Server exiting")
}
