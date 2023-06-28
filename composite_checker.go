package main

type CompositeChecker []Checker

func (c CompositeChecker) RequiredUpdate(
	pkg string,
	requirement VersionRequirement,
) (string, UpdateLevel, string, error) {
	for _, checker := range c {
		latest, level, url, err := checker.RequiredUpdate(pkg, requirement)

		if err != nil {
			return "", 0, "", err
		}

		if level > 0 {
			return latest, level, url, nil
		}
	}

	return "", 0, "", nil
}

// CreateCompositeChecker creates a new CompositeChecker with a PypiChecker as the first element.
// If config.Gitlab is not nil, a corresponding instance of GitlabChecker is appended to the CompositeChecker.
func CreateCompositeChecker(config *Config) CompositeChecker {
	checkers := CompositeChecker{&PypiChecker{}}

	if config != nil {
		if config.Gitlab != nil {
			checkers = append(checkers, &GitlabChecker{
				Config: *config.Gitlab,
			})
		}
	}

	return checkers
}
