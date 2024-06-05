package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"

	"github.com/gvidow/go-post-service/internal/entity"
)

func TestAddNewUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewPostgresRepo(pool)
	wantID := 6

	pool.ExpectQuery("INSERT INTO post").
		WithArgs("author1", "new", "good", true).
		WillReturnRows(
			pgxmock.NewRows([]string{"id"}).
				AddRow(wantID),
		)

	post := &entity.Post{
		Author:       "author1",
		Title:        "new",
		Content:      "good",
		AllowComment: true,
	}

	err = repo.AddPost(ctx, post)
	require.NoError(t, err)
	require.Equal(t, wantID, post.ID)
}

func TestAddNewComment(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewPostgresRepo(pool)
	wantID := 6

	pool.ExpectQuery("INSERT INTO comment").
		WithArgs("author1", "good", 4).
		WillReturnRows(
			pgxmock.NewRows([]string{"id"}).
				AddRow(wantID),
		)

	comment := &entity.Comment{
		Author:  "author1",
		Content: "good",
		Parent:  4,
	}

	err = repo.AddComment(ctx, comment)
	require.NoError(t, err)
	require.Equal(t, wantID, comment.ID)
}

func TestGetFeedPosts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewPostgresRepo(pool)
	time := pgtype.Timestamptz{Time: time.Now()}

	pool.ExpectQuery("SELECT id, author, title, content, allow_comment, created_at FROM post").
		WithArgs(5, 6).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "author", "title", "content", "allow_comment", "created_at"}).
				AddRow(6, "a", "b", "c", true, time),
		)

	feed, err := repo.GetFeedPosts(ctx, 5, 6)
	require.NoError(t, err)
	require.Equal(t, &entity.FeedPost{
		Cursor: 7,
		Posts: []*entity.Post{{
			ID:           6,
			Author:       "a",
			Title:        "b",
			Content:      "c",
			AllowComment: true,
			CreatedAt:    int(time.Time.UTC().UnixNano()),
		}},
	}, feed)
}

func TestGetComments(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewPostgresRepo(pool)
	time := pgtype.Timestamptz{Time: time.Now()}
	cfg := entity.QueryConfig{Limit: 7, Cursor: 1}

	pool.ExpectQuery("SELECT post_id, id, author, content, parent, depth, created_at FROM f").
		WithArgs(2, 3, cfg.Cursor, cfg.Limit+cfg.Cursor).
		WillReturnRows(
			pgxmock.NewRows([]string{"post_id", "id", "author", "content", "parent", "depth", "created_at"}).
				AddRow(2, 6, "a", "b", 8, 9, time).
				AddRow(2, 7, "a", "b", 8, 9, time).
				AddRow(3, 8, "a", "b", 8, 9, time),
		)

	batch, err := repo.GetComments(ctx, []int{2, 3}, cfg)
	require.NoError(t, err)
	require.Len(t, batch, 2)
	require.Equal(t, batch[2], []entity.Comment([]entity.Comment{
		{ID: 6, Author: "a", Content: "b", Parent: 8, Depth: 9, CreatedAt: int(time.Time.UTC().UnixNano())},
		{ID: 7, Author: "a", Content: "b", Parent: 8, Depth: 9, CreatedAt: int(time.Time.UTC().UnixNano())},
	}))
	require.Equal(t, batch[3], []entity.Comment([]entity.Comment{
		{ID: 8, Author: "a", Content: "b", Parent: 8, Depth: 9, CreatedAt: int(time.Time.UTC().UnixNano())},
	}))
}
