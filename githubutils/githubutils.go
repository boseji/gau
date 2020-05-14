// ####
// ##
// ## Boseji's Github Utilities
// ##
// ## SPDX-License-Identifier: GPL-2.0-only
// ##
// ## Copyright (C) 2020 Abhijit Bose <boseji@users.noreply.github.com>
// ##
// ####

package githubutils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/boseji/gau/sioutils"
)

const (
	githubAPIURL     = "https://api.github.com"
	repositoryAPI    = "repos"
	latestReleaseAPI = "releases/latest"
	latestTagFilter  = "\"tag_name\":"
)

func prepareRepoPath(apiPrefix, owner, repo, postfix string) (string, error) {
	uri := githubAPIURL + "/" + apiPrefix + "/"
	uri += owner + "/" + repo + "/" + postfix
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	// Success
	return u.String(), nil
}

func filterTagFromBody(body string) (string, error) {
	if !strings.Contains(body, latestTagFilter) {
		return "", fmt.Errorf("Could not find the specific tag")
	}
	// Get the Tag like "tag_name":"1.2.3" for first location
	block := strings.SplitN(body[strings.Index(body, latestTagFilter):], ",", 2)
	if len(block) != 2 {
		return "", fmt.Errorf("Incorrect Tag format in the body")
	}
	// Remove the "tag_name"
	tag := strings.Split(block[0], ":")
	if len(tag) != 2 {
		return "", fmt.Errorf("Incorrect Tag Name format in the body")
	}
	// filter out the quotes
	ver := strings.Trim(tag[1], "\"")
	// Upon Success
	return ver, nil
}

// ProcessHTTPRequest perform the execution of the http.Request and returns the
//  read back 'body' the 'Status Code' and 'Errors'
// Note: This function uses the Default HTTP Client
//       `http.DefaultClient` to perform the request
func ProcessHTTPRequest(request *http.Request) (string, int, error) {
	// Make the Request
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	// Make Sure the Close the Body
	defer resp.Body.Close()
	// Check the Status
	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, fmt.Errorf("Unknown Status : %s", resp.Status)
	}
	// Read and Pass On the Body
	_, body, err := sioutils.ReadAll(resp.Body)
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	// Upon Success
	return body, resp.StatusCode, nil
}

// MakeGetRequest performs the HTTP GET request and returns the
//  read back 'body' the 'Status Code' and 'Errors'
// Note: This function uses the Default HTTP Client
//       `http.DefaultClient` to perform the request
func MakeGetRequest(uri string) (string, int, error) {
	// Create Request
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	return ProcessHTTPRequest(req)
}

// MakeGetRequestCtx performs the HTTP GET request and returns the
//  read back 'body' the 'Status Code' and 'Errors' depending on
//  the supplied context.
// Note: This function uses the Default HTTP Client
//       `http.DefaultClient` to perform the request
func MakeGetRequestCtx(ctx context.Context, uri string) (string, int, error) {
	// Create Request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	return ProcessHTTPRequest(req)
}

// GetLatestRelease fetches the latest release version from 'https://github.com'
//  based on the supplied 'owner' and 'repo' repository information.
//  It either returns the latest release tag without any error or
//  produces an empty tag with the error.
//  The 'timeout' parameter helps to regulate the time spent in making
//   the API request.
func GetLatestRelease(owner, repo string, timeout time.Duration) (string, error) {
	// Create a Timeout based Context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create the API URI
	uri, err := prepareRepoPath(repositoryAPI, owner, repo, latestReleaseAPI)
	if err != nil {
		return "", err
	}

	// Make the Request
	body, _, err := MakeGetRequestCtx(ctx, uri)
	if err != nil {
		return "", err
	}

	// Get the Tag
	tag, err := filterTagFromBody(body)
	if err != nil {
		return "", err
	}

	// Upon Success
	return tag, nil
}
