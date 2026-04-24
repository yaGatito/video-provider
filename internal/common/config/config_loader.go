package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"video-provider/common/shared"

	"github.com/kelseyhightower/envconfig"
)

const jwtSec string = "JWT_SECRET"

type EnvSecureConfig struct {
	ApiPort string `required:"true" split_words:"true"`
	DbName  string `required:"true" split_words:"true"`
	DbPort  string `required:"true" split_words:"true"`
	DbHost  string `required:"true" split_words:"true"`
	DbPass  string `required:"true" split_words:"true"`
	DbUser  string `required:"true" split_words:"true"`
}

type JsonServiceConfig struct {
	PoolCons        int    `json:"pool_max_conns"`
	PoolConLifetime string `json:"pool_max_conn_lifetime"`
	SSLMode         string `json:"ssl_mode"`

	TokenExpTime time.Duration
	// The only usage of it should be time.ParseDuration func.
	TokenExpTimeRaw string `json:"token_exp,omitempty"`
}

type Config struct {
	JsonConf  JsonServiceConfig
	EnvConf   EnvSecureConfig
	JwtSecret []byte
}

func LoadConfig(svcPrefix string, jsonConfig []byte) (Config, error) {
	cm := Config{}

	var err error
	cm.EnvConf, err = loadEnvConfig(svcPrefix)
	if err != nil {
		return Config{}, err
	}

	cm.JsonConf, err = loadJsonConfig(jsonConfig)
	if err != nil {
		return Config{}, err
	}

	sec, ok := os.LookupEnv(jwtSec)
	if !ok || len(sec) == 0 {
		return Config{}, fmt.Errorf("failed to load env")
	}
	cm.JwtSecret = []byte(sec)

	return cm, nil
}

var nilJsonConfig = JsonServiceConfig{}

func loadJsonConfig(config []byte) (JsonServiceConfig, error) {
	var jsonConf JsonServiceConfig
	err := json.NewDecoder(bytes.NewReader(config)).Decode(&jsonConf)
	if err != nil {
		return nilJsonConfig, err
	}
	if jsonConf.TokenExpTimeRaw != "" {
		jsonConf.TokenExpTime, err = time.ParseDuration(jsonConf.TokenExpTimeRaw)
		if err != nil {
			return nilJsonConfig, err
		}
	}
	if jsonConf == nilJsonConfig {
		return nilJsonConfig, fmt.Errorf("empty env config")
	}

	return jsonConf, nil
}

var nilEnvConfig = EnvSecureConfig{}

func loadEnvConfig(svcPrefix string) (EnvSecureConfig, error) {
	var envConf EnvSecureConfig
	err := envconfig.Process(svcPrefix, &envConf)
	switch err := err.(type) {
	case nil:
		if envConf == nilEnvConfig {
			return EnvSecureConfig{}, fmt.Errorf("empty env config")
		}

		return envConf, nil

	case *envconfig.ParseError:
		return EnvSecureConfig{}, shared.NewError(shared.ErrInternal, fmt.Sprintf("failed to extract from env %s for field %s", err.KeyName, err.FieldName), err)

	default:
		return EnvSecureConfig{}, shared.NewError(shared.ErrInternal, "failed to parse config", err)
	}
}
