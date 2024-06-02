package usecase

import (
	"github.com/gvidow/go-post-service/internal/pkg/usecase/comment"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/post"
)

type Repository interface {
	comment.Repository
	post.Repository
}
