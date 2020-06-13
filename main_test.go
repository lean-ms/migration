package migration

import (
	"os"
	"os/exec"
	"testing"

	"github.com/lean-ms/database"
)

var goPath, _ = exec.LookPath("go")

func TestRunSuccessMigration(t *testing.T) {
	cmd := exec.Command(goPath, "run", "main.go", "-path=/tmp/lean-ms/test")
	cmd.Dir = "./examples/create_file"
	cmd.CombinedOutput()
	_, err := os.Stat("/tmp/lean-ms/test/migration.ok")
	if os.IsNotExist(err) {
		t.Error("Could not run migration")
	}
}

func TestRunFailedMigration(t *testing.T) {
	cmd := exec.Command(goPath, "run", "main.go", "-path", "/tmp/lean-ms/test")
	cmd.Dir = "./examples/failed_migration"
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Output should be fail")
	}
}

var dbConfigPath = "config/database.yml"

func TestCheckCurrentMigrationVersion(t *testing.T) {
	setup()
	dbConnection := database.CreateConnection(dbConfigPath)
	defer dbConnection.Close()
	if timestamp := GetCurrentVersion(dbConnection); len(timestamp) > 0 {
		t.Error("First migration timestamp should be an empty string")
	}
	SetCurrentVersion(dbConnection, "1234567890")
	if "1234567890" != GetCurrentVersion(dbConnection) {
		t.Error("Migration timestamp was not properly set")
	}
	tearDown()
}

func setup() {
	os.Setenv("LEANMS_ENV", "test")
	database.DropDatabase(dbConfigPath)
	database.CreateDatabase(dbConfigPath)
	database.CreateTable(dbConfigPath, &Migration{}, nil)
}

func tearDown() {
	database.DropDatabase(dbConfigPath)
}
