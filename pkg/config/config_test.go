package config_test

import (
	"os"
	"testing"

	"github.com/bbengfort/cosmos/pkg/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var testEnv = map[string]string{
	"COSMOS_MAINTENANCE":        "true",
	"COSMOS_LOG_LEVEL":          "debug",
	"COSMOS_CONSOLE_LOG":        "true",
	"COSMOS_BIND_ADDR":          ":443",
	"COSMOS_DATABASE_URL":       "postgres://localhost:5432/cosmos?sslmode=disable",
	"COSMOS_DATABASE_READ_ONLY": "true",
	"COSMOS_AUTH_TOKEN_KEYS":    "26eyytIJDwcp4ldjJTGVmsl3nEl:testdata/key1.pem,26eyzpOJMcQQCBWmIcS4mUscBPZ:testdata/key2.pem",
	"COSMOS_AUTH_AUDIENCE":      "cosmos.bengfort.com:443",
}

func TestConfig(t *testing.T) {
	// Set required environment variables and cleanup after
	prevEnv := curEnv()
	t.Cleanup(func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	})
	setEnv()

	conf, err := config.New()
	require.NoError(t, err)
	require.False(t, conf.IsZero())

	// Test configuration set from the environment
	require.Equal(t, true, conf.Maintenance)
	require.Equal(t, zerolog.DebugLevel, conf.GetLogLevel())
	require.Equal(t, true, conf.ConsoleLog)
	require.Equal(t, testEnv["COSMOS_BIND_ADDR"], conf.BindAddr)
	require.Equal(t, testEnv["COSMOS_DATABASE_URL"], conf.Database.URL)
	require.Equal(t, true, conf.Database.ReadOnly)
	require.Len(t, conf.Auth.TokenKeys, 2)
	require.Equal(t, testEnv["COSMOS_AUTH_AUDIENCE"], conf.Auth.Audience)
}

func TestRequiredConfig(t *testing.T) {
	required := []string{
		"COSMOS_DATABASE_URL",
		"COSMOS_AUTH_TOKEN_KEYS",
	}

	// Collect required environment variables and cleanup after
	prevEnv := curEnv(required...)
	cleanup := func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}
	t.Cleanup(cleanup)

	// Ensure that we've captured the complete set of required environment variables
	setEnv(required...)
	conf, err := config.New()
	require.NoError(t, err)

	// Ensure that each environment variable is required
	for _, envvar := range required {
		// Add all environment variables but the current one
		for _, key := range required {
			if key == envvar {
				os.Unsetenv(key)
			} else {
				setEnv(key)
			}
		}

		_, err := config.New()
		require.Errorf(t, err, "expected %q to be required but no error occurred", envvar)
	}

	// Test required configuration
	require.Len(t, conf.Auth.TokenKeys, 2)
}

// Returns the current environment for the specified keys, or if no keys are specified
// then returns the current environment for all keys in testEnv.
func curEnv(keys ...string) map[string]string {
	env := make(map[string]string)
	if len(keys) > 0 {
		for _, envvar := range keys {
			if val, ok := os.LookupEnv(envvar); ok {
				env[envvar] = val
			}
		}
	} else {
		for key := range testEnv {
			env[key] = os.Getenv(key)
		}
	}

	return env
}

// Sets the environment variable from the testEnv, if no keys are specified, then sets
// all environment variables from the test env.
func setEnv(keys ...string) {
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := testEnv[key]; ok {
				os.Setenv(key, val)
			}
		}
	} else {
		for key, val := range testEnv {
			os.Setenv(key, val)
		}
	}
}
