package ast

import (
	"testing"
)

func TestIf(t *testing.T) {
	ast := NewAST(QuickInnerInterrupt([]int{2, 3}))
	interpreter := NewInterpreter(ast)
	typeDefine(ast, 1, "int")

	reInitAST(ast, t)

	ast.Inst(0, "pass", ParamSize(0), false, PassInst)
	ast.Inst(1, "+", ParamSize(2), true, func(ds ...DelayData) Data {
		a, ok := ds[0]().(int)
		if !ok {
			ast.Int(0, "The Left Value of '+' is not int")
		}
		b, ok := ds[1]().(int)
		if !ok {
			ast.Int(0, "The Right Value of '+' is not int")
		}
		return a + b
	})
	ast.Inst(2, "print", nil, false, func(ds ...DelayData) Data {
		for _, d := range ds {
			ast.IntV(1, d())
		}
		return nil
	})
	ast.Inst(3, "block", nil, false, PassInst)
	ast.Inst(4, "if", ParamRange(2, 3), false, func(ds ...DelayData) Data {
		condition, ok := ds[0]().(int)
		if !ok {
			ast.Int(0, "The condition of 'if' is not int")
		}
		if condition == 0 {
			// false
			if len(ds) == 3 {
				ds[2]()
			}
		} else {
			// true
			ds[1]()
		}
		return nil
	})

	root := ast.Root(4) // If

	a := root.Child(-1)
	a.Data(1) // true
	//a.Data(0) // false
	A := root.Child(3) // then block
	a = A.Child(2)     // print
	b := a.Child(1)    // +
	b.Child(-1).Data(2)
	b.Child(-1).Data(3)
	a = A.Child(2) // print
	b = a.Child(1) // +
	b.Child(-1).Data(-2)
	b.Child(-1).Data(3)

	a = root.Child(2) // false print
	b = a.Child(1)    // +
	b.Child(-1).Data(2)
	b.Child(-1).Data(-3)

	interpreter.Run()
}

func TestLoop(t *testing.T) {
	scope := NewScope()
	ast := NewAST(QuickInnerInterrupt([]int{2, 3}))
	interpreter := NewInterpreter(ast)
	typeDefine(ast, 1, "int")
	typeDefine(ast, 2, "string")
	reInitAST(ast, t)

	ast.Inst(0, "print", nil, false, func(ds ...DelayData) Data {
		for _, d := range ds {
			ast.IntV(1, d())
		}
		return nil
	})
	ast.Inst(1, "save", ParamSize(2), false, func(ds ...DelayData) Data {
		k, ok := ds[0]().(string)
		if !ok {
			ast.Int(0, "The Key of 'save' is not string")
		}
		v, ok := ds[1]().(int)
		if !ok {
			ast.Int(0, "The Value of 'save' is not int")
		}
		scope.Put(k, v)
		return nil
	})
	ast.Inst(2, "load", ParamSize(1), false, func(ds ...DelayData) Data {
		k, ok := ds[0]().(string)
		if !ok {
			ast.Int(0, "The Key of 'load' is not string")
		}
		return scope.Get(k)
	})
	ast.Inst(3, "block", nil, false, func(ds ...DelayData) Data {
		for _, d := range ds {
			if d != nil {
				d()
			}
		}
		return nil
	})
	ast.Inst(4, "if", ParamRange(2, 3), false, func(ds ...DelayData) Data {
		condition, ok := ds[0]().(int)
		if !ok {
			ast.Int(0, "The condition of 'if' is not int")
		}
		if condition == 0 {
			// false
			if len(ds) == 3 {
				ds[2]()
			}
		} else {
			// true
			ds[1]()
		}
		return nil
	})
	//ast.Inst(5, "label", ParamSize(1), false, func(ds ...DelayData) Data {
	//	ip := interpreter.IP
	//	k, ok := ds[0]().(string)
	//	if !ok {
	//		ast.Int(0, "The Key of 'label' is not string")
	//	}
	//	scope.Put("lable "+k, ip)
	//	return nil
	//})
	//ast.Inst(6, "goto", ParamSize(1), false, func(ds ...DelayData) Data {
	//	k, ok := ds[0]().(string)
	//	if !ok {
	//		ast.Int(0, "The Key of 'goto' is not string")
	//	}
	//	ip := scope.Get("lable " + k).(*Node)
	//	interpreter.IP = ip
	//	t.Log("goto", ip)
	//	return nil
	//})
	ast.Inst(7, "inc", ParamSize(1), false, func(ds ...DelayData) Data {
		k, ok := ds[0]().(string)
		if !ok {
			ast.Int(0, "The Left Value of 'inc' is not string")
		}
		v, ok := scope.Get(k).(int)
		if !ok {
			ast.Int(0, "The Value in scope is not int")
		}
		scope.Put(k, v+1)
		return nil
	})
	ast.Inst(8, "!=", ParamSize(2), false, func(ds ...DelayData) Data {
		a, ok := ds[0]().(int)
		if !ok {
			ast.Int(0, "The Left Value of '!=' is not int")
		}
		b, ok := ds[1]().(int)
		if !ok {
			ast.Int(0, "The Right Value of '!=' is not int")
		}
		if a != b {
			//t.Log(a, "!=", b)
			return 1
		}
		//t.Log(a, "==", b)
		return 0
	})
	ast.Inst(9, "for", ParamSize(4), false, func(ds ...DelayData) Data {
		//ds[0]()
		//for ds[1]() == 1 {
		//	ds[3]()
		//	ds[2]()
		//}
		for ds[0](); ds[1]() == 1; ds[2]() {
			ds[3]()
		}
		return 0
	})

	root := ast.Root(3) // <block>

	//i := root.Child(1) // i = 0
	//i.Child(-2).Data("i")
	//i.Child(-1).Data(0)
	//root.Child(5).Child(-2).Data("loop")       // :loop
	//root.Child(0).Child(2).Child(-2).Data("i") // print(<load>i)
	//root.Child(7).Child(-2).Data("i")          // i++
	//IF := root.Child(4)                        // if
	//condition := IF.Child(8)                   // !=
	//condition.Child(2).Child(-2).Data("i")
	//condition.Child(-1).Data(5)
	//IF.Child(6).Child(-2).Data("loop") // then goto :loop

	FOR := root.Child(9)
	i := FOR.Child(1) // i = 0
	i.Child(-2).Data("i")
	i.Child(-1).Data(0)
	condition := FOR.Child(8) // i != 5
	condition.Child(2).Child(-2).Data("i")
	condition.Child(-1).Data(5)
	FOR.Child(7).Child(-2).Data("i")        // i++
	B := FOR.Child(3)                       // <block>
	B.Child(0).Child(2).Child(-2).Data("i") // print(<load>i)

	interpreter.Run()
}
