package health

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Checker struct {
	DB *sql.DB
}

func (c *Checker) Check(w http.ResponseWriter, r *http.Request) {
	if err := c.DB.Ping(); err != nil {
		logrus.Warn("Database not ready")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Not Ready")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ready")
}
