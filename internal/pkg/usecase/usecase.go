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
	postUsecase := post.NewPostUsecase(repo)
	return &usecase{
		PostUsecase: postUsecase,
		CommentUsecase: comment.NewCommentUsecase(
			repo,
			func(p post.RequestPermission) (bool, error) {
				return postUsecase.IsAllowCommenting(p)
			},
		),
	}
}
