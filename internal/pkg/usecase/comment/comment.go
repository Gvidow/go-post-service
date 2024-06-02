package comment

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
)

type CommentUsecase struct {
}

func (c *CommentUsecase) WriteReply(ctx context.Context, comment *entity.Comment) error {
	return nil
}

func (c *CommentUsecase) WriteComment(ctx context.Context, comment *entity.Comment) error {
	return nil
}

func (c *CommentUsecase) GetReplies(ctx context.Context, limit, cursor, depth int) (*entity.FeedComment, error) {
	return nil, nil
}

func (c *CommentUsecase) GetComments(ctx context.Context, postIds []int, limit, cursor, depth int) (entity.BatchComments, error) {
	return nil, nil
}
