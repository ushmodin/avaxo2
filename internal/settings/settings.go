package settings

import (
	ini "gopkg.in/ini.v1"
)

var (
	// AgentSettings config for agent mode
	AgentSettings = struct {
		Listen   string `ini:"listen"`
		Keyfile  string `ini:"keyfile"`
		Certfile string `ini:"certfile"`
	}{
		Listen:   "0.0.0.0:8843",
		Keyfile:  "avaxo2.key",
		Certfile: "avaxo2.crt",
	}
	// ConfigPath path to config file
	ConfigPath string
)

// InitSettings initialize application settings
func InitSettings() error {
	cfg := ini.Empty()
	if err := cfg.Append(ConfigPath); err != nil {
		return err
	}
	if err := readVars(cfg); err != nil {
		return err
	}
	return nil
}

func readVars(cfg *ini.File) error {
	if err := cfg.Section("agent").MapTo(&AgentSettings); err != nil {
		return err
	}
	return nil
}
