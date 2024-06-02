package delivery

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
)

type Usecase interface {
	postUsecase
	commentUsecase
}
type postUsecase interface {
	PublishPost(context.Context, *entity.Post) error
	ProhibitCommenting(ctx context.Context, author string, postId int) error
	AllowCommenting(ctx context.Context, author string, postId int) error
	GetFeedPosts(ctx context.Context, limit, cursor int) (*entity.FeedComment, error)
	GetPost(ctx context.Context, postId int) (*entity.Post, error)
	SubscribeOnPost(ctx context.Context, postId int) (<-chan entity.NotifyComment, error)
}

type commentUsecase interface {
	WriteReply(context.Context, *entity.Comment) error
	WriteComment(context.Context, *entity.Comment) error
	GetReplies(ctx context.Context, limit, cursor, depth int) (*entity.FeedComment, error)
}
