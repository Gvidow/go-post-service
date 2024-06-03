package memory

import (
	"sync/atomic"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/pkg/sync"
)

type postItem struct {
	ID           uint64
	Author       string
	Title        string
	Content      string
	AllowComment atomic.Bool
	CreatedAt    int
	Comments     sync.Slice[*commentNode]
}

func newPostItem(post *entity.Post) *postItem {
	resPost := &postItem{
		ID:        uint64(post.ID),
		Author:    post.Author,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
	}
	resPost.AllowComment.Store(post.AllowComment)
	return resPost
}

func (p *postItem) dto() *entity.Post {
	return &entity.Post{
		ID:           int(p.ID),
		Author:       p.Author,
		Title:        p.Title,
		Content:      p.Content,
		AllowComment: p.AllowComment.Load(),
		CreatedAt:    p.CreatedAt,
	}
}

type commentNode struct {
	ID        uint64
	Author    string
	Content   string
	PostID    uint64
	Parent    uint64
	Depth     int
	CreatedAt int
	Replies   sync.Slice[*commentNode]
}

func newCommentNode(comment *entity.Comment) *commentNode {
	return &commentNode{
		ID:        uint64(comment.ID),
		Author:    comment.Author,
		Content:   comment.Content,
		Parent:    uint64(comment.Parent),
		Depth:     comment.Depth,
		CreatedAt: comment.CreatedAt,
	}
}

func (c *commentNode) dto() *entity.Comment {
	return &entity.Comment{
		ID:        int(c.ID),
		Author:    c.Author,
		Content:   c.Content,
		Parent:    int(c.Parent),
		Depth:     c.Depth,
		CreatedAt: c.CreatedAt,
	}
}
