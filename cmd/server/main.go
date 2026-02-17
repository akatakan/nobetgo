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
		&core.Schedule{},
		&core.Attendance{},
	)
	if err != nil {
		slog.Error("Cannot migrate database", "error", err)
		os.Exit(1)
	}
	// Partial unique index: enforce uniqueness only on non-empty emails
	db.Exec("DROP INDEX IF EXISTS idx_employees_email")
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_employees_email_unique ON employees (email) WHERE email != '' AND deleted_at IS NULL")
	slog.Info("Database migration completed")

	// Initialize Layers - Department
	departmentRepo := repositories.NewDepartmentRepository(db)
	departmentService := services.NewDepartmentService(departmentRepo)
	departmentHandler := handlers.NewDepartmentHandler(departmentService)

	// Initialize Layers - Title
	titleRepo := repositories.NewTitleRepository(db)
	titleService := services.NewTitleService(titleRepo)
	titleHandler := handlers.NewTitleHandler(titleService)

	// Initialize Layers - Employee
	employeeRepo := repositories.NewEmployeeRepository(db)
	employeeService := services.NewEmployeeService(employeeRepo, departmentRepo, titleRepo)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

	// Initialize Layers - ShiftType
	shiftTypeRepo := repositories.NewShiftTypeRepository(db)
	shiftTypeService := services.NewShiftTypeService(shiftTypeRepo)
	shiftTypeHandler := handlers.NewShiftTypeHandler(shiftTypeService)

	// Initialize Layers - Schedule
	scheduleRepo := repositories.NewScheduleRepository(db)
	schedulerService := services.NewSchedulerService(scheduleRepo, employeeRepo, shiftTypeRepo)
	scheduleHandler := handlers.NewScheduleHandler(schedulerService)

	// Initialize Layers - Attendance
	attendanceRepo := repositories.NewAttendanceRepository(db)
	timekeepingService := services.NewTimekeepingService(attendanceRepo, scheduleRepo)
	attendanceHandler := handlers.NewAttendanceHandler(timekeepingService)

	router := gin.New()

	// Request logging middleware (replaces gin.Default logger)
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
		departments := api.Group("/departments")
		{
			departments.POST("", departmentHandler.CreateDepartment)
			departments.GET("", departmentHandler.GetAllDepartments)
			departments.GET("/:id", departmentHandler.GetDepartment)
			departments.PUT("/:id", departmentHandler.UpdateDepartment)
			departments.DELETE("/:id", departmentHandler.DeleteDepartment)
		}

		titles := api.Group("/titles")
		{
			titles.POST("", titleHandler.CreateTitle)
			titles.GET("", titleHandler.GetAllTitles)
			titles.GET("/:id", titleHandler.GetTitle)
			titles.PUT("/:id", titleHandler.UpdateTitle)
			titles.DELETE("/:id", titleHandler.DeleteTitle)
		}

		employees := api.Group("/employees")
		{
			employees.POST("", employeeHandler.CreateEmployee)
			employees.GET("", employeeHandler.GetAllEmployees)
			employees.POST("/import", employeeHandler.ImportEmployees) // Added import route
			employees.GET("/:id", employeeHandler.GetEmployee)
			employees.PUT("/:id", employeeHandler.UpdateEmployee)
			employees.DELETE("/:id", employeeHandler.DeleteEmployee)
		}

		shiftTypes := api.Group("/shift-types")
		{
			shiftTypes.POST("", shiftTypeHandler.CreateShiftType)
			shiftTypes.GET("", shiftTypeHandler.GetAllShiftTypes)
			shiftTypes.GET("/:id", shiftTypeHandler.GetShiftType)
			shiftTypes.PUT("/:id", shiftTypeHandler.UpdateShiftType)
			shiftTypes.DELETE("/:id", shiftTypeHandler.DeleteShiftType)
		}

		schedules := api.Group("/schedules")
		{
			schedules.POST("/generate", scheduleHandler.GenerateSchedule)
			schedules.POST("", scheduleHandler.CreateSingleSchedule) // Manual add
			schedules.PUT("/:id", scheduleHandler.UpdateSchedule)
			schedules.GET("", scheduleHandler.GetSchedule)
			schedules.DELETE("/clear", scheduleHandler.ClearSchedule)
			schedules.DELETE("/:id", scheduleHandler.DeleteSchedule) // Manual delete
		}

		attendance := api.Group("/attendance")
		{
			attendance.POST("", attendanceHandler.LogAttendance)
			attendance.PUT("/:id", attendanceHandler.UpdateAttendance) // Manual update
			attendance.GET("/reports", attendanceHandler.GetPayrollReport)
		}
	}

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	slog.Info("Server starting", "address", addr, "mode", cfg.Server.Mode)
	if err := router.Run(addr); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
