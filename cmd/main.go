package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/royanqodri/Absensi/config"
	"github.com/royanqodri/Absensi/controller"
	"github.com/royanqodri/Absensi/database"
	"github.com/royanqodri/Absensi/repository"
	"github.com/royanqodri/Absensi/service"
)

var (
	// Repository
	tAbsensiRepo repository.TAbsensiRepository
	tUserRepo    repository.TUserRepository

	// Service
	tAbsensiService service.TAbsensiService
	tUserService    service.TUserService

	// Controller
	tAbsensiController controller.TAbsensiController
	tUserController    controller.TUserController
)

func main() {
	// Init Config
	initConfig()

	// Init Database
	initDatabase()

	// Run Migration
	initMigration()

	// Init Repository
	initRepository()

	// Init Service
	initService()

	// Init Controller
	initController()

	// Init Router & Start Server
	runServer()
}

func initConfig() {
	cfg := config.InitConfig()
	if cfg == nil {
		log.Fatalf("❌ failed to initialize config")
	}
	log.Println("✅ Config initialized successfully")
}

func initDatabase() {
	// Panggil ConnectDatabase
	database.ConnectDatabase()

	if database.DBConn == nil {
		log.Fatalf("❌ Database connection is nil")
	}

	log.Println("✅ Database ready")
}

// TAMBAHKAN FUNGSI INI UNTUK MIGRASI
func initMigration() {
	log.Println("🔄 Running database migration...")
	database.MigrateDatabase()
	log.Println("✅ Database migration completed")
}

func initRepository() {
	tAbsensiRepo = repository.NewTAbsensiRepository()
	tUserRepo = repository.NewTUserRepository()
}

func initService() {
	tAbsensiService = service.NewTAbsensiService(tAbsensiRepo)
	tUserService = service.NewTUserService(tUserRepo)
}

func initController() {
	tAbsensiController = controller.NewTAbsensiController(tAbsensiService)
	tUserController = controller.NewTUserController(tUserService)
}

func runServer() {
	e := echo.New()

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))
	e.Use(middleware.Recover())

	// EndPoint
	api := e.Group("/api/v1")
	api.GET("/absensi", tAbsensiController.GetByParams)
	api.POST("/absensi", tAbsensiController.Post)
	api.GET("/user", tUserController.GetByParams)
	api.POST("/user", tUserController.Post)

	// Health Check
	e.GET("/health", func(c echo.Context) error {
		if database.DBConn == nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status": "unhealthy",
				"error":  "database not connected",
			})
		}

		sqlDB, err := database.DBConn.DB()
		if err != nil || sqlDB.Ping() != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status": "unhealthy",
				"error":  "database connection lost",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"db":     "connected",
		})
	})

	// Get port from config
	cfg := config.Get()
	port := cfg.ServerPort
	if port == 0 {
		port = 8080
	}
	portStr := fmt.Sprintf(":%d", port)

	// Start server with graceful shutdown
	go func() {
		log.Printf("🚀 Server running on %s", portStr)
		log.Printf("📊 Environment: %s", cfg.ServiceEnv)
		if err := e.Start(portStr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("❌ server shutdown error: %v", err)
	}
	log.Println("✅ Server shutdown gracefully")
}
