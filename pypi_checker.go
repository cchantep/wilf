package main

// PypiChecker is a struct that represents a PyPI checker.
type PypiChecker struct {
	PythonRequirement VersionRequirement
}

// RequiredUpdate checks if a package requires an update.
// It returns the current version of the package, the update level,
// the home URL of the package, and an error (if any).
func (c PypiChecker) RequiredUpdate(
	pkg string,
	requirement VersionRequirement,
) (string, UpdateLevel, string, error) {
	info, err := GetProjectInfo(pkg)

	if err != nil {
		return "", 0, "", err
	}

	if info == nil {
		return "", 0, "", nil
	}

	if info.RequiresPython != "" && len(c.PythonRequirement) > 0 {
		pkgPythonReq, err := ParseVersionRequirement(info.RequiresPython)

		if err != nil {
			return "", 0, "", err
		}

		if !AreCompatibles(c.PythonRequirement, pkgPythonReq) {
			// Python version are not compatible,
			// so package itself should not be updated
			return info.Version, 0, info.HomeURL, nil
		}
	}

	if !ShouldUpdate(requirement, info.Version) {
		return info.Version, 0, info.HomeURL, nil
	}

	lvl, err := CreateUpdateLevel(requirement, info.Version)

	if err != nil {
		return "", 0, "", err
	}

	return info.Version, lvl, info.HomeURL, nil
}
