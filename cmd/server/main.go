package main

import (
	"fmt"
	"log"

	"github.com/akatakan/nobetgo/config"
	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/database"
	"github.com/akatakan/nobetgo/internal/handlers"
	"github.com/akatakan/nobetgo/internal/repositories"
	"github.com/akatakan/nobetgo/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	// Auto Migration
	err = db.AutoMigrate(&core.Employee{}, &core.ShiftType{}, &core.Schedule{}, &core.Attendance{})
	if err != nil {
		log.Fatalf("cannot migrate database: %v", err)
	}

	// Initialize Layers - Employee
	employeeRepo := repositories.NewEmployeeRepository(db)
	employeeService := services.NewEmployeeService(employeeRepo)
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
	timekeepingService := services.NewTimekeepingService(attendanceRepo)
	attendanceHandler := handlers.NewAttendanceHandler(timekeepingService)

	router := gin.Default()

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
		employees := api.Group("/employees")
		{
			employees.POST("", employeeHandler.CreateEmployee)
			employees.GET("", employeeHandler.GetAllEmployees)
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
			schedules.PUT("/:id", scheduleHandler.UpdateSchedule)
			// schedules.GET("", scheduleHandler.GetSchedule)
		}

		attendance := api.Group("/attendance")
		{
			attendance.POST("", attendanceHandler.LogAttendance)
			attendance.GET("/reports", attendanceHandler.GetPayrollReport)
		}
	}

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s...", addr)
	router.Run(addr)
}
