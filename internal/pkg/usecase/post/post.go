package post

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
)

type PostUsecase struct {
}

func (p *PostUsecase) PublishPost(ctx context.Context, post *entity.Post) error {
	return nil
}

func (p *PostUsecase) ProhibitCommenting(ctx context.Context, author string, postId int) error {
	return nil
}

func (p *PostUsecase) AllowCommenting(ctx context.Context, author string, postId int) error {
	return nil
}

func (p *PostUsecase) GetFeedPosts(ctx context.Context, limit, cursor int) (*entity.FeedPost, error) {
	return nil, nil
}

func (p *PostUsecase) GetPost(ctx context.Context, postId int) (*entity.Post, error) {
	return nil, nil
}

func (p *PostUsecase) SubscribeOnPost(ctx context.Context, postId int) (<-chan entity.NotifyComment, error) {
	return nil, nil
}
