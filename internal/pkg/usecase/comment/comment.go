package comment

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/notify"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/post"
)

var errCommentsNotAllow = errors.New("the post is closed from comments")

type CommentUsecase struct {
	repo     Repository
	notifier notify.Notifier

	isAllow func(post.RequestPermission) (bool, error)
}

func NewCommentUsecase(
	repo Repository,
	notifier notify.Notifier,
	check func(post.RequestPermission) (bool, error),
) *CommentUsecase {
	return &CommentUsecase{
		repo:     repo,
		notifier: notifier,
		isAllow:  check,
	}
}

func (c *CommentUsecase) WriteReply(ctx context.Context, comment *entity.Comment) error {
	if err := c.checkPermission(post.RequestPermission{
		Entity: post.ReplyToComment(comment.Parent),
		Ctx:    ctx,
	}); err != nil {
		return err
	}

	if err := c.repo.AddReply(ctx, comment); err != nil {
		return errors.Wrap(err, "add reply to repository")
	}

	go c.notifier.PublishComment(ctx, comment, comment.Parent)

	return nil
}

func (c *CommentUsecase) WriteComment(ctx context.Context, comment *entity.Comment) error {
	if err := c.checkPermission(post.RequestPermission{
		Entity: post.CommentToPost(comment.Parent),
		Ctx:    ctx,
	}); err != nil {
		return err
	}

	if err := c.repo.AddComment(ctx, comment); err != nil {
		return errors.Wrap(err, "add comment to repository")
	}

	go c.notifier.PublishComment(ctx, comment, comment.Parent)

	return nil
}

func (c *CommentUsecase) GetReplies(ctx context.Context, commentId, limit, cursor, depth int) (*entity.FeedComment, error) {
	queryCfg := entity.QueryConfig{
		Limit:  limit,
		Cursor: cursor,
		Depth:  depth,
	}

	return c.repo.GetReplies(ctx, commentId, queryCfg)
}

func (c *CommentUsecase) checkPermission(r post.RequestPermission) error {
	if ok, err := c.isAllow(r); err != nil {
		return errors.WrapFail(err, "checking permission to leave comments")
	} else if !ok {
		return errors.WithType(errCommentsNotAllow, errors.CommentsAreProhibited)
	}
	return nil
}
