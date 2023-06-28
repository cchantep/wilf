package main

import "fmt"

type UpdateLevel uint

const (
	Patch UpdateLevel = 1
	Minor UpdateLevel = 2
	Major UpdateLevel = 3
)

func (l UpdateLevel) String() string {
	if l == Patch {
		return "patch"
	} else if l == Minor {
		return "minor"
	} else if l == Major {
		return "major"
	} else if l == 0 {
		return "<none>"
	}

	return "<undefined>"
}

// ParseUpdateLevel parses a string representation of an UpdateLevel
// and returns the corresponding UpdateLevel.
// If the string representation is not valid, an error is returned.
func ParseUpdateLevel(s string) (UpdateLevel, error) {
	switch s {
	case "patch":
		return Patch, nil
	case "minor":
		return Minor, nil
	case "major":
		return Major, nil
	default:
		return 0, fmt.Errorf("invalid UpdateLevel: %s", s)
	}
}
