package migration

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/lean-ms/database"
)

var dbConfigPath = "config/database.yml"
var baseVersion = 1234567890

type MigrationTestCase struct {
	command           []string
	version           int
	isMigrationFnOk   bool
	expectedUpCount   int
	expectedDownCount int
	expectedVersion   int
	description       string
}

type MigrationTestMock struct {
	upCount   int
	downCount int
}

func (m *MigrationTestMock) upFn() error {
	m.upCount++
	return nil
}

func (m *MigrationTestMock) downFn() error {
	m.downCount++
	return nil
}

var testCases = []MigrationTestCase{
	{
		command:           []string{"lean-ms", "migrate", "-config=config/database.yml"},
		version:           baseVersion,
		isMigrationFnOk:   false,
		description:       "First migration (with problems)",
		expectedUpCount:   0,
		expectedDownCount: 0,
		expectedVersion:   -1,
	},
	{
		command:           []string{"lean-ms", "migrate", "-config=config/database.yml"},
		version:           baseVersion,
		isMigrationFnOk:   true,
		description:       "First migration",
		expectedUpCount:   1,
		expectedDownCount: 0,
		expectedVersion:   baseVersion,
	},
	{
		command:           []string{"lean-ms", "migrate", "-config=config/database.yml"},
		version:           baseVersion,
		isMigrationFnOk:   true,
		description:       "Running same migration twice",
		expectedUpCount:   1,
		expectedDownCount: 0,
		expectedVersion:   baseVersion,
	},
	{
		command:           []string{"lean-ms", "migrate", "-rollback", "-config=config/database.yml"},
		version:           baseVersion - 1,
		isMigrationFnOk:   true,
		description:       "Rolling back with wrong version",
		expectedUpCount:   1,
		expectedDownCount: 0,
		expectedVersion:   baseVersion,
	},
	{
		command:           []string{"lean-ms", "migrate", "-rollback", "-config=config/database.yml"},
		version:           baseVersion,
		isMigrationFnOk:   false,
		description:       "Rolling back with errors",
		expectedUpCount:   1,
		expectedDownCount: 0,
		expectedVersion:   baseVersion,
	},
	{
		command:           []string{"lean-ms", "migrate", "-rollback", "-config=config/database.yml"},
		version:           baseVersion,
		isMigrationFnOk:   true,
		description:       "Rolling back correctly",
		expectedUpCount:   1,
		expectedDownCount: 1,
		expectedVersion:   -1,
	},
	{
		command:           []string{"lean-ms", "migrate", "-rollback", "-config=config/database.yml"},
		version:           baseVersion,
		isMigrationFnOk:   true,
		description:       "Rolling back from initial state",
		expectedUpCount:   1,
		expectedDownCount: 1,
		expectedVersion:   -1,
	},
	{
		command:           []string{"lean-ms", "migrate", "-config=config/database.yml"},
		version:           1111,
		isMigrationFnOk:   true,
		description:       "Migrating forward again",
		expectedUpCount:   2,
		expectedDownCount: 1,
		expectedVersion:   1111,
	},
	{
		command:           []string{"lean-ms", "migrate", "-config=config/database.yml"},
		version:           1110,
		isMigrationFnOk:   true,
		description:       "Migrating with lower version than actual",
		expectedUpCount:   2,
		expectedDownCount: 1,
		expectedVersion:   1111,
	},
	{
		command:           []string{"lean-ms", "migrate", "-config=config/database.yml"},
		version:           1115,
		isMigrationFnOk:   true,
		description:       "Migrating up again correctly",
		expectedUpCount:   3,
		expectedDownCount: 1,
		expectedVersion:   1115,
	},
	{
		command:           []string{"lean-ms", "migrate", "-rollback", "-config=config/database.yml"},
		version:           1115,
		isMigrationFnOk:   true,
		description:       "Rolling back once",
		expectedUpCount:   3,
		expectedDownCount: 2,
		expectedVersion:   1111,
	},

	{
		command:           []string{"lean-ms", "migrate", "-rollback", "-config=config/database.yml"},
		version:           1111,
		isMigrationFnOk:   true,
		description:       "Rolling back twice",
		expectedUpCount:   3,
		expectedDownCount: 3,
		expectedVersion:   -1,
	},
}

func problematicFn() error {
	return errors.New("injected error")
}

func TestSuccessfulMigrations(t *testing.T) {
	setupMigrationDB()
	m := &MigrationTestMock{}
	for _, testCase := range testCases {
		if testCase.isMigrationFnOk {
			Run(m.upFn, m.downFn, testCase.version, testCase.command...)
		} else {
			Run(problematicFn, problematicFn, testCase.version, testCase.command...)
		}
		if err := checkErrors(testCase, *m); err != nil {
			t.Error(err)
		}
	}
	tearDownMigrationDB()
}

func checkErrors(testCase MigrationTestCase, m MigrationTestMock) error {
	currentVersion := GetCurrentVersion(dbConfigPath)
	if m.upCount != testCase.expectedUpCount {
		return errors.New(fmt.Sprintf("%s. Unexpected up count", testCase.description))
	} else if m.downCount != testCase.expectedDownCount {
		return errors.New(fmt.Sprintf("%s. Unexpected down count", testCase.description))
	} else if currentVersion != testCase.expectedVersion {
		return errors.New(fmt.Sprintf("%s. Unexpected version", testCase.description))
	}
	return nil
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
