package main

import (
	"fmt"
	"os"

	"github.com/lean-ms/migration"
)

func Up() error {
	fmt.Println("running up migration")
	return nil
}

func Down() error {
	fmt.Println("running down migration")
	return nil
}

func main() {
	version := 20200613113048
	migration.Run(Up, Down, version, os.Args...)
}
