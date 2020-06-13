package other

import (
	"os"
	"os/exec"
	"testing"

	"github.com/go-pg/pg/v10/orm"
	"github.com/lean-ms/database"
)

type Migration struct {
	ID        int64
	Timestamp string
}

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

func getLastMigrationTimestamp(dbConnection *database.DbConnection) string {
	var migrations []Migration
	dbConnection.Database.Model(&migrations).Order("id DESC").Limit(1).Select()
	if len(migrations) == 0 {
		return ""
	}
	return migrations[0].Timestamp
}

func setupMigrationTable(dbConnection *database.DbConnection) {
	model := (*Migration)(nil)
	dbConnection.Database.CreateTable(model, &orm.CreateTableOptions{
		IfNotExists: true,
	})
}

func TestCheckCurrentMigrationVersion(t *testing.T) {
	dbConnection := database.CreateConnection("config/database.yml")
	defer dbConnection.Close()
	timestamp := getLastMigrationTimestamp(dbConnection)
	if len(timestamp) > 0 {
		t.Error("Should not find any migration")
	}
}
