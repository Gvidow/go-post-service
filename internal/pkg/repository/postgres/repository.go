package postgres

import (
	"context"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"github.com/gvidow/go-post-service/internal/pkg/usecase"
)

var _ usecase.Repository = (*postgresRepo)(nil)

type postgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(pool *pgxpool.Pool) *postgresRepo {
	return &postgresRepo{pool}
}

func (p *postgresRepo) GetComments(ctx context.Context, postIds []int, cfg entity.QueryConfig) (
	entity.BatchComments,
	error,
) {
	res := make(entity.BatchComments, len(postIds))
	_ = res

	return nil, nil
}

func (p *postgresRepo) GetReplies(ctx context.Context, commentId int, cfg entity.QueryConfig) (
	*entity.FeedComment,
	error,
) {
	rows, err := p.pool.Query(ctx, SelectFeedReplies, commentId, cfg.Depth, cfg.Limit, cfg.Cursor)
	if err != nil {
		return nil, errors.WrapFail(err, "select replies")
	}
	defer rows.Close()

	for rows.Next() {

	}
	comments := make([]*entity.Comment, 0)
	return &entity.FeedComment{Comments: comments, Cursor: cfg.Cursor + cfg.Limit}, nil
}

func (m *postgresRepo) GetPostById(ctx context.Context, postId int) (*entity.Post, error) {
	post, err := m.getPostUniversalId(ctx, SelectPostById, postId)
	if err != nil {
		return nil, errors.WrapFail(err, "get by post id")
	}
	return post, nil
}

func (m *postgresRepo) GetPostByComment(ctx context.Context, commentId int) (*entity.Post, error) {
	post, err := m.getPostUniversalId(ctx, SelectPostByComment, commentId)
	if err != nil {
		return nil, errors.WrapFail(err, "get by comment id")
	}
	return post, nil
}

func (m *postgresRepo) AddComment(ctx context.Context, comment *entity.Comment) error {
	return nil
}

func (m *postgresRepo) AddReply(ctx context.Context, comment *entity.Comment) error {
	return nil
}

func (m *postgresRepo) AddPost(ctx context.Context, post *entity.Post) error {
	return errors.WrapFail(
		m.pool.QueryRow(
			ctx,
			InsertNewPost,
			post.Author,
			post.Title,
			post.Content,
			post.AllowComment,
		).Scan(&post.ID),
		"add new post to storage",
	)
}

func (m *postgresRepo) GetFeedPosts(ctx context.Context, limit, cursor int) (*entity.FeedPost, error) {
	return nil, nil
}

func (m *postgresRepo) SetPermAddComments(ctx context.Context, postId int, allow bool) error {
	if _, err := m.pool.Exec(ctx, UpdateCommentingPermission, allow, postId); err != nil {
		return errors.WrapFailf(err, "update post(id=%d) commenting permission", postId)
	}
	return nil
}

func (m *postgresRepo) getPostUniversalId(ctx context.Context, query string, id int) (*entity.Post, error) {
	post := &entity.Post{}
	createdAt := pgtype.Timestamptz{}

	if err := errors.WrapFailf(
		m.pool.QueryRow(ctx, query, id).
			Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.AllowComment, &createdAt),
		"get post(id=%d) from storage", id,
	); errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithType(err, errors.PostNotFound)
	} else if err != nil {
		return nil, err
	}

	post.CreatedAt = int(createdAt.Time.UTC().UnixNano())
	return post, nil
}
