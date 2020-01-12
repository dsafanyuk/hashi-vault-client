package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/mitchellh/cli"
)

// Interface to easily add new authentication backends.
type AuthBackend interface {
	Ask() (*http.Request, error)
	// AskGithub() (*http.Request, error)
}

type LDAPAuth struct {
	ui cli.Ui
}
type GithubAuth struct {
	ui cli.Ui
}

func (l LDAPAuth) Ask() (*http.Request, error) {

	username, err := l.ui.Ask("Username:")
	if err != nil {
		return new(http.Request), err
	}

	password, err := l.ui.AskSecret("Password:")
	if err != nil {
		return new(http.Request), err
	}

	body := []byte(fmt.Sprintf(`{"password":"%s"}`, password))
	url := fmt.Sprintf("%v/v1/auth/%s/login/%s", ComposeUrl(), cfg.AuthMethod, username)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	return req, nil
}

func (g GithubAuth) Ask() (*http.Request, error) {

	githubOrg := cfg.GithubOrg
	if len(githubOrg) < 0 {
		org, err := g.ui.Ask("What's your github org?")
		if err != nil {
			return new(http.Request), err
		}
		githubOrg = org
	}

	githubPAT := cfg.GithubPAT
	if len(githubOrg) < 0 {
		token, err := g.ui.AskSecret("Github PAT:")
		if err != nil {
			return new(http.Request), err
		}
		githubPAT = token
	}

	body := []byte(fmt.Sprintf(`{"token":"%s"}`, githubPAT))
	url := fmt.Sprintf("%v/v1/auth/%s_%s/login", ComposeUrl(), cfg.AuthMethod, githubOrg)
	fmt.Println(url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return new(http.Request), err
	}
	return req, nil
}
