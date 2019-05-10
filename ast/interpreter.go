package ast

import (
	"strconv"
)

// Interpreter for AST
type Interpreter struct {
	ast *AST
	opt bool // auto change `2+3` -> `5`

	IP *Node
}

// NewInterpreter for AST
func NewInterpreter(ast *AST) *Interpreter {
	if ast == nil {
		return nil
	}
	return &Interpreter{
		ast: ast,
		opt: false,
	}
}

// Opt Open
func (i *Interpreter) Opt() error {
	if i.ast.CheckType == nil {
		return ErrOptNeedCheckType
	}
	i.opt = true
	return nil
}

// Run Interpreter
func (i *Interpreter) Run() Data {
	if !i.ast.CheckInnerInterrupt() {
		panic(ErrInnerInterruptCheckFail)
	}
	i.IP = i.ast.root
	if !i.IP.Reset {
		return i.run()
	}
	defer i.ast.InterruptHande(nil)
	return i.run()
}

func (i *Interpreter) run() Data {
	ip := i.IP
	defer func() {
		i.IP = ip
	}()

	// name := i.ast.instName[ip.op]
	limit := i.ast.instParam[ip.op]
	opt := i.ast.instOpt[ip.op]
	fn := i.ast.instSet[ip.op]
	if fn == nil {
		return ip.data
	}
	cs := ip.children
	l := len(cs)
	if limit != nil &&
		((limit.Min > -1 && l < limit.Min) ||
			(limit.Max > -1 && l > limit.Max)) {
		panic(Interrupt{
			code: i.ast.innerInterrupt.OverRange,
			data: "len=" + strconv.Itoa(l) + " over (" + strconv.Itoa(limit.Min) + ".." + strconv.Itoa(limit.Max) + ")",
		})
	}
	var d Data
	if l > 0 {
		ds := make([]DelayData, 0, l)
		//for i.IP = cs[0]; i.IP != nil; i.IP = i.IP.Next() {
		//	if i.IP.Reset {
		//		ds = append(ds, func() Data {
		//			defer i.ast.InterruptHande(nil)
		//			return i.run()
		//		})
		//	} else {
		//		ds = append(ds, func() Data {
		//			return i.run()
		//		})
		//	}
		//}
		for _, IP := range cs {
			ip := IP
			if i.IP.Reset {
				ds = append(ds, func() Data {
					i.IP = ip
					defer i.ast.InterruptHande(nil)
					return i.run()
				})
			} else {
				ds = append(ds, func() Data {
					i.IP = ip
					return i.run()
				})
			}
		}
		d = fn(ds...)
	} else {
		d = fn()
	}
	if opt && i.opt {
		ip.op = i.ast.CheckType(d)
		ip.children = nil
		ip.next = nil
		ip.data = d
		ip.Reset = false
	}
	return d
}
