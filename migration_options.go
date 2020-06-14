package migration

import "fmt"

type MigrationOptions struct {
	IsRollback bool
	ConfigPath string
	Version    int
}

func (m *MigrationOptions) String() string {
	return fmt.Sprintf("{IsRollback: %v, Version: %v, ConfigPath: %v}", m.IsRollback, m.Version, m.ConfigPath)
}
