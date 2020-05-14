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
	"os"
	"strings"
	"testing"
	"time"

	"github.com/boseji/gau/sioutils"
	"golang.org/x/sync/errgroup"
)

func TestTagReadTrialErrGroup(t *testing.T) {
	t.Log("Running the Test using Error Group")
	timeout := time.Second * 10
	owner := "boseji"
	repo := "gau"

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Create the ErrGroup
	eg, _ := errgroup.WithContext(ctx)
	// One Buffer so that we can return from function
	data := make(chan string, 1)

	eg.Go(func() error {
		uri, err := prepareRepoPath(repositoryAPI, owner, repo, latestReleaseAPI)
		if err != nil {
			return err
		}

		body, status, err := MakeGetRequest(uri)
		if err != nil {
			return err
		}
		fmt.Println("Fetch for ", uri, ":", status)
		// Send through the Output Channel
		data <- body
		return nil
	})

	if err := eg.Wait(); err != nil {
		t.Errorf("\nFailed to Fetch URL\n= %s", err)
		return
	}

	// Receive the Data
	s := <-data

	tag, err := filterTagFromBody(s)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(tag)
}

func TestGetLatestRelease(t *testing.T) {
	// Open the Release File
	f, err := os.Open("../RELEASE")
	if err != nil {
		t.Log(err)
		return
	}
	_, releaseVersion, err := sioutils.ReadAll(f)
	if err != nil {
		t.Log(err)
		return
	}
	releaseVersion = strings.TrimSpace(releaseVersion)

	// Test Cases
	type args struct {
		owner   string
		repo    string
		timeout time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"Positive Test Case1",
			args{
				owner:   "boseji",
				repo:    "gau",
				timeout: 10 * time.Second,
			},
			releaseVersion,
			false,
		},
		{
			"Negative Test Case1 URI Error",
			args{
				owner:   "boseji\n",
				repo:    "gau",
				timeout: 10 * time.Second,
			},
			"",
			true,
		},
		{
			"Negative Test Case2 No Releases",
			args{
				owner:   "boseji",
				repo:    "dotfiles",
				timeout: 10 * time.Second,
			},
			"",
			true,
		},
		{
			"Negative Test Case3 Invalid Repo",
			args{
				owner:   "boseji",
				repo:    "dotfiles1",
				timeout: 10 * time.Second,
			},
			"",
			true,
		},
		{
			"Negative Test Case4 API Failure",
			args{
				owner:   "boseji",
				repo:    "dotfiles1?",
				timeout: 10 * time.Second,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestRelease(tt.args.owner, tt.args.repo, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLatestRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}
