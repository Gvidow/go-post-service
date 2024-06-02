package main

import (
	"context"
	stdlog "log"

	"github.com/gvidow/go-post-service/internal/app"
	"github.com/gvidow/go-post-service/pkg/logger"
)

func main() {
	log, err := logger.New(logger.WithRFC3339TimeEncoder())
	if err != nil {
		stdlog.Fatal(err)
	}
	defer log.Sync()

	if err = app.Main(context.Background(), log); err != nil {
		log.Error(err.Error())
	}
}
