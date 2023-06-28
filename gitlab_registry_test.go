package main

import (
	"testing"

	"golang.org/x/mod/semver"
)

func TestGetPromptlibProjectInfo(t *testing.T) {
	// Create a GitlabRegistryConfig struct with the URL
	// and token for a public project
	config := GitlabRegistryConfig{
		ProjectApiPackagesUrl: "https://gitlab.com/api/v4/projects/46678122/packages",
		PrivateToken:          "",
	}

	// Test a package that exists in the project
	projectInfo, err := GetGitlabProjectInfo(config, "promptlib")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if projectInfo == nil {
		t.Errorf("Expected projectInfo to be non-nil")
	}

	if projectInfo.Name != "promptlib" {
		t.Errorf("Expected projectInfo.Name to be 'promptlib', got '%s'", projectInfo.Name)
	}

	if semver.Compare(projectInfo.Version, "v0.1.3") < 0 {
		t.Errorf("Expected projectInfo.Version to be >= '0.1.3', got '%s'", projectInfo.Version)
	}

	// Test a package that does not exist in the project
	projectInfo, err = GetGitlabProjectInfo(config, "non-existent-package")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if projectInfo != nil {
		t.Errorf("Expected projectInfo to be nil")
	}
}

func TestGitlabGetNonExistentPackageInfo(t *testing.T) {
	// Create a GitlabRegistryConfig struct with the URL
	// and token for a public project
	config := GitlabRegistryConfig{
		ProjectApiPackagesUrl: "https://gitlab.com/api/v4/projects/46678122/packages",
		PrivateToken:          "",
	}

	// Test a package that does not exist in the project
	projectInfo, err := GetGitlabProjectInfo(config, "non-existent-package")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if projectInfo != nil {
		t.Errorf("Expected projectInfo to be nil")
	}
}

func TestGitlabGetPackageInfoUnauthorized(t *testing.T) {
	// Create a GitlabRegistryConfig struct with the URL
	// and token for a public project
	config := GitlabRegistryConfig{
		ProjectApiPackagesUrl: "https://gitlab.com/api/v4/projects/46678122/packages",
		PrivateToken:          "NOT_AUTHORIZED",
	}

	_, err := GetGitlabProjectInfo(config, "promptlib")

	if err == nil || err.Error() != "Project information not found in the JSON response: 401 Unauthorized" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestLoadValidGitlabRegistryConfig(t *testing.T) {
	configFilePath := "resources/valid-gitlab-config.toml"
	expectedConfig := GitlabRegistryConfig{
		ProjectApiPackagesUrl: "https://gitlab.com/api/v4/projects/12345678/packages",
		PrivateToken:          "YOUR_PRIVATE_TOKEN",
	}

	config, err := LoadGitlabRegistryConfig(configFilePath)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if config.ProjectApiPackagesUrl != expectedConfig.ProjectApiPackagesUrl {
		t.Errorf("Expected ProjectApiPackagesUrl to be '%s', got '%s'", expectedConfig.ProjectApiPackagesUrl, config.ProjectApiPackagesUrl)
	}

	if config.PrivateToken != expectedConfig.PrivateToken {
		t.Errorf("Expected PrivateToken to be '%s', got '%s'", expectedConfig.PrivateToken, config.PrivateToken)
	}
}
