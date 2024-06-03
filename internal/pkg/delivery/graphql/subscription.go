package graphql

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
)

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) SubscribeOnPost(ctx context.Context, postID int) (<-chan *entity.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.TimeServeSubscribers)

	chanNotifier, err := r.usecase.SubscribeOnPost(ctx, postID)
	if err != nil {
		cancel()
		return nil, r.makeResponseErrorAndLog(errors.Wrap(err, "fail subscribe"))
	}

	chanComment := make(chan *entity.Comment)
	go func() {
		defer cancel()

		for msg := range chanNotifier {
			if msg.Err != nil {
				r.log.Error(msg.Err.Error())
				continue
			}

			chanComment <- msg.Comment
		}
	}()

	return chanComment, nil
}
