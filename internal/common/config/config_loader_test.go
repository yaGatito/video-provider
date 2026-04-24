package config_test

import (
	"os"
	"testing"
	"time"
	"video-provider/common/config"

	"github.com/stretchr/testify/assert"
)

const (
	apiPort string = "9090"
	dbPort  string = "9070"
	dbHost  string = "localghost"
	dbName  string = "test"
	dbUser  string = "testUser"
	dbPass  string = "testPass"
)

func TestLoadConfig_Success(t *testing.T) {
	os.Clearenv()

	// Set up environment variables
	os.Setenv("JWT_SECRET", "testjwtsecret")

	os.Setenv("TEST_API_PORT", apiPort)
	os.Setenv("TEST_DB_PORT", dbPort)
	os.Setenv("TEST_DB_HOST", dbHost)
	os.Setenv("TEST_DB_NAME", dbName)
	os.Setenv("TEST_DB_USER", dbUser)
	os.Setenv("TEST_DB_PASS", dbPass)

	// Create a valid JSON configuration file
	jsonConfig := `{
		"pool_max_conns": 10,
		"pool_max_conn_lifetime": "30m",
		"ssl_mode": "disable",
		"token_exp": "1h30m"
	}`

	// Call LoadConfig
	config, err := config.LoadConfig("test", []byte(jsonConfig))
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Verify the configuration
	assert.Equal(t, []byte("testjwtsecret"), config.JwtSecret)

	assert.Equal(t, 10, config.JsonConf.PoolCons)
	assert.Equal(t, "30m", config.JsonConf.PoolConLifetime)
	assert.Equal(t, 90*time.Minute, config.JsonConf.TokenExpTime)
	assert.Equal(t, "disable", config.JsonConf.SSLMode)

	assert.Equal(t, apiPort, config.EnvConf.ApiPort)
	assert.Equal(t, dbPort, config.EnvConf.DbPort)
	assert.Equal(t, dbHost, config.EnvConf.DbHost)
	assert.Equal(t, dbName, config.EnvConf.DbName)
	assert.Equal(t, dbUser, config.EnvConf.DbUser)
	assert.Equal(t, dbPass, config.EnvConf.DbPass)
}

func TestLoadConfig_FailedJwtSecret(t *testing.T) {
	os.Clearenv()

	// Set up environment variables
	os.Unsetenv("JWT_SECRET")

	os.Setenv("TEST_API_PORT", apiPort)
	os.Setenv("TEST_DB_PORT", dbPort)
	os.Setenv("TEST_DB_HOST", dbHost)
	os.Setenv("TEST_DB_NAME", dbName)
	os.Setenv("TEST_DB_USER", dbUser)
	os.Setenv("TEST_DB_PASS", dbPass)

	// Create a valid JSON configuration file
	jsonConfig := `{
		"pool_max_conns": 10,
		"pool_max_conn_lifetime": "30m",
		"ssl_mode": "disable",
		"token_exp": "1h"
	}`

	// Call LoadConfig
	config, err := config.LoadConfig("test", []byte(jsonConfig))
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to load env")
	assert.Empty(t, config)
}

func TestLoadConfig_EmptyJwtSecret(t *testing.T) {
	os.Clearenv()

	// Set up environment variables
	os.Setenv("JWT_SECRET", "")

	os.Setenv("TEST_API_PORT", apiPort)
	os.Setenv("TEST_DB_PORT", dbPort)
	os.Setenv("TEST_DB_HOST", dbHost)
	os.Setenv("TEST_DB_NAME", dbName)
	os.Setenv("TEST_DB_USER", dbUser)
	os.Setenv("TEST_DB_PASS", dbPass)

	// Create a valid JSON configuration file
	jsonConfig := `{
		"pool_max_conns": 10,
		"pool_max_conn_lifetime": "30m",
		"ssl_mode": "disable",
		"token_exp": "1h"
	}`

	// Call LoadConfig
	config, err := config.LoadConfig("test", []byte(jsonConfig))
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to load env")
	assert.Empty(t, config)
}

func TestLoadConfig_FailedEnvConfig(t *testing.T) {
	os.Clearenv()

	// Set up environment variables
	os.Setenv("JWT_SECRET", "testjwtsecret")

	os.Setenv("TEST_1API_PORT", apiPort)

	// Create a valid JSON configuration file
	jsonConfig := `{
		"pool_max_conns": 10,
		"pool_max_conn_lifetime": "30m",
		"ssl_mode": "disable",
		"token_exp": "1h"
	}`

	// Call LoadConfig
	config, err := config.LoadConfig("test", []byte(jsonConfig))
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to parse config: required key TEST_API_PORT missing value")
	assert.Empty(t, config)
}

func TestLoadConfig_EmptyEnvConfig(t *testing.T) {
	os.Clearenv()

	// Set up environment variables
	os.Setenv("JWT_SECRET", "testjwtsecret")

	os.Setenv("TEST_API_PORT", "")
	os.Setenv("TEST_DB_PORT", "")
	os.Setenv("TEST_DB_HOST", "")
	os.Setenv("TEST_DB_NAME", "")
	os.Setenv("TEST_DB_USER", "")
	os.Setenv("TEST_DB_PASS", "")

	// Create a valid JSON configuration file
	jsonConfig := `{
		"pool_max_conns": 10,
		"pool_max_conn_lifetime": "30m",
		"ssl_mode": "disable",
		"token_exp": "1h"
	}`

	// Call LoadConfig
	config, err := config.LoadConfig("test", []byte(jsonConfig))
	assert.Error(t, err)
	assert.EqualError(t, err, "empty env config")
	assert.Empty(t, config)
}

func TestLoadConfig_FailedJsonConfig(t *testing.T) {
	os.Clearenv()

	// Set up environment variables
	os.Setenv("JWT_SECRET", "testjwtsecret")

	os.Setenv("TEST_API_PORT", apiPort)
	os.Setenv("TEST_DB_PORT", dbPort)
	os.Setenv("TEST_DB_HOST", dbHost)
	os.Setenv("TEST_DB_NAME", dbName)
	os.Setenv("TEST_DB_USER", dbUser)
	os.Setenv("TEST_DB_PASS", dbPass)

	// Create a valid JSON configuration file
	jsonConfig := `{
		"pool_max_conns": 10,
		"pool_max_conn_lifetime": "30m",
		"ssl_mode": "disable",
		"token_exp": "1h44"
	}`

	// Call LoadConfig
	config, err := config.LoadConfig("test", []byte(jsonConfig))
	assert.Error(t, err)
	assert.Empty(t, config)
}

func TestLoadConfig_EmptyRequiredAttrJsonConfig(t *testing.T) {
	os.Clearenv()

	// Set up environment variables
	os.Setenv("JWT_SECRET", "testjwtsecret")

	os.Setenv("TEST_API_PORT", apiPort)
	os.Setenv("TEST_DB_PORT", dbPort)
	os.Setenv("TEST_DB_HOST", dbHost)
	os.Setenv("TEST_DB_NAME", dbName)
	os.Setenv("TEST_DB_USER", dbUser)
	os.Setenv("TEST_DB_PASS", dbPass)

	// Create a valid JSON configuration file
	jsonConfig := `{
		"pool_max_conns": 0,
		"pool_max_conn_lifetime": "",
		"ssl_mode": "",
		"token_exp": ""
	}`

	// Call LoadConfig
	config, err := config.LoadConfig("test", []byte(jsonConfig))
	assert.Error(t, err)
	assert.Empty(t, config)
}

func TestLoadConfig_EmptyTokenAttrJsonConfig(t *testing.T) {
	os.Clearenv()

	// Set up environment variables
	os.Setenv("JWT_SECRET", "testjwtsecret")

	os.Setenv("TEST_API_PORT", apiPort)
	os.Setenv("TEST_DB_PORT", dbPort)
	os.Setenv("TEST_DB_HOST", dbHost)
	os.Setenv("TEST_DB_NAME", dbName)
	os.Setenv("TEST_DB_USER", dbUser)
	os.Setenv("TEST_DB_PASS", dbPass)

	// Create a valid JSON configuration file
	jsonConfig := `{
		"pool_max_conns": 11,
		"pool_max_conn_lifetime": "32m",
		"ssl_mode": "enabled"
	}`

	// Call LoadConfig
	config, err := config.LoadConfig("test", []byte(jsonConfig))
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Verify the configuration
	assert.Equal(t, []byte("testjwtsecret"), config.JwtSecret)

	assert.Equal(t, 11, config.JsonConf.PoolCons)
	assert.Equal(t, "32m", config.JsonConf.PoolConLifetime)
	assert.Equal(t, "enabled", config.JsonConf.SSLMode)
	assert.Equal(t, time.Duration(0), config.JsonConf.TokenExpTime)

	assert.Equal(t, apiPort, config.EnvConf.ApiPort)
	assert.Equal(t, dbPort, config.EnvConf.DbPort)
	assert.Equal(t, dbHost, config.EnvConf.DbHost)
	assert.Equal(t, dbName, config.EnvConf.DbName)
	assert.Equal(t, dbUser, config.EnvConf.DbUser)
	assert.Equal(t, dbPass, config.EnvConf.DbPass)
}
