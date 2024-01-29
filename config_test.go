package main

import (
	"testing"
)

func TestLoadSettings(t *testing.T) {
	// Load the settings from the fixture file
	settings, err := LoadSettings("resources/valid-settings.toml")

	if err != nil {
		t.Fatalf("Failed to load settings: %v", err)
	}

	// Verify the loaded settings match the expected values
	if !settings.CheckDevPackages {
		t.Errorf("Expected CheckDevPackages to be true, but got false")
	}

	if len(settings.ExcludedPackages) != 2 || settings.ExcludedPackages[0] != "pkg1" || settings.ExcludedPackages[1] != "pkg2" {
		t.Errorf("Expected ExcludedPackages to be [\"pkg1\", \"pkg2\"], but got %v", settings.ExcludedPackages)
	}

	if settings.UpdateLevel != Major {
		t.Errorf("Expected UpdateLevel to be Major, but got %v", settings.UpdateLevel)
	}
}

func TestLoadSettingsOnlyConfig(t *testing.T) {
	path := "resources/valid-settings.toml"
	config, err := LoadConfig(path)

	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}

	if config.Settings == nil {
		t.Errorf("Expected settings to be loaded, but got nil")
	}

	if config.Gitlab != nil {
		t.Errorf("Expected Gitlab to be nil, but got %v", config.Gitlab)
	}

	settings := config.Settings

	// Verify the loaded settings match the expected values
	if !settings.CheckDevPackages {
		t.Errorf("Expected CheckDevPackages to be true, but got false")
	}

	if len(settings.ExcludedPackages) != 2 || settings.ExcludedPackages[0] != "pkg1" || settings.ExcludedPackages[1] != "pkg2" {
		t.Errorf("Expected ExcludedPackages to be [\"pkg1\", \"pkg2\"], but got %v", settings.ExcludedPackages)
	}

	if settings.UpdateLevel != Major {
		t.Errorf("Expected UpdateLevel to be Major, but got %v", settings.UpdateLevel)
	}
}

func TestLoadConfigWithGitlabRegistry(t *testing.T) {
	path := "resources/valid-gitlab-config.toml"
	config, err := LoadConfig(path)

	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}

	if config.Settings == nil {
		t.Errorf("Expected settings to be loaded, but got nil")
	}

	if config.Gitlab == nil {
		t.Errorf("Expected Gitlab registry config to be loaded, but got %v", config.Gitlab)
	}

	settings := config.Settings

	// Verify the loaded settings match the expected values
	if settings.CheckDevPackages {
		t.Errorf("Expected CheckDevPackages to be false, but got true")
	}

	if len(settings.ExcludedPackages) != 0 {
		t.Errorf("Expected empty ExcludedPackages, but got %v", settings.ExcludedPackages)
	}

	if settings.UpdateLevel != Minor {
		t.Errorf("Expected UpdateLevel to be Minor, but got %v", settings.UpdateLevel)
	}
}
