package uvm

import (
	"fmt"
	"testing"
	"time"
)

var (
	NOP,
	PUSH, POP,
	ADD,
	HALT byte
	DEBUG byte
)

const (
	MemSize = 1024
)

var uvm = NewVM(func(v *UVM) {
	v.InitReg(1, 0, 0, 1)
	v.InitInst(3)
	var op *Op

	NOP, op = v.Inst("NOP")
	op.Handle = func() {
		time.Sleep(1)
	}

	PUSH, op = v.Inst("PUSH")
	op.Num = 8
	op.Seek = func(bs []byte) {
		v.Reg64[0] -= 8
		if v.Reg64[0] < 0 {
			panic(ErrOutOfMem)
		}
		top := v.Reg64[0]
		for i := uint64(0); i < 8; i++ {
			v.Mem[top+i] = bs[i]
		}
	}
	// op.Handle = func() {}

	POP, op = v.Inst("POP")
	op.Handle = func() {
		v.Reg64[0] += 8
		if v.Reg64[0]+8 >= MemSize {
			panic(ErrOutOfMem)
		}
	}

	ADD, op = v.Inst("ADD")
	op.Handle = func() {
		top := v.Reg64[0]
		if top+16 > MemSize {
			panic(ErrOutOfMem)
		}
		ans := bsUint64(v.Mem[top:top+8]) + bsUint64(v.Mem[top+8:top+16])
		bs := uint64Bs(ans)
		v.Reg64[0] -= 8
		top = v.Reg64[0]
		for i := uint64(0); i < 8; i++ {
			v.Mem[top+i] = bs[i]
		}
	}

	DEBUG, op = v.Inst("DEBUG")
	op.Handle = func() {
		top := v.Reg64[0]
		fmt.Println(bsUint64(v.Mem[top : top+8]))
	}

	HALT, op = v.Inst("HALT")
	op.Handle = func() {
		v.Stop()
	}

	v.Mem = make([]byte, MemSize) // 1k mem
	v.InitMem = func() {
		v.Reg64[0] = MemSize - 1

		bs := []byte{}
		bs = append(bs, PUSH)
		bs = append(bs, uint64Bs(100)...)
		bs = append(bs, PUSH)
		bs = append(bs, uint64Bs(260)...)
		bs = append(bs, ADD)
		bs = append(bs, DEBUG)
		bs = append(bs, PUSH)
		bs = append(bs, uint64Bs(2445)...)
		bs = append(bs, ADD)
		bs = append(bs, DEBUG)
		bs = append(bs, HALT)

		copy(v.Mem, bs)
	}
})

func TestVM(t *testing.T) {
	uvm.Start(0)
	time.Sleep(time.Second * 3)
}

func bsUint64(bs []byte) uint64 {
	return uint64(bs[0])<<56 +
		uint64(bs[1])<<48 +
		uint64(bs[2])<<40 +
		uint64(bs[3])<<32 +
		uint64(bs[4])<<24 +
		uint64(bs[5])<<16 +
		uint64(bs[6])<<8 +
		uint64(bs[7])
}

func uint64Bs(n uint64) []byte {
	bs := make([]byte, 8)
	bs[0] = byte((n & uint64(0xff00000000000000)) >> 56)
	bs[1] = byte((n & uint64(0x00ff000000000000)) >> 48)
	bs[2] = byte((n & uint64(0x0000ff0000000000)) >> 40)
	bs[3] = byte((n & uint64(0x000000ff00000000)) >> 32)
	bs[4] = byte((n & uint64(0x00000000ff000000)) >> 24)
	bs[5] = byte((n & uint64(0x0000000000ff0000)) >> 16)
	bs[6] = byte((n & uint64(0x000000000000ff00)) >> 8)
	bs[7] = byte((n & uint64(0x00000000000000ff)) >> 0)
	return bs
}
