package errors

import "errors"

type TypeError interface {
	error

	Type() Type
}

type typingError struct {
	error

	t Type
}

func (err typingError) Type() Type {
	return err.t
}

func (err typingError) Unwrap() error {
	return err.error
}

func WithType(err error, t Type) error {
	return &typingError{err, t}
}

func New(text string) error {
	return errors.New(text)
}
