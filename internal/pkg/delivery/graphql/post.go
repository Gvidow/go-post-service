package graphql

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/loader"
	"github.com/gvidow/go-post-service/pkg/slices"
)

type postResolver struct{ *Resolver }

func (r *postResolver) Comments(
	ctx context.Context,
	obj *entity.Post,
	limit int,
	cursor int,
	depth int,
) (*entity.FeedComment, error) {

	comments, newCursor, err := loader.GetComments(ctx, obj.ID, limit, cursor, depth)
	if err != nil {
		return nil, r.makeResponseErrorAndLog(err)
	}

	return &entity.FeedComment{
		Comments: slices.ToPointers(comments),
		Cursor:   newCursor,
	}, nil
}
