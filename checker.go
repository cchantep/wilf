package main

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/mod/semver"
)

type Checker interface {
	RequiredUpdate(
		pkg string,
		requirement VersionRequirement,
	) (string, UpdateLevel, string, error)
}

type DependencyKind int

const (
	DevDependency DependencyKind = iota
	RunDependency
)

func (k DependencyKind) String() string {
	if k == DevDependency {
		return "dev"
	}

	return "runtime"
}

func ContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// ReportUpdates reports updates for the given dependencies.
// It returns a boolean indicating whether there is at least one update available and an error if any.
func ReportUpdates(
	dependencies Dependencies,
	kind DependencyKind,
	minLevel UpdateLevel,
	excludedPackages []string,
	checker Checker,
	reporter UpdateReporter,
	out io.Writer,
) (bool, error) {
	atLeastOneUpdate := false

	for pkg, requirement := range dependencies {
		ts := time.Now()
		ver, lvl, url, err := checker.RequiredUpdate(pkg, requirement)

		if err != nil {
			return false, err
		}

		if lvl == 0 {
			log.Debugf("no update available for %s: '%s'", pkg, ver)

			continue
		}

		fatal := lvl >= minLevel

		if fatal && !ContainsString(excludedPackages, pkg) {
			atLeastOneUpdate = true
		}

		reporter.Report(
			pkg,
			requirement,
			ver,
			lvl,
			kind,
			url,
			fatal,
			excludedPackages,
			time.Since(ts).Seconds(),
			out,
		)
	}

	return atLeastOneUpdate, nil
}

// ShouldUpdate checks if the given requirement should be updated to the latest version.
// It returns a boolean indicating whether the requirement should be updated or not.
// It loops through each constraint in the requirement and checks if the latest version
// matches the constraint.
// If the latest version does not match any of the constraints, it returns true.
// If the latest version matches all of the constraints, it returns false.
// If the requirement contains a wildcard constraint, it returns false.
func ShouldUpdate(
	requirement VersionRequirement,
	latest string,
) bool {
	for _, constraint := range requirement {
		op := constraint[0]

		if op == "*" {
			return false
		}

		if !MatchConstraint(latest, constraint) {
			return true
		}
	}

	return false
}

func AreCompatibles(a, b VersionRequirement) bool {
	for _, ac := range a {
		ao := ac[0]

		if ao == "*" {
			continue
		}

		for _, bc := range b {
			bo := bc[0]

			if bo == "*" {
				continue
			}

			if bo == "==" && MatchConstraint(bc[1], ac) {
				continue
			}

			if ao == "==" && MatchConstraint(ac[1], bc) {
				continue
			}

			if !MatchConstraint(bc[1], ac) || !MatchConstraint(ac[1], bc) {
				return false
			}
		}
	}

	return true
}

// MatchConstraint checks if the latest version matches the given constraint.
// It returns a boolean indicating whether the latest version matches the constraint or not.
// It supports the following operators from PEP-404 (https://peps.python.org/pep-0440):
// - ===: arbitrary equality
// - ~: version matching
// - !~: version not matching
// - <=: less than or equal to
// - <: less than
// - !=: not equal to
// - ~=: compatible release
// - ==: equal to
// - >=: greater than or equal to
// - >: greater than
// It uses the semver package to compare versions.
// If the constraint is invalid, it logs a warning and returns false.
func MatchConstraint(
	latest string,
	constraint VersionConstraint,
) bool {
	op := constraint[0]
	ver := constraint[1]

	if op == "===" {
		// arbitrary equality
		return ver == latest
	}

	normalizedLatest := latest

	if !semver.IsValid(latest) {
		normalizedLatest = NormalizeNonStandardVersion(latest)
	}

	if op == "~" || op == "!~" {
		// '==' + version matching
		re, err := CompileMatching(ver)

		if err != nil {
			log.Warnf("Invalid version matching: %s", ver)

			return false
		}

		matches := re.MatchString(normalizedLatest)

		if op == "~" {
			return matches
		} else {
			return !matches
		}
	}

	// ---

	normalized := ver

	if !semver.IsValid(ver) {
		// Normalize non standard version
		normalized = NormalizeNonStandardVersion(ver)
	}

	c := semver.Compare(normalizedLatest, normalized)

	if op == "<=" {
		return c <= 0
	}

	if op == "<" {
		return c < 0
	}

	if op == "!=" {
		return c != 0
	}

	if op == "~=" {
		// See https://peps.python.org/pep-0440/#compatible-release
		return semver.Major(normalizedLatest) == semver.Major(normalized) && c >= 0
	}

	if op == "==" {
		return c == 0
	}

	if op == ">=" {
		return c >= 0
	}

	if op == ">" {
		return c > 0
	}

	return false
}

// NormalizeNonStandardVersion normalizes a non-standard version to a standard one.
// It takes a version string as input and returns the first three segments
// of the version string separated by a period.
// For example, "1.2.3.4" would be normalized to "1.2.3".
// This function is used to normalize non-standard version strings before
// comparing them with standard version strings.
func NormalizeNonStandardVersion(ver string) string {
	return strings.Join(strings.Split(ver, ".")[0:3], ".")
}

// CompileMatching compiles a regular expression pattern that matches the given string.
// It replaces any asterisk (*) with the regular expression pattern ".*".
// It returns a compiled regular expression and an error if any.
func CompileMatching(matching string) (*regexp.Regexp, error) {
	pattern := strings.ReplaceAll(matching, "*", ".*")

	return regexp.Compile(pattern)
}

// CreateUpdateLevel creates an update level based on
// the given version requirement and latest version.
//
// It returns an UpdateLevel and an error if any.
// It loops through each constraint in the requirement
// and finds the highest version that matches the constraint.
// It then compares the highest version with the latest version
// to determine the update level.
// If the latest version is a major version ahead of the highest version,
// it returns Major.
// If the latest version is a minor version ahead of the highest version,
// it returns Minor.
// If the latest version is a patch version ahead of the highest version,
// it returns Patch.
// If the latest version is the same as the highest version, it returns 0.
// If the requirement is empty, it returns an error.
func CreateUpdateLevel(
	requirement VersionRequirement,
	latest string,
) (UpdateLevel, error) {
	if len(requirement) == 0 {
		return 0, fmt.Errorf("missing requirement")
	}

	// ---

	maxVer := "v0.0.0"

	for _, constraint := range requirement {
		spec := constraint[1]

		if spec == "*" {
			continue
		}

		spec = strings.ReplaceAll(spec, "*", "0")

		normalized := spec

		if !semver.IsValid(spec) {
			// Normalize non standard version
			normalized = NormalizeNonStandardVersion(spec)
		}

		if semver.Compare(maxVer, normalized) < 0 {
			maxVer = normalized
		}
	}

	if maxVer == "v0.0.0" {
		return 0, nil
	}

	if semver.Major(latest) != semver.Major(maxVer) {
		return Major, nil
	}

	if semver.MajorMinor(latest) != semver.MajorMinor(maxVer) {
		return Minor, nil
	}

	if semver.Compare(latest, maxVer) != 0 {
		return Patch, nil
	}

	return 0, nil
}
