package graphql

import (
	"context"
	"fmt"

	"github.com/gvidow/go-post-service/internal/entity"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) PublishPost(ctx context.Context, author string, title string, content string, allowComment bool) (*entity.Post, error) {
	panic(fmt.Errorf("not implemented: PublishPost - publishPost"))
}

func (r *mutationResolver) AddCommentToPost(ctx context.Context, author string, postID int, content string) (*entity.Comment, error) {
	panic(fmt.Errorf("not implemented: AddCommentToPost - addCommentToPost"))
}

func (r *mutationResolver) AddCommentToComment(ctx context.Context, author string, commentID int, content string) (*entity.Comment, error) {
	panic(fmt.Errorf("not implemented: AddCommentToComment - addCommentToComment"))
}

func (r *mutationResolver) ProhibitWritingComments(ctx context.Context, author string, postID int) (bool, error) {
	panic(fmt.Errorf("not implemented: ProhibitWritingComments - prohibitWritingComments"))
}

func (r *mutationResolver) AllowWritingComments(ctx context.Context, author string, postID int) (bool, error) {
	panic(fmt.Errorf("not implemented: AllowWritingComments - allowWritingComments"))
}
