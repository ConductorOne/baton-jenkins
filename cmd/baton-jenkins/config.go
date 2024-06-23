package main

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/cli"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig  `mapstructure:",squash"` // Puts the base config options in the same place as the connector options
	JenkinsUsername string                   `mapstructure:"jenkins-username" description:"Username of administrator used to connect to the Jenkins API."`
	JenkinsPassword string                   `mapstructure:"jenkins-password" description:"Application password used to connect to the Jenkins API."`
	JenkinstBaseUrl string                   `mapstructure:"jenkins-baseurl" description:"Jenkins. example http://localhost:8080." defaultValue:"http://localhost:8080"`
	JenkinsToken    string                   `mapstructure:"jenkins-token" description:"HTTP access tokens in Jenkins"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.JenkinstBaseUrl == "" {
		return fmt.Errorf("jenkins-baseurl must be provided")
	}

	if cfg.JenkinsUsername == "" {
		return fmt.Errorf("jenkins-username must be provided")
	}

	if cfg.JenkinsToken == "" {
		if cfg.JenkinsUsername == "" || cfg.JenkinsPassword == "" {
			return fmt.Errorf("either bitbucketdc-token or (bitbucketdc-username/bitbucketdc-password) must be provided")
		}
	}

	if cfg.JenkinsToken != "" && cfg.JenkinsUsername != "" && cfg.JenkinsPassword != "" {
		return fmt.Errorf("bitbucketdc-token, and (bitbucketdc-username/bitbucketdc-password) cannot be provided simultaneously")
	}

	return nil
}
