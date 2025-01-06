package handlers

import (
	"log/slog"
	"net/http"
	"point-system-api/internal/database"

	"github.com/gin-gonic/gin"
)

// HandleMigrate returns a gin.HandlerFunc that runs the database migrations.
// It calls the MigrateDB method on the database package and returns the result as a JSON response.
// @Summary Run database migrations
// @Description Runs the database migrations
func HandleMigrate(c *gin.Context) {
	err := database.MigrateDB()
	if err != nil {
		slog.Error("Error migrating database", "err", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "Database migrated successfully")
}

// HelloWorldHandler returns a gin.HandlerFunc that responds to GET /
// requests with a JSON message of "Hello World".
// @Summary Hello World
// @Description Returns a hello world message
// @Produce json
// @Success 200 {object} map[string]string
func HelloWorldHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
}

// HealthHandler returns a gin.HandlerFunc that returns the health status of the database.
// It calls the Health method on the provided database.Service and returns the result as a JSON response.
// @Summary Health check
// @Description Checks if the server is healthy
// @Produce plain
// @Success 200 {string} string "Server is healthy"
// @Router /client-auth/health [get]
func HealthHandler(dbService database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		healthStatus := dbService.Health()
		c.JSON(http.StatusOK, healthStatus)
	}
}
