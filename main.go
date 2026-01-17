package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/wabtcdi/user_service/cmd"
)

func main() {
	configName := flag.String("config", "local", "config file to use (local or cloud)")
	flag.Parse()

	if err := run(*configName, sql.Open, &cmd.RealStarter{}); err != nil {
		log.Fatal(err)
	}
}

func run(configName string, openFunc cmd.DBOpener, starter cmd.ServerStarter) error {
	err := os.Chdir("cmd")
	if err != nil {
		return err
	}
	return cmd.Init(configName, openFunc, starter)
}
