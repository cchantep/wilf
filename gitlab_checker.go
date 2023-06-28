package main

// GitlabChecker represents a struct that holds the configuration for a GitLab registry.
type GitlabChecker struct {
	Config GitlabRegistryConfig
}

// RequiredUpdate checks if a package requires an update and returns the current version,
// update level, home URL and error (if any).
func (c GitlabChecker) RequiredUpdate(
	pkg string,
	requirement VersionRequirement,
) (string, UpdateLevel, string, error) {
	info, err := GetGitlabProjectInfo(c.Config, pkg)

	if err != nil {
		return "", 0, "", err
	}

	if info == nil {
		return "", 0, "", nil
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
