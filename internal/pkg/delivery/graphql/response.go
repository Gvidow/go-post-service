package graphql

import "github.com/gvidow/go-post-service/internal/pkg/errors"

var _unknowResponseError = &responseError{"internal server error"}

type responseError struct {
	Message string
}

func (r *responseError) Error() string {
	return r.Message
}

func MakeResponseError(err error) error {
	var typeErr errors.TypeError

	if ok := errors.As(err, &typeErr); !ok {
		return _unknowResponseError
	}

	switch typeErr.Type() {
	case errors.CommentNotFound:
		return &responseError{"the comment was not found"}
	case errors.PostNotFound:
		return &responseError{"the post was not found"}
	case errors.NotPermission:
		return &responseError{"there are no rights to perform the action"}
	case errors.CommentsAreProhibited:
		return &responseError{"it is forbidden to leave comments under the post"}
	default:
		return _unknowResponseError
	}
}
