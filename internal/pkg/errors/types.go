package errors

type Type uint8

const (
	_ Type = iota

	TypePostNotFound
	TypeCommentNotFound

	TypeInvalidComment
	TypeCommentsAreProhibited

	TypeNotPermission

	TypeUnknow
)
