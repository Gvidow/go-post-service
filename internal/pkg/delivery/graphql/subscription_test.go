package graphql

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/delivery/mock"
	"github.com/gvidow/go-post-service/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestSubscribeSuccess(t *testing.T) {
	defer goleak.VerifyNone(t)

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

	postId := 26
	comments := []*entity.Comment{{ID: 5}, {ID: 9}, {ID: 3}, {ID: 242}}

	chanNotifier := make(chan entity.NotifyComment, len(comments))

	usecaseMock.EXPECT().SubscribeOnPost(gomock.Any(), postId).Return(chanNotifier, nil)

	chanComment, err := resolver.Subscription().SubscribeOnPost(ctx, postId)
	require.NoError(t, err)

	for _, comment := range comments {
		chanNotifier <- entity.NotifyComment{Comment: comment}
	}
	close(chanNotifier)

	received := 0
	timer := time.After(200 * time.Millisecond)
	for ; received < len(comments); received++ {
		select {
		case comment := <-chanComment:
			require.Equal(t, comments[received], comment)
		case <-timer:
			t.Fatal("didn't have time to write it down")
		}
	}

	time.Sleep(50 * time.Millisecond)
	select {
	case _, ok := <-chanComment:
		require.False(t, ok, "channel is not closed")
	default:
		t.Fatal("channel is not closed")
	}
}
