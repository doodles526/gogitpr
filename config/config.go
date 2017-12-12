package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("base_url", "https://api.github.com")
	viper.SetDefault("application_name", "gogitpr")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("print", false)

	viper.SetConfigName("gogitpr") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // path to look for the config file in
	viper.ReadInConfig()     // Find and read the config file

	viper.SetEnvPrefix("gitpr")
	viper.AutomaticEnv()
}

type Config struct {
	// BaseURL is the URL to use for the Github API
	// if not set, defaults to https://api.github.com
	BaseURL string

	// GithubToken allows you to specify a token in case of accessing
	// private resources
	GithubToken string

	// ApplicationName allows you to specify an application name
	// in case you want to change the default - used in querying the
	// github API
	ApplicationName string

	// GithubOrg is which github organization to populate DB from
	GithubOrg string

	// GithubUser is which github user to populate DB from
	GithubUser string

	PrintResult bool

	// Logger instance
	Logger *logrus.Logger
}

func NewConfig() (*Config, error) {
	logger := logrus.New()
	logger.Level = getLogLevel()

	cfg := &Config{
		BaseURL:         viper.GetString("base_url"),
		GithubToken:     viper.GetString("github_token"),
		ApplicationName: viper.GetString("application_name"),
		GithubOrg:       viper.GetString("github_org"),
		GithubUser:      viper.GetString("github_user"),
		PrintResult:     viper.GetBool("print"),
		Logger:          logger,
	}

	return cfg, nil
}

func getLogLevel() logrus.Level {
	level, err := logrus.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		return logrus.InfoLevel
	}

	return level
}
