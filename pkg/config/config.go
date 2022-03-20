package config

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type Config struct {
	Maintenance bool         `split_words:"true" default:"false"`
	LogLevel    LevelDecoder `split_words:"true" default:"info"`
	ConsoleLog  bool         `split_words:"true" default:"false"`
	BindAddr    string       `split_words:"true" default:"10001"`
	Database    DatabaseConfig
	Auth        AuthConfig
	processed   bool
}

type DatabaseConfig struct {
	URL      string `split_words:"true" required:"true"`
	ReadOnly bool   `split_words:"true" default:"false"`
}

type AuthConfig struct {
	Audience  string            `split_words:"true" default:"localhost:10001"`
	TokenKeys map[string]string `split_words:"true"`
}

func New() (conf Config, err error) {
	if err = envconfig.Process("cosmos", &conf); err != nil {
		return Config{}, err
	}

	if err = conf.Validate(); err != nil {
		return Config{}, err
	}

	conf.processed = true
	return conf, nil
}

func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}

func (c Config) IsZero() bool {
	return !c.processed
}

// Mark a manually constructed as processed as long as it is validated.
func (c Config) Mark() (Config, error) {
	if err := c.Validate(); err != nil {
		return c, err
	}
	c.processed = true
	return c, nil
}

func (c Config) Validate() (err error) {
	if err = c.Auth.Validate(); err != nil {
		return err
	}
	return nil
}

func (c AuthConfig) Validate() error {
	if len(c.TokenKeys) == 0 {
		return errors.New("invalid configuration: at least one token key is required")
	}
	return nil
}
