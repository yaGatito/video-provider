package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Api Api
	Db  Db
}

type Api struct {
	Name string
	Host string
	Port string
}

type Db struct {
	Name     string
	Host     string
	Port     string
	MaxConns string
	Vendor   string
	Version  string
	Url      string
}

const (
	congfigFolder    = "./config"
	secretPgUser     = "POSTGRES_USER"
	secretPgPassword = "POSTGRES_PWD"
)

// GetURL must be called only after setup OS env variables.
func (db Db) GetURL() string {
	// vendor user pwd port db_name db_max_conn
	url := strings.Replace(db.Url, "$1", db.Vendor, 1)
	url = strings.Replace(url, "$2", os.Getenv(secretPgUser), 1)
	url = strings.Replace(url, "$3", os.Getenv(secretPgPassword), 1)
	url = strings.Replace(url, "$4", db.Host, 1)
	url = strings.Replace(url, "$5", db.Port, 1)
	url = strings.Replace(url, "$6", db.Name, 1)
	url = strings.Replace(url, "$7", db.MaxConns, 1)
	return url
}

func ParseFromFS(name string) (Config, error) {
	file, err := os.ReadFile(fmt.Sprintf("%s/%s.yml", congfigFolder, name))
	if err != nil {
		return Config{}, fmt.Errorf("failed to load file %w", err)
	}
	return ParseConfig(file)
}

func ParseConfig(cfg []byte) (Config, error) {
	var appConfig Config
	err := yaml.Unmarshal(cfg, &appConfig)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshall app.yml config: %w", err)
	}

	return appConfig, nil
}
