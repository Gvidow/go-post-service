package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_timeoutConnect = time.Minute
	_maxConnDB      = 500
	_schemaDB       = "post_service"
)

func NewPoolConnectPG(ctx context.Context) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, _timeoutConnect)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dsn())
	if err != nil {
		return nil, errors.WrapFail(err, "parse config")
	}

	cfg.MaxConns = _maxConnDB
	cfg.ConnConfig.RuntimeParams["search_path"] = _schemaDB

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "init pool connect by config")
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, errors.Wrap(err, "ping db")
	}

	return pool, nil
}

func dsn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
}
