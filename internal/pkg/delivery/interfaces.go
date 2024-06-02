package delivery

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
)

type Usecase interface {
	PublishPost(context.Context, *entity.Post) error
	WriteComment(context.Context, int, *entity.Comment) error
	WriteReply(context.Context, int, *entity.Comment) error
}
