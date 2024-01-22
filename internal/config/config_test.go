package config

import (
	"os"
	"testing"

	"github.com/quadev-ltd/qd-common/pkg/config"
	"github.com/stretchr/testify/assert"
)

const (
	MockConfigPath = "./"
)

func TestLoad(t *testing.T) {
	t.Run("Load_Should_Show_File_Values_If_No_Env_Vars", func(t *testing.T) {
		// Setup
		cfg := &Config{}
		os.Setenv(config.AppEnvironmentKey, "test")
		os.Setenv(config.VerboseKey, "false")

		defer os.Unsetenv(config.AppEnvironmentKey)
		defer os.Unsetenv(config.VerboseKey)

		err := cfg.Load(MockConfigPath)
		assert.NoError(t, err, "expected no error from Load")

		// Assertions
		assert.Equal(t, "QuaDevEmailTest", cfg.App)
		assert.Equal(t, "localhost", cfg.SMTP.Host)
		assert.Equal(t, "1111", cfg.SMTP.Port)
		assert.Equal(t, "test@test.com", cfg.SMTP.Username)
		assert.Equal(t, "test_password", cfg.SMTP.Password)
		assert.Equal(t, "localhost", cfg.GRPC.Host)
		assert.Equal(t, "3333", cfg.GRPC.Port)
		assert.Equal(t, true, cfg.TLSEnabled)

		assert.False(t, cfg.Verbose)
		assert.Equal(t, "test", cfg.Environment)
	})

	t.Run("Load_Should_Show_Env_Vars_Values", func(t *testing.T) {
		// Setup
		cfg := &Config{}
		os.Setenv(config.AppEnvironmentKey, "test")
		os.Setenv(config.VerboseKey, "false")
		os.Setenv("TEST_ENV_APP", "QuaDevEmailTest_env")
		os.Setenv("TEST_ENV_SMTP_HOST", "localhost_env")
		os.Setenv("TEST_ENV_SMTP_PORT", "1111_env")
		os.Setenv("TEST_ENV_SMTP_USERNAME", "test_env@test.com")
		os.Setenv("TEST_ENV_SMTP_PASSWORD", "test_password_env")
		os.Setenv("TEST_ENV_GRPC_HOST", "localhost_env")
		os.Setenv("TEST_ENV_GRPC_PORT", "3333_env")

		defer os.Unsetenv(config.AppEnvironmentKey)
		defer os.Unsetenv(config.VerboseKey)
		defer os.Unsetenv("TEST_ENV_SMTP_HOST")
		defer os.Unsetenv("TEST_ENV_SMTP_PORT")
		defer os.Unsetenv("TEST_ENV_SMTP_USERNAME")
		defer os.Unsetenv("TEST_ENV_SMTP_PASSWORD")
		defer os.Unsetenv("TEST_ENV_GRPC_HOST")
		defer os.Unsetenv("TEST_ENV_GRPC_PORT")

		err := cfg.Load(MockConfigPath)
		assert.NoError(t, err, "expected no error from Load")

		// Assertions
		assert.Equal(t, "QuaDevEmailTest_env", cfg.App)
		assert.Equal(t, "localhost_env", cfg.SMTP.Host)
		assert.Equal(t, "1111_env", cfg.SMTP.Port)
		assert.Equal(t, "test_env@test.com", cfg.SMTP.Username)
		assert.Equal(t, "test_password_env", cfg.SMTP.Password)
		assert.Equal(t, "localhost_env", cfg.GRPC.Host)
		assert.Equal(t, "3333_env", cfg.GRPC.Port)

		assert.False(t, cfg.Verbose)
		assert.Equal(t, "test", cfg.Environment)
	})

}
