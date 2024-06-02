package usecase

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
)

type Repository interface {
	GetComments(ctx context.Context, postIds []int, cfg entity.QueryConfig) (entity.BatchComments, error)
}
