package main

import (
	"fmt"
	"io"
	"strings"

	color "github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

type ColorizedTableReporter struct {
	Version string
}

func (r ColorizedTableReporter) Before(out io.Writer) {
	fmt.Fprint(out, "-- ")
	color.New(color.Bold).Fprintf(out, "wilf v%s", r.Version)
	fmt.Fprint(out, " --\n\n")

	color.New(color.FgBlue).Fprint(out, "info")
	fmt.Fprintf(out, " Color legend:\n ")

	color.New(color.FgRed).Fprint(out, "<red>")
	fmt.Fprintf(out, "    : Major Update backward-incompatible updates\n ")

	color.New(color.FgYellow).Fprint(out, "<yellow>")
	fmt.Fprintf(out, " : Minor Update backward-compatible features\n ")

	color.New(color.FgGreen).Fprint(out, "<green>")
	fmt.Fprintln(out, "  : Patch Update backward-compatible bug fixes")
	fmt.Fprintln(out)

	underline := color.New(color.Underline)

	underline.Fprint(out, "Package")
	fmt.Fprint(out, "         ")

	underline.Fprint(out, "Wanted")
	fmt.Fprint(out, "          ")

	underline.Fprint(out, "Latest")
	fmt.Fprint(out, "      ")

	underline.Fprint(out, "Package type")
	fmt.Fprint(out, "  ")

	underline.Fprintln(out, "Details")
	fmt.Fprintln(out)
}

func (r ColorizedTableReporter) Report(
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

	pc := color.New(color.FgHiBlack, color.Bold)

	if updateLevel == Major {
		pc = color.New(color.FgRed)
	} else if updateLevel == Minor {
		pc = color.New(color.FgYellow)
	} else if updateLevel == Patch {
		pc = color.New(color.FgGreen)
	}

	pc.Fprintf(out, "%-14.14s", packageName)
	fmt.Fprint(out, "\t")

	fmt.Fprintf(out, "%-12.12s", strings.Join(reqs, ", "))
	fmt.Fprint(out, "\t")

	pc.Add(color.Bold).Fprintf(out, "%-10.10s", latestVersion)
	fmt.Fprint(out, "  ")

	kind := "dev"

	if dependencyKind == RunDependency {
		kind = "runtime"
	}

	fmt.Fprintf(out, "%-12.12s", kind)
	fmt.Fprint(out, "  ")

	fmt.Fprintf(out, "%s; %s\n", packageName, packageUrl)

	return nil
}

func (r ColorizedTableReporter) After(out io.Writer) {
	fmt.Fprintln(out)
}
