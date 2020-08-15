package config

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DB DB
}

type DB struct {
	Host     string
	Port     int
	DBName   string
	SSLMode  string
	User     string
	Password string
}

var conf *Config

func GetDSN(filePath string) (string, error) {
	c, err := load(filePath)
	if err != nil {
		return "", err
	}

	d := c.DB
	var dsn string
	dsn = fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s user=%s", d.Host, d.Port, d.DBName, d.SSLMode, d.User)

	if d.Password != "" {
		dsn += fmt.Sprintf(" password=%s", d.Password)
	}

	return dsn, nil
}

func load(filePath string) (*Config, error) {
	if conf != nil {
		return conf, nil
	}

	conf = &Config{}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return conf, errors.WithStack(err)
	}

	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		return conf, errors.WithStack(err)
	}
	return conf, nil
}
