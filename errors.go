package uvm

import (
	"errors"
)

// errors
var (
	ErrOutOfMem   = errors.New("Out of Mem")
	ErrTooManyOps = errors.New("Too many Ops")
	ErrUnknowOp   = errors.New("Unknow Op")
)
