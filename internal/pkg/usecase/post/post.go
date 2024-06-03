package post

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
)

var errNotAllowed = errors.New("action not allowed")

type PostUsecase struct {
	repo Repository
}

func NewPostUsecase(repo Repository) *PostUsecase {
	return &PostUsecase{repo}
}

func (p *PostUsecase) PublishPost(ctx context.Context, post *entity.Post) error {
	return errors.WrapFail(p.repo.AddPost(ctx, post), "add post to repository")
}

func (p *PostUsecase) ProhibitCommenting(ctx context.Context, author string, postId int) error {
	if err := p.checkIsAuthor(ctx, author, postId); err != nil {
		return errors.Wrap(err, "failed check permission prohibit")
	}

	return errors.WrapFail(p.repo.SetPermAddComments(ctx, postId, false), "prohibit commenting")
}

func (p *PostUsecase) AllowCommenting(ctx context.Context, author string, postId int) error {
	if err := p.checkIsAuthor(ctx, author, postId); err != nil {
		return errors.Wrap(err, "failed check permission allow")
	}

	return errors.WrapFail(p.repo.SetPermAddComments(ctx, postId, true), "allow commenting")
}

func (p *PostUsecase) GetFeedPosts(ctx context.Context, limit, cursor int) (*entity.FeedPost, error) {
	return p.repo.GetFeedPosts(ctx, limit, cursor)
}

func (p *PostUsecase) GetPost(ctx context.Context, postId int) (*entity.Post, error) {
	return p.repo.GetPostById(ctx, postId)
}

func (p *PostUsecase) SubscribeOnPost(ctx context.Context, postId int) (<-chan entity.NotifyComment, error) {
	return nil, nil
}

func (p *PostUsecase) checkIsAuthor(ctx context.Context, author string, postId int) error {
	switch post, err := p.repo.GetPostById(ctx, postId); {
	case err != nil:
		return errors.WrapFail(err, "get post author")
	case post.Author != author:
		return errors.WithType(errNotAllowed, errors.NotPermission)
	}
	return nil
}
