package app

import (
	"context"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"github.com/gvidow/go-post-service/internal/api/server"
	"github.com/gvidow/go-post-service/internal/app/config"
	"github.com/gvidow/go-post-service/internal/pkg/delivery/graphql"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"github.com/gvidow/go-post-service/internal/pkg/middleware"
	"github.com/gvidow/go-post-service/internal/pkg/repository/memory"
	"github.com/gvidow/go-post-service/internal/pkg/repository/postgres"
	"github.com/gvidow/go-post-service/internal/pkg/usecase"
	"github.com/gvidow/go-post-service/pkg/logger"
)

func Main(ctx context.Context, log *logger.Logger) error {
	godotenv.Load()

	ctx, cancel := WithGracefulShutdown(ctx)
	defer cancel()

	cfg, err := config.Parse()
	if err != nil {
		return errors.WrapFail(err, "parse application config")
	}

	log.Info("parse config", logger.String("config", cfg.String()))

	var repo usecase.Repository
	switch cfg.Repository {
	case config.Postgres:
		pool, err := NewPoolConnectPG(ctx)
		if err != nil {
			return errors.WrapFail(err, "open connect to db")
		}
		defer pool.Close()

		repo = postgres.NewPostgresRepo(pool)

	case config.Memory:
		repo = memory.NewMemoryRepo()
	}

	resolver := graphql.NewResolver(log, usecase.NewUsecase(log, repo))

	server := server.NewServer(resolver, cfg)
	server.Handler = middleware.WithLoaders(repo, server.Handler)
	server.Handler = middleware.WithLogger(log, server.Handler)
	server.Handler = middleware.RequestID(server.Handler)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Error(err.Error())
		}
	}()

	log.Info("start server", logger.String("address", server.Addr))
	err = server.ListenAndServe()

	wg.Wait()
	return errors.Wrap(err, "serve server")
}
