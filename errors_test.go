package gserrors

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewf(t *testing.T) {
	s := "high level error message"
	err := Newf(errors.New("origin error message"), "error info: %s", s)
	fmt.Println(err)
}

func TestRequire(t *testing.T) {
	a := 1
	b := 2
	c := a + b
	Require(c == 3, "")
}
