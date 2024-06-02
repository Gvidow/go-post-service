package usecase

import (
	"github.com/gvidow/go-post-service/internal/pkg/delivery"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/comment"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/post"
)

var _ delivery.Usecase = (*usecase)(nil)

type usecase struct {
	*post.PostUsecase
	*comment.CommentUsecase
}

func NewUsecase(repo Repository) *usecase {
	return &usecase{
		PostUsecase:    nil,
		CommentUsecase: nil,
	}
}
