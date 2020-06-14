package migration

import (
	"os"
	"testing"

	"github.com/lean-ms/database"
)

func TestConfigFileExists(t *testing.T) {
	_, err := os.Stat("config/database.yml")
	if os.IsNotExist(err) {
		t.Error("Database config file was not found")
	}
}

func TestDatabaseConnection(t *testing.T) {
	config := "config/database.yml"
	database.CreateDatabase(config)
	dbConnection := database.CreateConnection(config)
	defer dbConnection.Close()
	_, err := dbConnection.Database.Exec("SELECT 1")
	if err != nil {
		t.Errorf("Could not connect to database: %s", err.Error())
	}
	database.DropDatabase(config)
}
