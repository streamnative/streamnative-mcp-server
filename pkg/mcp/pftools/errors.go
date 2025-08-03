package pftools

import "errors"

var (
	ErrFunctionNotFound = errors.New("function not found")
	ErrNotOurMessage    = errors.New("not our message")
)
