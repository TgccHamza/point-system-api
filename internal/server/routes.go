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

	// Employee Workday routes
	workdayHandler := handlers.NewEmployeeWorkdayHandler(s.employeeWorkdayService)
	r.POST("/employee-workdays", workdayHandler.CreateEmployeeWorkday)
	r.GET("/employee-workdays/:id", workdayHandler.GetEmployeeWorkdayByID)
	r.GET("/employees/workdays/:employee_id", workdayHandler.GetEmployeeWorkdaysByEmployeeID)
	r.PUT("/employee-workdays/:id", workdayHandler.UpdateEmployeeWorkday)
	r.DELETE("/employee-workdays/:id", workdayHandler.DeleteEmployeeWorkday)

	// User routes
	userHandler := handlers.NewUserHandler(s.userService)
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users/:id", userHandler.GetUserByID)
	r.GET("/users/username/:username", userHandler.GetUserByUsername)
	r.GET("/users", userHandler.ListUsers)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)
	r.POST("/users/authenticate", userHandler.AuthenticateUser)

	// Company routes
	companyHandler := handlers.NewCompanyHandler(s.companyService)
	r.POST("/companies", companyHandler.CreateCompany)
	r.GET("/companies/:id", companyHandler.GetCompanyByID)
	r.GET("/companies", companyHandler.ListCompanies)
	r.PUT("/companies/:id", companyHandler.UpdateCompany)
	r.DELETE("/companies/:id", companyHandler.DeleteCompany)

	// Employee routes
	employeeHandler := handlers.NewEmployeeHandler(s.employeeService)
	r.POST("/employees", employeeHandler.CreateEmployee)
	r.GET("/employees/:id", employeeHandler.GetEmployeeByID)
	r.GET("/employees/by-company/:id", employeeHandler.GetEmployeesByCompanyID)
	r.PUT("/employees/:id", employeeHandler.UpdateEmployee)
	r.DELETE("/employees/:id", employeeHandler.DeleteEmployee)

	// Hello World endpoint
	r.GET("/", s.HelloWorldHandler)

	s.httpServer.Handler = r
	return r
}

// HelloWorldHandler returns a simple "Hello World" message.
func (s *Server) HelloWorldHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
}
