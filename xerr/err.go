package xerr

import (
	"github.com/pkg/errors"
)

// stackTracer is an interface that provides access to the error's stack trace
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// New creates a new error with the given message
func New(msg string) error {
	return errors.New(msg)
}

// Errorf creates a new error with the given formatted message
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// WithMessage adds a message to an existing error
func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

// WithMessagef adds a formatted message to an existing error
func WithMessagef(err error, message string, args ...interface{}) error {
	return errors.WithMessagef(err, message, args...)
}

// Wrap wraps an error with a message
func Wrap(err error, msg string) error {
	return errors.Wrap(err, msg)
}

// Wrapf wraps an error with a formatted message
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// Cause returns the underlying cause of the error
func Cause(err error) error {
	return errors.Cause(err)
}

// Is checks if the error matches the target error
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As checks if the error can be converted to the target type
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// WithStack adds a stack trace to an error if it doesn't already have one
func WithStack(err error) error {
	_, ok := err.(stackTracer)
	if ok {
		return err
	}

	return errors.WithStack(err)
}
