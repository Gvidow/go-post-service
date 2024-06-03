package notify

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
)

type Notifier interface {
	PublishComment(ctx context.Context, comment *entity.Comment, postId int)
	RegistryChanNotifier(ctx context.Context, ch chan<- entity.NotifyComment, postId int) error
}
