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

func TestPublishPost(t *testing.T) {
	tests := []struct {
		name         string
		wantErr      error
		wantMessage  string
		author       string
		title        string
		content      string
		allowComment bool
	}{
		{
			name:         "success publish",
			wantErr:      nil,
			author:       "author",
			title:        "example",
			content:      "For example ...",
			allowComment: true,
		},
		{
			name:         "fail publish with unknow type error",
			wantErr:      errors.New("fail"),
			wantMessage:  "internal server error",
			author:       "bad",
			allowComment: true,
		},
		{
			name:         "fail publish with typing error",
			wantErr:      errors.WithType(errors.New("fail"), errors.TypeNotPermission),
			wantMessage:  "there are no rights to perform the action",
			author:       "bad",
			allowComment: true,
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

			wantPost := &entity.Post{
				Author:       tt.author,
				Title:        tt.title,
				Content:      tt.content,
				AllowComment: tt.allowComment,
			}

			usecaseMock.EXPECT().PublishPost(ctx, wantPost).Return(tt.wantErr)

			actualPost, err := resolver.Mutation().PublishPost(
				ctx,
				wantPost.Author,
				wantPost.Title,
				wantPost.Content,
				wantPost.AllowComment,
			)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.IsType(t, &responseError{}, err)
				require.Equal(t, tt.wantMessage, err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, wantPost, actualPost)
		})
	}
}

func TestAddCommentToPost(t *testing.T) {
	tests := []struct {
		name        string
		wantErr     error
		wantMessage string
		author      string
		postId      int
		content     string
	}{
		{
			name:    "success write",
			wantErr: nil,
			author:  "author",
			content: "For example ...",
			postId:  5,
		},
		{
			name:        "fail write comment with typing error",
			wantErr:     errors.WithType(errors.New("fail"), errors.TypeUnknow),
			wantMessage: "internal server error",
			author:      "bad",
			postId:      34,
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

			wantComment := &entity.Comment{
				Author:  tt.author,
				Parent:  tt.postId,
				Content: tt.content,
				Depth:   1,
			}

			usecaseMock.EXPECT().WriteComment(ctx, wantComment).Return(tt.wantErr)

			actualPost, err := resolver.Mutation().AddCommentToPost(
				ctx,
				wantComment.Author,
				tt.postId,
				wantComment.Content,
			)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.IsType(t, &responseError{}, err)
				require.Equal(t, tt.wantMessage, err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, wantComment, actualPost)
		})
	}
}

func TestAddCommentToComment(t *testing.T) {
	tests := []struct {
		name        string
		wantErr     error
		wantMessage string
		author      string
		postId      int
		content     string
	}{
		{
			name:    "success write",
			wantErr: nil,
			author:  "author",
			content: "For example ...",
			postId:  5,
		},
		{
			name:        "fail write reply with typing error",
			wantErr:     errors.WithType(errors.New("fail"), errors.TypeUnknow),
			wantMessage: "internal server error",
			author:      "bad",
			postId:      34,
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

			wantComment := &entity.Comment{
				Author:  tt.author,
				Parent:  tt.postId,
				Content: tt.content,
			}

			usecaseMock.EXPECT().WriteReply(ctx, wantComment).Return(tt.wantErr)

			actualPost, err := resolver.Mutation().AddCommentToComment(
				ctx,
				wantComment.Author,
				tt.postId,
				wantComment.Content,
			)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.IsType(t, &responseError{}, err)
				require.Equal(t, tt.wantMessage, err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, wantComment, actualPost)
		})
	}
}

func TestProhibitWritingComments(t *testing.T) {
	tests := []struct {
		name        string
		wantErr     error
		wantMessage string
		author      string
		postId      int
		wantStatus  bool
	}{
		{
			name:       "success prohibit",
			wantErr:    nil,
			author:     "author",
			postId:     52,
			wantStatus: true,
		},
		{
			name:        "fail prohibit",
			wantErr:     errors.WithType(errors.New("fail"), errors.TypeCommentsAreProhibited),
			wantMessage: "it is forbidden to leave comments under the post",
			author:      "bad",
			postId:      34,
			wantStatus:  false,
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

			usecaseMock.EXPECT().ProhibitCommenting(ctx, tt.author, tt.postId).Return(tt.wantErr)

			actualStatus, err := resolver.Mutation().ProhibitWritingComments(
				ctx,
				tt.author,
				tt.postId,
			)

			require.Equal(t, tt.wantStatus, actualStatus)

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

func TestAllowWritingComments(t *testing.T) {
	tests := []struct {
		name        string
		wantErr     error
		wantMessage string
		author      string
		postId      int
		wantStatus  bool
	}{
		{
			name:       "status change passed without error",
			wantErr:    nil,
			author:     "author",
			postId:     52,
			wantStatus: true,
		},
		{
			name:        "status change failed",
			wantErr:     errors.New("fail"),
			wantMessage: "internal server error",
			author:      "bad",
			postId:      34,
			wantStatus:  false,
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

			usecaseMock.EXPECT().AllowCommenting(ctx, tt.author, tt.postId).Return(tt.wantErr)

			actualStatus, err := resolver.Mutation().AllowWritingComments(
				ctx,
				tt.author,
				tt.postId,
			)

			require.Equal(t, tt.wantStatus, actualStatus)

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
