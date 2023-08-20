package main

import (
	"testing"
)

func TestParseArguments(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expected          CommandArguments
		expectedReporters []string
		err               bool
	}{
		/* TODO
		{
			name: "no arguments",
			args: []string{},
			err:  true,
		},
		{
			name: "missing pipfile argument",
			args: []string{"-v"},
			err:  true,
		},
		{
			name: "invalid option",
			args: []string{"-x", "Pipfile"},
			err:  true,
		},
		{
			name: "multiple pipfiles",
			args: []string{"Pipfile1", "Pipfile2"},
			err:  true,
		},
		{
			name: "verbose flag",
			args: []string{
				"-v", "Pipfile",
				"-r", "monochrome-table",
				"-r", "junit",
			},
			expected: CommandArguments{
				Verbose:   true,
				Pipfile:   "Pipfile",
				Reporters: "monochrome-table, junit",
			},
		},
		{
			name: "config flag",
			args: []string{"-c", "config.toml", "Pipfile"},
			expected: CommandArguments{
				Config:    "config.toml",
				Pipfile:   "Pipfile",
				Reporters: "colorized-table",
			},
		},*/
		{
			name: "pipfile argument",
			args: []string{"Pipfile"},
			expected: CommandArguments{
				Pipfile:   "Pipfile",
				Reporters: "colorized-table",
			},
			expectedReporters: []string{"colorized-table"},
		},
		/*TODO{
			name: "print usage argument",
			args: []string{"-h"},
			expected: CommandArguments{
				PrintUsage: true,
			},
		},*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args

			actual, reporting, err := ParseArguments(args)

			if tt.err && err == nil {
				t.Errorf("expected an error but got none")
			}

			if !tt.err && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if actual != tt.expected {
				t.Errorf("expected %v but got %v", tt.expected, actual)
			}

			if len(reporting) != len(tt.expectedReporters) {
				t.Errorf("expected %d reporters but got %d", len(tt.expectedReporters), len(reporting))
			}

			for _, r := range reporting {
				if !ContainsString(tt.expectedReporters, r.Reporter.ReporterName()) {
					t.Errorf("unexpected reporter: %s", r.Reporter.ReporterName())
				}
			}
		})
	}
}
