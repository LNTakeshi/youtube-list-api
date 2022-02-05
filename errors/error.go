package errors

import "github.com/morikuni/failure"

const (
	ErrBadRequest  failure.StringCode = "-400"
	ErrBadUrl      failure.StringCode = "-1000"
	ErrFetchUrl    failure.StringCode = "-1001"
	ErrAdd         failure.StringCode = "-1002"
	ErrInvalidTime failure.StringCode = "-1003"
)
