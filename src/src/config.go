package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Token     string `yaml:"token"`
	TLS       bool   `yaml:"tls"`
	VerifyTLS bool   `yaml:"verify_tls"`
	// Name of authentication method
	AuthMethod string `yaml:"auth_method"`
	// Type of the authentication backend
	AuthBackend string `yaml:"auth_backend"`
	// Github org
	GithubOrg string `yaml:"github_org"`
	// Github Personal Access Token
	GithubPAT string `yaml:"github_pat"`
	Path      string
}

func LoadConfig() (Config, error) {

	cfg = Config{
		Host:        "127.0.0.1",
		Port:        8200,
		Token:       "password",
		TLS:         true,
		VerifyTLS:   true,
		AuthMethod:  "token",
		AuthBackend: "token",
	}

	var err error

	cfg.Path, err = GetConfigPath()
	if err != nil {
		return cfg, err
	}

	file, err := os.Stat(cfg.Path)
	if err != nil {
		return cfg, err
	}

	// Ensure that the config file is only readable by the user.
	// And not by his group or others (-rwx------)
	cfgFilePerm := file.Mode().String()
	if !strings.HasSuffix(cfgFilePerm, "------") {
		return cfg, fmt.Errorf("Your config file %q is accessible by others.\nYou can fix this by issuing:\n\n  $ chmod 700 %q\n", cfg.Path, cfg.Path)
	}

	content, err := ioutil.ReadFile(cfg.Path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func ComposeUrl() string {

	protocol := "http"
	if cfg.TLS {
		protocol = "https"
	}

	return fmt.Sprintf("%v://%v:%v", protocol, cfg.Host, cfg.Port)
}

// Update the token in the configuration file
func UpdateConfigToken(token string) error {

	// Reauthenticate against Vault and update in-memory config
	vc.SetToken(token)
	vc.Auth()
	cfg.Token = token

	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	token_found := false

	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "token:") {
			lines[i] = "token: " + token
			token_found = true
		}
	}

	if !token_found {
		lines = append(lines, "token: "+token)

	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0600)
	if err != nil {
		return err
	}

	return nil
}

func GetConfigPath() (string, error) {

	path := os.Getenv("VAULT_CLIENT_CONFIG")

	if path != "" {
		path, err := filepath.Abs(path)
		if err != nil {
			return "", errors.New("Unable to determine absolute path to config specified in $VAULT_CLIENT_CONFIG")
		}
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", errors.New("Unable to determine user home to locate ~/.vaultrc")
	}

	return usr.HomeDir + "/.vaultrc", nil
}
