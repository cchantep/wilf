package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"
)

type JUnitReporter struct {
	Version      string
	StartTime    time.Time
	DevTestSuite JUnitTestSuite
	RunTestSuite JUnitTestSuite
}

type testSuites struct {
	XMLName  xml.Name         `xml:"testsuites"`
	Name     string           `xml:"name,attr"`
	Tests    int              `xml:"tests,attr"`
	Failures int              `xml:"failures,attr"`
	Errors   int              `xml:"errors,attr"`
	Skipped  int              `xml:"skipped,attr"`
	Time     float64          `xml:"time,attr"`
	Suites   []JUnitTestSuite `xml:"testsuite"`
}

type JUnitTestSuite struct {
	Name       string `xml:"name,attr"`
	_timestamp time.Time
	Timestamp  string          `xml:"timestamp,attr"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Errors     int             `xml:"errors,attr"`
	Skipped    int             `xml:"skipped,attr"`
	Time       float64         `xml:"time,attr"`
	TestCases  []JUnitTestCase `xml:"testcase"`
}

type JUnitTestCase struct {
	Name      string        `xml:"name,attr"`
	Failure   *JUnitFailure `xml:"failure,omitempty"`
	Skipped   *JUnitSkipped `xml:"skipped,omitempty"`
	Time      float64       `xml:"time,attr"`
	Timestamp string        `xml:"timestamp,attr"`
}

type JUnitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Text    string `xml:",chardata"`
}

type JUnitSkipped struct {
	Message string `xml:"message,attr"`
	Text    string `xml:",chardata"`
}

func (r *JUnitReporter) Before(out io.Writer) {
	r.StartTime = time.Now()

	ts := time.Now()

	r.DevTestSuite = JUnitTestSuite{
		Name:       "dev",
		_timestamp: ts,
		TestCases:  []JUnitTestCase{},
		Tests:      0,
		Failures:   0,
		Errors:     0,
		Skipped:    0,
		Time:       0.0,
	}

	r.RunTestSuite = JUnitTestSuite{
		Name:       "run",
		_timestamp: ts,
		TestCases:  []JUnitTestCase{},
		Tests:      0,
		Failures:   0,
		Errors:     0,
		Skipped:    0,
		Time:       0.0,
	}
}

func (r *JUnitReporter) Report(
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
	// Select the appropriate testSuite
	testSuite := &r.RunTestSuite

	if dependencyKind == DevDependency {
		testSuite = &r.DevTestSuite
	}

	// Prepare the testCase representation
	testCase := JUnitTestCase{
		Name:      fmt.Sprintf("%s %s", packageName, updateLevel),
		Time:      timeSec,
		Timestamp: time.Now().Format("2006-01-02T15:04:05"),
	}

	if ContainsString(excludedPackages, packageName) {
		testCase.Skipped = &JUnitSkipped{
			Message: fmt.Sprintf(
				"package '%s' is excluded", packageName),
			Text: fmt.Sprintf("Package '%s' is excluded by configuration: [%s]", packageName, strings.Join(excludedPackages, ", ")),
		}

		testSuite.TestCases = append(testSuite.TestCases, testCase)

		return nil
	}

	// ---

	if fatal {
		msg := fmt.Sprintf("%s %s is outdated. Latest version is %s", packageName, updateLevel, latestVersion)

		testCase.Failure = &JUnitFailure{
			Message: msg,
			Type:    "error",
			Text:    msg,
		}
	}

	testSuite.TestCases = append(testSuite.TestCases, testCase)

	return nil
}

func finalizeTestSuite(suite *JUnitTestSuite) {
	ts := suite._timestamp

	suite.Tests = len(suite.TestCases)
	suite.Time = time.Since(ts).Seconds()
	suite.Timestamp = ts.Format("2006-01-02T15:04:05")

	for _, testCase := range suite.TestCases {
		if testCase.Failure != nil {
			suite.Failures++
		} else if testCase.Skipped != nil {
			suite.Skipped++
		} else {
			suite.Errors++
		}
	}
}

func (r *JUnitReporter) After(out io.Writer) {
	finalizeTestSuite(&r.DevTestSuite)
	finalizeTestSuite(&r.RunTestSuite)

	testSuites := testSuites{
		Name:     fmt.Sprintf("wilf v%s", r.Version),
		Tests:    r.DevTestSuite.Tests + r.RunTestSuite.Tests,
		Failures: r.DevTestSuite.Failures + r.RunTestSuite.Failures,
		Errors:   r.DevTestSuite.Errors + r.RunTestSuite.Errors,
		Skipped:  r.DevTestSuite.Skipped + r.RunTestSuite.Skipped,
		Time:     time.Since(r.StartTime).Seconds(),
		Suites:   []JUnitTestSuite{r.DevTestSuite, r.RunTestSuite},
	}

	xmlHeader := []byte(xml.Header)
	xmlBody, err := xml.MarshalIndent(testSuites, "", "  ")

	if err != nil {
		panic(err)
	}

	xmlOutput := append(xmlHeader, xmlBody...)

	fmt.Fprintln(out, string(xmlOutput))
}
