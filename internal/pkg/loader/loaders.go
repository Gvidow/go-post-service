package loader

import (
	"context"
	"time"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"github.com/gvidow/go-post-service/internal/pkg/usecase"
	"github.com/vikstrous/dataloadgen"
)

const _timeWait = 10 * time.Millisecond

type Loaders struct {
	CommentsLoader *dataloadgen.Loader[int, []entity.Comment]
}

func NewLoaders(repo usecase.Repository) *Loaders {
	ur := &commentsGetter{repo}
	return &Loaders{
		CommentsLoader: dataloadgen.NewLoader(ur.getBatchComments, dataloadgen.WithWait(_timeWait)),
	}
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(LoadersKey).(*Loaders)
}

func GetComments(ctx context.Context, postId, limit, cursor, depth int) (
	comments []entity.Comment,
	newCursor int,
	err error,
) {
	loaders := For(ctx)
	queryCfg := entity.QueryConfig{
		Limit:  limit,
		Cursor: cursor,
		Depth:  depth,
	}

	comments, err = loaders.CommentsLoader.Load(context.WithValue(ctx, _queryConfigKey, queryCfg), postId)
	if err != nil {
		return nil, 0, errors.WrapFail(err, "load from context loader")
	}

	return comments, cursor + len(comments), nil
}
