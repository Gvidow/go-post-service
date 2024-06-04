package usecase

import (
	"github.com/gvidow/go-post-service/internal/pkg/usecase/comment"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/notify"
	"github.com/gvidow/go-post-service/internal/pkg/usecase/post"
)

//go:generate mockgen -destination=./mock/repository_mock.go -package=mock -source=interfaces.go

type Repository interface {
	comment.Repository
	post.Repository
}

type Notifier interface {
	notify.Notifier
}
