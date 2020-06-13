package migration

import (
	"log"

	"github.com/go-pg/pg/orm"
	"github.com/lean-ms/database"
)

// Migration model.
// Timestamp is used as an ID that avoids merging problems
type Migration struct {
	ID        int64
	Timestamp string
}

func setupMigrationTable(dbConnection *database.DbConnection) {
	model := (*Migration)(nil)
	dbConnection.Database.CreateTable(model, &orm.CreateTableOptions{
		IfNotExists: true,
	})
}

// GetCurrentVersion returns the Timestamp of latest migration
func GetCurrentVersion(dbConnection *database.DbConnection) string {
	var migrations []Migration
	dbConnection.Database.Model(&migrations).Order("id DESC").Limit(1).Select()
	if len(migrations) == 0 {
		return ""
	}
	return migrations[0].Timestamp
}

// SetCurrentVersion creates a new migration with a given stimestamp
func SetCurrentVersion(dbConnection *database.DbConnection, timestamp string) {
	err := dbConnection.Database.Insert(&Migration{Timestamp: timestamp})
	if err != nil {
		log.Fatal("Error in SetCurrentVersion: " + err.Error())
	}
}
