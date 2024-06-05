package memory

import (
	"context"
	"testing"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/pkg/sync"
	"github.com/stretchr/testify/require"
)

func TestGetTreeCommentsForManyPosts(t *testing.T) {
	memPost := map[uint64]*postItem{
		1: {ID: 1, Comments: sync.NewSlice([]*commentNode{
			{ID: 2, Replies: sync.NewSlice([]*commentNode{
				{ID: 4, Replies: sync.NewSlice([]*commentNode{})},
			})},
			{ID: 3, Replies: sync.NewSlice([]*commentNode{})},
		})},
		2: {ID: 2, Comments: sync.NewSlice([]*commentNode{
			{ID: 5, Replies: sync.NewSlice([]*commentNode{
				{ID: 6, Replies: sync.NewSlice([]*commentNode{
					{ID: 8, Replies: sync.NewSlice([]*commentNode{
						{ID: 9, Replies: sync.NewSlice([]*commentNode{})},
					})},
				})},
				{ID: 8, Replies: sync.NewSlice([]*commentNode{})},
			})},
			{ID: 7, Replies: sync.NewSlice([]*commentNode{})},
		})},
	}

	mem := NewMemoryRepo()
	mem.Post = sync.NewMap(memPost)
	mem.postLastIndex = uint64(mem.Post.Len())

	wantComments := []int{4, 8, 6}
	posts := []int{1, 2}
	batch, err := mem.GetComments(context.Background(), posts, entity.QueryConfig{Limit: 5, Depth: 2, Cursor: 2})
	require.NoError(t, err)
	require.NotEmpty(t, batch)

	var actualComments []int
	for _, postsId := range posts {
		for _, comment := range batch[postsId] {
			actualComments = append(actualComments, comment.ID)
		}
	}

	require.Equal(t, wantComments, actualComments)
}
