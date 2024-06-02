package loader

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"github.com/gvidow/go-post-service/internal/pkg/usecase"
)

type commentsGetter struct {
	repo usecase.Repository
}

func (c *commentsGetter) getBatchComments(ctx context.Context, postIds []int) ([][]entity.Comment, []error) {
	cfg := ctx.Value(_queryConfigKey).(entity.QueryConfig)
	comments, err := c.repo.GetComments(ctx, postIds, cfg)

	if err != nil {
		return make([][]entity.Comment, len(postIds)), multiplyError(
			errors.Wrap(err, "get comments from repository"),
			len(postIds),
		)
	}

	res := make([][]entity.Comment, len(postIds))
	for ind, postId := range postIds {
		res[ind] = comments[postId]
	}

	return res, nil
}

func multiplyError(err error, n int) []error {
	errs := make([]error, n)
	for i := range errs {
		errs[i] = err
	}
	return errs
}
