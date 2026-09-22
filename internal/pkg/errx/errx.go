package errx

import "errors"

// Sentinel domain errors. Callers wrap them with detail using
// fmt.Errorf("%w: ...", ...) and transports map them with errors.Is.
var (
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrTodoNotFound        = errors.New("todo not found")
	ErrTodoLimitReached    = errors.New("todo limit reached for this session")
	ErrSessionLimitReached = errors.New("too many active sessions, try again later")
	ErrUnauthenticated     = errors.New("unauthenticated")
)
