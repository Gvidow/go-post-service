package graphql

import (
	"context"
	"time"

	"github.com/gvidow/go-post-service/internal/api/graph"
	"github.com/gvidow/go-post-service/internal/pkg/delivery"
	"github.com/gvidow/go-post-service/internal/pkg/middleware"
	"github.com/gvidow/go-post-service/pkg/logger"
)

var DefaultTimeServeSubscribers = time.Hour

var _ graph.ResolverRoot = (*Resolver)(nil)

type Resolver struct {
	log     *logger.Logger
	usecase delivery.Usecase

	TimeServeSubscribers time.Duration
}

func NewResolver(log *logger.Logger, u delivery.Usecase) *Resolver {
	return &Resolver{
		log:     log,
		usecase: u,

		TimeServeSubscribers: DefaultTimeServeSubscribers,
	}
}

func (r *Resolver) Mutation() graph.MutationResolver         { return &mutationResolver{r} }
func (r *Resolver) Post() graph.PostResolver                 { return &postResolver{r} }
func (r *Resolver) Query() graph.QueryResolver               { return &queryResolver{r} }
func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }

func (r *Resolver) makeResponseErrorAndLog(ctx context.Context, err error) error {
	if err != nil {
		res := MakeResponseError(err)
		r.getLog(ctx).Error(err.Error(), logger.String("response", res.Message))
		return res
	}
	return nil
}

func (r *Resolver) getLog(ctx context.Context) *logger.Logger {
	if log, ok := ctx.Value(middleware.Logger).(*logger.Logger); ok {
		return log
	}
	return r.log
}
