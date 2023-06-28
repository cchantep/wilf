package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestTextReporter(t *testing.T) {
	testCases := []struct {
		name           string
		reporter       TextReporter
		packageName    string
		requirement    VersionRequirement
		latestVersion  string
		updateLevel    UpdateLevel
		dependencyKind DependencyKind
		packageUrl     string
		expectedOutput string
	}{
		{
			name: "Test case 1",
			reporter: TextReporter{
				MessageBefore: "",
				Pattern:       "%s %s %s %s %s %.0f %s",
				MessageAfter:  "",
			},
			packageName:    "github.com/test/package",
			requirement:    VersionRequirement{{">=", "1.0.0"}},
			latestVersion:  "2.0.0",
			updateLevel:    Major,
			dependencyKind: DevDependency,
			packageUrl:     "https://github.com/test/package",
			expectedOutput: "github.com/test/package >=1.0.0 2.0.0 major dev 0 https://github.com/test/package",
		},
		{
			name: "Test case 2",
			reporter: TextReporter{
				MessageBefore: "",
				Pattern:       "%s %s %s %s %s %.0f %s\n",
				MessageAfter:  "",
			},
			packageName:    "github.com/test/package",
			requirement:    VersionRequirement{{">=", "1.0.0"}, {"<", "2.0.0"}},
			latestVersion:  "1.5.0",
			updateLevel:    Patch,
			dependencyKind: RunDependency,
			packageUrl:     "https://github.com/test/package",
			expectedOutput: "github.com/test/package >=1.0.0, <2.0.0 1.5.0 patch runtime 0 https://github.com/test/package\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tc.reporter.Report(
				tc.packageName,
				tc.requirement,
				tc.latestVersion,
				tc.updateLevel,
				tc.dependencyKind,
				tc.packageUrl,
				false,
				[]string{},
				0,
				&buf,
			)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if strings.Compare(buf.String(), tc.expectedOutput) != 0 {
				t.Fatalf("expected output: %s, but got: %s", tc.expectedOutput, buf.String())
			}
		})
	}

	// ---

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%sExcluded", tc.name), func(t *testing.T) {
			var buf bytes.Buffer
			err := tc.reporter.Report(
				tc.packageName,
				tc.requirement,
				tc.latestVersion,
				tc.updateLevel,
				tc.dependencyKind,
				tc.packageUrl,
				false,
				[]string{tc.packageName},
				0,
				&buf,
			)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if buf.String() != "" {
				t.Fatalf("unexpected output: %s", buf.String())
			}
		})
	}
}

func TestMonochromeTableReporter(t *testing.T) {
	var buf bytes.Buffer

	reporter := MonochromeTableReporter("1.0.0")

	reporter.Before(&buf)

	reporter.Report("github.com/user/repo",
		VersionRequirement{{">=", "1.0.0"}}, "1.2.3",
		Patch, RunDependency, "https://github.com/user/repo", false, []string{}, 0, &buf)

	reporter.After(&buf)

	// Package name "github.com/user/repo" is truncated to "github.com/use" because of the width of the terminal
	expected := "-- wilf v1.0.0 --\nPackage         Wanted          Latest      Package type  Details\n" +
		"github.com/use\t>=1.0.0     \t1.2.3       runtime       patch for github.com/user/repo; https://github.com/user/repo\n"

	if buf.String() != expected {
		t.Errorf("Unexpected output:\nExpected: %s\nGot     : %s", expected, buf.String())
	}
}
