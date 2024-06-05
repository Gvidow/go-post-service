package comment

import (
	"context"
	"unicode/utf8"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/notify"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/post"
)

var (
	errCommentsNotAllow = errors.New("the post is closed from comments")
	errVeryLongContent  = errors.New("the comment is too long")
)

const MaxLenComment = 2000

type CommentUsecase struct {
	repo     Repository
	notifier notify.Notifier

	post postGetter
}

func NewCommentUsecase(
	repo Repository,
	notifier notify.Notifier,
	post postGetter,
) *CommentUsecase {
	return &CommentUsecase{
		repo:     repo,
		notifier: notifier,
		post:     post,
	}
}

func (c *CommentUsecase) WriteReply(ctx context.Context, comment *entity.Comment) error {
	return errors.Wrap(c.writeComment(post.Request{
		Entity: post.ReplyToComment(comment.Parent),
		Ctx:    ctx,
	}, comment, c.repo.AddReply), "write reply")
}

func (c *CommentUsecase) WriteComment(ctx context.Context, comment *entity.Comment) error {
	return errors.Wrap(c.writeComment(post.Request{
		Entity: post.CommentToPost(comment.Parent),
		Ctx:    ctx,
	}, comment, c.repo.AddComment), "write comment")
}

func (c *CommentUsecase) GetReplies(ctx context.Context, commentId, limit, cursor, depth int) (*entity.FeedComment, error) {
	queryCfg := entity.QueryConfig{
		Limit:  limit,
		Cursor: cursor,
		Depth:  depth,
	}

	return c.repo.GetReplies(ctx, commentId, queryCfg)
}

func (c *CommentUsecase) writeComment(
	req post.Request,
	comment *entity.Comment,
	addInRepository func(context.Context, *entity.Comment) error,
) error {
	if err := isValid(comment); err != nil {
		return errors.WrapFail(err, "check content")
	}

	post, err := c.post.GetPostByEntity(req)
	if err != nil {
		return errors.WrapFail(err, "get parent post")
	}

	if !post.AllowComment {
		return errors.WithType(errCommentsNotAllow, errors.TypeCommentsAreProhibited)
	}

	if err = addInRepository(req.Ctx, comment); err != nil {
		return errors.Wrap(err, "add to repository")
	}

	go c.notifier.PublishComment(req.Ctx, comment, post.ID)

	return nil
}

func isValid(comment *entity.Comment) error {
	if utf8.RuneCountInString(comment.Content) > MaxLenComment {
		return errors.WithType(errVeryLongContent, errors.TypeInvalidComment)
	}
	return nil
}
