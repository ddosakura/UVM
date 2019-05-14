package uvm

import (
	"fmt"
)

// DefaultErrHandler for UVM
func DefaultErrHandler() {
	e := recover()
	if e != nil {
		panic(fmt.Errorf("VM Error: %v", e))
	}
}
