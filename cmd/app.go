package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/wabtcdi/user_service/cmd/health"
	"github.com/wabtcdi/user_service/cmd/log"

	"github.com/gorilla/mux"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

type Pinger interface {
	Ping() error
}

type DBOpener func(driverName, dataSourceName string) (*sql.DB, error)

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

func connectDatabase(cfg Config, opener DBOpener) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
	db, err := opener("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, "../database/migrations"); err != nil {
		return nil, fmt.Errorf("failed to run goose migrations: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logrus.Info("Successfully connected to PostgreSQL!")
	return db, nil
}

func startServer(cfg Config, db *sql.DB, starter ServerStarter) error {
	r := createRouter(cfg, db)
	addr := getAddr(cfg)
	logrus.Infof("Starting server on %s", addr)
	return starter.Start(addr, r)
}

func createRouter(cfg Config, db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	checker := &health.Checker{DB: db}

	r.HandleFunc(cfg.Server.LivenessPath, livenessHandler).Methods("GET")
	r.HandleFunc(cfg.Server.ReadinessPath, checker.Check).Methods("GET")

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
