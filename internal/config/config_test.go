package config

import (
	"os"
	"testing"

	"github.com/gustavo-m-franco/qd-common/pkg/config"

	"github.com/stretchr/testify/assert"
)

const (
	MockConfigPath = "./"
)

func TestLoad(t *testing.T) {
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

	assert.False(t, cfg.Verbose)
	assert.Equal(t, "test", cfg.Environment)
}
