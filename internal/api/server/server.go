package server

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/gvidow/go-post-service/internal/api/graph"
	"github.com/gvidow/go-post-service/internal/app/config"
)

func NewServer(resolver graph.ResolverRoot, cfg *config.Config) *http.Server {
	mux := http.NewServeMux()

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AddTransport(transport.Websocket{})
	srv.AddTransport(transport.POST{})

	srv.Use(extension.Introspection{})

	mux.Handle("/query", srv)

	return &http.Server{Handler: mux, Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}
}
