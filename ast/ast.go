package ast

import (
	"strconv"
)

// DelayData for Inst
type DelayData func() Data

// Inst Action
type Inst func(...DelayData) Data

// PassInst Do Nothing
var PassInst = func(ds ...DelayData) Data {
	for _, d := range ds {
		if d != nil {
			d()
		}
	}
	return nil
}

// Data for AST
type Data interface{}

// InterruptHander for AST
type InterruptHander func(interface{}) interface{}

// ParamLimit for Inst
type ParamLimit struct {
	Min int
	Max int
}

// ParamSize Build ParamLimit a..a
func ParamSize(size int) *ParamLimit {
	return &ParamLimit{
		Min: size,
		Max: size,
	}
}

// ParamRange Build ParamLimit a..b
func ParamRange(Min int, Max int) *ParamLimit {
	return &ParamLimit{
		Min,
		Max,
	}
}

// AST info
type AST struct {
	DataType       map[int]string
	CheckType      func(interface{}) int
	interruptType  map[int]string
	interrupt      map[int]InterruptHander
	innerInterrupt *InnerInterrupt

	instName  map[int]string
	instParam map[int]*ParamLimit
	instOpt   map[int]bool
	instSet   map[int]Inst

	root *Node
}

// InnerInterrupt for Interpreter
type InnerInterrupt struct {
	OverRange         int
	InterruptNotFound int
}

// QuickInnerInterrupt Builder
func QuickInnerInterrupt(oii []int) *InnerInterrupt {
	return &InnerInterrupt{
		OverRange:         oii[0],
		InterruptNotFound: oii[1],
	}
}

// CheckInnerInterrupt not nil
func (a *AST) CheckInnerInterrupt() bool {
	return !(a.interrupt[a.innerInterrupt.OverRange] == nil ||
		a.interrupt[a.innerInterrupt.InterruptNotFound] == nil)
}

// NewAST for Interpreter
func NewAST(innerInterrupt *InnerInterrupt) *AST {
	return &AST{
		DataType:       make(map[int]string),
		interruptType:  make(map[int]string),
		interrupt:      make(map[int]InterruptHander),
		innerInterrupt: innerInterrupt,

		instName:  make(map[int]string),
		instParam: make(map[int]*ParamLimit),
		instOpt:   make(map[int]bool),
		instSet:   make(map[int]Inst),
	}
}

// Inst Builder
func (a *AST) Inst(op int, name string, param *ParamLimit, opt bool, fn Inst) {
	a.instName[op] = name
	a.instParam[op] = param
	a.instOpt[op] = opt
	a.instSet[op] = fn
}

// Interrupt Builder
func (a *AST) Interrupt(code int, name string, fn InterruptHander) {
	a.interruptType[code] = name
	a.interrupt[code] = fn
}

// Interrupt for AST
type Interrupt struct {
	code int
	data interface{}
}

// Int Call Interrupt
func (a *AST) Int(code int, data interface{}) {
	if a.interrupt[code] == nil {
		panic(Interrupt{
			code: a.innerInterrupt.InterruptNotFound,
			data: "INT " + strconv.Itoa(code),
		})
	}
	panic(Interrupt{
		code,
		data,
	})
}

// InterruptHande for Interpreter
func (a *AST) InterruptHande(fn InterruptHander) {
	e := recover()
	if e == nil {
		return
	}
	i, ok := e.(Interrupt)
	if ok {
		defer a.InterruptHande(nil)
		a.interrupt[i.code](i.data)
	} else if fn != nil {
		fn(e)
	} else {
		panic(e)
	}
}

// IntV Call Interrupt will return data
func (a *AST) IntV(code int, data interface{}) (v interface{}) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		i, ok := e.(Interrupt)
		if ok {
			if a.interrupt[code] == nil {
				panic(Interrupt{
					code: a.innerInterrupt.InterruptNotFound,
					data: "INT " + strconv.Itoa(code),
				})
			}
			v = a.IntV(i.code, i.data)
		} else {
			panic(e)
		}
	}()
	return a.interrupt[code](data)
}
