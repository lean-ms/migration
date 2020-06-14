package migration

import (
	"os"
	"os/exec"
	"testing"
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
