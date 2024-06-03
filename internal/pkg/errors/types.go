package errors

type Type uint8

const (
	_ Type = iota

	PostNotFound
	CommentNotFound

	Unknow
)
