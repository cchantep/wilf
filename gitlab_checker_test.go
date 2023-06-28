package main

import (
	"testing"
)

func TestGitlabChecker(t *testing.T) {
	checker := GitlabChecker{Config: GitlabRegistryConfig{
		ProjectApiPackagesUrl: "https://gitlab.com/api/v4/projects/46678122/packages",
		PrivateToken:          "",
	}}

	// Test existing project with required major update
	pkg := "promptlib"
	req := VersionRequirement{
		VersionConstraint{">=", "v2.26.0"},
	}
	expectedVersion := "v0.1.3"
	expectedLevel := Major
	expectedUrl := "https://gitlab.com/gitlab-org/modelops/applied-ml/code-suggestions/prompt-library/-/packages/16106475"
	version, level, url, err := checker.RequiredUpdate(pkg, req)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if version != expectedVersion {
		t.Errorf("Expected version %s, but got %s", expectedVersion, version)
	}

	if url != expectedUrl {
		t.Errorf("Expected url %s, but got %s", expectedUrl, url)
	}

	if level != expectedLevel {
		t.Errorf("Expected level %s, but got %v", expectedLevel, level)
	}
}

func TestGitlabCheckerNonExisting(t *testing.T) {
	checker := GitlabChecker{Config: GitlabRegistryConfig{
		ProjectApiPackagesUrl: "https://gitlab.com/api/v4/projects/46678122/packages",
		PrivateToken:          "",
	}}

	pkg := "non_existing_project"
	req := VersionRequirement{
		VersionConstraint{">=", "v1.0.0"},
	}
	version, level, url, err := checker.RequiredUpdate(pkg, req)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if version != "" {
		t.Errorf("Expected empty version, but got %s", version)
	}

	if url != "" {
		t.Errorf("Expected empty url, but got %s", url)
	}

	if level != 0 {
		t.Errorf("Expected level 0, but got %v", level)
	}
}
