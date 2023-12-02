package config

import (
	"fmt"
	"time"

	"github.com/bbengfort/cosmos/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/rotationalio/confire"
	"github.com/rs/zerolog"
)

const Prefix = "cosmos"

type Config struct {
	Maintenance  bool                `default:"false" desc:"sets the server to maintenance mode if true"`
	BindAddr     string              `split_words:"true" default:":8888" desc:"the ip address and port to bind the server to"`
	Mode         string              `default:"release" desc:"one of debug, test, or release"`
	LogLevel     logger.LevelDecoder `split_words:"true" default:"info" desc:"the verbosity of logging"`
	ConsoleLog   bool                `split_words:"true" default:"false" desc:"human readable instead of json logging"`
	AllowOrigins []string            `split_words:"true" default:"http://localhost:3000" desc:"origin of website accessing API"`
	Database     DatabaseConfig      `desc:"database configuration"`
	Auth         AuthConfig          `desc:"authentication and claims issuer configuration"`
	processed    bool                // set when the config is properly processed from the environment
}

type DatabaseConfig struct {
	URL      string `default:"postgres://localhost:5432/cosmos?sslmode=disable" required:"true" desc:"specify the connection to the database via a DSN"`
	ReadOnly bool   `split_words:"true" default:"false" desc:"open the database in readonly mode"`
	Testing  bool   `default:"false" desc:"if set to true, opens a sql mock rather than an actual db connection"`
}

type AuthConfig struct {
	Keys            map[string]string `desc:"a map of key id to key path on disk"`
	Audience        string            `default:"http://localhost:3000" desc:"value for the aud jwt claim"`
	Issuer          string            `default:"http://localhost:3000" desc:"value for the iss jwt claim"`
	CookieDomain    string            `split_words:"true" default:"localhost" desc:"limit the cookies to the specified domain (same as allowed origins)"`
	AccessTokenTTL  time.Duration     `split_words:"true" default:"24h" desc:"the amount of time before an access token expires"`
	RefreshTokenTTL time.Duration     `split_words:"true" default:"48h" desc:"the amount of time before a refresh token expires"`
	TokenOverlap    time.Duration     `split_words:"true" default:"-1h" desc:"the amount of overlap between the access and refresh token"`
}

func New() (conf Config, err error) {
	if err = confire.Process(Prefix, &conf); err != nil {
		return Config{}, err
	}

	if err = conf.Validate(); err != nil {
		return Config{}, err
	}

	conf.processed = true
	return conf, nil
}

// Returns true if the config has not been correctly processed from the environment.
func (c Config) IsZero() bool {
	return !c.processed
}

// Custom validations are added here, particularly validations that require one or more
// fields to be processed before the validation occurs.
// NOTE: ensure that all nested config validation methods are called here.
func (c Config) Validate() (err error) {
	if c.Mode != gin.ReleaseMode && c.Mode != gin.DebugMode && c.Mode != gin.TestMode {
		return fmt.Errorf("invalid configuration: %q is not a valid gin mode", c.Mode)
	}
	return nil
}

func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}
