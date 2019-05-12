package state

import (
	"fmt"
)

// OpCode for VM
type OpCode byte

// InstFn for MicorInst
type InstFn func(v *VM)

// Base Op Code
const (
	HALT OpCode = iota // 终止

	INPUT // 存放输入数据到内存
	PRINT // 输出内存中的数据到屏幕

	LOAD  // 加载内存数据到累加寄存器
	STORE // 存放累加寄存器到内存

	ADD // 累加寄存器数据加上内存数据
	SUB // 累加寄存器数据减去内存数据

	JMP // 转移
	JZ  // 累加寄存器数据为零则转移
	JNZ // 累加寄存器数据不为零则转移
	JS  // 累加寄存器数据为负则转移
	JNS // 累加寄存器数据为正则转移

	// MUL // 乘法
	// DIV // 除法
	// MOD // 求余
)

func (v *VM) initOp() *VM {
	v.opName = []string{
		"HALT",  // 终止
		"INPUT", // 存放输入数据到内存
		"PRINT", // 输出内存中的数据到屏幕
		"LOAD",  // 加载内存数据到累加寄存器
		"STORE", // 存放累加寄存器到内存
		"ADD",   // 累加寄存器数据加上内存数据
		"SUB",   // 累加寄存器数据减去内存数据
		"JMP",   // 转移
		"JZ",    // 累加寄存器数据为零tring移
		"JNZ",   // 累加寄存器数据不为tring转移
		"JS",    // 累加寄存器数据为负tring移
		"JNS",   // 累加寄存器数据为正则转移
	}
	v.opFn = make(map[byte]InstFn)
	return v
}

func (v *VM) run() {
	defer func() {
		e := recover()
		if e != nil {
			panic(fmt.Errorf("VM Error: %v", e))
		}
	}()
	// pretty.Println(v.mem)
	for {
		v.ir = OpCode(v.mem[v.ip])
		// fmt.Println("IR", v.ir)
		v.ip++
		switch v.ir {
		case HALT:
		case INPUT, PRINT,
			LOAD, STORE,
			ADD, SUB,
			JMP, JZ, JNZ, JS, JNS:
			v.op = int(v.loadMem32(v.ip))
			v.ip += 4
			// fmt.Println("OP", v.op)
		}
		switch v.ir {
		case HALT:
			v.Stop()
		case INPUT:
			fmt.Scanf("%d", &v.mem[v.op])
		case PRINT:
			fmt.Printf("%d\n", v.mem[v.op])
		case LOAD:
			v.acc = int(v.loadMem32(v.op))
			// fmt.Println("LOAD ACC=", v.acc&0xffffffff)
		case STORE:
			// fmt.Println("STORE ACC=", v.acc&0xffffffff)
			v.storeMem32(v.op)
		case ADD:
			// fmt.Println("ADD ACC=", v.acc&0xffffffff)
			v.acc += int(v.loadMem32(v.op))
			// fmt.Println("ACC=", v.acc&0xffffffff)
			v.acc &= 0xffffffff
		case SUB:
			// fmt.Println("SUB ACC=", v.acc&0xffffffff)
			v.acc -= int(v.loadMem32(v.op))
			fmt.Println("ACC=", v.acc&0xffffffff)
			v.acc &= 0xffffffff
		case JMP:
			v.ip = v.op
		case JZ:
			if v.acc == 0 {
				v.ip = v.op
			}
		case JNZ:
			if v.acc != 0 {
				v.ip = v.op
			}
		case JS:
			//if v.acc < 0 {
			//	v.ip = v.op
			//}
			if v.acc < 0 || v.acc > 0x7fffffff {
				v.ip = v.op
			}
		case JNS:
			//if v.acc > 0 {
			//	v.ip = v.op
			//}
			if v.acc != 0 && v.acc != 0x80000000 &&
				!(v.acc < 0 || v.acc > 0x7fffffff) {
				v.ip = v.op
			}
		}

		v.mutex.RLock()
		if !v.running {
			v.mutex.RUnlock()
			break
		}
		v.mutex.RUnlock()
	}
}

func (v *VM) loadMem32(addr int) int32 {
	return int32(v.mem[addr])<<24 +
		int32(v.mem[addr+1])<<16 +
		int32(v.mem[addr+2])<<8 +
		int32(v.mem[addr+3])
}

func (v *VM) storeMem32(addr int) {
	v.mem[addr+0] = byte((v.acc & 0xff000000) >> 24)
	v.mem[addr+1] = byte((v.acc & 0x00ff0000) >> 16)
	v.mem[addr+2] = byte((v.acc & 0x0000ff00) >> 8)
	v.mem[addr+3] = byte((v.acc & 0x000000ff) >> 0)
	// fmt.Println(v.acc, v.acc&0xff000000)
	// fmt.Println(v.mem[addr+0], v.mem[addr+1], v.mem[addr+2], v.mem[addr+3])
}
