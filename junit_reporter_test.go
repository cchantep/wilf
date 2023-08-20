package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestBefore(t *testing.T) {
	testTime := time.Now()

	// Create a new JUnitReporter
	r := &JUnitReporter{}

	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Call the Before function
	r.Before(&buf)

	if r.ReporterName() != "junit" {
		t.Errorf("Expected ReporterName to be 'junit', got '%s'", r.ReporterName())
	}

	// Check that the Start time is set
	if r.StartTime.IsZero() {
		t.Error("Expected Start time to be set")
	}

	if r.StartTime.Compare(testTime) <= 0 {
		t.Errorf("Expected Start time to be after %s, got %s", testTime, r.StartTime)
	}

	// Check that the DevTestSuite is initialized correctly
	if r.DevTestSuite.Name != "dev" {
		t.Errorf("Expected DevTestSuite Name to be 'dev', got '%s'", r.DevTestSuite.Name)
	}

	if r.DevTestSuite._timestamp != r.StartTime {
		t.Errorf("Expected DevTestSuite _timestamp to be %s, got %s", r.StartTime, r.DevTestSuite._timestamp)
	}

	if r.DevTestSuite.Tests != 0 {
		t.Errorf("Expected DevTestSuite Tests to be 0, got %d", r.DevTestSuite.Tests)
	}

	if r.DevTestSuite.Failures != 0 {
		t.Errorf("Expected DevTestSuite Failures to be 0, got %d", r.DevTestSuite.Failures)
	}

	if r.DevTestSuite.Errors != 0 {
		t.Errorf("Expected DevTestSuite Errors to be 0, got %d", r.DevTestSuite.Errors)
	}

	if r.DevTestSuite.Skipped != 0 {
		t.Errorf("Expected DevTestSuite Skipped to be 0, got %d", r.DevTestSuite.Skipped)
	}

	if r.DevTestSuite.Time != 0.0 {
		t.Errorf("Expected DevTestSuite Time to be 0.0, got %f", r.DevTestSuite.Time)
	}

	if len(r.DevTestSuite.TestCases) != 0 {
		t.Errorf("Expected DevTestSuite TestCases to be empty, got %d", len(r.DevTestSuite.TestCases))
	}

	// Check that the RunTestSuite is initialized correctly
	if r.RunTestSuite.Name != "run" {
		t.Errorf("Expected RunTestSuite Name to be 'run', got '%s'", r.RunTestSuite.Name)
	}

	if r.RunTestSuite._timestamp != r.StartTime {
		t.Errorf("Expected RunTestSuite _timestamp to be %s, got %s", r.StartTime, r.RunTestSuite._timestamp)
	}

	if r.RunTestSuite.Tests != 0 {
		t.Errorf("Expected RunTestSuite Tests to be 0, got %d", r.RunTestSuite.Tests)
	}

	if r.RunTestSuite.Failures != 0 {
		t.Errorf("Expected RunTestSuite Failures to be 0, got %d", r.RunTestSuite.Failures)
	}

	if r.RunTestSuite.Errors != 0 {
		t.Errorf("Expected RunTestSuite Errors to be 0, got %d", r.RunTestSuite.Errors)
	}

	if r.RunTestSuite.Skipped != 0 {
		t.Errorf("Expected RunTestSuite Skipped to be 0, got %d", r.RunTestSuite.Skipped)
	}

	if r.RunTestSuite.Time != 0.0 {
		t.Errorf("Expected RunTestSuite Time to be 0.0, got %f", r.RunTestSuite.Time)
	}

	if len(r.RunTestSuite.TestCases) != 0 {
		t.Errorf("Expected RunTestSuite TestCases to be empty, got %d", len(r.RunTestSuite.TestCases))
	}

	// Check that the output is empty
	if buf.String() != "" {
		t.Errorf("Expected output to be empty, got '%s'", buf.String())
	}
}

func TestReport(t *testing.T) {
	r := &JUnitReporter{}
	out := &bytes.Buffer{}

	// Call the Report function with some sample data
	err := r.Report(
		"mypackage",
		VersionRequirement{{"<", "1.0.0"}},
		"1.0.0",
		Major,
		RunDependency,
		"https://mypackage.com",
		false,
		[]string{},
		1.23,
		out,
	)

	// Check that the function returned no error
	if err != nil {
		t.Errorf("Report returned an error: %v", err)
	}

	// Check that the output buffer contains some expected text
	if len(out.String()) != 0 {
		t.Errorf("Report produced unexpected output: %v", out.String())
	}

	if len(r.RunTestSuite.TestCases) != 1 {
		t.Errorf("Expected RunTestSuite TestCases to be empty, got %d", len(r.RunTestSuite.TestCases))
	}

	if len(r.DevTestSuite.TestCases) != 0 {
		t.Errorf("Expected DevTestSuite TestCases to be empty, got %d", len(r.DevTestSuite.TestCases))
	}
}

func TestSecondsSince(t *testing.T) {
	now := time.Now()
	from := now.Add(-time.Second * 10)
	expected := float64(10)

	if got := SecondsSince(from); got != expected {
		t.Errorf("SecondsSince(%v) = %v, want %v", from, got, expected)
	}
}

func TestAfter(t *testing.T) {
	ts := time.Now().Add(time.Second * -5)
	formattedTs := ts.Format("2006-01-02T15:04:05")

	// Create a new JUnitReporter instance
	reporter := &JUnitReporter{
		Version:   "1.0.0",
		StartTime: ts,
		DevTestSuite: JUnitTestSuite{
			_timestamp: ts,
			Name:       "dev",
			TestCases: []JUnitTestCase{
				{
					Name:      "test1",
					Timestamp: formattedTs,
					Time:      2.0,
				},
				{
					Name:      "test2",
					Timestamp: formattedTs,
					Failure: &JUnitFailure{
						Message: "test2 failed",
						Type:    "error",
						Text:    "test2 error message",
					},
					Time: 3.0,
				},
			},
		},
		RunTestSuite: JUnitTestSuite{
			Name:       "run",
			_timestamp: ts,
			TestCases: []JUnitTestCase{
				{
					Name:      "test3",
					Timestamp: formattedTs,
					Time:      3.0,
				},
				{
					Name:      "test4",
					Timestamp: formattedTs,
					Skipped: &JUnitSkipped{
						Message: "test4 skipped",
						Text:    "test4 skipped message",
					},
					Time: 4.0,
				},
			},
		},
	}

	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Call the After function
	reporter.After(&buf)

	// Check the output
	expected := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<testsuites name="wilf v1.0.0" tests="4" failures="1" errors="2" skipped="1" time="5">
  <testsuite name="dev" timestamp="%[1]s" tests="2" failures="1" errors="1" skipped="0" time="5">
    <testcase name="test1" time="2" timestamp="%[1]s"></testcase>
    <testcase name="test2" time="3" timestamp="%[1]s">
      <failure message="test2 failed" type="error">test2 error message</failure>
    </testcase>
  </testsuite>
  <testsuite name="run" timestamp="%[1]s" tests="2" failures="0" errors="1" skipped="1" time="5">
    <testcase name="test3" time="3" timestamp="%[1]s"></testcase>
    <testcase name="test4" time="4" timestamp="%[1]s">
      <skipped message="test4 skipped">test4 skipped message</skipped>
    </testcase>
  </testsuite>
</testsuites>
`, formattedTs)
	if buf.String() != expected {
		t.Errorf("unexpected output:\n%s", buf.String())
	}
}
