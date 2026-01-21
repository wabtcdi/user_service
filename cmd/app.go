package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/wabtcdi/user_service/cmd/health"
	"github.com/wabtcdi/user_service/cmd/log"
	"github.com/wabtcdi/user_service/handlers"
	"github.com/wabtcdi/user_service/repository"
	"github.com/wabtcdi/user_service/service"

	"github.com/gorilla/mux"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Pinger interface {
	Ping() error
}

type DBOpener func(dsn string) (*gorm.DB, error)

type ServerStarter interface {
	Start(addr string, handler http.Handler) error
}

type RealStarter struct{}

func (r *RealStarter) Start(addr string, handler http.Handler) error {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	errChan := make(chan error, 1)

	// Start server in goroutine
	go func() {
		errChan <- server.ListenAndServe()
	}()

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Wait for either server error or timeout
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		// Timeout occurred, shutdown server gracefully
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		server.Shutdown(shutdownCtx)
		return fmt.Errorf("server start timeout after 10 seconds")
	}
}

func Init(configName string, opener DBOpener, starter ServerStarter) error {
	path := "../resources/" + configName + ".yaml"
	cfg, err := loadConfiguration(path)
	if err != nil {
		return err
	}
	log.Configure(cfg.Logging.Level, cfg.Logging.Format)
	db, err := connectDatabase(cfg, opener)
	if err != nil {
		return err
	}
	err = startServer(cfg, db, starter)
	if err != nil {
		return err
	}
	return nil
}

func loadConfiguration(path string) (Config, error) {
	var cfg Config
	err := LoadConfig(&cfg, path)
	if err != nil {
		return Config{}, err
	}
	logrus.Info("Configuration loaded successfully")
	return cfg, nil
}

func connectDatabase(cfg Config, opener DBOpener) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
	db, err := opener(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Get underlying sql.DB for connection pooling and Goose migrations
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	maxOpenConns := 25
	maxIdleConns := 5
	connMaxLifetime := 5 * time.Minute

	if cfg.Resources.Threads > 0 {
		maxOpenConns = cfg.Resources.Threads * 2
		maxIdleConns = cfg.Resources.Threads / 2
		if maxIdleConns < 2 {
			maxIdleConns = 2
		}
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	logrus.Infof("Connection pool configured: max_open=%d, max_idle=%d, max_lifetime=%v",
		maxOpenConns, maxIdleConns, connMaxLifetime)

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(sqlDB, "../database/migrations"); err != nil {
		return nil, fmt.Errorf("failed to run goose migrations: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logrus.Info("Successfully connected to PostgreSQL!")
	return db, nil
}

func startServer(cfg Config, db *gorm.DB, starter ServerStarter) error {
	r := createRouter(cfg, db)
	addr := getAddr(cfg)
	logrus.Infof("Starting server on %s", addr)
	return starter.Start(addr, r)
}

func createRouter(cfg Config, db *gorm.DB) *mux.Router {
	r := mux.NewRouter()

	// Health checks
	checker := &health.Checker{DB: db}
	r.HandleFunc(cfg.Server.LivenessPath, livenessHandler).Methods("GET")
	r.HandleFunc(cfg.Server.ReadinessPath, checker.Check).Methods("GET")

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(db)
	accessLevelRepo := repository.NewPostgresAccessLevelRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, accessLevelRepo)
	accessLevelService := service.NewAccessLevelService(accessLevelRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	accessLevelHandler := handlers.NewAccessLevelHandler(accessLevelService)

	// User routes
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
	r.HandleFunc("/users/{id}/access-levels", userHandler.AssignAccessLevels).Methods("POST")
	r.HandleFunc("/users/{id}/access-levels", userHandler.GetUserAccessLevels).Methods("GET")

	// Authentication routes
	r.HandleFunc("/auth/login", userHandler.Login).Methods("POST")

	// Access level routes
	r.HandleFunc("/access-levels", accessLevelHandler.CreateAccessLevel).Methods("POST")
	r.HandleFunc("/access-levels", accessLevelHandler.ListAccessLevels).Methods("GET")
	r.HandleFunc("/access-levels/{id}", accessLevelHandler.GetAccessLevel).Methods("GET")

	return r
}

func getAddr(cfg Config) string {
	return fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
}

func livenessHandler(w http.ResponseWriter, _ *http.Request) {
	logrus.Debug("Liveness check requested")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
