package post

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/mock"
	"github.com/stretchr/testify/require"
)

func TestPublishPost(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock.NewMockRepository(ctrl)

	usecase := NewPostUsecase(mockRepository, nil)

	post := &entity.Post{
		ID:           25,
		Author:       "apple",
		Content:      "test",
		AllowComment: true,
	}

	mockRepository.EXPECT().AddPost(ctx, post).Return(nil)

	err := usecase.PublishPost(ctx, post)
	require.NoError(t, err)
}

func TestAuthorUpdateAllowCommenting(t *testing.T) {
	tests := []struct {
		name   string
		author string
		postId int
		allow  bool
	}{
		{name: "ProhibitCommenting", author: "bear", postId: 45, allow: false},
		{name: "ProhibitCommenting", author: "bear", postId: 45, allow: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepository := mock.NewMockRepository(ctrl)

			usecase := NewPostUsecase(mockRepository, nil)

			mockRepository.EXPECT().GetPostById(ctx, tt.postId).Return(&entity.Post{ID: tt.postId, Author: tt.author}, nil)

			mockRepository.EXPECT().SetPermAddComments(ctx, tt.postId, tt.allow).Return(nil)
			var err error
			if tt.allow {
				err = usecase.AllowCommenting(ctx, tt.author, tt.postId)
			} else {
				err = usecase.ProhibitCommenting(ctx, tt.author, tt.postId)
			}

			require.NoError(t, err)
		})
	}
}

func TestNoAuthorUpdateAllowCommenting(t *testing.T) {
	tests := []struct {
		name   string
		author string
		postId int
		allow  bool
	}{
		{name: "ProhibitCommenting", author: "bear", postId: 45, allow: false},
		{name: "ProhibitCommenting", author: "bear", postId: 45, allow: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepository := mock.NewMockRepository(ctrl)

			usecase := NewPostUsecase(mockRepository, nil)

			mockRepository.EXPECT().GetPostById(ctx, tt.postId).Return(&entity.Post{ID: tt.postId, Author: "!" + tt.author}, nil)

			var err error
			if tt.allow {
				err = usecase.AllowCommenting(ctx, tt.author, tt.postId)
			} else {
				err = usecase.ProhibitCommenting(ctx, tt.author, tt.postId)
			}

			require.ErrorIs(t, err, errNotAllowed)
		})
	}
}

func TestGetPost(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock.NewMockRepository(ctrl)

	usecase := NewPostUsecase(mockRepository, nil)

	wantPost := &entity.Post{
		ID:           25,
		Author:       "apple",
		Content:      "test",
		AllowComment: true,
	}

	mockRepository.EXPECT().GetPostById(ctx, wantPost.ID).Return(wantPost, nil)

	actualPost, err := usecase.GetPost(ctx, wantPost.ID)
	require.NoError(t, err)
	require.Equal(t, wantPost, actualPost)
}

func TestSubscribeOnPost(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock.NewMockRepository(ctrl)
	mockNotifier := mock.NewMockNotifier(ctrl)

	usecase := NewPostUsecase(mockRepository, mockNotifier)

	post := &entity.Post{
		ID:           25,
		Author:       "apple",
		Content:      "test",
		AllowComment: true,
	}

	mockRepository.EXPECT().GetPostById(ctx, post.ID).Return(post, nil)

	wantChannelNotifier := make(chan entity.NotifyComment)
	mockNotifier.
		EXPECT().
		RegistryChanNotifier(
			ctx,
			gomock.AssignableToTypeOf((chan<- entity.NotifyComment)(wantChannelNotifier)),
			post.ID,
		).
		Return(nil)

	actualChannelNotifier, err := usecase.SubscribeOnPost(ctx, post.ID)
	require.NoError(t, err)
	require.IsType(t, actualChannelNotifier, (<-chan entity.NotifyComment)(wantChannelNotifier))
}
