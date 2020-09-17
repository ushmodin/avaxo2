package settings

import (
	"errors"

	ini "gopkg.in/ini.v1"
)

var (
	cfg *ini.File
	// MinionSettings config for minion mode
	MinionSettings = struct {
		Listen   string `ini:"listen"`
		Keyfile  string `ini:"keyfile"`
		Certfile string `ini:"certfile"`
		Cafile   string `ini:"cafile"`
	}{
		Listen:   "0.0.0.0:8843",
		Keyfile:  "avaxo2.key",
		Certfile: "avaxo2.crt",
		Cafile:   "ca.crt",
	}

	// GruSettings config for gru mode
	GruSettings = struct {
		Keyfile  string `ini:"keyfile"`
		Certfile string `ini:"certfile"`
		Cafile   string `ini:"cafile"`
	}{
		Keyfile:  "avaxo2.key",
		Certfile: "avaxo2.crt",
		Cafile:   "ca.crt",
	}
	// ConfigPath path to config file
	ConfigPath string
)

// MinionAddress info about minion
type MinionAddress struct {
	Host string `ini:"host"`
}

// InitSettings initialize application settings
func InitSettings() error {
	cfg = ini.Empty()
	if err := cfg.Append(ConfigPath); err != nil {
		return err
	}
	if err := readVars(cfg); err != nil {
		return err
	}
	return nil
}

func readVars(cfg *ini.File) error {
	if err := cfg.Section("minion").MapTo(&MinionSettings); err != nil {
		return err
	}

	if err := cfg.Section("gru").MapTo(&GruSettings); err != nil {
		return err
	}
	return nil
}

// GetMinionAddress read minion settings for gru
func GetMinionAddress(name string) (MinionAddress, error) {
	var ma MinionAddress
	if err := cfg.Section(name).MapTo(&ma); err != nil {
		return MinionAddress{}, err
	}
	if ma.Host == "" {
		return MinionAddress{}, errors.New("minion settings not found")
	}
	return ma, nil
}
