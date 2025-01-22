package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"point-system-api/internal/database"
	"point-system-api/internal/handlers"
	"point-system-api/internal/repositories"
	"point-system-api/internal/services"
)

// Server represents the HTTP server and its dependencies.
type Server struct {
	httpServer             *http.Server
	port                   int
	db                     database.Service
	userService            services.UserService
	companyService         services.CompanyService
	employeeService        services.EmployeeService
	workDayService         services.WorkDayService
	rawAttendanceService   services.RawAttendanceService
	employeeWorkDayService services.EmployeeWorkDayService
	attendanceService      services.AttendanceService
	deviceService          services.DeviceService
}

// NewServer creates a new instance of the Server.
func NewServer() *Server {
	handlers.StartManager()
	// Load configuration
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// Initialize database
	db := database.New()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db.GetDB())
	companyRepo := repositories.NewCompanyRepository(db.GetDB())
	employeeRepo := repositories.NewEmployeeRepository(db.GetDB())
	workDayRepo := repositories.NewWorkDayRepository(db.GetDB())
	rawAttendanceRepo := repositories.NewRawAttendanceRepository(db.GetDB())
	employeeWorkDayRepo := repositories.NewEmployeeWorkDayRepository(db.GetDB())
	deviceRepo := repositories.NewDeviceRepository(db.GetDB())
	attendanceRepo := repositories.NewAttendanceRepository(db.GetDB())

	// Initialize services
	userService := services.NewUserService(userRepo)
	companyService := services.NewCompanyService(companyRepo)
	employeeService := services.NewEmployeeService(employeeRepo, userService)
	workDayService := services.NewWorkDayService(workDayRepo)
	rawAttendanceService := services.NewRawAttendanceService(rawAttendanceRepo)
	employeeWorkDayService := services.NewEmployeeWorkDayService(workDayRepo, rawAttendanceRepo, employeeWorkDayRepo)
	attendanceService := services.NewAttendanceService(deviceRepo, attendanceRepo)
	deviceService := services.NewDeviceService(deviceRepo)

	// Create the HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      nil, // Will be set in RegisterRoutes
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &Server{
		httpServer:             httpServer,
		port:                   port,
		db:                     db,
		userService:            userService,
		companyService:         companyService,
		employeeService:        employeeService,
		workDayService:         workDayService,
		rawAttendanceService:   rawAttendanceService,
		employeeWorkDayService: employeeWorkDayService,
		attendanceService:      attendanceService,
		deviceService:          deviceService,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	// Register routes
	s.RegisterRoutes()

	// Start the server
	log.Printf("Server started on port %d", s.port)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	// Shutdown the HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %v", err)
	}

	log.Println("Server shutdown complete")
	return nil
}
