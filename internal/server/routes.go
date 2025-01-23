package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"point-system-api/internal/handlers"

	"github.com/gin-contrib/cors"
)

// RegisterRoutes sets up all the routes for the application.
func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Referer", "Accept", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	r.GET("/health", handlers.HealthHandler(s.db))

	// EmployeeWorkDay routes
	employeeWorkDayHandler := handlers.NewEmployeeWorkDayHandler(s.employeeWorkDayService)
	r.POST("/employee-workdays/generate/:workDayID", employeeWorkDayHandler.GenerateEmployeeWorkDay)
	r.PUT("/employee-workdays/:id", employeeWorkDayHandler.UpdateEmployeeWorkDay)

	// RawAttendance routes
	rawAttendanceHandler := handlers.NewRawAttendanceHandler(s.rawAttendanceService)
	r.POST("/raw-attendances", rawAttendanceHandler.CreateRawAttendance)
	r.POST("/raw-attendances/bulk", rawAttendanceHandler.CreateManyRawAttendances)
	r.GET("/raw-attendances/:id", rawAttendanceHandler.GetRawAttendanceByID)
	r.GET("/raw-attendances/work-day/:workDayID", rawAttendanceHandler.GetRawAttendancesByWorkDayID)
	r.PUT("/raw-attendances/:id", rawAttendanceHandler.UpdateRawAttendance)
	r.DELETE("/raw-attendances/:id", rawAttendanceHandler.DeleteRawAttendance)

	// WorkDay routes
	workDayHandler := handlers.NewWorkDayHandler(s.workDayService)
	r.POST("/workdays", workDayHandler.CreateWorkDay)
	r.GET("/workdays/:id", workDayHandler.GetWorkDayByID)
	r.GET("/workdays", workDayHandler.ListWorkDays)
	r.PUT("/workdays/:id", workDayHandler.UpdateWorkDay)
	r.DELETE("/workdays/:id", workDayHandler.DeleteWorkDay)

	// User routes
	userHandler := handlers.NewUserHandler(s.userService)
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users/:id", userHandler.GetUserByID)
	r.GET("/users/username/:username", userHandler.GetUserByUsername)
	r.GET("/users", userHandler.ListUsers)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)
	r.POST("/users/authenticate", userHandler.AuthenticateUser)
	r.GET("/users/select", userHandler.ListUsersForSelect)

	// Company routes
	companyHandler := handlers.NewCompanyHandler(s.companyService)
	r.POST("/companies", companyHandler.CreateCompany)
	r.GET("/companies/:id", companyHandler.GetCompanyByID)
	r.GET("/companies", companyHandler.ListCompanies)
	r.PUT("/companies/:id", companyHandler.UpdateCompany)
	r.DELETE("/companies/:id", companyHandler.DeleteCompany)
	r.GET("/companies/select", companyHandler.ListCompaniesForSelect)

	// Employee routes
	employeeHandler := handlers.NewEmployeeHandler(s.employeeService)
	r.POST("/employees", employeeHandler.CreateEmployee)
	r.GET("/employees/:id", employeeHandler.GetEmployeeByID)
	r.GET("/employees/by-company/:id", employeeHandler.GetEmployeesByCompanyID)
	r.GET("/employees", employeeHandler.FetchEmployees)
	r.PUT("/employees/:id", employeeHandler.UpdateEmployee)
	r.DELETE("/employees/:id", employeeHandler.DeleteEmployee)

	// Hello World endpoint
	r.GET("/", s.HelloWorldHandler)

	attendanceHandler := handlers.NewAttendanceHandler(s.attendanceService)
	r.POST("/process-hex", attendanceHandler.CreateAttendanceLog)
	r.GET("/attendance-logs", attendanceHandler.ListAttendanceLogs)
	r.GET("/attendance-logs/:id", attendanceHandler.GetAttendanceLogByID)
	deviceHandler := handlers.NewDeviceHandler(s.deviceService)
	RegisterDeviceRoutes(r, deviceHandler)

	r.GET("/ws", handlers.ServeWs)
	s.httpServer.Handler = r
	return r
}

// HelloWorldHandler returns a simple "Hello World" message.
func (s *Server) HelloWorldHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
}

func RegisterDeviceRoutes(router *gin.Engine, deviceHandler *handlers.DeviceHandler) {
	devices := router.Group("/devices")
	{
		devices.GET("", deviceHandler.GetAllDevices)       // Retrieve all devices (with filters)
		devices.GET("/:id", deviceHandler.GetDeviceByID)   // Retrieve a single device
		devices.POST("", deviceHandler.CreateDevice)       // Create a new device
		devices.PUT("/:id", deviceHandler.UpdateDevice)    // Update a device
		devices.DELETE("/:id", deviceHandler.DeleteDevice) // Delete a device
	}
}
