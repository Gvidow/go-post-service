package comment

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/post"
)

type Repository interface {
	GetReplies(ctx context.Context, commentId int, cfg entity.QueryConfig) (*entity.FeedComment, error)
	GetComments(ctx context.Context, postIds []int, cfg entity.QueryConfig) (entity.BatchComments, error)
	AddComment(context.Context, *entity.Comment) error
	AddReply(context.Context, *entity.Comment) error
}

type postGetter interface {
	GetPostByEntity(r post.Request) (*entity.Post, error)
}
