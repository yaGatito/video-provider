package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"video-provider/common/shared"

	"github.com/kelseyhightower/envconfig"
)

const (
	dbUserEnv = "POSTGRES_USER"
	dbPassEnv = "POSTGRES_PASSWORD"
)

type Config struct {
	ApiPort         string `required:"true" split_words:"true"`
	ApiMaxDbCons    int    `ignored:"true"`
	ApiMaxDbConLife string `ignored:"true"`
	ApiSslModCon    string `ignored:"true"`

	DbName string `required:"true" split_words:"true"`
	DbPort string `required:"true" split_words:"true"`
	DbHost string `required:"true" split_words:"true"`
	PgPass string `ignored:"true"`
	PgUser string `ignored:"true"`
}

type serviceConfig struct {
	PoolCons        int    `json:"pool_max_conns"`
	PoolConLifetime string `json:"pool_max_conn_lifetime"`
	SSLMode         string `json:"ssl_mode"`
}

func LoadConfig(svcPrefix string) (Config, error) {
	var c Config

	file, err := os.ReadFile(fmt.Sprintf("./config/%s_config.json", svcPrefix))
	if err != nil {
		return Config{}, err
	}

	var sc serviceConfig
	err = json.NewDecoder(bytes.NewReader(file)).Decode(&sc)
	if err != nil {
		return Config{}, err
	}

	err = envconfig.Process(svcPrefix, &c)
	switch err := err.(type) {
	case nil:
		var ok bool
		c.PgPass, ok = os.LookupEnv(dbPassEnv)
		if !ok {
			return Config{}, shared.NewError(shared.ErrInternal, "failed to lookup env:"+dbPassEnv, err)
		}
		c.PgUser, ok = os.LookupEnv(dbUserEnv)
		if !ok {
			return Config{}, shared.NewError(shared.ErrInternal, "failed to lookup env:"+dbUserEnv, err)
		}
		c.ApiMaxDbCons = sc.PoolCons
		c.ApiMaxDbConLife = sc.PoolConLifetime
		c.ApiSslModCon = sc.SSLMode
		return c, nil

	case *envconfig.ParseError:
		return Config{}, shared.NewError(shared.ErrInternal, fmt.Sprintf("failed to extract from env %s for field %s", err.KeyName, err.FieldName), err)

	default:
		return Config{}, shared.NewError(shared.ErrInternal, "failed to parse config", err)
	}
}
