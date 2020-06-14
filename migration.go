package migration

import (
	"fmt"
	"log"

	"github.com/lean-ms/database"
	"github.com/lean-ms/database/models"
)

// Migration model.
// Version is used as an ID. Timestamp version is recomended to avoid merging problems
type Migration struct {
	ID      int64
	Version int
}

// GetCurrentVersion returns the Timestamp of latest migration
func GetCurrentVersion(dbConfigPath string) int {
	dbConnection := database.CreateConnection(dbConfigPath)
	defer dbConnection.Close()
	var migrations []Migration
	dbConnection.Database.Model(&migrations).Order("id DESC").Limit(1).Select()
	if len(migrations) == 0 {
		return -1
	}
	return migrations[0].Version
}

// SetCurrentVersion creates a new migration with a given a version number
func SetCurrentVersion(dbConfigPath string, version int) error {
	database.CreateTable(dbConfigPath, new(Migration), &models.CreateTableOptions{
		IfNotExists: true,
	})
	dbConnection := database.CreateConnection(dbConfigPath)
	defer dbConnection.Close()
	return dbConnection.Database.Insert(&Migration{Version: version})
}

// RollbackVersion removes last version
func RollbackVersion(dbConfigPath string) error {
	dbConnection := database.CreateConnection(dbConfigPath)
	defer dbConnection.Close()
	migration := new(Migration)
	err := dbConnection.Database.Model(migration).Last()
	if err != nil {
		return err
	}
	return dbConnection.Database.Delete(migration)
}

type migrateFn func() error

// func Timestamp() string {
// 	timestamp := time.Now().Format(time.RFC3339)
// 	a, _ := regexp.Compile("\\..*")
// 	timestamp = a.ReplaceAllLiteralString(timestamp, "")
// 	a, _ = regexp.Compile("[^\\d]")
// 	timestamp = a.ReplaceAllLiteralString(timestamp, "")
// 	return timestamp[:14]
// }

func printableVersion(version int) string {
	if version < 0 {
		return "Empty"
	}
	return fmt.Sprintf("%v", version)
}

func Run(upFn migrateFn, downFn migrateFn, opts *Options) {
	log.Printf("Starting migration. Options: %s\n", opts.String())
	currentVersion := GetCurrentVersion(opts.ConfigPath)
	log.Printf("Current version is %s\n", printableVersion(currentVersion))
	if opts.IsRollback && runRollback(opts.Version, currentVersion, downFn) {
		RollbackVersion(opts.ConfigPath)
	} else if !opts.IsRollback && runForward(opts.Version, currentVersion, upFn) {
		SetCurrentVersion(opts.ConfigPath, opts.Version)
	}
	currentVersion = GetCurrentVersion(opts.ConfigPath)
	log.Printf("Finished. Version now is %s\n", printableVersion(currentVersion))
}

func runForward(version int, currentVersion int, upFn migrateFn) bool {
	if version <= currentVersion {
		log.Printf("Doing nothing. New version %v is not higher than %v", version, currentVersion)
		return false
	}
	err := upFn()
	if err != nil {
		log.Printf("Error running forward migration: %s", err.Error())
		return false
	}
	return true
}

func runRollback(version int, currentVersion int, downFn migrateFn) bool {
	if currentVersion < 0 {
		log.Println("Doing nothing. Cannot rollback from empty state")
		return false
	}
	if version != currentVersion {
		log.Printf("Doing nothing. Version to rollback is %d and requested was %d", currentVersion, version)
		return false
	}
	err := downFn()
	if err != nil {
		log.Printf("Error rolling back: %s", err.Error())
		return false
	}
	return true
}
