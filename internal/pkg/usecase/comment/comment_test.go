package comment

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/comment/mock"
	mocknotify "github.com/gvidow/go-post-service/internal/pkg/usecase/mock"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/post"
	"github.com/stretchr/testify/require"
)

func TestXxx(t *testing.T) {
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
