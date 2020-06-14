package migration

import "fmt"

// Options model exposes migration parameters
type Options struct {
	IsRollback bool
	ConfigPath string
	Version    int
}

func (m *Options) String() string {
	return fmt.Sprintf("{IsRollback: %v, Version: %v, ConfigPath: %v}", m.IsRollback, m.Version, m.ConfigPath)
}
