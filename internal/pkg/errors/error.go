package errors

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

func NewWithType(err error, t Type) error {
	return &typingError{err, t}
}
