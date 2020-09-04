package settings

import (
	ini "gopkg.in/ini.v1"
)

var (
	ConfigPath string
	cfg        *ini.File
)

// InitSettings initialize application settings
func InitSettings() error {
	cfg = ini.Empty()
	if err := cfg.Append(ConfigPath); err != nil {
		return err
	}
	return nil
}
