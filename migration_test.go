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
	options           *MigrationOptions
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
		options:           &MigrationOptions{IsRollback: false, Version: baseVersion, ConfigPath: "config/database.yml"},
		isMigrationFnOk:   false,
		description:       "First migration (with problems)",
		expectedUpCount:   0,
		expectedDownCount: 0,
		expectedVersion:   -1,
	},
	{
		options:           &MigrationOptions{IsRollback: false, Version: baseVersion, ConfigPath: "config/database.yml"},
		isMigrationFnOk:   true,
		description:       "First migration",
		expectedUpCount:   1,
		expectedDownCount: 0,
		expectedVersion:   baseVersion,
	},
	{
		options:           &MigrationOptions{IsRollback: false, Version: baseVersion, ConfigPath: "config/database.yml"},
		isMigrationFnOk:   true,
		description:       "Running same migration twice",
		expectedUpCount:   1,
		expectedDownCount: 0,
		expectedVersion:   baseVersion,
	},
	{
		options:           &MigrationOptions{IsRollback: true, Version: baseVersion - 1, ConfigPath: "config/database.yml"},
		isMigrationFnOk:   true,
		description:       "Rolling back with wrong version",
		expectedUpCount:   1,
		expectedDownCount: 0,
		expectedVersion:   baseVersion,
	},
	{
		options:           &MigrationOptions{IsRollback: true, Version: baseVersion, ConfigPath: "config/database.yml"},
		isMigrationFnOk:   false,
		description:       "Rolling back with errors",
		expectedUpCount:   1,
		expectedDownCount: 0,
		expectedVersion:   baseVersion,
	},
	{
		options:           &MigrationOptions{IsRollback: true, Version: baseVersion, ConfigPath: "config/database.yml"},
		isMigrationFnOk:   true,
		description:       "Rolling back correctly",
		expectedUpCount:   1,
		expectedDownCount: 1,
		expectedVersion:   -1,
	},
	{
		options:           &MigrationOptions{IsRollback: true, Version: baseVersion, ConfigPath: "config/database.yml"},
		isMigrationFnOk:   true,
		description:       "Rolling back from initial state",
		expectedUpCount:   1,
		expectedDownCount: 1,
		expectedVersion:   -1,
	},
	{
		options:           &MigrationOptions{IsRollback: false, Version: 1111, ConfigPath: "config/database.yml"},
		isMigrationFnOk:   true,
		description:       "Migrating forward again",
		expectedUpCount:   2,
		expectedDownCount: 1,
		expectedVersion:   1111,
	},
	{
		options:           &MigrationOptions{IsRollback: false, Version: 1110, ConfigPath: "config/database.yml"},
		isMigrationFnOk:   true,
		description:       "Migrating with lower version than actual",
		expectedUpCount:   2,
		expectedDownCount: 1,
		expectedVersion:   1111,
	},
	{
		options:           &MigrationOptions{IsRollback: false, Version: 1115, ConfigPath: "config/database.yml"},
		isMigrationFnOk:   true,
		description:       "Migrating up again correctly",
		expectedUpCount:   3,
		expectedDownCount: 1,
		expectedVersion:   1115,
	},
	{
		options:           &MigrationOptions{IsRollback: true, Version: 1115, ConfigPath: "config/database.yml"},
		isMigrationFnOk:   true,
		description:       "Rolling back once",
		expectedUpCount:   3,
		expectedDownCount: 2,
		expectedVersion:   1111,
	},

	{
		options:           &MigrationOptions{IsRollback: true, Version: 1111, ConfigPath: "config/database.yml"},
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
			Run(m.upFn, m.downFn, testCase.options)
		} else {
			Run(problematicFn, problematicFn, testCase.options)
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
