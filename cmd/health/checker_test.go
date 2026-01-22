package health

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

func TestChecker_Check_Success(t *testing.T) {
	// Use sqlmock with successful ping
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer mockDB.Close()

	// Expect ping to succeed
	mock.ExpectPing()

	// Wrap the sql.DB in a gorm.DB
	dialector := &mockDialector{sqlDB: mockDB}
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		DisableAutomaticPing:   true, // Don't ping during Open
	})
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}

	checker := &Checker{DB: gormDB}

	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	checker.Check(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Ready"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChecker_Check_PingError(t *testing.T) {
	// Use sqlmock to simulate a ping error
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer mockDB.Close()

	// Expect ping to fail
	mock.ExpectPing().WillReturnError(errors.New("ping failed"))

	// Wrap the sql.DB in a gorm.DB
	dialector := &mockDialector{sqlDB: mockDB}
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		DisableAutomaticPing:   true, // Don't ping during Open
	})
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}

	checker := &Checker{DB: gormDB}

	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	checker.Check(rr, req)

	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusServiceUnavailable)
	}

	expected := "Not Ready"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChecker_Check_MultipleRequests(t *testing.T) {
	// Use sqlmock with multiple successful pings
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer mockDB.Close()

	// Expect multiple pings
	for i := 0; i < 5; i++ {
		mock.ExpectPing()
	}

	// Wrap the sql.DB in a gorm.DB
	dialector := &mockDialector{sqlDB: mockDB}
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		DisableAutomaticPing:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}

	checker := &Checker{DB: gormDB}

	// Make multiple requests to ensure consistency
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", "/ready", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		checker.Check(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("request %d: handler returned wrong status code: got %v want %v",
				i, status, http.StatusOK)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChecker_Check_AlternatingResults(t *testing.T) {
	// Test alternating success and failure
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer mockDB.Close()

	// Expect alternating ping results
	mock.ExpectPing()
	mock.ExpectPing().WillReturnError(errors.New("temporary failure"))
	mock.ExpectPing()

	// Wrap the sql.DB in a gorm.DB
	dialector := &mockDialector{sqlDB: mockDB}
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		DisableAutomaticPing:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}

	checker := &Checker{DB: gormDB}

	// First request should succeed
	req1, _ := http.NewRequest("GET", "/ready", nil)
	rr1 := httptest.NewRecorder()
	checker.Check(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Errorf("First request: expected 200, got %d", rr1.Code)
	}

	// Second request should fail
	req2, _ := http.NewRequest("GET", "/ready", nil)
	rr2 := httptest.NewRecorder()
	checker.Check(rr2, req2)
	if rr2.Code != http.StatusServiceUnavailable {
		t.Errorf("Second request: expected 503, got %d", rr2.Code)
	}

	// Third request should succeed again
	req3, _ := http.NewRequest("GET", "/ready", nil)
	rr3 := httptest.NewRecorder()
	checker.Check(rr3, req3)
	if rr3.Code != http.StatusOK {
		t.Errorf("Third request: expected 200, got %d", rr3.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// mockDialector is a custom GORM dialector for testing with sqlmock
type mockDialector struct {
	sqlDB *sql.DB
}

func (d *mockDialector) Name() string {
	return "sqlmock"
}

func (d *mockDialector) Initialize(db *gorm.DB) error {
	db.ConnPool = d.sqlDB
	return nil
}

func (d *mockDialector) Migrator(db *gorm.DB) gorm.Migrator {
	return nil
}

func (d *mockDialector) DataTypeOf(*schema.Field) string {
	return ""
}

func (d *mockDialector) DefaultValueOf(*schema.Field) clause.Expression {
	return clause.Expr{}
}

func (d *mockDialector) BindVarTo(clause.Writer, *gorm.Statement, interface{}) {
}

func (d *mockDialector) QuoteTo(clause.Writer, string) {
}

func (d *mockDialector) Explain(sql string, vars ...interface{}) string {
	return sql
}
