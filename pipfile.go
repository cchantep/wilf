package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"golang.org/x/mod/semver"
)

type VersionConstraint = [2]string
type VersionRequirement = []VersionConstraint
type Dependencies = map[string]VersionRequirement

type Pipfile struct {
	RuntimeDependencies   Dependencies
	DevDependencies       Dependencies
	RequiresPythonVersion VersionRequirement
}

func ParsePipfile(reader io.Reader) (Pipfile, error) {
	pipfile := Pipfile{
		RuntimeDependencies:   make(Dependencies),
		DevDependencies:       make(Dependencies),
		RequiresPythonVersion: VersionRequirement{},
	}

	scanner := bufio.NewScanner(reader)
	var currentSection string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Check for comment and remove it
		if commentIndex := strings.Index(line, "#"); commentIndex != -1 {
			line = line[:commentIndex]
		}

		// Check for empty line or comment
		if line == "" {
			continue
		}

		// Check for section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.ToLower(line[1 : len(line)-1])
			continue
		}

		if currentSection == "[source]" {
			continue
		}

		// Parse key-value pairs
		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			return Pipfile{}, fmt.Errorf("invalid line format: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if currentSection == "requires" {
			if key == "python_version" {
				value = strings.Trim(value, `"'`)
				versionReq, err := ParseVersionRequirement(value)

				if err != nil {
					return Pipfile{}, err
				}

				pipfile.RequiresPythonVersion = versionReq

				continue
			}
		}

		if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
			// Value is a dictionary, parse it
			dict, err := ParsePythonDict(value)

			if err != nil {
				return Pipfile{}, fmt.Errorf("error parsing dictionary: %v: %s", err, value)
			}

			if version, ok := dict["version"]; ok {
				value = version
			} else if _, ok := dict["path"]; ok {
				// Ignore missing version if path is set to "."
				continue
			} else {
				return Pipfile{}, fmt.Errorf("Field 'version' not found in specification: %s", key)
			}
		} else {
			// Remove surrounding quotes from the value
			value = strings.Trim(value, `"'`)
		}

		versionReq, err := ParseVersionRequirement(value)

		if err != nil {
			return Pipfile{}, err
		}

		switch currentSection {
		case "packages":
			pipfile.RuntimeDependencies[key] = versionReq
		case "dev-packages":
			pipfile.DevDependencies[key] = versionReq
		default:
			log.Debugf("Ignoring unknown section '%s'\n", currentSection)
		}
	}

	if err := scanner.Err(); err != nil {
		return Pipfile{}, err
	}

	return pipfile, nil
}

func ParsePythonDict(input string) (map[string]string, error) {
	chars := []rune(input[1 : len(input)-1])

	inQuote := false
	inArray := false

	for i, ch := range chars {
		if ch == '\'' || ch == '"' {
			inQuote = !inQuote
		} else if ch == '[' && !inArray {
			inArray = true
		} else if ch == ']' && inArray {
			inArray = false
		} else if !inQuote && !inArray && ch == ',' {
			if i+1 < len(chars) {
				chars[i] = '\r'
				chars[i+1] = '\n'
			} else {
				chars[i] = '\n'
			}
		}
	}

	input = string(chars)

	result := make(map[string]interface{})

	// Unmarshal the input using yaml
	_, err := toml.Decode(input, &result)

	if err != nil {
		return nil, err
	}

	dict := make(map[string]string)

	for key, val := range result {
		dict[key] = fmt.Sprintf("%v", val)
	}

	return dict, nil
}

// TODO: Aside from '*' support V.*
func ParseVersionRequirement(input string) (VersionRequirement, error) {
	operators := []string{"<=", "<", "!=", "==", ">=", ">", "~=", "==="}
	var verReq VersionRequirement

NEXT_SPEC:
	for _, spec := range strings.Split(input, ",") {
		spec := strings.TrimSpace(spec)

		if spec == "*" {
			verReq = append(verReq, VersionConstraint{"*", "*"})
			continue
		}

		for _, operator := range operators {
			if strings.HasPrefix(spec, operator) {
				remaining := strings.TrimPrefix(spec, operator)

				if remaining == "" {
					return VersionRequirement{}, fmt.Errorf("missing version: %s", spec)
				}

				ver := fmt.Sprintf("v%s", remaining)

				n := strings.ReplaceAll(ver, "*", "0")

				if n != ver {
					// Was normalized as matching expr
					if !IsValidVersion(n) {
						return VersionRequirement{}, fmt.Errorf("invalid version: %s", spec)
					}

					if operator == "==" {
						operator = "~"
					} else if operator == "!=" {
						operator = "!~"
					} else {
						return VersionRequirement{}, fmt.Errorf("invalid version matching: %s", spec)
					}
				} else if !IsValidVersion(ver) {
					return VersionRequirement{}, fmt.Errorf("invalid version: %s", remaining)
				}

				verReq = append(
					verReq,
					VersionConstraint{operator, ver},
				)

				continue NEXT_SPEC
			}
		}

		ver := fmt.Sprintf("v%s", spec)

		if !IsValidVersion(ver) {
			return VersionRequirement{}, fmt.Errorf("invalid version: %s", spec)
		}

		verReq = append(verReq, VersionConstraint{"==", ver})
	}

	return verReq, nil
}

func IsValidVersion(version string) bool {
	// Check non-standard version like zstd 1.5.5.1 (https://pypi.org/project/zstd/)
	// which is not supported by semver
	if match, _ := regexp.MatchString(`^v\d+\.\d+\.\d+(\.\d+)*$`, version); match {
		return true
	}

	return semver.IsValid(version)
}
