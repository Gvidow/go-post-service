package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"github.com/gvidow/go-post-service/internal/pkg/usecase"
)

var (
	_ usecase.Repository = (*postgresRepo)(nil)
	_ PgxPoolIface       = (*pgxpool.Pool)(nil)
)

type postgresRepo struct {
	pool PgxPoolIface
}

func NewPostgresRepo(pool PgxPoolIface) *postgresRepo {
	return &postgresRepo{pool}
}

func (p *postgresRepo) GetComments(ctx context.Context, postIds []int, cfg entity.QueryConfig) (
	entity.BatchComments,
	error,
) {
	if len(postIds) == 0 {
		return nil, nil
	}

	query, args, err := sq.Select(
		"post_id", "id", "author", "content", "parent", "depth", "created_at",
	).
		From("f").
		Where("rank > ?", cfg.Cursor).
		Where("rank <= ?", cfg.Cursor+cfg.Limit).
		PrefixExpr(
			sq.Select(
				"path[1] AS post_id",
				"id",
				"author",
				"content",
				"path[icount(path)] AS parent",
				"icount(path) AS depth",
				"created_at",
				"RANK() OVER (PARTITION BY path[1] ORDER BY path + intset(id) + (SELECT MAX(id) + 1 FROM comment) DESC) AS rank",
			).
				From("comment").
				Where(sq.Eq{"path[1]": postIds}).
				Where("? < 0 OR icount(path) <= ?", cfg.Depth, cfg.Depth).
				Prefix("WITH f AS (").
				Suffix(")"),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query select")
	}

	res := make(entity.BatchComments, len(postIds))

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.WrapFail(err, "select feed comments for many posts from storage")
	}
	defer rows.Close()

	var (
		createdAt  pgtype.Timestamptz
		comment    entity.Comment
		prevPostID int
		postID     int

		comments = make([]entity.Comment, 0, cfg.Limit)
	)
	for rows.Next() {
		comment = entity.Comment{}
		if err = rows.Scan(
			&postID,
			&comment.ID,
			&comment.Author,
			&comment.Content,
			&comment.Parent,
			&comment.Depth,
			&createdAt,
		); err != nil {
			return nil, errors.Wrap(err, "scan selected comments")
		}

		comment.CreatedAt = int(createdAt.Time.UTC().UnixNano())
		switch {
		case postID != prevPostID && prevPostID != 0:
			res[prevPostID] = comments
			comments = make([]entity.Comment, 0, cfg.Limit)
			fallthrough
		case prevPostID == 0:
			prevPostID = postID
		}

		comments = append(comments, comment)
	}
	res[prevPostID] = comments

	return res, nil
}

func (p *postgresRepo) GetReplies(ctx context.Context, commentId int, cfg entity.QueryConfig) (
	*entity.FeedComment,
	error,
) {
	rows, err := p.pool.Query(ctx, SelectFeedReplies, commentId, cfg.Depth, cfg.Limit, cfg.Cursor)
	if err != nil {
		return nil, errors.WrapFail(err, "select feed replies from storage")
	}
	defer rows.Close()

	comments := make([]*entity.Comment, 0, cfg.Limit)
	createdAt := pgtype.Timestamptz{}
	for rows.Next() {
		comment := &entity.Comment{}
		if err = rows.Scan(
			&comment.ID,
			&comment.Author,
			&comment.Content,
			&comment.Parent,
			&comment.Depth,
			&createdAt,
		); err != nil {
			return nil, errors.Wrap(err, "scan selected replies")
		}
		comment.CreatedAt = int(createdAt.Time.UTC().UnixNano())
		comments = append(comments, comment)
	}

	return &entity.FeedComment{Comments: comments, Cursor: cfg.Cursor + len(comments)}, nil
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
	return errors.WrapFail(
		m.pool.QueryRow(
			ctx,
			InsertNewComment,
			comment.Author,
			comment.Content,
			comment.Parent,
		).Scan(&comment.ID),
		"add new comment to storage",
	)
}

func (m *postgresRepo) AddReply(ctx context.Context, comment *entity.Comment) error {
	return errors.WrapFail(
		m.pool.QueryRow(
			ctx,
			InsertNewReply,
			comment.Author,
			comment.Content,
			comment.Parent,
		).Scan(&comment.ID),
		"add new reply to storage",
	)
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
	rows, err := m.pool.Query(ctx, SelectFeedPosts, limit, cursor)
	if err != nil {
		return nil, errors.WrapFail(err, "select feed post from storage")
	}
	defer rows.Close()

	posts := make([]*entity.Post, 0, limit)
	createdAt := pgtype.Timestamptz{}
	for rows.Next() {
		post := &entity.Post{}
		if err = rows.Scan(
			&post.ID,
			&post.Author,
			&post.Title,
			&post.Content,
			&post.AllowComment,
			&createdAt,
		); err != nil {
			return nil, errors.Wrap(err, "scan selected feed post")
		}
		post.CreatedAt = int(createdAt.Time.UTC().UnixNano())
		posts = append(posts, post)
	}

	return &entity.FeedPost{Posts: posts, Cursor: cursor + len(posts)}, nil
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

	if err := errors.WrapFail(
		m.pool.QueryRow(ctx, query, id).
			Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.AllowComment, &createdAt),
		"get entity from storage",
	); errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithType(err, errors.TypePostNotFound)
	} else if err != nil {
		return nil, err
	}

	post.CreatedAt = int(createdAt.Time.UTC().UnixNano())
	return post, nil
}
