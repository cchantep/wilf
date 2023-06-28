package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	version string = "0" // Specified as build time: -ldflags '-X main.version=...'
)

func main() {
	args := os.Args[1:]

	commandArgs, reportings, err := ParseArguments(args)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		PrintUsage()
		os.Exit(1)
		return
	}

	if commandArgs.PrintUsage {
		PrintUsage()
		os.Exit(0)
		return
	}

	if commandArgs.Verbose {
		log.SetLevel(log.DebugLevel)
		log.Debugf("Command arguments: %s", commandArgs)
	}

	var config *Config = nil

	if commandArgs.Config != "" {
		config, err = LoadConfig(commandArgs.Config)

		if err != nil {
			fmt.Fprintf(os.Stderr, "fails to load configuration '%s': %s",
				commandArgs.Config,
				err.Error(),
			)

			os.Exit(2)

			return
		}
	}

	settings := DefaultSettings()

	if config != nil {
		settings = *config.Settings
	}

	file, err := os.Open(commandArgs.Pipfile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "fails to open Pipfile '%s': %s",
			commandArgs.Pipfile, err.Error())

		os.Exit(3)

		return
	}

	defer file.Close()

	pipfile, err := ParsePipfile(file)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
		return
	}

	checker := CreateCompositeChecker(config)

	reportUpdates := func(
		dependencies Dependencies,
		kind DependencyKind,
	) (bool, error) {
		globalFatal := false

		for _, reporting := range reportings {
			fatal, err := ReportUpdates(
				dependencies,
				kind,
				settings.UpdateLevel,
				settings.ExcludedPackages,
				checker,
				reporting.Reporter,
				reporting.Output,
			)

			if err != nil {
				return fatal, err
			}

			if fatal {
				globalFatal = true
			}
		}

		return globalFatal, nil
	}

	for _, reporting := range reportings {
		reporting.Reporter.Before(reporting.Output)
	}

	log.Debugln("Checking runtime Dependencies ...")

	requiresUpdates, err := reportUpdates(
		pipfile.RuntimeDependencies,
		RunDependency,
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(5)
		return
	}

	if settings.CheckDevPackages {
		log.Debugln("Checking dev dependencies ...")

		reqDevUpdates, err := reportUpdates(
			pipfile.DevDependencies,
			DevDependency,
		)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(6)
			return
		}

		if reqDevUpdates {
			requiresUpdates = true
		}
	}

	for _, reporting := range reportings {
		reporting.Reporter.After(reporting.Output)
	}

	if !requiresUpdates {
		log.Debugf("no updates required")

		os.Exit(0)
		return
	}

	os.Exit(5)
}
