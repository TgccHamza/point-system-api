package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"point-system-api/internal/models"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	GetDB() *gorm.DB
}

func (s *service) GetDB() *gorm.DB {
	return s.db
}

type service struct {
	db *gorm.DB
}

var (
	dbname     =  os.Getenv("BLUEPRINT_DB_DATABASE")
	password   =  os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   =  os.Getenv("BLUEPRINT_DB_USERNAME")
	port       =  os.Getenv("BLUEPRINT_DB_PORT")
	host       =  os.Getenv("BLUEPRINT_DB_HOST")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	// Opening a driver typically will not attempt to connect to the database.
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname))
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	log.Printf("Initializing Gorm...")
	gormDb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		log.Printf("Failed to initialize Gorm: %v", err)
		return nil
	}

	log.Printf("Database connected successfully")

	dbInstance = &service{
		db: gormDb,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	sqldb, err := s.db.DB()
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database
	err = sqldb.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("%s", fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := sqldb.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dbname)
	sqldb, err := s.db.DB()
	if err != nil {
		log.Fatal(err)
	}
	return sqldb.Close()
}

func MigrateDB() error {
	err := dbInstance.db.AutoMigrate(
		&models.User{},
		&models.Employee{},
		&models.Company{},
		&models.WorkDay{},
		&models.RawAttendance{},
		&models.EmployeeWorkDay{},
		&models.AttendanceLog{},
		&models.Device{},
	)
	if err != nil {
		log.Printf("Database migration failed: %v", err)
		return err
	}

	log.Printf("Database migrated successfully")
	return nil
}

func InitializeViewDB() error {
	user_daily_checkin_checkout := `CREATE OR REPLACE VIEW user_daily_checkin_checkout AS 
	WITH DailyPunchData AS (
            SELECT 
                attendance_logs.user_id AS user_id,
                CAST(attendance_logs.timestamp AS DATE) AS date,
                MIN(CASE WHEN (attendance_logs.system_punch = 'IN') THEN attendance_logs.timestamp END) AS checkin,
                MAX(CASE WHEN (attendance_logs.system_punch = 'OUT') THEN attendance_logs.timestamp END) AS last_out_punch,
                MAX(attendance_logs.timestamp) AS last_punch_of_day,
                MAX(CASE WHEN (attendance_logs.system_punch = 'IN') THEN attendance_logs.timestamp END) AS last_in_punch
            FROM attendance_logs 
            GROUP BY attendance_logs.user_id, CAST(attendance_logs.timestamp AS DATE)
        ), 
        NextDayPunch AS (
            SELECT 
                c.user_id AS user_id,
                c.date AS date,
                MIN(n.timestamp) AS next_out_punch
            FROM DailyPunchData c 
            LEFT JOIN attendance_logs n 
                ON ((c.user_id = n.user_id) 
                AND (CAST(n.timestamp AS DATE) = (c.date + INTERVAL 1 DAY)) 
                AND (n.system_punch = 'OUT'))
            GROUP BY c.user_id, c.date
       ) 



        SELECT 
            d.user_id AS user_id,
            d.date AS date,
            d.checkin AS checkin,
            (CASE 
                WHEN ((d.last_punch_of_day = d.last_in_punch) AND (nd.next_out_punch IS NOT NULL)) 
                THEN nd.next_out_punch 
                WHEN (d.last_punch_of_day <> d.last_in_punch) 
                THEN d.last_out_punch 
                ELSE NULL 
            END) AS checkout
        FROM DailyPunchData d 
        LEFT JOIN NextDayPunch nd 
            ON ((d.user_id = nd.user_id) AND (d.date = nd.date))
        ORDER BY d.user_id, d.date`

	if err := dbInstance.db.Exec(user_daily_checkin_checkout).Error; err != nil {
		return fmt.Errorf("failed to create user_daily_checkin_checkout view: %w", err)
	}

	return nil
}
