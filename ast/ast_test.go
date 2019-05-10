package ast

import (
	"testing"

	"github.com/kr/pretty"
)

var ast *AST

func typeDefine(ast *AST, t int, name string) {
	ast.DataType[t] = name
	ast.Inst(-t, name, ParamSize(0), false, nil)
}

func init() {
	ast = NewAST(QuickInnerInterrupt([]int{2, 3}))
	typeDefine(ast, 1, "int")

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
		pretty.Println(ds)
		for _, d := range ds {
			ast.IntV(1, d())
		}
		return nil
	})
}

func reInitAST(ast *AST, t *testing.T) {
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
}

func TestAST(t *testing.T) {
	reInitAST(ast, t)
	root := ast.Root(2)

	a := root.Child(1)
	a.Child(-1).Data(2)
	a.Child(-1).Data(3)

	pretty.Println(ast)

	interpreter := NewInterpreter(ast)
	interpreter.Run()
}

func TestOpt(t *testing.T) {
	reInitAST(ast, t)
	root := ast.Root(2)

	a := root.Child(1)
	a.Child(-1).Data(2)
	a.Child(-1).Data(3)

	pretty.Println(ast)

	interpreter := NewInterpreter(ast)
	ast.CheckType = func(v interface{}) int {
		return -1
	}
	err := interpreter.Opt()
	if err != nil {
		t.Fatal(err)
	}

	interpreter.Run()
	pretty.Println(ast)

	interpreter.Run()
	pretty.Println(ast)
}

func TestStrSaveToInt(t *testing.T) {
	reInitAST(ast, t)
	root := ast.Root(2)

	a := root.Child(1)
	a.Child(-1).Data(2)
	a.Child(-1).Data("H")

	//pretty.Println(ast)

	interpreter := NewInterpreter(ast)

	interpreter.Run()
}
