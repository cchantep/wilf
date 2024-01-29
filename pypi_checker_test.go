package main

import (
	"testing"
)

func TestPypiRequiredUpdate(t *testing.T) {
	pypiChecker := PypiChecker{}

	tests := []struct {
		pkg           string
		requirement   VersionRequirement
		pythonVersion VersionRequirement
		expectedVer   string
		expectedLvl   UpdateLevel
		expectedUrl   string
		expectedError error
	}{
		// Test case for an existing package with no update required
		{
			pkg: "requests",
			requirement: VersionRequirement{
				VersionConstraint{">=", "v2.26.0"},
			},
			expectedVer:   "v2.31.0",
			expectedLvl:   0,
			expectedUrl:   "https://requests.readthedocs.io",
			expectedError: nil,
		},
		// Test case for an existing package with minor update required
		{
			pkg: "numpy",
			requirement: VersionRequirement{
				VersionConstraint{"<", "v1.24.0"},
			},
			expectedVer:   "v1.26.3",
			expectedLvl:   Minor,
			expectedUrl:   "https://numpy.org",
			expectedError: nil,
		},
		// Test case for an existing package with major update required
		{
			pkg: "pytest",
			requirement: VersionRequirement{
				VersionConstraint{"~=", "v6.0"},
			},
			expectedVer:   "v8.0.0",
			expectedLvl:   Major,
			expectedUrl:   "https://docs.pytest.org/en/latest/",
			expectedError: nil,
		},
		// Test case with python version incompatible with update
		{
			pkg: "pytest",
			requirement: VersionRequirement{
				VersionConstraint{"~=", "v6.0"},
			},
			pythonVersion: VersionRequirement{
				VersionConstraint{"==", "v3.7"},
			},
			expectedVer:   "v8.0.0",
			expectedLvl:   0,
			expectedUrl:   "https://docs.pytest.org/en/latest/",
			expectedError: nil,
		},
		// Test case for a non-existing package
		{
			pkg: "nonexistent-package",
			requirement: VersionRequirement{
				VersionConstraint{">=", "v1.0.0"},
			},
			expectedVer:   "",
			expectedLvl:   0,
			expectedUrl:   "",
			expectedError: nil,
		},
	}

	for _, test := range tests {
		pypiChecker.PythonRequirement = test.pythonVersion

		ver, lvl, url, err := pypiChecker.RequiredUpdate(test.pkg, test.requirement)

		if (err == nil && test.expectedError != nil) ||
			(err != nil && test.expectedError == nil) ||
			(err != nil &&
				err.Error() != test.expectedError.Error()) {
			t.Errorf("For package %s, expected error: %v, but got: %v", test.pkg, test.expectedError, err)
		}

		if ver != test.expectedVer {
			t.Errorf("For package %s, expected version: %s, but got: %s", test.pkg, test.expectedVer, ver)
		}

		if url != test.expectedUrl {
			t.Errorf("For package %s, expected url: %s, but got: %s", test.pkg, test.expectedUrl, url)
		}

		if lvl != test.expectedLvl {
			t.Errorf("For package %s, expected update level: %d, but got: %d", test.pkg, test.expectedLvl, lvl)
		}
	}
}
