package app

import (
	"context"

	"github.com/gvidow/go-post-service/internal/api/server"
	"github.com/gvidow/go-post-service/internal/pkg/delivery/graphql"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"github.com/gvidow/go-post-service/internal/pkg/middleware"
	"github.com/gvidow/go-post-service/internal/pkg/repository/memory"
	"github.com/gvidow/go-post-service/internal/pkg/usecase"
	"github.com/gvidow/go-post-service/pkg/logger"
)

func Main(ctx context.Context, log *logger.Logger) error {
	ctx, cancel := WithGracefulShutdown(ctx)
	defer cancel()

	pool, err := NewPoolConnectPG(ctx)
	if err != nil {
		return errors.WrapFail(err, "open connect to db")
	}
	defer pool.Close()

	repo := memory.NewMemoryRepo()

	resolver := graphql.NewResolver(log, usecase.NewUsecase(repo))

	server := server.NewServer(resolver)
	server.Handler = middleware.WithLoaders(repo, server.Handler)
	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()
	server.ListenAndServe()

	return nil
}
