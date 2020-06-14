package migration

import (
	"flag"
	"log"
	"regexp"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/lean-ms/database"
)

// Migration model.
// Version is used as an ID. Timestamp is recomended to avoid merging problems
type Migration struct {
	ID      int64
	Version int
}

func setupMigrationTable(dbConnection *database.DbConnection) {
	model := (*Migration)(nil)
	dbConnection.Database.CreateTable(model, &orm.CreateTableOptions{
		IfNotExists: true,
	})
}

// GetCurrentVersion returns the Timestamp of latest migration
func GetCurrentVersion(dbConnection *database.DbConnection) int {
	var migrations []Migration
	dbConnection.Database.Model(&migrations).Order("id DESC").Limit(1).Select()
	if len(migrations) == 0 {
		return -1
	}
	return migrations[0].Version
}

// SetCurrentVersion creates a new migration with a given a version number
func SetCurrentVersion(dbConnection *database.DbConnection, version int) error {
	err := dbConnection.Database.Insert(&Migration{Version: version})
	return err
}

// RollbackVersion removes last version
func RollbackVersion(dbConnection *database.DbConnection) error {
	migration := new(Migration)
	err := dbConnection.Database.Model(migration).Last()
	if err != nil {
		return err
	}
	dbConnection.Database.Delete(migration)
	return nil
}

type migrateFn func() error

func Timestamp() string {
	timestamp := time.Now().Format(time.RFC3339)
	a, _ := regexp.Compile("\\..*")
	timestamp = a.ReplaceAllLiteralString(timestamp, "")
	a, _ = regexp.Compile("[^\\d]")
	timestamp = a.ReplaceAllLiteralString(timestamp, "")
	return timestamp[:14]
}

func Run(upFn migrateFn, downFn migrateFn, version int, args ...string) {
	isRollback := getIsRollbackFromCli(args...)
	dbConnection := database.CreateConnection(dbConfigPath)
	defer dbConnection.Close()
	currentVersion := GetCurrentVersion(dbConnection)
	if isRollback && runRollback(version, currentVersion, downFn) {
		RollbackVersion(dbConnection)
	} else if !isRollback && runForward(version, currentVersion, upFn) {
		SetCurrentVersion(dbConnection, version)
	}
	log.Println("Finished")
}

func getIsRollbackFromCli(args ...string) bool {
	cmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	isRollback := cmd.Bool("rollback", false, "migrate one version behind")
	cmd.Parse(args[2:])
	return *isRollback
}

func runForward(version int, currentVersion int, upFn migrateFn) bool {
	log.Println("Starting migration...")
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
	log.Println("Starting rollback...")
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
