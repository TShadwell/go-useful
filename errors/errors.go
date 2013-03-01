package errors

import (
	"errors"
	"runtime"
	"strconv"
)

const maxStackLen = 20

var _ error = Error{}

type Error struct {
	Stack []Frame
	Child error
}

type Frame struct {
	*runtime.Func
	Pc uintptr
}

func (e Error) Error() string {
	out := e.Child.Error() + " in:\n"
	stackLen := len(e.Stack)
	for i, ed := 0, stackLen; i < ed; i++ {
		this := e.Stack[i]
		file, line := this.FileLine(this.Pc)
		name := this.Name()
		out += file + ": " + strconv.Itoa(line) + " (" + name + ")"
		if i < ed-1 {
			out += "\n"
		}
	}
	if stackLen == maxStackLen {
		out += "\n-- stack limit reached --\n"
	}
	return out
}

/*
	Extends an error to include information about
	function invocations 'frames' long.

	Frames defaults to 1 (just callee).
*/
func Extend(e error) Error {

	stk := make([]uintptr, maxStackLen)
	num := runtime.Callers(3, stk)
	callStack := make([]Frame, num)
	for i, ed := 0, num; i < ed; i++ {
		callStack[i] = Frame{
			Func: runtime.FuncForPC(stk[i]),
			Pc:   stk[i],
		}
	}
	if v, ok := e.(Error); ok {
		return Error{
			Child: v.Child,
			Stack: callStack,
		}
	}
	return Error{
		Child: e,
		Stack: callStack,
	}
}

/*
	Creates a new error that returns error 's' as well as 
	function frame stack 'frames' long.

	Frames defaults to 1 (just callee).
*/
func New(s string) Error {
	return Extend(errors.New(s))
}