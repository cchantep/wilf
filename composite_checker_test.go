package main

import (
	"errors"
	"testing"
)

type mockChecker struct {
	pkg           string
	latestVersion string
	updateLevel   UpdateLevel
	url           string
	err           error
}

func (m mockChecker) RequiredUpdate(
	pkg string,
	requirement VersionRequirement,
) (string, UpdateLevel, string, error) {
	if pkg != m.pkg {
		return "", 0, "", nil
	}

	return m.latestVersion, m.updateLevel, m.url, m.err
}

func TestCompositeChecker(t *testing.T) {
	checker1 := mockChecker{
		pkg:           "test-pkg1",
		latestVersion: "v1.2.3",
		updateLevel:   Patch,
		url:           "http://foo/pkg1",
		err:           nil,
	}
	checker2 := mockChecker{
		pkg:           "test-pkg2",
		latestVersion: "v2.0.0",
		updateLevel:   Major,
		url:           "http://foo/pkg2",
		err:           nil,
	}
	checker3 := mockChecker{
		pkg:           "test-pkg3",
		latestVersion: "",
		updateLevel:   0,
		url:           "",
		err:           errors.New("failed to get latest version"),
	}

	compositeChecker := CompositeChecker{checker1, checker2, checker3}

	testCases := []struct {
		pkg            string
		requirement    VersionRequirement
		expectedLatest string
		expectedLevel  UpdateLevel
		expectedUrl    string
		expectedError  error
	}{
		{
			pkg:            "test-pkg1",
			requirement:    VersionRequirement{VersionConstraint{">=", "v1.2.2"}},
			expectedLatest: "v1.2.3",
			expectedLevel:  Patch,
			expectedUrl:    "http://foo/pkg1",
			expectedError:  nil,
		},
		{
			pkg:            "test-pkg2",
			requirement:    VersionRequirement{VersionConstraint{">=", "v3.0.0"}},
			expectedLatest: "v2.0.0",
			expectedLevel:  Major,
			expectedUrl:    "http://foo/pkg2",
			expectedError:  nil,
		},
		{
			pkg:            "test-pkg3",
			requirement:    VersionRequirement{VersionConstraint{">=", "v1.2.2"}},
			expectedLatest: "",
			expectedLevel:  0,
			expectedUrl:    "",
			expectedError:  errors.New("failed to get latest version"),
		},
	}

	for _, tc := range testCases {
		latest, level, url, err := compositeChecker.RequiredUpdate(tc.pkg, tc.requirement)

		if latest != tc.expectedLatest {
			t.Errorf("Expected latest version to be %s, but got %s", tc.expectedLatest, latest)
		}

		if level != tc.expectedLevel {
			t.Errorf("Expected update level to be %d, but got %d", tc.expectedLevel, level)
		}

		if url != tc.expectedUrl {
			t.Errorf("Expected url to be %s, but got %s", tc.expectedUrl, url)
		}

		if (err == nil && tc.expectedError != nil) || (err != nil && tc.expectedError == nil) || (err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error()) {
			t.Errorf("Expected error to be %v, but got %v", tc.expectedError, err)
		}
	}
}
