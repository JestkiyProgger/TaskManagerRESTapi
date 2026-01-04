package task

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidTitle = errors.New("invalid title")
)
