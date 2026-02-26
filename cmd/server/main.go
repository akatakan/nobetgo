package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/akatakan/nobetgo/config"
	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/database"
	"github.com/akatakan/nobetgo/internal/handlers"
	"github.com/akatakan/nobetgo/internal/logger"
	"github.com/akatakan/nobetgo/internal/repositories"
	"github.com/akatakan/nobetgo/internal/services"
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
		&core.Notification{}, // <-- new notification model
	)
	if err != nil {
		slog.Error("Cannot migrate database", "error", err)
		os.Exit(1)
	}
	// Partial unique index: enforce uniqueness only on non-empty emails
	db.Exec("DROP INDEX IF EXISTS idx_employees_email")
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_employees_email_unique ON employees (email) WHERE email != '' AND deleted_at IS NULL")
	slog.Info("Database migration completed")

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

	// Schedule
	scheduleRepo := repositories.NewScheduleRepository(db)
	schedulerService := services.NewSchedulerService(scheduleRepo, employeeRepo, shiftTypeRepo)
	scheduleHandler := handlers.NewScheduleHandler(schedulerService)

	// TimeEntry (Puantaj)
	timeEntryRepo := repositories.NewTimeEntryRepository(db)
	timekeepingService := services.NewTimekeepingService(timeEntryRepo, scheduleRepo)
	timeEntryHandler := handlers.NewTimeEntryHandler(timekeepingService)

	// Leave (İzin)
	leaveRepo := repositories.NewLeaveRepository(db)
	leaveService := services.NewLeaveService(leaveRepo)
	leaveHandler := handlers.NewLeaveHandler(leaveService)

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

	// ===== Router =====

	router := gin.New()

	// Request logging middleware
	router.Use(logger.GinLoggerMiddleware())
	router.Use(gin.Recovery())

	// CORS Middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
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
		// Department
		departments := api.Group("/departments")
		{
			departments.POST("", departmentHandler.CreateDepartment)
			departments.GET("", departmentHandler.GetAllDepartments)
			departments.GET("/:id", departmentHandler.GetDepartment)
			departments.PUT("/:id", departmentHandler.UpdateDepartment)
			departments.DELETE("/:id", departmentHandler.DeleteDepartment)
		}

		// Title
		titles := api.Group("/titles")
		{
			titles.POST("", titleHandler.CreateTitle)
			titles.GET("", titleHandler.GetAllTitles)
			titles.GET("/:id", titleHandler.GetTitle)
			titles.PUT("/:id", titleHandler.UpdateTitle)
			titles.DELETE("/:id", titleHandler.DeleteTitle)
		}

		// Employee
		employees := api.Group("/employees")
		{
			employees.POST("", employeeHandler.CreateEmployee)
			employees.GET("", employeeHandler.GetAllEmployees)
			employees.POST("/import", employeeHandler.ImportEmployees)
			employees.GET("/:id", employeeHandler.GetEmployee)
			employees.PUT("/:id", employeeHandler.UpdateEmployee)
			employees.DELETE("/:id", employeeHandler.DeleteEmployee)
		}

		// ShiftType (Çalışma Tipleri)
		shiftTypes := api.Group("/shift-types")
		{
			shiftTypes.POST("", shiftTypeHandler.CreateShiftType)
			shiftTypes.GET("", shiftTypeHandler.GetAllShiftTypes)
			shiftTypes.GET("/:id", shiftTypeHandler.GetShiftType)
			shiftTypes.PUT("/:id", shiftTypeHandler.UpdateShiftType)
			shiftTypes.DELETE("/:id", shiftTypeHandler.DeleteShiftType)
		}

		// Schedule (Nöbet Çizelgesi)
		schedules := api.Group("/schedules")
		{
			schedules.POST("/generate", scheduleHandler.GenerateSchedule)
			schedules.POST("", scheduleHandler.CreateSingleSchedule)
			schedules.PUT("/:id", scheduleHandler.UpdateSchedule)
			schedules.GET("", scheduleHandler.GetSchedule)
			schedules.DELETE("/clear", scheduleHandler.ClearSchedule)
			schedules.DELETE("/:id", scheduleHandler.DeleteSchedule)
		}

		// TimeEntry (Puantaj — Giriş/Çıkış)
		timeEntries := api.Group("/time-entries")
		{
			timeEntries.POST("/clock-in", timeEntryHandler.ClockIn)
			timeEntries.POST("/clock-out", timeEntryHandler.ClockOut)
			timeEntries.POST("", timeEntryHandler.CreateTimeEntry)
			timeEntries.GET("", timeEntryHandler.ListTimeEntries)
			timeEntries.GET("/:id", timeEntryHandler.GetTimeEntry)
			timeEntries.PUT("/:id", timeEntryHandler.UpdateTimeEntry)
			timeEntries.DELETE("/:id", timeEntryHandler.DeleteTimeEntry)
		}

		// Leave (İzin Yönetimi)
		leaves := api.Group("/leaves")
		{
			leaves.POST("", leaveHandler.RequestLeave)
			leaves.GET("", leaveHandler.ListLeaves)
			leaves.GET("/balance", leaveHandler.GetLeaveBalance)
			leaves.GET("/:id", leaveHandler.GetLeave)
			leaves.POST("/:id/approve", leaveHandler.ApproveLeave)
			leaves.POST("/:id/reject", leaveHandler.RejectLeave)
		}

		// LeaveType (İzin Türleri)
		leaveTypes := api.Group("/leave-types")
		{
			leaveTypes.POST("", leaveHandler.CreateLeaveType)
			leaveTypes.GET("", leaveHandler.GetAllLeaveTypes)
			leaveTypes.GET("/:id", leaveHandler.GetLeaveType)
			leaveTypes.PUT("/:id", leaveHandler.UpdateLeaveType)
			leaveTypes.DELETE("/:id", leaveHandler.DeleteLeaveType)
		}

		// Overtime (Mesai Hesaplama)
		overtime := api.Group("/overtime")
		{
			overtime.GET("/calculate", overtimeHandler.CalculateOvertime)
			overtime.GET("/summary", overtimeHandler.GetDepartmentSummary)
		}

		// OvertimeRule (Mesai Kuralları)
		overtimeRules := api.Group("/overtime-rules")
		{
			overtimeRules.POST("", overtimeHandler.CreateRule)
			overtimeRules.GET("", overtimeHandler.GetAllRules)
			overtimeRules.GET("/:id", overtimeHandler.GetRule)
			overtimeRules.PUT("/:id", overtimeHandler.UpdateRule)
			overtimeRules.DELETE("/:id", overtimeHandler.DeleteRule)
		}

		// PublicHoliday (Resmi Tatiller)
		holidays := api.Group("/public-holidays")
		{
			holidays.POST("", overtimeHandler.CreateHoliday)
			holidays.GET("", overtimeHandler.GetHolidays)
			holidays.PUT("/:id", overtimeHandler.UpdateHoliday)
			holidays.DELETE("/:id", overtimeHandler.DeleteHoliday)
		}

		// Approval (Onay Mekanizması)
		approvals := api.Group("/approvals")
		{
			approvals.GET("/pending", approvalHandler.GetPendingApprovals)
			approvals.POST("/time-entry/:id/approve", approvalHandler.ApproveTimeEntry)
			approvals.POST("/time-entry/:id/reject", approvalHandler.RejectTimeEntry)
		}

		// AuditLog (Denetim İzi)
		auditLogs := api.Group("/audit-logs")
		{
			auditLogs.GET("", approvalHandler.GetAuditLogs)
		}

		// Reports (Raporlama)
		reports := api.Group("/reports")
		{
			reports.GET("/work-hours", reportingHandler.GetWorkHoursReport)
			reports.GET("/absences", reportingHandler.GetAbsenceReport)
			reports.GET("/employee-summary", reportingHandler.GetEmployeeSummary)
			reports.GET("/trends", reportingHandler.GetTrendAnalysis)
			reports.GET("/cost-analysis", reportingHandler.GetCostAnalysis)
		}

		// Notifications (Bildirimler)
		notifications := api.Group("/notifications")
		{
			notifications.GET("/unread", notificationHandler.GetUnread)
			notifications.POST("/:id/read", notificationHandler.MarkAsRead)
			notifications.POST("/read-all", notificationHandler.MarkAllAsRead)
		}
	}

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	slog.Info("Server starting", "address", addr, "mode", cfg.Server.Mode)
	if err := router.Run(addr); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
