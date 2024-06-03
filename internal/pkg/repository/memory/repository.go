package memory

import (
	"context"
	"log"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/usecase"
)

var _ usecase.Repository = (*memoryRepo)(nil)

type memoryRepo struct{}

func NewMemoryRepo() *memoryRepo {
	return &memoryRepo{}
}

func (m *memoryRepo) GetComments(ctx context.Context, postIds []int, cfg entity.QueryConfig) (
	entity.BatchComments,
	error,
) {
	log.Println("YES", postIds, cfg)
	return map[int][]entity.Comment{
		78: {{Content: "this 78"}},
		4:  {{Content: "this 4"}},
	}, nil
}

func (m *memoryRepo) GetReplies(ctx context.Context, commentId int, cfg entity.QueryConfig) (
	*entity.FeedComment,
	error,
) {
	return nil, nil
}

func (m *memoryRepo) GetPostById(ctx context.Context, id int) (*entity.Post, error) {
	return nil, nil
}

func (m *memoryRepo) GetPostByComment(ctx context.Context, id int) (*entity.Post, error) {
	return nil, nil
}

func (m *memoryRepo) AddComment(context.Context, *entity.Comment) error {
	return nil
}

func (m *memoryRepo) AddReply(context.Context, *entity.Comment) error {
	return nil
}

func (m *memoryRepo) AddPost(context.Context, *entity.Post) error {
	return nil
}

func (m *memoryRepo) GetFeedPosts(ctx context.Context, limit, cursor int) (*entity.FeedPost, error) {
	return nil, nil
}

func (m *memoryRepo) SetPermAddComments(ctx context.Context, postId int, allow bool) error {
	return nil
}
