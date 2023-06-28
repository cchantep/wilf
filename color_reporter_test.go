package main

import (
	"bytes"
	"testing"

	color "github.com/fatih/color"
)

func TestColorizedTableReporterBeforeAfter(t *testing.T) {
	var buf bytes.Buffer

	reporter := ColorizedTableReporter{Version: "1.2.3"}

	reporter.Before(&buf)

	var expected bytes.Buffer

	expected.WriteString("-- ")
	color.New(color.Bold).Fprint(&expected, "wilf v1.2.3")
	expected.WriteString(" --\n\n")

	color.New(color.FgBlue).Fprint(&expected, "info")
	expected.WriteString(" Color legend:\n ")

	color.New(color.FgRed).Fprint(&expected, "<red>")
	expected.WriteString("    : Major Update backward-incompatible updates\n ")

	color.New(color.FgYellow).Fprint(&expected, "<yellow>")
	expected.WriteString(" : Minor Update backward-compatible features\n ")

	color.New(color.FgGreen).Fprint(&expected, "<green>")
	expected.WriteString("  : Patch Update backward-compatible bug fixes\n\n")

	underline := color.New(color.Underline)

	underline.Fprint(&expected, "Package")
	expected.WriteString("         ")

	underline.Fprint(&expected, "Wanted")
	expected.WriteString("          ")

	underline.Fprint(&expected, "Latest")
	expected.WriteString("      ")

	underline.Fprint(&expected, "Package type")
	expected.WriteString("  ")

	underline.Fprintln(&expected, "Details")
	expected.WriteString("\n")

	if buf.String() != expected.String() {
		t.Errorf("Unexpected output:\nExpected: %s\nGot     : %s", expected.String(), buf.String())
	}

	buf.Reset()
	expected.Reset()

	reporter.After(&buf)

	if buf.String() != "\n" {
		t.Errorf("Unexpected output:\nExpected: <endline>\nGot     : [%s]", buf.String())
	}
}

func TestColorizedTableReporterReportMajorDev(t *testing.T) {
	var buf bytes.Buffer

	reporter := ColorizedTableReporter{Version: "1.2.3"}

	var expected bytes.Buffer

	reporter.Report(
		"github.com/test/package",
		VersionRequirement{{">=", "1.0.0"}},
		"2.0.0",
		Major,
		DevDependency,
		"https://github.com/test/package",
		false,
		[]string{},
		0,
		&buf,
	)

	pc := color.New(color.FgRed)

	pc.Fprint(&expected, "github.com/tes")
	expected.WriteString("\t")

	expected.WriteString(">=1.0.0     \t")

	pc.Add(color.Bold).Fprint(&expected, "2.0.0     ")
	expected.WriteString("  ")

	expected.WriteString("dev           github.com/test/package; https://github.com/test/package\n")

	if buf.String() != expected.String() {
		t.Errorf("Unexpected output:\nExpected: %s\nGot     : %s", expected.String(), buf.String())
	}
}

func TestColorizedTableReporterReportExcludedMajorDev(t *testing.T) {
	var buf bytes.Buffer

	reporter := ColorizedTableReporter{Version: "1.2.3"}

	reporter.Report(
		"github.com/test/package",
		VersionRequirement{{">=", "1.0.0"}},
		"2.0.0",
		Major,
		DevDependency,
		"https://github.com/test/package",
		false,
		[]string{"github.com/test/package"},
		0,
		&buf,
	)

	if buf.String() != "" {
		t.Errorf("Unexpected output:\nExpected: <empty>\nGot     : %s", buf.String())
	}
}

func TestColorizedTableReporterMinorRun(t *testing.T) {
	var buf bytes.Buffer

	reporter := ColorizedTableReporter{Version: "1.2.3"}

	var expected bytes.Buffer

	reporter.Report(
		"github.com/foo/package",
		VersionRequirement{{">=", "1.0.0"}},
		"1.1.0",
		Minor,
		RunDependency,
		"https://github.com/foo/package",
		false,
		[]string{},
		0,
		&buf,
	)

	pc := color.New(color.FgYellow)

	pc.Fprint(&expected, "github.com/foo")
	expected.WriteString("\t")

	expected.WriteString(">=1.0.0     \t")

	pc.Add(color.Bold).Fprint(&expected, "1.1.0     ")
	expected.WriteString("  ")

	expected.WriteString("runtime       github.com/foo/package; https://github.com/foo/package\n")

	if buf.String() != expected.String() {
		t.Errorf("Unexpected output:\nExpected: %s\nGot     : %s", expected.String(), buf.String())
	}
}

func TestColorizedTableReporterPatchRun(t *testing.T) {
	var buf bytes.Buffer

	reporter := ColorizedTableReporter{Version: "1.2.3"}

	var expected bytes.Buffer

	reporter.Report(
		"bar",
		VersionRequirement{{">=", "3.4"}},
		"3.4.5",
		Patch,
		RunDependency,
		"https://github.com/bar/package",
		false,
		[]string{},
		0,
		&buf,
	)

	pc := color.New(color.FgGreen)

	pc.Fprint(&expected, "bar           ")
	expected.WriteString("\t")

	expected.WriteString(">=3.4       \t")

	pc.Add(color.Bold).Fprint(&expected, "3.4.5     ")
	expected.WriteString("  ")

	expected.WriteString("runtime       bar; https://github.com/bar/package\n")

	if buf.String() != expected.String() {
		t.Errorf("Unexpected output:\nExpected: %s\nGot     : %s", expected.String(), buf.String())
	}
}
