package usecase

import (
	"github.com/gvidow/go-post-service/internal/pkg/delivery"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/comment"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/notify"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/post"
	"github.com/gvidow/go-post-service/pkg/logger"
)

var _ delivery.Usecase = (*usecase)(nil)

type usecase struct {
	*post.PostUsecase
	*comment.CommentUsecase
}

func NewUsecase(log *logger.Logger, repo Repository) *usecase {
	notifier := notify.NewNotifier(log)
	postUsecase := post.NewPostUsecase(repo, notifier)

	return &usecase{
		PostUsecase: postUsecase,
		CommentUsecase: comment.NewCommentUsecase(
			repo,
			notifier,
			postUsecase,
		),
	}
}
