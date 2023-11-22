package config

import (
	"fmt"

	"github.com/bbengfort/cosmos/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/rotationalio/confire"
	"github.com/rs/zerolog"
)

const Prefix = "cosmos"

type Config struct {
	Maintenance  bool                `default:"false" desc:"sets the server to maintenance mode if true"`
	BindAddr     string              `split_words:"true" default:":8888" desc:"the ip address and port to bind the server to""`
	Mode         string              `default:"release" desc:"one of debug, test, or release"`
	LogLevel     logger.LevelDecoder `split_words:"true" default:"info" desc:"the verbosity of logging"`
	ConsoleLog   bool                `split_words:"true" default:"false" desc:"human readable instead of json logging"`
	AllowOrigins []string            `split_words:"true" default:"http://localhost:8888" desc:"origin of website accessing API"`
	processed    bool                // set when the config is properly processed from the environment
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
