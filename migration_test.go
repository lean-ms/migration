package migration

import (
	"os"
	"testing"

	"github.com/lean-ms/database"
)

type MigrationRunTest struct {
	upCount   int
	downCount int
}

func (m *MigrationRunTest) upFn() error {
	m.upCount++
	return nil
}

func (m *MigrationRunTest) downFn() error {
	m.downCount++
	return nil
}

var dbConfigPath = "config/database.yml"
var baseVersion = 1234567890

var testCases = []struct {
	command           []string
	version           int
	expectedUpCount   int
	expectedDownCount int
	expectedVersion   int
	description       string
}{
	{
		command:           []string{"lean-ms", "migrate"},
		version:           baseVersion,
		description:       "First migration",
		expectedUpCount:   1,
		expectedDownCount: 0,
		expectedVersion:   baseVersion,
	},
	{
		command:           []string{"lean-ms", "migrate"},
		version:           baseVersion,
		description:       "Running same migration twice",
		expectedUpCount:   1,
		expectedDownCount: 0,
		expectedVersion:   baseVersion,
	},
	{
		command:           []string{"lean-ms", "migrate", "-rollback"},
		version:           baseVersion - 1,
		description:       "Rolling back with wrong version",
		expectedUpCount:   1,
		expectedDownCount: 0,
		expectedVersion:   baseVersion,
	},
	{
		command:           []string{"lean-ms", "migrate", "-rollback"},
		version:           baseVersion,
		description:       "Rolling back correctly",
		expectedUpCount:   1,
		expectedDownCount: 1,
		expectedVersion:   -1,
	},
	{
		command:           []string{"lean-ms", "migrate", "-rollback"},
		version:           baseVersion,
		description:       "Rolling back from initial state",
		expectedUpCount:   1,
		expectedDownCount: 1,
		expectedVersion:   -1,
	},
	{
		command:           []string{"lean-ms", "migrate"},
		version:           1111,
		description:       "Migrating forward again",
		expectedUpCount:   2,
		expectedDownCount: 1,
		expectedVersion:   1111,
	},
	{
		command:           []string{"lean-ms", "migrate"},
		version:           1110,
		description:       "Migrating with lower version than actual",
		expectedUpCount:   2,
		expectedDownCount: 1,
		expectedVersion:   1111,
	},
	{
		command:           []string{"lean-ms", "migrate"},
		version:           1115,
		description:       "Migrating up again correctly",
		expectedUpCount:   3,
		expectedDownCount: 1,
		expectedVersion:   1115,
	},
	{
		command:           []string{"lean-ms", "migrate", "-rollback"},
		version:           1115,
		description:       "Rolling back once",
		expectedUpCount:   3,
		expectedDownCount: 2,
		expectedVersion:   1111,
	},

	{
		command:           []string{"lean-ms", "migrate", "-rollback"},
		version:           1111,
		description:       "Rolling back twice",
		expectedUpCount:   3,
		expectedDownCount: 3,
		expectedVersion:   -1,
	},
}

func TestUpMigration(t *testing.T) {
	setupMigrationDB()
	dbConnection := database.CreateConnection(dbConfigPath)
	defer dbConnection.Close()
	m := &MigrationRunTest{}
	for _, testCase := range testCases {
		Run(m.upFn, m.downFn, testCase.version, testCase.command...)
		currentVersion := GetCurrentVersion(dbConnection)
		if m.upCount != testCase.expectedUpCount {
			t.Errorf("%s. Unexpected up count", testCase.description)
		} else if m.downCount != testCase.expectedDownCount {
			t.Errorf("%s. Unexpected down count", testCase.description)
		} else if currentVersion != testCase.expectedVersion {
			t.Errorf("%s. Unexpected version", testCase.description)
		}
	}
	tearDownMigrationDB()
}

func setupMigrationDB() {
	os.Setenv("LEANMS_ENV", "test")
	database.DropDatabase(dbConfigPath)
	database.CreateDatabase(dbConfigPath)
	database.CreateTable(dbConfigPath, &Migration{}, nil)
}

func tearDownMigrationDB() {
	database.DropDatabase(dbConfigPath)
}
