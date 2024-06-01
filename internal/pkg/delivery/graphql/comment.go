package graphql

import (
	"context"
	"fmt"

	"github.com/gvidow/go-post-service/internal/entity"
)

type commentResolver struct{ *Resolver }

func (r *commentResolver) Replies(ctx context.Context, obj *entity.Comment, limit int, after int) (*entity.FeedComment, error) {
	panic(fmt.Errorf("not implemented: Replies - replies"))
}
