package graphql

import (
	"context"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) PublishPost(
	ctx context.Context,
	author string,
	title string,
	content string,
	allowComment bool,
) (*entity.Post, error) {

	post := &entity.Post{
		Author:       author,
		Title:        title,
		Content:      content,
		AllowComment: allowComment,
	}

	if err := errors.WrapFail(
		r.usecase.PublishPost(ctx, post),
		"publish post",
	); err != nil {
		return nil, err
	}

	return post, nil
}

func (r *mutationResolver) AddCommentToPost(
	ctx context.Context,
	author string,
	postID int,
	content string,
) (*entity.Comment, error) {

	comment := &entity.Comment{
		Author:  author,
		Content: content,
		Parent:  postID,
		Depth:   1,
	}

	if err := errors.WrapFail(
		r.usecase.WriteComment(ctx, comment),
		"write comment to post",
	); err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *mutationResolver) AddCommentToComment(
	ctx context.Context,
	author string,
	commentID int,
	content string,
) (*entity.Comment, error) {

	comment := &entity.Comment{
		Author:  author,
		Content: content,
		Parent:  commentID,
	}

	if err := errors.WrapFail(
		r.usecase.WriteReply(ctx, comment),
		"write reply to comment",
	); err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *mutationResolver) ProhibitWritingComments(ctx context.Context, author string, postID int) (bool, error) {
	if err := errors.WrapFail(
		r.usecase.ProhibitCommenting(ctx, author, postID),
		"prohibit commenting",
	); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) AllowWritingComments(ctx context.Context, author string, postID int) (bool, error) {
	if err := errors.WrapFail(
		r.usecase.AllowCommenting(ctx, author, postID),
		"allow commenting",
	); err != nil {
		return false, err
	}

	return true, nil
}
