package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"testing"
)

func TestParsePythonDict(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]string
		err      bool
	}{
		{
			input: `{version = "==4.4.0", extras = ["srv", "tls"]}`,
			expected: map[string]string{
				"extras":  `[srv tls]`,
				"version": "==4.4.0",
			},
			err: false,
		},
		{
			input: `{field1 = "value1", field2 = "value2", field3 = {nested = "value"}}`,
			expected: map[string]string{
				"field1": "value1",
				"field2": "value2",
				"field3": `map[nested:value]`,
			},
			err: false,
		},
		{
			input:    `{}`,
			expected: map[string]string{},
			err:      false,
		},
		{
			input:    `{field1: "value1", field2}`,
			expected: nil,
			err:      true,
		},
	}

	for _, test := range tests {
		result, err := ParsePythonDict(test.input)

		if (err != nil) != test.err {
			t.Errorf("Expected error: %v, got error: %v",
				test.err, err)
		}

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Expected result: %v, got: %v",
				test.expected, result)
		}
	}
}

func TestParsePipfile(t *testing.T) {
	log.SetLevel(log.ErrorLevel)

	tests := []struct {
		path     string
		expected Pipfile
	}{
		{
			path: "resources/valid1.pipfile",
			expected: Pipfile{
				RuntimeDependencies: Dependencies{
					"requests": VersionRequirement{
						VersionConstraint{
							"==",
							"v2.26.0",
						},
					},
					"numpy": VersionRequirement{
						VersionConstraint{
							">=",
							"v1.21.0",
						},
						VersionConstraint{
							"<",
							"v1.22.0",
						},
					},
				},
				DevDependencies: Dependencies{
					"pytest": VersionRequirement{
						VersionConstraint{
							"~=",
							"v6.0",
						},
					},
					"black": VersionRequirement{
						VersionConstraint{"*", "*"},
					},
				},
				RequiresPythonVersion: VersionRequirement{},
			},
		},
		{
			path: "resources/valid2.pipfile",
			expected: Pipfile{
				RuntimeDependencies:   make(Dependencies),
				DevDependencies:       make(Dependencies),
				RequiresPythonVersion: VersionRequirement{},
			},
		},
		{
			path: "resources/valid3.pipfile",
			expected: Pipfile{
				RuntimeDependencies: Dependencies{
					"envyaml": VersionRequirement{
						VersionConstraint{
							"==",
							"v1.10.211231",
						},
					},
					"requests": VersionRequirement{
						VersionConstraint{
							"==",
							"v2.31.0",
						},
					},
					"python-dateutil": VersionRequirement{
						VersionConstraint{
							"==",
							"v2.8.2",
						},
					},
					"pymongo": VersionRequirement{
						VersionConstraint{
							"==",
							"v4.4.0",
						},
					},
				},
				DevDependencies: Dependencies{
					"ipython": VersionRequirement{
						VersionConstraint{"*", "*"},
					},
					"coverage": VersionRequirement{
						VersionConstraint{"*", "*"},
					},
					"flake8": VersionRequirement{
						VersionConstraint{"*", "*"},
					},
					"flake8-import-order": VersionRequirement{
						VersionConstraint{"*", "*"},
					},
					"flake8_formatter_junit_xml": VersionRequirement{
						VersionConstraint{"*", "*"},
					},
					"mypy": VersionRequirement{
						VersionConstraint{"*", "*"},
					},
					"pytest": VersionRequirement{
						VersionConstraint{"*", "*"},
					},
					"junit-xml": VersionRequirement{
						VersionConstraint{"*", "*"},
					},
					"types-requests": VersionRequirement{
						VersionConstraint{"*", "*"},
					},
					"types-python-dateutil": VersionRequirement{
						VersionConstraint{"*", "*"},
					},
					"pymongo-stubs": VersionRequirement{
						VersionConstraint{"*", "*"},
					},
					"lorem": VersionRequirement{
						VersionConstraint{
							"==",
							"v1.3.3",
						},
					},
					"numpy": VersionRequirement{
						VersionConstraint{
							"==",
							"v1.23.3",
						},
					},
				},
				RequiresPythonVersion: VersionRequirement{
					VersionConstraint{
						">=",
						"v3.8",
					},
				},
			},
		},
	}

	for _, test := range tests {
		file, err := os.Open(test.path)

		if err != nil {
			t.Errorf("Error occurred while loading fixture '%s': %v", test.path, err)
		}

		result, err := ParsePipfile(file)

		if err != nil {
			t.Errorf("Error occurred while parsing Pipfile '%s': %v",
				test.path, err)
		}

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Expected result for '%s': %v, got: %v",
				test.path, test.expected, result,
			)
		}
	}
}

func TestParseVersionRequirement(t *testing.T) {
	tests := []struct {
		input    string
		expected VersionRequirement
		err      error
	}{
		{
			input: "<=5",
			expected: VersionRequirement{
				VersionConstraint{"<=", "v5"},
			},
			err: nil,
		},
		{
			input:    "==value",
			expected: VersionRequirement{},
			err:      fmt.Errorf("invalid version: value"),
		},
		{
			input: "==1.*",
			expected: VersionRequirement{
				VersionConstraint{"~", "v1.*"},
			},
			err: nil,
		},
		{
			input: "==1.2.*",
			expected: VersionRequirement{
				VersionConstraint{"~", "v1.2.*"},
			},
			err: nil,
		},
		{
			input:    ">=1.2.*",
			expected: VersionRequirement{},
			err:      fmt.Errorf("invalid version matching: >=1.2.*"),
		},
		{
			input:    "==1.A.*",
			expected: VersionRequirement{},
			err:      fmt.Errorf("invalid version: 1.A.*"),
		},
		{
			input: ">100",
			expected: VersionRequirement{
				VersionConstraint{">", "v100"},
			},
			err: nil,
		},
		{
			input: "!=10",
			expected: VersionRequirement{
				VersionConstraint{"!=", "v10"},
			},
			err: nil,
		},
		{
			input: ">=3.14",
			expected: VersionRequirement{
				VersionConstraint{">=", "v3.14"},
			},
			err: nil,
		},
		{
			input: ">=1.21.0, <1.22.0",
			expected: VersionRequirement{
				VersionConstraint{">=", "v1.21.0"},
				VersionConstraint{"<", "v1.22.0"},
			},
			err: nil,
		},
		{
			input:    "~=",
			expected: VersionRequirement{},
			err:      fmt.Errorf("missing version: ~="),
		},
		{
			input:    "some value",
			expected: VersionRequirement{},
			err:      fmt.Errorf("invalid version: some value"),
		},
		{
			input: "*",
			expected: VersionRequirement{
				VersionConstraint{"*", "*"},
			},
			err: nil,
		},
		{
			input: "3.7",
			expected: VersionRequirement{
				VersionConstraint{"==", "v3.7"},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		req, err := ParseVersionRequirement(test.input)

		if test.err == nil && err != nil {
			t.Errorf("Unexpected parse error: %v", err)
		}

		if test.err != nil && err == nil {
			t.Errorf("Expected parse error: %v, got: %v",
				test.err, err)

			continue
		}

		if len(test.expected) != len(req) {
			t.Errorf("Expected requirement: %v, got: %v",
				test.expected, req)

			continue
		}

		for i, c := range test.expected {
			xop := c[0]
			rop := req[i][0]

			if xop != rop {
				t.Errorf("Expected operator: %v, got: %v",
					xop, rop)
			}

			xver := c[1]
			rver := req[i][1]

			if xver != rver {
				t.Errorf("Expected version: %v, got: %v",
					xver, rver)
			}
		}
	}
}
