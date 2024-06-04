package post

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
)

type Request struct {
	Entity object
	Ctx    context.Context
}

type (
	requestComment int
	requestReply   int
)

func CommentToPost(id int) requestComment {
	return requestComment(id)
}

func ReplyToComment(id int) requestReply {
	return requestReply(id)
}

func (r requestComment) Id() int { return int(r) }
func (r requestComment) entity() {}

func (r requestReply) Id() int { return int(r) }
func (r requestReply) entity() {}

func (p *PostUsecase) GetPostByEntity(r Request) (*entity.Post, error) {
	var (
		post *entity.Post
		err  error
	)

	switch obj := r.Entity.(type) {
	case requestComment:
		post, err = p.repo.GetPostById(r.Ctx, obj.Id())
	case requestReply:
		post, err = p.repo.GetPostByComment(r.Ctx, obj.Id())
	}

	if err != nil {
		return nil, errors.WrapFail(err, "get post for check permission")
	}

	return post, nil
}
