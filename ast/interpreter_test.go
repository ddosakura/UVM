package ast

import (
	"testing"
)

func TestIf(t *testing.T) {
	ast := NewAST(QuickInnerInterrupt([]int{2, 3}))
	interpreter := NewInterpreter(ast)
	typeDefine(ast, 1, "int")

	ast.Interrupt(0, "Type Error", func(e interface{}) interface{} {
		t.Log("Type Error:", e)
		return nil
	})
	ast.Interrupt(1, "Output", func(e interface{}) interface{} {
		t.Log("Output:", e)
		return nil
	})
	ast.Interrupt(2, "OverRange", func(e interface{}) interface{} {
		t.Log("OverRange:", e)
		return nil
	})
	ast.Interrupt(3, "InterruptNotFound", func(e interface{}) interface{} {
		t.Log("InterruptNotFound:", e)
		return nil
	})

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
			ast.Int(0, "The Left Value of '+' is not int")
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
	a.Data(0)         // false
	a = root.Child(2) // then print
	b := a.Child(1)   // +
	b.Child(-1).Data(2)
	b.Child(-1).Data(3)
	a = root.Child(2) // false print
	b = a.Child(1)    // +
	b.Child(-1).Data(2)
	b.Child(-1).Data(-3)

	interpreter.Run()
}

func TestLoop(t *testing.T) {

}
