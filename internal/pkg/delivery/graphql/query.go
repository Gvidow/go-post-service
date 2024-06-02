package graphql

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) Posts(ctx context.Context, limit int, cursor int) (*entity.FeedPost, error) {
	feed, err := r.usecase.GetFeedPosts(ctx, limit, cursor)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	return feed, nil
}

func (r *queryResolver) GetPost(ctx context.Context, postID int) (*entity.Post, error) {
	post, err := r.usecase.GetPost(ctx, postID)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	return post, nil
}

func (r *queryResolver) Replies(
	ctx context.Context,
	commentID int,
	limit int,
	cursor int,
	depth int,
) (*entity.FeedComment, error) {

	feed, err := r.usecase.GetReplies(ctx, commentID, limit, cursor, depth)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	return feed, err
}
