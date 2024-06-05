package config

import (
	"flag"
	"fmt"
)

type repository string

const DefaultPathConfig = "configs/config.yaml"

var cfgFlag = flag.String("config", DefaultPathConfig, "the path to the config")

const (
	Memory   repository = "memory"
	Postgres repository = "postgres"
)

type Config struct {
	Host       string
	Port       int
	Repository repository
}

func (cfg Config) String() string {
	return fmt.Sprintf("{host=%s, port=%d, repository=%s}", cfg.Host, cfg.Port, cfg.Repository)
}
