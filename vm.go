package uvm

import (
	"sync"
)

// VM interface
type VM interface {
	Start(ip uint)
	Stop()
}

// UVM core
type UVM struct {
	mutex      *sync.RWMutex
	running    bool
	ErrHandler func()
	IP         uint // 指令计数器

	Reg8  []byte   // 8 位寄存器 (byte, alias for uint8, [0~255]), 需要至少一个作为指令寄存器
	Reg16 []uint16 // 16 位寄存器
	Reg32 []uint32 // 32 位寄存器
	Reg64 []uint64 // 64 位寄存器

	OpName []string        // 指令名
	OpCode map[string]byte // 指令号
	Ops    []*Op           // 指令集

	Mem     []byte // 内存
	InitMem func() // 内存初始化
}

// NewVM for UVM
func NewVM(fn func(*UVM)) VM {
	v := &UVM{
		mutex:      new(sync.RWMutex),
		running:    false,
		ErrHandler: DefaultErrHandler,
	}
	fn(v)
	return v
}

// InitReg for UVM
func (v *UVM) InitReg(r8, r16, r32, r64 int) {
	if r8 > 0 {
		v.Reg8 = make([]byte, r8)
	}
	if r16 > 0 {
		v.Reg16 = make([]uint16, r16)
	}
	if r32 > 0 {
		v.Reg32 = make([]uint32, r32)
	}
	if r64 > 0 {
		v.Reg64 = make([]uint64, r64)
	}
}

// InitInst for UVM
func (v *UVM) InitInst(num int) {
	v.OpName = make([]string, 0, num)
	v.OpCode = make(map[string]byte, num)
	v.Ops = make([]*Op, 0, num)
}

// Inst Builder for UVM
func (v *UVM) Inst(name string) (CODE byte, OP *Op) {
	l := len(v.OpName)
	if l > 255 {
		panic(ErrTooManyOps)
	}
	CODE = byte(l)
	OP = &Op{}

	v.OpName = append(v.OpName, name)
	v.OpCode[name] = CODE
	v.Ops = append(v.Ops, OP)
	return
}

// Start UVM
func (v *UVM) Start(ip uint) {
	v.mutex.Lock()
	if v.running {
		v.mutex.Unlock()
		return
	}
	v.running = true
	v.InitMem()
	v.mutex.Unlock()

	v.IP = ip

	go v.run()
}

// Stop UVM
func (v *UVM) Stop() {
	defer v.mutex.Unlock()
	v.mutex.Lock()
	if v.running {
		v.running = false
	}
}
