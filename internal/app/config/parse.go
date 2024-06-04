package config

import (
	"flag"
	"fmt"
	"slices"

	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"go.uber.org/config"
)

var (
	ErrBadListRepositories = errors.New("bad list repositories")
	ErrInvalidRepository   = errors.New("invalid selected repository")
)

func Parse() (*Config, error) {
	flag.Parse()

	yml, err := config.NewYAML(config.File(*cfgFlag))
	if err != nil {
		return nil, errors.WrapFail(err, "parse yaml config")
	}

	cfg := &Config{}

	serverCfg := yml.Get("app.server")

	err = errors.Join(
		serverCfg.Get("host").Populate(&cfg.Host),
		serverCfg.Get("port").Populate(&cfg.Port),
		checkListRepositories(yml),
		yml.Get("app.repository").Populate(&cfg.Repository),
	)

	if err != nil {
		return nil, errors.WrapFail(err, "parse config")
	} else if cfg.Repository != Memory && cfg.Repository != Postgres {
		return nil, ErrInvalidRepository
	}

	return cfg, nil
}

func checkListRepositories(yml *config.YAML) error {
	var repositories []repository
	err := yml.Get("repositories").Populate(&repositories)
	if err != nil {
		return errors.WrapFail(err, "parse repositories")
	}

	fmt.Println(repositories, len(repositories))

	if len(repositories) != 2 ||
		!slices.Contains(repositories, Memory) ||
		!slices.Contains(repositories, Postgres) {

		return ErrBadListRepositories
	}

	return nil
}
