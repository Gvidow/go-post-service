package graphql

import (
	"context"
	"fmt"

	"github.com/gvidow/go-post-service/internal/entity"
)

type postResolver struct{ *Resolver }

func (r *postResolver) Comments(ctx context.Context, obj *entity.Post, limit int, after int) (*entity.FeedComment, error) {
	panic(fmt.Errorf("not implemented: Comments - comments"))
}
