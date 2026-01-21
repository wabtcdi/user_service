package health

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Checker struct {
	DB *gorm.DB
}

func (c *Checker) Check(w http.ResponseWriter, r *http.Request) {
	sqlDB, err := c.DB.DB()
	if err != nil {
		logrus.Warn("Failed to get underlying database")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Not Ready")
		return
	}

	if err := sqlDB.Ping(); err != nil {
		logrus.Warn("Database not ready")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Not Ready")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ready")
}
