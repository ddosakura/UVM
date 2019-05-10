package ast

import (
	"errors"
)

// errors
var (
	ErrOverChildrenCap         = errors.New("The size of children is over size")
	ErrOptNeedCheckType        = errors.New("Interpreter Opt Need Set CheckType in AST")
	ErrInnerInterruptCheckFail = errors.New("InnerInterrupt of AST Error in Check")
)
