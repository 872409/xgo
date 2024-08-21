package xerr

import (
	"errors"
	//"google.golang.org/grpc/status"
)

func AsCodeError(err error) *CodeError {
	codeError := new(CodeError)
	if errors.As(err, &codeError) {
		return codeError
	}
	return nil
}

func UnwrapCodeError(err error) *CodeError {

	if e, ok := err.(*CodeError); ok {
		return e
	}

	switch x := err.(type) {
	case interface{ Unwrap() error }:
		err = x.Unwrap()
		if err == nil {
			return nil
		}
		return UnwrapCodeError(err)

	case interface{ Unwrap() []error }:
		for _, err := range x.Unwrap() {
			if e := UnwrapCodeError(err); e != nil {
				return e
			}
		}
		return nil
	default:
		//fmt.Println("err.(type)", x)
		return nil
	}
}

func AsCodeErrorCase(err error, call func(err *CodeError)) *CodeError {
	codeError := AsCodeError(err)

	if codeError != nil {
		call(codeError)
		return codeError
	}
	return nil
}

//func FromRPCError(err error) error {
//	if err == nil {
//		return nil
//	}
//
//	if e, ok := status.FromError(err); ok {
//		code := fmt.Sprintf("%d", e.Code())
//		err = NewCodeError(code, e.Message()).Join(err)
//	}
//	return err
//}

func Join(err ...error) error {
	return errors.Join(err...)
}
