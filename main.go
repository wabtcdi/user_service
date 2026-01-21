package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/wabtcdi/user_service/cmd"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	configName := flag.String("config", "local", "config file to use (local or cloud)")
	flag.Parse()

	if err := run(*configName, gormOpener, &cmd.RealStarter{}); err != nil {
		log.Fatal(err)
	}
}

func gormOpener(dsn string) (*gorm.DB, error) {
	// Configure GORM logger for query monitoring
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Log queries slower than 200ms
			LogLevel:                  logger.Warn,            // Log slow queries and errors
			IgnoreRecordNotFoundError: true,                   // Don't log ErrRecordNotFound
			Colorful:                  true,                   // Enable color output
		},
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
}

func run(configName string, openFunc cmd.DBOpener, starter cmd.ServerStarter) error {
	err := os.Chdir("cmd")
	if err != nil {
		return err
	}
	return cmd.Init(configName, openFunc, starter)
}
