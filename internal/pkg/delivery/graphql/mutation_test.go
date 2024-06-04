package graphql

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/delivery/mock"
	"github.com/gvidow/go-post-service/pkg/logger"
)

func TestPublishPost(t *testing.T) {
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
		Author:       "author",
		Title:        "example",
		Content:      "For example ...",
		AllowComment: true,
	}

	usecaseMock.EXPECT().PublishPost(ctx, wantPost).Return(nil)

	actualPost, err := resolver.Mutation().PublishPost(ctx, wantPost.Author, wantPost.Title, wantPost.Content, wantPost.AllowComment)
	require.NoError(t, err)
	require.Equal(t, wantPost, actualPost)
}
