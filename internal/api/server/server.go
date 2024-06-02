package server

import (
	"net/http"

	h "github.com/99designs/gqlgen/handler"

	// "github.com/99designs/gqlgen/graphql/playground"
	"github.com/gvidow/go-post-service/internal/api/graph"
)

func NewServer(resolver graph.ResolverRoot) *http.Server {
	mux := http.NewServeMux()

	srv := h.GraphQL(
		graph.NewExecutableSchema(graph.Config{Resolvers: resolver}),
	)

	// handler.New(nil).AddTransport

	mux.Handle("/query", srv)
	// mux.Handle("/", h.Playground("GraphQL playground", "/query"))

	return &http.Server{Handler: mux, Addr: ":8080"}
}
