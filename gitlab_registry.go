package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
)

// GitlabRegistryConfig represents the configuration
// for accessing Gitlab's registry API.
type GitlabRegistryConfig struct {
	ProjectApiPackagesUrl string `toml:"project_api_packages_url"`
	PrivateToken          string `toml:"private_token"`
}

type gitlabProjectLinks struct {
	WebPath string `json:"web_path"`
}

type gitlabProjectInfo struct {
	Name    string             `json:"name"`
	Version string             `json:"version"`
	Links   gitlabProjectLinks `json:"_links"`
}

// LoadGitlabRegistryConfig loads the Gitlab registry configuration from a TOML file.
// It takes a file path as input and returns a pointer to a GitlabRegistryConfig struct and an error.
func LoadGitlabRegistryConfig(path string) (*GitlabRegistryConfig, error) {
	var config struct {
		Gitlab GitlabRegistryConfig `toml:"gitlab"`
	}

	_, err := toml.DecodeFile(path, &config)

	if err != nil {
		return nil, err
	}

	return &config.Gitlab, nil
}

// GetGitlabProjectInfo retrieves information about a project from Gitlab's registry API.
// It takes a GitlabRegistryConfig struct and a package name as input.
// It returns a pointer to a ProjectInfo struct and an error.
// If the package is not found, it returns nil and no error.
// If the package is found but there are multiple projects with the same name,
// it returns an error.
// If there is an error while retrieving the package information,
// it returns an error.
func GetGitlabProjectInfo(
	gitlabConfig GitlabRegistryConfig,
	packageName string,
) (*ProjectInfo, error) {
	url := fmt.Sprintf(
		"%s?package_type=pypi&package_name=%s",
		gitlabConfig.ProjectApiPackagesUrl,
		packageName,
	)

	// Create a new HTTP request with the Gitlab API URL
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	// Set the Gitlab API token in the request header
	if gitlabConfig.PrivateToken != "" {
		request.Header.Set("PRIVATE-TOKEN", gitlabConfig.PrivateToken)
	}

	// Send the HTTP request
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into a slice of ProjectInfo structs
	var projectInfo []gitlabProjectInfo

	err = json.Unmarshal(body, &projectInfo)

	if err != nil {
		var jsonResp struct {
			Message string `json:"message"`
		}

		err2 := json.Unmarshal(body, &jsonResp)

		if err2 != nil {
			return nil, err
		}

		if err2 == nil && jsonResp.Message == "Not Found" {
			return nil, nil
		}

		errMsg := err.Error()

		if err2 == nil {
			errMsg = jsonResp.Message
		}

		return nil, errors.New(fmt.Sprintf("Project information not found in the JSON response: %s", errMsg))
	}

	// ---

	// Check if the response is empty
	if len(projectInfo) == 0 {
		return nil, nil
	}

	// Check if the response contains more than one element
	if len(projectInfo) > 1 {
		return nil, errors.New("The Gitlab API indicates more than one project with the same name")
	}

	// ---

	urlParts := strings.SplitAfterN(gitlabConfig.ProjectApiPackagesUrl, "/", 4)
	homeUrl := fmt.Sprintf("%s%s",
		strings.Join(urlParts[0:3], ""),
		strings.TrimPrefix(projectInfo[0].Links.WebPath, "/"),
	)

	// Return the first element of the slice
	return &ProjectInfo{
		Name:    projectInfo[0].Name,
		Version: fmt.Sprintf("v%s", projectInfo[0].Version),
		Summary: "",
		HomeURL: homeUrl,
	}, nil
}
