package state

import (
	"fmt"
	"sync"
)

// VM of state
type VM struct {
	mutex   *sync.RWMutex
	running bool

	memMutex *sync.RWMutex
	memSize  int

	mem []byte // 内存
	ip  int    // 指令计数器
	ir  OpCode // 指令寄存器
	op  int    // 操作码寄存器
	//addr int    // 内存地址寄存器
	acc int // 累加寄存器 accumulator

	opName []string
	//opArgs map[OpCode]bool
	opFn map[OpCode]InstFn
}

// NewVM of state
func NewVM(memSize int) *VM {
	return (&VM{
		mutex:   new(sync.RWMutex),
		running: false,
		memSize: memSize,
	}).initOp()
}

// MicroInstruction can custom op
func (v *VM) MicroInstruction(name string, fn InstFn) OpCode {
	code := OpCode(len(v.opName))
	v.opName = append(v.opName, name)
	//v.opArgs[code] = args
	v.opFn[code] = fn
	return code
}

// Load Mem by addr
func (v *VM) Load(addr, length int) []byte {
	data := make([]byte, 0, length)
	pos := addr + length
	if l := len(v.mem); pos > l {
		pos = l
	}
	data = append(data, v.mem[addr:pos]...)
	return data
}

// Store Mem by addr
func (v *VM) Store(addr int, data []byte) {
	pos := addr + len(data)
	if l := len(v.mem); pos > l {
		pos = l
	}
	// v.mem = append(v.mem[:addr], data[:pos-addr]..., v.mem[pos:]...)
	for p := addr; p < pos; p++ {
		v.mem[p] = data[p-addr]
	}
}

// Start VM
func (v *VM) Start(sd SD, f, t int) {
	v.mutex.Lock()
	if v.running {
		v.mutex.Unlock()
		return
	}
	v.running = true
	v.mem = make([]byte, v.memSize)
	data := sd.Data(f, t)
	l := len(data)
	if l > v.memSize {
		l = v.memSize
	}
	copy(v.mem, data[:l])
	v.memMutex = new(sync.RWMutex)
	v.mutex.Unlock()

	v.ip = 0
	v.ir = 0
	v.op = 0
	//v.addr = 0
	v.acc = 0

	fmt.Println("VM Start\n---------")
	go v.run()
}

// Stop VM
func (v *VM) Stop() {
	defer v.mutex.Unlock()
	v.mutex.Lock()
	if v.running {
		v.running = false
		fmt.Println("\n---------\nVM Stop")
		return
	}
	fmt.Println("VM Stop")
}
