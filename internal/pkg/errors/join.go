package errors

import (
	"strings"
)

type joinedError struct {
	errs []error
}

func (j *joinedError) Error() string {
	msgs := make([]string, len(j.errs))
	for ind, err := range j.errs {
		msgs[ind] = err.Error()
	}
	return strings.Join(msgs, "; ")
}

func Join(errs ...error) error {
	errsReal := make([]error, 0)

	for _, err := range errs {
		if err != nil {
			errsReal = append(errsReal, err)
		}
	}

	if len(errsReal) == 0 {
		return nil
	}

	return &joinedError{errsReal}
}
