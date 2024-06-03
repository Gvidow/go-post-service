package postgres

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"github.com/gvidow/go-post-service/internal/pkg/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ usecase.Repository = (*postgresRepo)(nil)

type postgresRepo struct {
	pool pgxpool.Pool
}

func NewPostgresRepo(pool pgxpool.Pool) *postgresRepo {
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

func (m *postgresRepo) GetPostById(ctx context.Context, id int) (*entity.Post, error) {
	post := &entity.Post{ID: id}

	if err := errors.WrapFail(
		m.pool.QueryRow(ctx, SelectPostById, id).
			Scan(&post.Author, &post.Title, &post.Content, &post.AllowComment, &post.CreatedAt),
		"get post by id from storage",
	); err != nil {
		return nil, err
	}

	return post, nil
}

func (m *postgresRepo) GetPostByComment(ctx context.Context, id int) (*entity.Post, error) {
	return nil, nil
}

func (m *postgresRepo) AddComment(ctx context.Context, comment *entity.Comment) error {
	return nil
}

func (m *postgresRepo) AddReply(ctx context.Context, comment *entity.Comment) error {
	return nil
}

func (m *postgresRepo) AddPost(_ context.Context, post *entity.Post) error {
	return nil
}

func (m *postgresRepo) GetFeedPosts(ctx context.Context, limit, cursor int) (*entity.FeedPost, error) {
	return nil, nil
}

func (m *postgresRepo) SetPermAddComments(ctx context.Context, postId int, allow bool) error {
	return nil
}
