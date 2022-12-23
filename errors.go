// @file:  errors.go
// @author:caibo
// @email: caibo923@gmail.com
// @desc:  在标准错误信息基础上添加自定义错误信息、函数栈、文件名、行数，方便定位错误

package gserrors

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
)

var (
	ErrRequire  = errors.New("precondition error")
	ErrAssert   = errors.New("assert error")
	ErrEnsure   = errors.New("ensure error")
	ErrOverload = errors.New("service overload")
)

type GSError interface {
	error
	Stack() string
	Origin() error
	NewOrigin(error)
}

type errorHost struct {
	origin  error
	stack   string
	message string
}

func (err *errorHost) Error() string {
	if err.message != "" {
		if err.origin != nil {
			return fmt.Sprintf(
				"%s\nbacktrace:\n%sbacktrace error:\n%s",
				err.message,
				err.stack,
				err.origin.Error(),
			)
		}
		return fmt.Sprintf("%s\nbacktrace:\n%s", err.message, err.stack)
	}
	if err.origin != nil {
		return fmt.Sprintf(
			"%s\nbacktrace:\n%s",
			err.origin.Error(),
			err.stack,
		)
	}
	return fmt.Sprintf("<unknown error>\n%s", err.stack)
}

func (err *errorHost) Stack() string {
	return err.stack
}

func (err *errorHost) Origin() error {
	return err.origin
}

func (err *errorHost) NewOrigin(target error) {
	err.origin = target
}

func stack() []byte {
	var buff bytes.Buffer
	for skip := 2; ; skip++ {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		buff.WriteString(fmt.Sprintf("\tfile = %s, line = %d\n", file, line))
	}
	return buff.Bytes()
}

func Panic(err error) {
	panic(New(err))
}

func Panicf(err error, fmtstring string, args ...interface{}) {
	panic(Newf(err, fmtstring, args...))
}

func New(err error) GSError {
	return &errorHost{
		origin: err,
		stack:  string(stack()),
	}
}

func Newf(err error, fmtstring string, args ...interface{}) GSError {
	return &errorHost{
		origin:  err,
		stack:   string(stack()),
		message: fmt.Sprintf(fmtstring, args...),
	}
}

func Require(status bool, fmtstring string, args ...interface{}) {
	if !status {
		Panicf(ErrRequire, fmtstring, args...)
	}
}

func Assert(status bool, fmtstring string, args ...interface{}) {
	if !status {
		Panicf(ErrAssert, fmtstring, args...)
	}
}

func Ensure(condition func() bool, fmtstring string, args ...interface{}) {
	if !condition() {
		Panicf(ErrEnsure, fmtstring, args...)
	}
}
