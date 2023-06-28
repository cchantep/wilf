package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// CommandArguments represents the parsed command line arguments.
type CommandArguments struct {
	Verbose    bool
	Config     string
	Pipfile    string
	PrintUsage bool
	Reporters  string // comma separated list of reporters
}

type Reporting struct {
	Reporter UpdateReporter
	Output   io.Writer
}

// ParseArguments parses the command line arguments and returns the parsed arguments,
// an update reporter, an output writer and an error.
func ParseArguments(args []string) (CommandArguments, []Reporting, error) {
	var verbose bool
	var config string
	var pipfile string
	var printUsage bool
	var reporters []string

	if len(args) == 0 {
		return CommandArguments{}, []Reporting{},
			fmt.Errorf("please provide the path to the Pipfile as an argument")
	}

	for i := 0; i < len(args); i++ {
		if args[i] == "-v" {
			verbose = true
		} else if args[i] == "-c" {
			if i+1 >= len(args) {
				return CommandArguments{}, []Reporting{},
					fmt.Errorf("error: -c option requires a value")
			}

			config = args[i+1]

			i++
		} else if args[i] == "-h" {
			printUsage = true
		} else if args[i] == "-r" {
			if i+1 >= len(args) {
				return CommandArguments{}, []Reporting{},
					fmt.Errorf("error: -r option requires a value")
			}

			reporters = append(reporters, args[i+1])

			i++
		} else if strings.HasPrefix(args[i], "-") {
			return CommandArguments{}, []Reporting{}, fmt.Errorf("error: invalid option %s", args[i])
		} else if pipfile != "" {
			return CommandArguments{}, []Reporting{}, fmt.Errorf("error: only one Pipfile can be provided")
		} else {
			pipfile = args[i]
		}
	}

	if printUsage {
		return CommandArguments{
			PrintUsage: true,
		}, []Reporting{}, nil
	}

	if pipfile == "" {
		return CommandArguments{}, []Reporting{},
			fmt.Errorf("please provide the path to the Pipfile as an argument")
	}

	rl := len(reporters)

	if rl == 0 {
		reporters = []string{"colorized-table"}
	}

	var updateReporters []Reporting

	if rl == 0 {
		updateReporters = []Reporting{}
	} else {
		updateReporters = make([]Reporting, rl)

		for i, reporter := range reporters {
			updateReporter, output, err := createReporter(reporter)

			if err != nil {
				return CommandArguments{}, []Reporting{}, err
			}

			updateReporters[i] = Reporting{
				Reporter: updateReporter,
				Output:   output,
			}
		}
	}

	return CommandArguments{
		Verbose:    verbose,
		Config:     config,
		Pipfile:    pipfile,
		PrintUsage: printUsage,
		Reporters:  strings.Join(reporters, ", "),
	}, updateReporters, nil
}

func (args CommandArguments) String() string {
	return fmt.Sprintf("{Verbose: %v, Config: '%s', Pipfile: '%s', Reporter: '%s'}",
		args.Verbose,
		args.Config,
		args.Pipfile,
		args.Reporters,
	)
}

func PrintUsage() {
	fmt.Println("Usage: wilf [OPTIONS] /path/to/Pipefile")
	fmt.Println("Options:")
	fmt.Println("  -c FILE      Use FILE as the configuration file")
	fmt.Println("  -h           Print this help message and exit")
	fmt.Println("  -r REPORTER  Use REPORTER as the reporter. It can be specified multi time to apply multiple reporters. Valid options are monochrome-table, colorized-table (default), junit or junit:/path/to/output/junit.xml")
	fmt.Println("  -v           Enable verbose output")
}

func createReporter(reporter string) (UpdateReporter, io.Writer, error) {
	validReporters := []string{"monochrome-table", "colorized-table", "junit"}

	if strings.HasPrefix(reporter, "junit:") {
		path := strings.TrimPrefix(reporter, "junit:")

		_, err := os.Stat(path)

		if err != nil && !os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("invalid path for junit reporter: %s: %s", path, err)
		}

		out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

		if err != nil {
			return nil, nil, err
		}

		return &JUnitReporter{Version: version}, out, nil
	}

	for _, r := range validReporters {
		if r != reporter {
			continue
		}

		// ---

		switch r {
		case "monochrome-table":
			return MonochromeTableReporter(version), os.Stdout, nil

		case "colorized-table":
			return &ColorizedTableReporter{Version: version}, os.Stdout, nil

		case "junit":
			return &JUnitReporter{Version: version}, os.Stdout, nil

		default:
			return nil, nil, fmt.Errorf("unknown reporter: %s", reporter)
		}
	}

	return nil, nil, fmt.Errorf("invalid reporter: %s", reporter)
}
