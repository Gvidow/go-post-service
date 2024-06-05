package memory

import (
	"context"
	stdsync "sync"
	"testing"

	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/internal/pkg/errors"
	"github.com/gvidow/go-post-service/pkg/sync"
	"github.com/stretchr/testify/require"
)

func TestAddComments(t *testing.T) {
	tests := []struct {
		name        string
		memPost     sync.Map[uint64, *postItem]
		memComment  sync.Map[uint64, *commentNode]
		comment     *entity.Comment
		funcAdd     func(*memoryRepo, context.Context, *entity.Comment) error
		wantTypeErr errors.Type
	}{
		{
			name:        "add coment to empty memory",
			memPost:     sync.NewMap(make(map[uint64]*postItem)),
			memComment:  sync.NewMap(make(map[uint64]*commentNode)),
			comment:     &entity.Comment{Content: "new comment", Parent: 12},
			funcAdd:     (*memoryRepo).AddComment,
			wantTypeErr: errors.TypePostNotFound,
		},
		{
			name:        "add reply to empty memory",
			memPost:     sync.NewMap(make(map[uint64]*postItem)),
			memComment:  sync.NewMap(make(map[uint64]*commentNode)),
			comment:     &entity.Comment{Content: "new comment", Parent: 42},
			funcAdd:     (*memoryRepo).AddReply,
			wantTypeErr: errors.TypeCommentNotFound,
		},
		{
			name: "adding a comment to an existing post",
			memPost: sync.NewMap(map[uint64]*postItem{
				1: {ID: 1, Content: "i am exists"},
				2: {ID: 1, Content: "i am exists too"},
			}),
			memComment: sync.NewMap(make(map[uint64]*commentNode)),
			comment:    &entity.Comment{Content: "new comment", Parent: 2},
			funcAdd:    (*memoryRepo).AddComment,
		},
		{
			name: "adding a comment to an existing comment",
			memPost: sync.NewMap(map[uint64]*postItem{
				1: {ID: 1, Content: "i am exists"},
				2: {ID: 2, Content: "i am exists too"},
			}),
			memComment: sync.NewMap(map[uint64]*commentNode{
				1: {ID: 1, Content: "i am exists too", PostID: 2, Parent: 2},
				2: {ID: 2, Content: "i am exists too", PostID: 2, Parent: 2},
				3: {ID: 3, Content: "i am exists too", PostID: 2, Parent: 2},
			}),
			comment: &entity.Comment{Content: "new comment", Parent: 3},
			funcAdd: (*memoryRepo).AddReply,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := &memoryRepo{
				Post:             tt.memPost,
				postLastIndex:    uint64(tt.memPost.Len()),
				Comment:          tt.memComment,
				commentLastIndex: uint64(tt.memComment.Len()),
			}

			actualErr := tt.funcAdd(mem, context.Background(), tt.comment)
			if tt.wantTypeErr == 0 {
				require.NoError(t, actualErr)
				return
			}

			var actualTypeErr errors.TypeError
			require.ErrorAs(t, actualErr, &actualTypeErr)
			require.Equal(t, tt.wantTypeErr, actualTypeErr.Type(), "the wrong type of error")
		})
	}
}

func TestParallelAddPost(t *testing.T) {
	mem := NewMemoryRepo()
	countAdd := 100000
	wg := stdsync.WaitGroup{}

	wg.Add(countAdd)
	for i := 0; i < countAdd; i++ {
		go func() {
			defer wg.Done()
			mem.AddPost(context.Background(), &entity.Post{Content: "new post"})
		}()
	}
	wg.Wait()

	require.Equal(t, countAdd, mem.Post.Len(), "posts are lost")
}

func TestParallelAddComment(t *testing.T) {
	mem := NewMemoryRepo()
	countAdd := 100000
	wg := stdsync.WaitGroup{}

	wg.Add(countAdd)
	for i := 0; i < countAdd; i++ {
		go func() {
			defer wg.Done()
			post := &entity.Post{Content: "new post"}
			mem.AddPost(context.Background(), post)
			mem.AddComment(context.Background(), &entity.Comment{Parent: post.ID, Content: "new comment"})
		}()
	}
	wg.Wait()

	require.Equal(t, countAdd, mem.Comment.Len(), "posts are lost")
}
