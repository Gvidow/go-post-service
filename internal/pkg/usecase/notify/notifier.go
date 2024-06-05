package notify

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/gvidow/go-post-service/internal/entity"
	"github.com/gvidow/go-post-service/pkg/logger"
)

type notifier struct {
	log *logger.Logger
	hub subscribersHub
	rw  sync.RWMutex
}

func NewNotifier(log *logger.Logger) *notifier {
	return &notifier{
		log: log,
		hub: make(subscribersHub),
	}
}

func (n *notifier) PublishComment(_ context.Context, comment *entity.Comment, postId int) {
	n.rw.RLock()
	defer n.rw.RUnlock()

	subscribers := n.hub[postId]

	for _, client := range subscribers {
		select {
		case client.transport <- entity.NotifyComment{Comment: comment}:
			n.log.Sugar().Infof("send notify comment(id=%d) to subscriber client(id=%s)", comment.ID, client.id)
		case <-client.ctx.Done():
			n.log.Info("the context was canceled when sending")
		}
	}
}

func (n *notifier) RegistryChanNotifier(ctx context.Context, ch chan<- entity.NotifyComment, postId int) error {
	client := &client{
		id:        uuid.New(),
		transport: ch,
		postId:    postId,
		ctx:       ctx,
	}

	n.log.Sugar().Infof("registry notifier client(id = %s) for post(id = %d)", client.id, postId)

	n.subscribeClient(client)
	n.registryUnsubscribe(client)

	return nil
}

func (n *notifier) subscribeClient(cl *client) {
	n.rw.Lock()
	defer n.rw.Unlock()

	hub, ok := n.hub[cl.postId]
	if !ok {
		n.hub[cl.postId] = map[uuid.UUID]*client{cl.id: cl}
	} else {
		hub[cl.id] = cl
	}
}

func (n *notifier) registryUnsubscribe(cl *client) {
	go func() {
		defer close(cl.transport)
		<-cl.ctx.Done()

		n.rw.Lock()
		delete(n.hub[cl.postId], cl.id)
		n.rw.Unlock()
	}()
}
