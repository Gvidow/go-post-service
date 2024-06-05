package notify

import (
	"context"

	"github.com/google/uuid"
	"github.com/gvidow/go-post-service/internal/entity"
)

type client struct {
	id        uuid.UUID
	transport chan<- entity.NotifyComment
	postId    int
	ctx       context.Context
}

type subscribersHub map[int]map[uuid.UUID]*client
