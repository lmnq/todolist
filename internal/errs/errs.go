package errs

import "errors"

var (
	ErrEmptyTitle   = errors.New("title is required!")
	ErrTooLongTitle = errors.New("too long title")

	ErrEmptyDate   = errors.New("date is required!")
	ErrInvalidDate = errors.New("invalid date")
)
