package graphql

import (
	"context"
	"fmt"

	"github.com/gvidow/go-post-service/internal/entity"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) Posts(ctx context.Context, limit int, cursor int) (*entity.FeedPost, error) {
	return &entity.FeedPost{Posts: []*entity.Post{{ID: 4}, {ID: 78}}, Cursor: 67}, nil
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
