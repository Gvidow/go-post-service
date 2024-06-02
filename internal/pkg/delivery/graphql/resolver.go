package graphql

import (
	"github.com/gvidow/go-post-service/internal/api/graph"
)

var _ graph.ResolverRoot = (*Resolver)(nil)

type Resolver struct{}

func (r *Resolver) Mutation() graph.MutationResolver         { return &mutationResolver{r} }
func (r *Resolver) Post() graph.PostResolver                 { return &postResolver{r} }
func (r *Resolver) Query() graph.QueryResolver               { return &queryResolver{r} }
func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }
