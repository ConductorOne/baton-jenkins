package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configuration is the top-level configuration type.
type Configuration struct {
	IntegrationTestNumber string `mapstructure:"INTEGRATION_TEST_NUMBER"`
	BitbucketdcUsername   string `mapstructure:"BATON_BITBUCKETDC_USERNAME"`
	BitbucketdcPassword   string `mapstructure:"BATON_BITBUCKETDC_PASSWORD"`
	BitbucketdcBaseUrl    string `mapstructure:"BATON_BITBUCKETDC_BASE_URL"`
}

// LoadConfig reads configuration from environment variables
// and unmarshalls it into a Configuration instance.
func LoadConfig() (*Configuration, error) {
	var config *Configuration
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %v", err.Error())
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %v", err.Error())
	}

	return config, nil
}
