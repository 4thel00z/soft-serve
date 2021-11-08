package types

import "errors"

var (
	ErrDiffTooLong      = errors.New("diff is too long")
	ErrDiffFilesTooLong = errors.New("diff files are too long")
)

type ErrMsg struct {
	Error error
}
