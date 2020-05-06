package syntax

import (
	"errors"
	"fmt"
)

var (
	functionRedeclireError = errors.New("function is redecliress")
)

// An ErrorHandler is called for each error encountered reading a .go file.
type ErrorHandler func(err error)

// Error describes a syntax error. Error implements the error interface.
type Error struct {
	Pos Pos
	Msg string
}

func (err Error) Error() string {
	return fmt.Sprintf("%s: %s", err.Pos, err.Msg)
}
