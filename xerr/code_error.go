package xerr

import (
	"errors"
	"fmt"
	"io"
)

/**
常用通用固定错误
*/

type CodeError struct {
	code string
	msg  string
	//cause error
}

func (e *CodeError) WithStack(skips ...int) *withStack {
	var skip = 4
	if len(skips) > 0 {
		skip = skips[0]
	}
	return WithStack(e, skip)
}

func (e *CodeError) Append(format string, val ...any) *CodeError {
	msg := fmt.Sprintf("%s %s", e.msg, fmt.Sprintf(format, val...))
	return &CodeError{
		code: e.code,
		msg:  msg,
	}
}

func (e *CodeError) JoinF(format string, val ...any) error {
	return errors.Join(e, fmt.Errorf(format, val...))
}

func (e *CodeError) Join(errs ...error) error {
	return errors.Join(append([]error{e}, errs...)...)
	//return errors.Join(append(errs, e)...)
}

func (e *CodeError) Code() string {
	return e.code
}

func (e *CodeError) Msg() string {
	return e.msg
}

func (e *CodeError) Is(err error) bool {
	if xe := AsCodeError(err); xe != nil {
		return xe.code == e.code
	}
	return false
}
func Is(err1 error, err2 error) bool {
	if xe := AsCodeError(err1); xe != nil {
		return xe.Is(err2)
	}
	return false
}

//func (e *CodeError) Cause() error {
//	return e.cause
//}

//	func (e *CodeError) Unwrap() error {
//		return e.cause
//	}
func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%s，ErrMsg:%s", e.code, e.msg)
}

func (e *CodeError) Format(s fmt.State, verb rune) {
	_, _ = io.WriteString(s, e.Error())
	//
	//switch verb {
	//case 'v':
	//	if s.Flag('+') {
	//		_, _ = io.WriteString(s, e.OutError())
	//		return
	//	}
	//	fallthrough
	//case 's', 'q':
	//	_, _ = io.WriteString(s, e.OutError())
	//}
}
