package errors

import (
	"errors"
	"testing"
)

func makeError() error{
	return New("Error!")
}

func extendError() error{
	return Extend(errors.New("Le"))
}

func TestTrace(t *testing.T) {
	t.Log(makeError())
	t.Log(extendError())
}
