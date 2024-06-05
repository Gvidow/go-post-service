package post

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
)

type Repository interface {
	AddPost(context.Context, *entity.Post) error
	GetPostById(ctx context.Context, id int) (*entity.Post, error)
	GetPostByComment(ctx context.Context, id int) (*entity.Post, error)
	GetFeedPosts(ctx context.Context, limit, cursor int) (*entity.FeedPost, error)
	SetPermAddComments(ctx context.Context, postId int, allow bool) error
}

type object interface {
	Id() int
	entity()
}
