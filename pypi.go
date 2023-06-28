package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ProjectInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Summary string `json:"summary"`
	HomeURL string `json:"home_page"`
}

func GetProjectInfo(packageName string) (*ProjectInfo, error) {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", packageName)

	// Send GET request to the API
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into a ProjectInfo struct
	var projectInfo struct {
		Info ProjectInfo `json:"info"`
	}

	err = json.Unmarshal(body, &projectInfo)

	if err != nil {
		return nil, err
	}

	// Check if the 'info' field is nil
	if projectInfo.Info.Name != "" {
		projectInfo.Info.Version = fmt.Sprintf(
			"v%s", projectInfo.Info.Version)

		return &projectInfo.Info, nil
	}

	// ---

	var jsonResp struct {
		Message string `json:"message"`
	}

	err2 := json.Unmarshal(body, &jsonResp)

	if err2 == nil && jsonResp.Message == "Not Found" {
		return nil, nil
	}

	errMsg := err.Error()

	if err2 == nil {
		errMsg = jsonResp.Message
	}

	return nil, errors.New(fmt.Sprintf("Project information not found in the JSON response: %s", errMsg))
}
