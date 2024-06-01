package graphql

import (
	"context"
	"fmt"

	"github.com/gvidow/go-post-service/internal/entity"
)

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) SubscribeOnPost(ctx context.Context, postID int) (<-chan *entity.Comment, error) {
	panic(fmt.Errorf("not implemented: SubscribeOnPost - subscribeOnPost"))
}
