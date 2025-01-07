package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"point-system-api/internal/handlers"
)

// RegisterRoutes sets up all the routes for the application.
func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", handlers.HealthHandler(s.db))

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
