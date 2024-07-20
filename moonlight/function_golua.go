package moonlight

import (
	rt "github.com/arnodel/golua/runtime"
)

type GoFunctionFunc = rt.GoFunctionFunc

func (mlr *Runtime) CheckNArgs(c *GoCont, num int) error {
	return c.cont.CheckNArgs(num)
}

func (mlr *Runtime) Check1Arg(c *GoCont) error {
	return c.cont.CheckNArgs(1)
}

func (mlr *Runtime) StringArg(c *GoCont, num int) (string, error) {
	return c.cont.StringArg(num)
}

func (mlr *Runtime) ClosureArg(c *GoCont, num int) (*Closure, error) {
	return c.cont.ClosureArg(num)
}

func (mlr *Runtime) GoFunction(fun GoToLuaFunc) rt.GoFunctionFunc {
	return func(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
		gocont := GoCont{
			cont: c,
			thread: t,
		}
		return fun(mlr, &gocont)
	}
}

func (mlr *Runtime) Call1(val Value, args ...Value) (Value, error) {
	return rt.Call1(mlr.rt.MainThread(), val, args...)
}
