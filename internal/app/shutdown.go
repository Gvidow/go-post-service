package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var GracefulShutdownSignals = []os.Signal{
	syscall.SIGTERM,
	syscall.SIGQUIT,
	syscall.SIGINT,
}

func WithGracefulShutdown(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, GracefulShutdownSignals...)
}
