package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/imdario/mergo"
	"github.com/nzqpeace/rbac/cache"
	"github.com/nzqpeace/rbac/db"
)

type Logger struct {
	Output string `json:"output"`
	Level  string `json:"level"`
}

type HttpServerConfig struct {
	Address string `json:"address"`
}

type Config struct {
	Log   *Logger            `json:"log"`
	Redis *cache.RedisConfig `json:"redis"`
	Mongo *db.MgoConf        `json:"mongo"`
	Http  *HttpServerConfig  `json:"http_server"`
}

func DefaultConfig() *Config {
	return &Config{
		Log: &Logger{
			Output: "./logs/cowshed.log",
			Level:  "info",
		},
		Redis: cache.DefaultConfig(),
		Mongo: db.DefaultConf(),
		Http: &HttpServerConfig{
			Address: ":60001",
		},
	}
}

func loadConfig(filename string) (c *Config) {
	byts, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(byts, &c)
	}

	if err != nil {
		fmt.Printf("load config file '%s' fail, use default config instead\n", filename)
	}

	if c == nil {
		c = DefaultConfig()
	} else {
		mergo.Merge(c, DefaultConfig())
	}
	return
}
