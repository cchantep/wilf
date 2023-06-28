package main

import (
	"github.com/BurntSushi/toml"
)

type Settings struct {
	CheckDevPackages bool     `toml:"check_dev_packages"`
	ExcludedPackages []string `toml:"excluded_packages"`
	UpdateLevel      UpdateLevel
	UpdateLevelRepr  string `toml:"update_level"`
}

type Config struct {
	Settings *Settings
	Gitlab   *GitlabRegistryConfig
}

func DefaultSettings() Settings {
	return Settings{
		CheckDevPackages: false,
		ExcludedPackages: []string{},
		UpdateLevel:      Minor,
		UpdateLevelRepr:  "",
	}
}

// LoadSettings loads a Settings instance from a TOML file specified as path in arguments.
// It returns either any encountered error, or the successfully loaded Settings.
// LoadSettings loads a Settings instance from a TOML file specified as path in arguments.
// It returns either any encountered error, or the successfully loaded Settings.
//
// If the TOML file does not contain an `update_level` field, the `UpdateLevel` field of the returned
// `Settings` instance will be set to `Minor` by default.
func LoadSettings(path string) (*Settings, error) {
	var settings Settings

	if _, err := toml.DecodeFile(path, &settings); err != nil {
		return nil, err
	}

	if settings.UpdateLevelRepr != "" {
		level, err := ParseUpdateLevel(settings.UpdateLevelRepr)

		if err != nil {
			return nil, err
		}

		settings.UpdateLevel = level
	} else {
		settings.UpdateLevel = DefaultSettings().UpdateLevel
	}

	return &settings, nil
}

// LoadConfig loads a Config instance from a TOML file specified as path in arguments.
// It returns either any encountered error, or the successfully loaded Config.
func LoadConfig(path string) (*Config, error) {
	// Load settings from the TOML file
	settings, err := LoadSettings(path)

	if err != nil {
		return nil, err
	}

	gitlabConfig, err := LoadGitlabRegistryConfig(path)

	if err != nil {
		return nil, err
	}

	if gitlabConfig.ProjectApiPackagesUrl == "" {
		gitlabConfig = nil
	}

	// Return the Config instance with the loaded settings and an empty GitlabRegistryConfig
	return &Config{
		Settings: settings,
		Gitlab:   gitlabConfig,
	}, nil
}
