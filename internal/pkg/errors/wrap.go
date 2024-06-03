package errors

import (
	"errors"
	"fmt"
)

func Wrap(err error, wrapped string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", wrapped, err)
	}
	return nil
}

func Wrapf(err error, format string, a ...any) error {
	if err != nil {
		return fmt.Errorf("%s: %w", fmt.Sprintf(format, a...), err)
	}
	return nil
}

func WrapFail(err error, wrapped string) error {
	return Wrapf(err, "couldn't %s", wrapped)
}

func WrapFailf(err error, format string, a ...any) error {
	if err != nil {
		return fmt.Errorf("couldn't %s: %w", fmt.Sprintf(format, a...), err)
	}
	return nil
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
