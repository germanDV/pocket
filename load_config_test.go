package pocket

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfigFromEnv(t *testing.T) {
	t.Run("all_defaults", func(t *testing.T) {
		cleanEnv()
		type MyConfig struct {
			Env         string        `env:"ENV" default:"dev"`
			Port        int           `env:"PORT" default:"8080"`
			EnableDebug bool          `env:"DEBUG" default:"false"`
			Timeout     time.Duration `env:"TIMEOUT" default:"5s"`
		}

		myConfig, err := LoadConfigFromEnv[MyConfig]()
		AssertNil(t, err)
		AssertEqual(t, myConfig.Env, "dev")
		AssertEqual(t, myConfig.Port, 8080)
		AssertEqual(t, myConfig.EnableDebug, false)
		AssertEqual(t, myConfig.Timeout, 5*time.Second)
	})

	t.Run("from_env", func(t *testing.T) {
		cleanEnv()
		os.Setenv("ENV", "production")
		os.Setenv("DEBUG", "true")
		type MyConfig struct {
			Env         string `env:"ENV"`
			EnableDebug bool   `env:"DEBUG"`
		}

		myConfig, err := LoadConfigFromEnv[MyConfig]()
		AssertNil(t, err)
		AssertEqual(t, myConfig.Env, "production")
		AssertEqual(t, myConfig.EnableDebug, true)
	})

	t.Run("parses_durations", func(t *testing.T) {
		cleanEnv()
		os.Setenv("TIMEOUT_SECS", "23s")
		os.Setenv("TIMEOUT_MINS", "45m")
		type MyConfig struct {
			TimeoutS time.Duration `env:"TIMEOUT_SECS"`
			TimeoutM time.Duration `env:"TIMEOUT_MINS"`
		}

		myConfig, err := LoadConfigFromEnv[MyConfig]()
		AssertNil(t, err)
		AssertEqual(t, myConfig.TimeoutS, 23*time.Second)
		AssertEqual(t, myConfig.TimeoutM, 45*time.Minute)
	})

	t.Run("env_overrides_default", func(t *testing.T) {
		cleanEnv()
		os.Setenv("ENV", "production")
		type MyConfig struct {
			Env string `env:"ENV" default:"dev"`
		}

		myConfig, err := LoadConfigFromEnv[MyConfig]()
		AssertNil(t, err)
		AssertEqual(t, myConfig.Env, "production")
	})

	t.Run("errors_on_missing_env", func(t *testing.T) {
		cleanEnv()
		type MyConfig struct {
			Env  string `env:"ENV" default:"dev"`
			Port int    `env:"PORT"`
		}

		_, err := LoadConfigFromEnv[MyConfig]()
		AssertNotNil(t, err)
	})

	t.Run("errors_on_wrong_type_int", func(t *testing.T) {
		cleanEnv()
		os.Setenv("PORT", "hello")
		type MyConfig struct {
			Port int `env:"PORT"`
		}

		_, err := LoadConfigFromEnv[MyConfig]()
		AssertNotNil(t, err)
	})

	t.Run("errors_on_wrong_type_bool", func(t *testing.T) {
		cleanEnv()
		os.Setenv("DEBUG", "hello")
		type MyConfig struct {
			EnableDebug bool `env:"DEBUG"`
		}

		_, err := LoadConfigFromEnv[MyConfig]()
		AssertNotNil(t, err)
	})

	t.Run("errors_on_wrong_type_duration", func(t *testing.T) {
		cleanEnv()
		os.Setenv("TIMEOUT", "hello")
		type MyConfig struct {
			Timeout time.Duration `env:"TIMEOUT" default:"5s"`
		}

		_, err := LoadConfigFromEnv[MyConfig]()
		AssertNotNil(t, err)
	})
}

// cleanEnv removes all env vars used for testing.
func cleanEnv() {
	os.Unsetenv("FOO")
	os.Unsetenv("ENV")
	os.Unsetenv("PORT")
	os.Unsetenv("TIMEOUT")
}
