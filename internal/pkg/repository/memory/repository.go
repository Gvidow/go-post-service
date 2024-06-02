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
