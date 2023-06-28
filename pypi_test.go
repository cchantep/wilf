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
		Name:    "requests",
		Version: "v2.31.0",
		Summary: "Python HTTP for Humans.",
		HomeURL: "https://requests.readthedocs.io",
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
