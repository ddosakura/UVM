package uvm

import "fmt"

// Op for UVM
type Op struct {
	Num    uint
	Seek   func([]byte)
	Handle func()
}

func (v *UVM) run() {
	defer v.ErrHandler()
	for {
		if v.IP >= uint(len(v.Mem)) {
			panic(ErrOutOfMem)
		}
		v.Reg8[0] = v.Mem[v.IP]
		v.IP++

		if v.Reg8[0] >= byte(len(v.Ops)) {
			panic(ErrUnknowOp)
		}
		op := v.Ops[v.Reg8[0]]
		if op.Seek != nil && op.Num > 0 {
			if v.IP+op.Num > uint(len(v.Mem)) {
				fmt.Println(v.IP, op.Num, uint(len(v.Mem)))
				panic(ErrOutOfMem)
			}
			op.Seek(v.Mem[v.IP : v.IP+op.Num])
			v.IP += op.Num
		}

		if op.Handle != nil {
			op.Handle()
		}

		v.mutex.RLock()
		if !v.running {
			v.mutex.RUnlock()
			break
		}
		v.mutex.RUnlock()
	}
}
