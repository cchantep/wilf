package main

import (
	"reflect"
	"testing"
)

func TestGetProjectInfo(t *testing.T) {
	result, err := GetProjectInfo("requests")

	if err != nil {
		t.Errorf("Unexpected error occurred for 'requests': %v", err)
	}

	expected := &ProjectInfo{
		Name:           "requests",
		Version:        "v2.31.0",
		RequiresPython: ">=3.7",
		Summary:        "Python HTTP for Humans.",
		HomeURL:        "https://requests.readthedocs.io",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf(
			"Expected result for 'requests': %+v, got: %+v",
			expected, result)
	}

	// Test case: Get project info for an invalid package name
	result, err = GetProjectInfo("invalid-package-name")

	if err != nil || result != nil {
		t.Errorf("Expected result to be nil for 'invalid-package-name'")
	}
}

func TestIsValidVersion(t *testing.T) {
	testCases := []struct {
		version string
		want    bool
	}{
		{"v1.0.0", true},
		{"v1.0.0-alpha", true},
		{"v1.0.0-alpha.1", true},
		{"v1.0.0-0.3.7", true},
		{"v1.0.0-x.7.z.92", true},
		{"v1.0.0-rc.1+build.1", true},
		{"v1.0.0+0.3.7", true},
		{"v1.0.0-beta+exp.sha.5114f85", true},
		{"v1.5.5.1", true},
		{"v1.5.5.1-alpha", false},
		{"v1.5.5.1+build.1", false},
		{"v1.5.5.1-beta+exp.sha.5114f85", false},
	}

	for _, tc := range testCases {
		got := IsValidVersion(tc.version)
		if got != tc.want {
			t.Errorf("IsValidVersion(%q) = %v; want %v", tc.version, got, tc.want)
		}
	}
}
