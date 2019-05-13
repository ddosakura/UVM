package state

import (
	"testing"
	"time"
)

var (
	MUL OpCode
)

type sd struct {
}

func (s *sd) Data(f, t int) []byte {
	bs := make([]byte, 1024)
	copy(bs, []byte{
		byte(PRINT), 0, 0, 0x3, 0xf, // 000-0
		byte(JMP), 0, 0, 0, 30, // 000-5
		0, 1, 0, 0, // 001-0
		0, 0, 0, 1, // 001-4
		0, 0, 0, 2, // 001-8
		0, 0, 0, 4, // 002-2
		0, 0, 0, 0, // 002-6

		//byte(PRINT), 0, 0, 0, 10, // 003-0
		//byte(PRINT), 0, 0, 0, 11,
		//byte(PRINT), 0, 0, 0, 12,
		//byte(PRINT), 0, 0, 0, 13,
		//byte(PRINT), 0, 0, 0, 14,
		//byte(PRINT), 0, 0, 0, 15,
		//byte(PRINT), 0, 0, 0, 16,
		//byte(PRINT), 0, 0, 0, 17,
		//byte(PRINT), 0, 0, 0, 18,
		//byte(PRINT), 0, 0, 0, 19,
		//byte(PRINT), 0, 0, 0, 20,
		//byte(PRINT), 0, 0, 0, 21,
		//byte(PRINT), 0, 0, 0, 22,
		//byte(PRINT), 0, 0, 0, 23,
		//byte(PRINT), 0, 0, 0, 24,
		//byte(PRINT), 0, 0, 0, 25,
		//byte(PRINT), 0, 0, 0, 26,
		//byte(PRINT), 0, 0, 0, 27,
		//byte(PRINT), 0, 0, 0, 28,
		//byte(PRINT), 0, 0, 0, 29,

		byte(LOAD), 0, 0, 0, 14,
		byte(ADD), 0, 0, 0, 18,
		byte(SUB), 0, 0, 0, 22,
		byte(STORE), 0, 0, 0, 26,

		byte(PRINT), 0, 0, 0, 26,
		byte(PRINT), 0, 0, 0, 27,
		byte(PRINT), 0, 0, 0, 28,
		byte(PRINT), 0, 0, 0, 29,

		byte(LOAD), 0, 0, 0, 18,
		byte(MUL), 0, 0, 0, 22,
		byte(ADD), 0, 0, 0, 14,
		byte(STORE), 0, 0, 0, 26,

		byte(PRINT), 0, 0, 0, 26,
		byte(PRINT), 0, 0, 0, 27,
		byte(PRINT), 0, 0, 0, 28,
		byte(PRINT), 0, 0, 0, 29,

		// byte(LOAD), 0, 0, 0, 10,
		// byte(SUB), 0, 0, 0, 14,
		// byte(STORE), 0, 0, 0, 10,
		// byte(JNZ), 0, 0, 0, 30,
	})

	// byte(JMP), 0, 0, 0, 30, // 000-5
	bs[1009] = byte(JMP)
	bs[1013] = 30
	return bs
}

func TestVM(t *testing.T) {
	v := NewVM(1024) // 1k mem
	MUL = v.MicroInstruction("MUL", func(v *VM, acc, op int) int {
		return acc * v.LoadMem32(op)
	})
	v.Start(&sd{}, 0, 0)
	//time.Sleep(time.Second * 1)
	time.Sleep(time.Second * 3)
	v.Stop()
}
