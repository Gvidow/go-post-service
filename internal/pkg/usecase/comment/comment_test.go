package comment

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/comment/mock"
	mocknotify "github.com/gvidow/go-post-service/internal/pkg/usecase/mock"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/post"
	"github.com/stretchr/testify/require"
)

func TestWriteCommentUnderPermissionOfThisPost(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock.NewMockRepository(ctrl)
	mockPost := mock.NewMockpostGetter(ctrl)
	mockNotifier := mocknotify.NewMockNotifier(ctrl)

	usecase := NewCommentUsecase(mockRepository, mockNotifier, mockPost)

	comment := &entity.Comment{
		ID:      24,
		Author:  "man",
		Content: "ozon",
		Parent:  21,
		Depth:   17,
	}

	p := &entity.Post{
		ID:           84,
		Author:       "man",
		Content:      "ozon",
		AllowComment: true,
	}

	mockPost.EXPECT().GetPostByEntity(post.Request{Entity: post.CommentToPost(comment.Parent), Ctx: ctx}).Return(p, nil)
	mockRepository.EXPECT().AddComment(ctx, comment).Return(nil)
	mockNotifier.EXPECT().PublishComment(ctx, comment, p.ID)

	err := usecase.WriteComment(ctx, comment)
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)
}

func TestWriteCommentUnderForbiddingOfThisPost(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock.NewMockRepository(ctrl)
	mockPost := mock.NewMockpostGetter(ctrl)
	mockNotifier := mocknotify.NewMockNotifier(ctrl)

	usecase := NewCommentUsecase(mockRepository, mockNotifier, mockPost)

	comment := &entity.Comment{
		ID:      31,
		Author:  "man",
		Content: "ozon",
		Parent:  81,
		Depth:   17,
	}

	p := &entity.Post{
		ID:           8,
		Author:       "man",
		Content:      "ozon",
		AllowComment: false,
	}

	mockPost.EXPECT().GetPostByEntity(post.Request{Entity: post.CommentToPost(comment.Parent), Ctx: ctx}).Return(p, nil)

	err := usecase.WriteComment(ctx, comment)
	require.ErrorIs(t, err, errCommentsNotAllow)
}

func TestWriteReplyUnderPermissionOfThisPost(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock.NewMockRepository(ctrl)
	mockPost := mock.NewMockpostGetter(ctrl)
	mockNotifier := mocknotify.NewMockNotifier(ctrl)

	usecase := NewCommentUsecase(mockRepository, mockNotifier, mockPost)

	comment := &entity.Comment{
		ID:      24,
		Author:  "man",
		Content: "ozon",
		Parent:  21,
		Depth:   17,
	}

	p := &entity.Post{
		ID:           84,
		Author:       "man",
		Content:      "ozon",
		AllowComment: true,
	}

	mockPost.EXPECT().GetPostByEntity(post.Request{Entity: post.ReplyToComment(comment.Parent), Ctx: ctx}).Return(p, nil)
	mockRepository.EXPECT().AddReply(ctx, comment).Return(nil)
	mockNotifier.EXPECT().PublishComment(ctx, comment, p.ID)

	err := usecase.WriteReply(ctx, comment)
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)
}

func TestVeryLongPost(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock.NewMockRepository(ctrl)
	mockPost := mock.NewMockpostGetter(ctrl)
	mockNotifier := mocknotify.NewMockNotifier(ctrl)

	usecase := NewCommentUsecase(mockRepository, mockNotifier, mockPost)

	comment := &entity.Comment{
		ID:      24,
		Author:  "man",
		Content: strings.Repeat("a", MaxLenComment+1),
		Parent:  21,
		Depth:   17,
	}

	err := usecase.WriteReply(ctx, comment)
	require.ErrorIs(t, err, errVeryLongContent)
}

func TestGetReplies(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock.NewMockRepository(ctrl)
	mockPost := mock.NewMockpostGetter(ctrl)
	mockNotifier := mocknotify.NewMockNotifier(ctrl)

	usecase := NewCommentUsecase(mockRepository, mockNotifier, mockPost)

	var (
		commentId = 3
		limit     = 49
		cursor    = 17
		depth     = 5
	)

	cfg := entity.QueryConfig{
		Cursor: cursor,
		Limit:  limit,
		Depth:  depth,
	}

	wantFeed := &entity.FeedComment{
		Comments: []*entity.Comment{{ID: 16}},
		Cursor:   27,
	}

	mockRepository.EXPECT().GetReplies(ctx, commentId, cfg).Return(wantFeed, nil)

	actualFeed, err := usecase.GetReplies(ctx, commentId, limit, cursor, depth)
	require.NoError(t, err)
	require.Equal(t, wantFeed, actualFeed)
}
