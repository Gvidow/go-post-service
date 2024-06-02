package graphql

import (
	"github.com/gvidow/go-post-service/internal/api/graph"
	"github.com/gvidow/go-post-service/internal/pkg/delivery"
	"github.com/gvidow/go-post-service/pkg/logger"
)

var _ graph.ResolverRoot = (*Resolver)(nil)

type Resolver struct {
	log     *logger.Logger
	usecase delivery.Usecase
}

func NewResolver(log *logger.Logger, u delivery.Usecase) *Resolver {
	return &Resolver{
		log:     log,
		usecase: u,
	}
}

func (r *Resolver) Mutation() graph.MutationResolver         { return &mutationResolver{r} }
func (r *Resolver) Post() graph.PostResolver                 { return &postResolver{r} }
func (r *Resolver) Query() graph.QueryResolver               { return &queryResolver{r} }
func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }
