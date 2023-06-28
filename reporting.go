package main

import (
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
)

type UpdateReporter interface {
	Before(out io.Writer)

	// Report an update for a package.
	// Arguments:
	// - packageName: string representing the name of the package being reported.
	// - requirement: VersionRequirement representing the version requirement for the package.
	// - latestVersion: string representing the latest version available for the package.
	// - updateLevel: UpdateLevel representing the level of update available for the package.
	// - dependencyKind: DependencyKind representing the kind of dependency for the package.
	// - packageUrl: string representing the URL of the package.
	// - fatal: boolean indicating whether the package is a fatal dependency.
	// - excludedPackages: slice of strings representing the names of packages to be excluded from the report.
	// - timeSec: duration in seconds to check the package.
	// - out: io.Writer representing the output writer to which the report will be written.
	// Returns an error if any.
	Report(
		packageName string,
		requirement VersionRequirement,
		latestVersion string,
		updateLevel UpdateLevel,
		dependencyKind DependencyKind,
		packageUrl string,
		fatal bool,
		excludedPackages []string,
		timeSec float64,
		out io.Writer,
	) error

	After(out io.Writer)
}

// ---

type TextReporter struct {
	MessageBefore string
	Pattern       string
	MessageAfter  string
}

// Before is a method of the UpdateReporter interface. It writes the MessageBefore field of the TextReporter struct to the output writer.
func (r TextReporter) Before(out io.Writer) {
	fmt.Fprint(out, r.MessageBefore)
}

// Report is a method of the UpdateReporter interface. It formats the output of the report in a text format and writes it to the output writer.
func (r TextReporter) Report(
	packageName string,
	requirement VersionRequirement,
	latestVersion string,
	updateLevel UpdateLevel,
	dependencyKind DependencyKind,
	packageUrl string,
	fatal bool,
	excludedPackages []string,
	timeSec float64,
	out io.Writer,
) error {
	if ContainsString(excludedPackages, packageName) {
		log.Debugf("skipping package %s", packageName)

		return nil
	}

	// ---

	reqs := make([]string, len(requirement))

	for i, req := range requirement {
		reqs[i] = fmt.Sprintf("%s%s", req[0], req[1])
	}

	fmt.Fprintf(
		out,
		r.Pattern,
		packageName,
		strings.Join(reqs, ", "),
		latestVersion,
		updateLevel,
		dependencyKind,
		timeSec,
		packageUrl,
	)

	return nil
}

// After is a method of the UpdateReporter interface. It writes the MessageAfter field of the TextReporter struct to the output writer.
func (r TextReporter) After(out io.Writer) {
	fmt.Fprint(out, r.MessageAfter)
}

// ---

// MonochromeTableReporter returns a TextReporter that formats the output of the report in a monochrome table format.
// The function takes a version string as input and returns a TextReporter struct with MessageBefore, Pattern, and MessageAfter fields.
// The MessageBefore field contains a formatted string with the version number and column headers.
// The Pattern field contains a formatted string with placeholders for package name, wanted version, latest version, package type, and details.
// The MessageAfter field is an empty string.
func MonochromeTableReporter(version string) TextReporter {
	return TextReporter{
		MessageBefore: fmt.Sprintf("-- wilf v%s --\nPackage         Wanted          Latest      Package type  Details\n", version),
		Pattern:       "%-14.14[1]s\t%-12.12[2]s\t%-12.12[3]s%-12.12[5]s  %[4]s for %[1]s; %[7]s\n",
		MessageAfter:  "",
	}
}
