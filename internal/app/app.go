package app

import (
	"context"

	"github.com/gvidow/go-post-service/internal/api/server"
	"github.com/gvidow/go-post-service/internal/pkg/delivery/graphql"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
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

	resolver := graphql.Resolver{}

	server := server.NewServer(&resolver)
	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()
	server.ListenAndServe()

	return nil
}
