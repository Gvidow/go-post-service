package app

import (
	"context"

	"github.com/gvidow/go-post-service/pkg/logger"
)

func Main(ctx context.Context, log *logger.Logger) error {
	ctx, cancel := WithGracefulShutdown(ctx)
	defer cancel()

	_ = ctx

	return nil
}
