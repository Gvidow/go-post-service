package graphql

import "github.com/gvidow/go-post-service/internal/pkg/errors"

var _unknowResponseError = &responseError{"internal server error"}

type responseError struct {
	Message string
}

func (r *responseError) Error() string {
	return r.Message
}

func MakeResponseError(err error) *responseError {
	var typeErr errors.TypeError

	if ok := errors.As(err, &typeErr); !ok {
		return _unknowResponseError
	}

	switch typeErr.Type() {
	case errors.TypeCommentNotFound:
		return &responseError{"the comment was not found"}
	case errors.TypePostNotFound:
		return &responseError{"the post was not found"}
	case errors.TypeNotPermission:
		return &responseError{"there are no rights to perform the action"}
	case errors.TypeCommentsAreProhibited:
		return &responseError{"it is forbidden to leave comments under the post"}
	case errors.TypeInvalidComment:
		return &responseError{"the length of the comment does not fit into 2000 characters"}
	default:
		return _unknowResponseError
	}
}
