package main

import (
	"os"
	"time"

	"github.com/lean-ms/migration"
)

type User struct {
	ID        int64
	Name      string
	Emails    []string
	UpdatedAt time.Time
	CreatedAt time.Time
}

var config = "../../config/database.yml"

func Up() error {
	if err := database.CreateDabase(config); err != nil {
		return err
	}
	if err := database.CreateTable(config, new(User), nil); err != nil {
		return err
	}
	return nil
}

func Down() error {
	if err := database.DropTable(config, new(User), nil); err != nil {
		return err
	}
	if err := database.DropDabase(config); err != nil {
		return err
	}
	return nil
}

func main() {
	version := 20200613113048
	migration.Run(Up, Down, version, os.Args...)
}
