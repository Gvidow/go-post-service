package graphql

import (
	"context"
	"fmt"

	"github.com/gvidow/go-post-service/internal/entity"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) Posts(ctx context.Context, limit int, cursor int) (*entity.FeedPost, error) {
	panic(fmt.Errorf("not implemented: Posts - posts"))
}

func (r *queryResolver) GetPost(ctx context.Context, postID int) (*entity.Post, error) {
	panic(fmt.Errorf("not implemented: GetPost - getPost"))
}

func (r *queryResolver) Replies(
	ctx context.Context,
	commentID int,
	limit int,
	cursor int,
	depth int,
) (*entity.FeedComment, error) {

	panic(fmt.Errorf("not implemented: Replies - replies"))
}
