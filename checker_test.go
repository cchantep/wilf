package main

import (
	"fmt"
	"testing"
)

func TestMatchConstraint(t *testing.T) {
	latestVersion := "v1.2.3"

	tests := []struct {
		constraint VersionConstraint
		expected   bool
	}{
		// Test cases for valid constraints
		{
			constraint: VersionConstraint{"==", "v1.2.3.4"}, // non standard
			expected:   true,
		},
		{
			constraint: VersionConstraint{">=", "v1.0.0"},
			expected:   true,
		},
		{
			constraint: VersionConstraint{">", "v1.0.0"},
			expected:   true,
		},
		{
			constraint: VersionConstraint{"<=", "v1.2.3"},
			expected:   true,
		},
		{
			constraint: VersionConstraint{"<", "v1.2.4"},
			expected:   true,
		},
		{
			constraint: VersionConstraint{"==", "v1.2.3"},
			expected:   true,
		},
		{
			constraint: VersionConstraint{"!=", "v1.2.2"},
			expected:   true,
		},
		{
			constraint: VersionConstraint{"~=", "v1.2.0"},
			expected:   true,
		},
		{
			constraint: VersionConstraint{"~", "v1.*"},
			expected:   true,
		},
		{
			constraint: VersionConstraint{"~", "v1.2.*"},
			expected:   true,
		},
		{
			constraint: VersionConstraint{"!~", "v1.3.*"},
			expected:   true,
		},
		// Test cases for invalid constraints
		{
			constraint: VersionConstraint{"<", "v1.2.3"},
			expected:   false,
		},
		{
			constraint: VersionConstraint{">", "v1.2.3"},
			expected:   false,
		},
		{
			constraint: VersionConstraint{">=", "v2.0.0"},
			expected:   false,
		},
		{
			constraint: VersionConstraint{"<", "v1.0.0"},
			expected:   false,
		},
		{
			constraint: VersionConstraint{"<=", "v1.2.2"},
			expected:   false,
		},
		{
			constraint: VersionConstraint{"==", "v1.0.0"},
			expected:   false,
		},
		{
			constraint: VersionConstraint{"!=", "v1.2.3"},
			expected:   false,
		},
		{
			constraint: VersionConstraint{"~=", "v1.3.0"},
			expected:   false,
		},
		{
			constraint: VersionConstraint{"~", "v1.5.*"},
			expected:   false,
		},
		{
			constraint: VersionConstraint{"!~", "v1.2.*"},
			expected:   false,
		},
	}

	for _, test := range tests {
		result := MatchConstraint(latestVersion, test.constraint)

		if result != test.expected {
			t.Errorf("For constraint '%s', expected %v, but got %v", test.constraint, test.expected, result)
		}
	}
}

func TestShouldUpdate(t *testing.T) {
	latestVersion := "v1.2.3"

	tests := []struct {
		requirement VersionRequirement
		expected    bool
	}{
		// Test cases for requirements that require an update
		{VersionRequirement{
			VersionConstraint{"<", "v1.2.3"},
		}, true},
		{VersionRequirement{
			VersionConstraint{"<=", "v1.2.2"},
			VersionConstraint{"!=", "v1.2.3"},
		}, true},
		{VersionRequirement{
			VersionConstraint{">=", "v1.0.0"},
			VersionConstraint{"~=", "v1.3.0"},
		}, true},
		{VersionRequirement{
			VersionConstraint{"<", "v1.2.3.4"}, // non standard
		}, true},

		// Test cases for requirements that don't require an update
		{VersionRequirement{
			VersionConstraint{"==", latestVersion},
		}, false},
		{VersionRequirement{
			VersionConstraint{">=", "v1.2.2"},
		}, false},
		{VersionRequirement{
			VersionConstraint{"==", "v1.2.3"},
		}, false},
		{VersionRequirement{
			VersionConstraint{">=", "v1.2.3"},
		}, false},
		{VersionRequirement{
			VersionConstraint{">=", "v1.2.2"},
			VersionConstraint{"<=", "v2.0.0"},
		}, false},
		{VersionRequirement{
			VersionConstraint{">=", "v1.0.0"},
			VersionConstraint{"~=", "v1.2.0"},
		}, false},
		{VersionRequirement{
			VersionConstraint{"*"},
			VersionConstraint{">=", "v1.2.3"},
		}, false},
		{VersionRequirement{
			VersionConstraint{">", "v1.2.2.3"}, // non standard
		}, false},
	}

	for _, test := range tests {
		result := ShouldUpdate(test.requirement, latestVersion)

		if result != test.expected {
			t.Errorf("For requirement '%v', expected %v, but got %v", test.requirement, test.expected, result)
		}
	}

	nonStandardRequirement := VersionRequirement{
		VersionConstraint{"==", "v1.5.5.1"},
	}

	result := ShouldUpdate(nonStandardRequirement, "v1.5.5.1")

	if result != false {
		t.Errorf("For requirement '%v', expected %v, but got %v", nonStandardRequirement, false, result)
	}
}

func TestCreateUpdateLevel(t *testing.T) {
	tests := []struct {
		requirement VersionRequirement
		latest      string
		expected    UpdateLevel
		expectedErr error
	}{
		// Test cases for no requirement
		{VersionRequirement{}, "v1.2.3", 0, fmt.Errorf("missing requirement")},

		// Test cases for Major update level
		{VersionRequirement{{">=", "v1.0.0"}}, "v2.0.0", Major, nil},
		{VersionRequirement{{">=", "v1.0.0"}}, "v1.5.0", Minor, nil},
		{VersionRequirement{{">=", "v1.0.0.0"}}, "v1.5.0", Minor, nil},
		{VersionRequirement{{">=", "v1.0.0"}, {"<", "v2.0.0"}}, "v3.0.0", Major, nil},

		// Test cases for Minor update level
		{VersionRequirement{{">=", "v1.2.0"}}, "v1.2.3", Patch, nil},
		{VersionRequirement{{">=", "v1.2.3"}, {"<", "v1.3.0"}}, "v1.2.5", Minor, nil},
		{VersionRequirement{{"~=", "v1.2.0"}}, "v1.2.8", Patch, nil},

		// Test cases for no update required
		{VersionRequirement{{">=", "v1.2.3"}}, "v1.2.3", 0, nil},
		{VersionRequirement{{">=", "v1.0.0"}}, "v1.0.0", 0, nil},
		{VersionRequirement{{"~=", "v1.2.3"}}, "v1.2.3", 0, nil},
	}

	for _, test := range tests {
		result, err := CreateUpdateLevel(test.requirement, test.latest)
		if err != nil {
			if test.expectedErr == nil || err.Error() != test.expectedErr.Error() {
				t.Errorf("For requirement '%v' and latest version '%s', expected error '%v', but got '%v'", test.requirement, test.latest, test.expectedErr, err)
			}
			continue
		}

		if result != test.expected {
			t.Errorf("For requirement '%v' and latest version '%s', expected update level '%v', but got '%v'", test.requirement, test.latest, test.expected, result)
		}
	}
}

func TestNormalizeNonStandardVersion(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal version",
			input:    "1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "non-standard version",
			input:    "1.2.3.4",
			expected: "1.2.3",
		},
		{
			name:     "long non-standard version",
			input:    "1.2.3.4.5.6",
			expected: "1.2.3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := NormalizeNonStandardVersion(tc.input)
			if actual != tc.expected {
				t.Errorf("expected %s but got %s", tc.expected, actual)
			}
		})
	}
}
