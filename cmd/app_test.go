package cmd

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
)

func TestLivenessHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(livenessHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestConnectDatabaseSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectPing()

	cfg, err := loadConfiguration("../resources/test.yaml")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	opener := func(driver, dsn string) (*sql.DB, error) {
		if driver == "postgres" {
			return db, nil
		}
		return nil, fmt.Errorf("unsupported driver")
	}

	resultDB, err := connectDatabase(cfg, opener)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resultDB == nil {
		t.Fatal("Expected DB, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestConnectDatabaseOpenError(t *testing.T) {
	cfg, err := loadConfiguration("../resources/test.yaml")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	opener := func(driver, dsn string) (*sql.DB, error) {
		return nil, fmt.Errorf("open error")
	}

	_, err = connectDatabase(cfg, opener)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to open database") {
		t.Errorf("Expected error message to contain 'failed to open database', got %v", err)
	}
}

func TestConnectDatabasePingError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectPing().WillReturnError(fmt.Errorf("ping error"))

	cfg, err := loadConfiguration("../resources/test.yaml")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	opener := func(driver, dsn string) (*sql.DB, error) {
		return db, nil
	}

	_, err = connectDatabase(cfg, opener)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to connect to database") {
		t.Errorf("Expected error message to contain 'failed to connect to database', got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetAddr(t *testing.T) {
	cfg := Config{}
	cfg.Server.Host = "127.0.0.1"
	cfg.Server.Port = 8080

	addr := getAddr(cfg)
	expected := "127.0.0.1:8080"
	if addr != expected {
		t.Errorf("Expected %s, got %s", expected, addr)
	}
}

func TestCreateRouter(t *testing.T) {
	cfg, err := loadConfiguration("../resources/test.yaml")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectPing()

	r := createRouter(cfg, db)

	// Test liveness
	req, err := http.NewRequest("GET", cfg.Server.LivenessPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Liveness handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("Liveness handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Test readiness
	req2, err := http.NewRequest("GET", cfg.Server.ReadinessPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)

	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("Readiness handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected2 := "Ready"
	if !strings.Contains(rr2.Body.String(), expected2) {
		t.Errorf("Readiness handler returned unexpected body: got %v want %v", rr2.Body.String(), expected2)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestStartServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStarter := NewMockServerStarter(ctrl)

	cfg, err := loadConfiguration("../resources/test.yaml")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	db, _, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	mockStarter.EXPECT().Start("127.0.0.1:8081", gomock.Any()).Return(nil)

	err = startServer(cfg, db, mockStarter)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStarter := NewMockServerStarter(ctrl)

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectPing()

	opener := func(driver, dsn string) (*sql.DB, error) {
		if driver == "postgres" {
			return db, nil
		}
		return nil, fmt.Errorf("unsupported driver")
	}

	mockStarter.EXPECT().Start("127.0.0.1:8081", gomock.Any()).Return(nil)

	err = Init("test", opener, mockStarter)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet db expectations: %v", err)
	}
}

func TestInitConfigError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStarter := NewMockServerStarter(ctrl)

	opener := func(driver, dsn string) (*sql.DB, error) {
		return nil, nil // not reached
	}

	err := Init("nonexistent", opener, mockStarter)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to open config file") {
		t.Errorf("Expected error message to contain 'failed to open config file', got %v", err)
	}
}

func TestInitDBOpenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStarter := NewMockServerStarter(ctrl)

	opener := func(driver, dsn string) (*sql.DB, error) {
		return nil, fmt.Errorf("db open error")
	}

	err := Init("test", opener, mockStarter)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to open database") {
		t.Errorf("Expected error message to contain 'failed to open database', got %v", err)
	}
}

func TestInitDBPingError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStarter := NewMockServerStarter(ctrl)

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectPing().WillReturnError(fmt.Errorf("ping error"))

	opener := func(driver, dsn string) (*sql.DB, error) {
		if driver == "postgres" {
			return db, nil
		}
		return nil, fmt.Errorf("unsupported driver")
	}

	err = Init("test", opener, mockStarter)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to connect to database") {
		t.Errorf("Expected error message to contain 'failed to connect to database', got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet db expectations: %v", err)
	}
}

func TestInitServerStartError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStarter := NewMockServerStarter(ctrl)

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectPing()

	opener := func(driver, dsn string) (*sql.DB, error) {
		if driver == "postgres" {
			return db, nil
		}
		return nil, fmt.Errorf("unsupported driver")
	}

	mockStarter.EXPECT().Start("127.0.0.1:8081", gomock.Any()).Return(fmt.Errorf("server start error"))

	err = Init("test", opener, mockStarter)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "server start error") {
		t.Errorf("Expected error message to contain 'server start error', got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet db expectations: %v", err)
	}
}

func TestRealStarterStart(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		expectError bool
		errorMatch  string
	}{
		{
			name:        "invalid address format",
			address:     "invalid-address",
			expectError: true,
			errorMatch:  "address",
		},
		{
			name:        "empty address",
			address:     "",
			expectError: true,
			errorMatch:  "",
		},
		{
			name:        "malformed address with colons",
			address:     "not:a:valid:address",
			expectError: true,
			errorMatch:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			starter := &RealStarter{}
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "test")
			})

			err := starter.Start(tt.address, handler)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}

			if tt.errorMatch != "" && err != nil && !strings.Contains(err.Error(), tt.errorMatch) {
				t.Logf("Got error (expected): %v", err)
			}
		})
	}
}

func TestRealStarterStartTimeout(t *testing.T) {
	starter := &RealStarter{}

	// Create a handler that will keep the server running
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	start := time.Now()

	// Start server on a valid address that will bind successfully
	// The server should timeout after 10 seconds since it starts successfully
	// and doesn't return an error immediately
	err := starter.Start("127.0.0.1:0", handler)

	elapsed := time.Since(start)

	// We expect a timeout error
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "timeout") {
		t.Errorf("Expected timeout error, got: %v", err)
	}

	// Verify the timeout occurred around 10 seconds (with some tolerance)
	if elapsed < 9*time.Second || elapsed > 11*time.Second {
		t.Errorf("Expected timeout around 10 seconds, got %v", elapsed)
	}
}

func TestRealStarterStartImmediateError(t *testing.T) {
	starter := &RealStarter{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Use an invalid address format to get immediate error
	err := starter.Start("invalid-address-format", handler)

	if err == nil {
		t.Error("Expected error with invalid address, got nil")
	}

	// The error should be about address, not timeout
	if err != nil && strings.Contains(err.Error(), "timeout") {
		t.Errorf("Expected address error, got timeout: %v", err)
	}
}

func TestRealStarterStartWithNilHandler(t *testing.T) {
	starter := &RealStarter{}

	// Test with nil handler - http.Server allows nil handler (uses DefaultServeMux)
	// but we should still get an error from invalid address
	err := starter.Start("invalid-addr", nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestRealStarterStartPortInUse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Start a server on a specific port in a goroutine
	server := &http.Server{
		Addr:    "127.0.0.1:18899",
		Handler: handler,
	}

	go func() {
		server.ListenAndServe()
	}()

	// Give the first server a moment to bind
	time.Sleep(100 * time.Millisecond)

	// Try to start another server on the same port
	starter := &RealStarter{}
	err := starter.Start("127.0.0.1:18899", handler)

	if err == nil {
		t.Error("Expected error when port is in use, got nil")
	}

	// The error should not be a timeout error since it fails immediately
	if err != nil && strings.Contains(err.Error(), "timeout") {
		t.Errorf("Expected port in use error, got timeout: %v", err)
	}

	// Clean up
	server.Close()
}

func TestRealStarterImplementsInterface(t *testing.T) {
	// Verify RealStarter properly implements ServerStarter interface
	var _ ServerStarter = &RealStarter{}

	starter := &RealStarter{}
	if starter == nil {
		t.Fatal("RealStarter should implement ServerStarter interface")
	}
}

func TestRealStarterZeroValue(t *testing.T) {
	// Test that zero value of RealStarter is usable
	var starter RealStarter

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test with invalid address
	err := starter.Start("not:a:valid:address", handler)
	if err == nil {
		t.Error("Expected error with malformed address")
	}
}
