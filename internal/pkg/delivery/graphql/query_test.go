package graphql

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/delivery/mock"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"github.com/gvidow/go-post-service/pkg/logger"
)

func TestPosts(t *testing.T) {
	tests := []struct {
		name        string
		wantErr     error
		wantMessage string
		cursor      int
		limit       int
		feed        *entity.FeedPost
	}{
		{
			name:    "successful get feed posts",
			wantErr: nil,
			cursor:  15,
			limit:   40,
			feed:    &entity.FeedPost{Cursor: 55},
		},
		{
			name:        "fail get feed posts",
			wantErr:     errors.WithType(errors.New("fail"), errors.TypeInvalidComment),
			wantMessage: "the length of the comment does not fit into 2000 characters",
			limit:       15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			log, err := logger.New()
			if err != nil {
				t.Fatal(err)
			}

			usecaseMock := mock.NewMockUsecase(ctrl)

			resolver := Resolver{
				log:     log,
				usecase: usecaseMock,
			}

			usecaseMock.EXPECT().GetFeedPosts(ctx, tt.limit, tt.cursor).Return(tt.feed, tt.wantErr)

			actualFeed, err := resolver.Query().Posts(ctx, tt.limit, tt.cursor)

			require.Equal(t, tt.feed, actualFeed)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.IsType(t, &responseError{}, err)
				require.Equal(t, tt.wantMessage, err.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestGetPost(t *testing.T) {
	tests := []struct {
		name        string
		wantErr     error
		wantMessage string
		postId      int
		wantPost    *entity.Post
	}{
		{
			name:     "successful get one post",
			wantErr:  nil,
			postId:   62,
			wantPost: &entity.Post{ID: 62, Title: "ozon"},
		},
		{
			name:        "fail get post",
			wantErr:     errors.WithType(errors.New("fail"), errors.TypePostNotFound),
			wantMessage: "the post was not found",
			postId:      15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			log, err := logger.New()
			if err != nil {
				t.Fatal(err)
			}

			usecaseMock := mock.NewMockUsecase(ctrl)

			resolver := Resolver{
				log:     log,
				usecase: usecaseMock,
			}

			usecaseMock.EXPECT().GetPost(ctx, tt.postId).Return(tt.wantPost, tt.wantErr)

			actualPost, err := resolver.Query().GetPost(ctx, tt.postId)

			require.Equal(t, tt.wantPost, actualPost)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.IsType(t, &responseError{}, err)
				require.Equal(t, tt.wantMessage, err.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestReplies(t *testing.T) {
	tests := []struct {
		name        string
		wantErr     error
		wantMessage string
		commentId   int
		cursor      int
		limit       int
		depth       int
		feed        *entity.FeedComment
	}{
		{
			name:      "successful get replies",
			wantErr:   nil,
			cursor:    15,
			commentId: 17,
			limit:     40,
			depth:     14,
			feed:      &entity.FeedComment{Cursor: 144},
		},
		{
			name:        "fail get feed replies",
			wantErr:     errors.WithType(errors.New("fail"), errors.TypeCommentNotFound),
			wantMessage: "the comment was not found",
			commentId:   86,
			cursor:      245,
			depth:       15,
			limit:       150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			log, err := logger.New()
			if err != nil {
				t.Fatal(err)
			}

			usecaseMock := mock.NewMockUsecase(ctrl)

			resolver := Resolver{
				log:     log,
				usecase: usecaseMock,
			}

			usecaseMock.EXPECT().GetReplies(ctx, tt.commentId, tt.limit, tt.cursor, tt.depth).Return(tt.feed, tt.wantErr)

			actualFeed, err := resolver.Query().Replies(ctx, tt.commentId, tt.limit, tt.cursor, tt.depth)

			require.Equal(t, tt.feed, actualFeed)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.IsType(t, &responseError{}, err)
				require.Equal(t, tt.wantMessage, err.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}
