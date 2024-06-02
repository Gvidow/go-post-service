package post

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
)

type Repository interface {
	GetPostById(ctx context.Context, id int) (*entity.Post, error)
	GetPostByComment(ctx context.Context, id int) (*entity.Post, error)
}

type object interface {
	Id() int
	entity()
}
