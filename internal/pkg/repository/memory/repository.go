package memory

import (
	"context"
	"slices"
	stdsync "sync"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"github.com/gvidow/go-post-service/internal/pkg/usecase"
	"github.com/gvidow/go-post-service/pkg/sync"
)

var _ usecase.Repository = (*memoryRepo)(nil)

type memoryRepo struct {
	Post          sync.Map[uint64, *postItem]
	postLastIndex uint64

	Comment          sync.Map[uint64, *commentNode]
	commentLastIndex uint64

	// for uniq id and maintain the sorting of comments by id
	rwMutexComment stdsync.RWMutex

	// for uniq id and maintain the sorting of posts by id
	rwMutexPost stdsync.RWMutex
}

func NewMemoryRepo() *memoryRepo {
	return &memoryRepo{
		Post:    sync.NewMap(make(map[uint64]*postItem)),
		Comment: sync.NewMap(make(map[uint64]*commentNode)),
	}
}

func (m *memoryRepo) GetComments(_ context.Context, postIds []int, cfg entity.QueryConfig) (
	entity.BatchComments,
	error,
) {
	res := make(entity.BatchComments, len(postIds))

	for _, postId := range postIds {
		cmntsPointers, err := m.getComments(postId, cfg)
		if err != nil {
			return nil, errors.WrapFailf(err, "get comments for post with id = %d", postId)
		}
		comments := make([]entity.Comment, len(cmntsPointers))
		for ind := range comments {
			comments[ind] = *cmntsPointers[ind]
		}
		res[postId] = comments
	}

	return res, nil
}

func (m *memoryRepo) GetReplies(_ context.Context, commentId int, cfg entity.QueryConfig) (
	*entity.FeedComment,
	error,
) {
	comment, ok := m.Comment.Get(uint64(commentId))
	if !ok {
		return nil, errors.WithType(ErrNotFound, errors.CommentNotFound)
	}

	comments := getCommentsRecurse(comment, cfg.Limit+1, cfg.Depth)
	comments = comments[:len(comments)-1]
	slices.Reverse(comments)

	return &entity.FeedComment{Comments: comments, Cursor: cfg.Cursor + cfg.Limit}, nil
}

func (m *memoryRepo) GetPostById(_ context.Context, id int) (*entity.Post, error) {
	post, ok := m.Post.Get(uint64(id))
	if !ok {
		return nil, errors.WithType(ErrNotFound, errors.PostNotFound)
	}
	return post.dto(), nil
}

func (m *memoryRepo) GetPostByComment(_ context.Context, id int) (*entity.Post, error) {
	comment, ok := m.Comment.Get(uint64(id))
	if !ok {
		return nil, errors.WithType(ErrNotFound, errors.CommentNotFound)
	}

	post, _ := m.Post.Get(comment.ID)
	return post.dto(), nil
}

func (m *memoryRepo) AddComment(_ context.Context, comment *entity.Comment) error {
	postId := uint64(comment.Parent)
	post, ok := m.Post.Get(postId)
	if !ok {
		return errors.WithType(ErrNotFound, errors.PostNotFound)
	}

	m.rwMutexComment.Lock()
	defer m.rwMutexComment.Unlock()

	m.commentLastIndex++
	id := m.commentLastIndex
	comment.ID = int(id)
	node := newCommentNode(comment)
	node.PostID = postId
	m.Comment.Set(id, node)

	post.Comments.PushBack(node)
	return nil
}

func (m *memoryRepo) AddReply(_ context.Context, comment *entity.Comment) error {
	commentId := uint64(comment.Parent)
	node, ok := m.Comment.Get(commentId)
	if !ok {
		return errors.WithType(ErrNotFound, errors.CommentNotFound)
	}

	m.rwMutexComment.Lock()
	defer m.rwMutexComment.Unlock()

	m.commentLastIndex++
	id := m.commentLastIndex
	comment.ID = int(id)
	newNode := newCommentNode(comment)
	newNode.PostID = node.PostID
	m.Comment.Set(id, newNode)

	node.Replies.PushBack(newNode)
	return nil
}

func (m *memoryRepo) AddPost(_ context.Context, post *entity.Post) error {
	m.rwMutexPost.Lock()
	defer m.rwMutexPost.Unlock()

	m.postLastIndex++
	id := m.postLastIndex
	post.ID = int(id)
	m.Post.Set(id, newPostItem(post))
	return nil
}

func (m *memoryRepo) GetFeedPosts(_ context.Context, limit, cursor int) (*entity.FeedPost, error) {
	posts := make([]*entity.Post, 0, limit)

	m.rwMutexPost.RLock()
	maxID := m.postLastIndex
	m.rwMutexPost.RUnlock()

	if cursor > 0 {
		maxID = uint64(cursor)
	}
	minID := max(maxID-uint64(limit), 0)

	for ind := maxID; ind > minID; ind-- {
		post, _ := m.Post.Get(ind)
		posts = append(posts, post.dto())
	}

	return &entity.FeedPost{
		Posts:  posts,
		Cursor: int(minID),
	}, nil
}

func (m *memoryRepo) SetPermAddComments(_ context.Context, postId int, allow bool) error {
	post, ok := m.Post.Get(uint64(postId))
	if !ok {
		return errors.WithType(ErrNotFound, errors.PostNotFound)
	}

	post.AllowComment.Store(allow)
	return nil
}

func (m *memoryRepo) getComments(postId int, cfg entity.QueryConfig) (
	[]*entity.Comment,
	error,
) {
	post, ok := m.Post.Get(uint64(postId))
	if !ok {
		return nil, errors.WithType(ErrNotFound, errors.PostNotFound)
	}

	comments := getCommentsRecurse(&commentNode{Replies: post.Comments}, cfg.Limit+1, cfg.Depth)
	comments = comments[:len(comments)-1]
	slices.Reverse(comments)

	return comments, nil
}

func getCommentsRecurse(node *commentNode, limit int, depth int) []*entity.Comment {
	switch {
	case limit <= 0 || depth == 0:
		return nil
	case depth == 1:
		countReplies := node.Replies.Len()

		res := make([]*entity.Comment, 0, limit)

		for len(res)+1 < limit && countReplies > 0 {
			countReplies--
			res = append(res, node.Replies.Get(countReplies).dto())
		}

		return append(res, node.dto())
	}

	countReplies := node.Replies.Len()

	res := make([]*entity.Comment, 0, limit)

	for len(res)+1 < limit && countReplies > 0 {
		countReplies--
		res = append(res, getCommentsRecurse(node.Replies.Get(countReplies), limit-len(res)-1, depth-1)...)
	}

	return append(res, node.dto())
}
