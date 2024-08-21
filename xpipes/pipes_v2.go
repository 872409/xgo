package xpipes

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"main/comm/xerr"
	"sync"
)

type PipesFunV2 func(args ...any) (any, error)

func DoV2(fun ...PipesFunV2) (any, error) {
	var err error
	var out any
	NewV2(fun...).DoHandle(func(_out any, outErr error, isRecover bool) {
		err = outErr
		out = _out
	})
	return out, err
}

//func DoRecover(fun ...PipesFunV2) (error, bool) {
//	var err error
//	var _isRecover bool
//	New(fun...).DoHandle(func(outErr error, isRecover bool) {
//		err = outErr
//		_isRecover = isRecover
//	})
//	return err, _isRecover
//}

func NewV2(fun ...PipesFunV2) *PipesV2 {
	pipes := &PipesV2{}

	if len(fun) > 0 {
		pipes.Fill(fun...)
	}
	return pipes
}

type PipesV2 struct {
	debug bool
	//recover    bool
	stopped    bool
	isRecover  bool
	currentFun interface{}
	logger     logx.Logger
	//isDone    bool
	OutError error
	once     sync.Once
	fun      []PipesFunV2
}

func (r *PipesV2) Fill(fun ...PipesFunV2) *PipesV2 {
	r.fun = append(r.fun, fun...)
	return r
}

func (r *PipesV2) Debug(debug bool) *PipesV2 {
	r.debug = debug
	return r
}

//func (r *PipesV2) Recover(recover bool) *PipesV2 {
//	r.recover = recover
//	return r
//}

func (r *PipesV2) Off() {
	r.stopped = true
}

func (r *PipesV2) Do() (out any, outErr error, isRecover bool) {
	r.DoHandle(func(_out any, _outErr error, _isRecover bool) {
		outErr = _outErr
		out = _out
		isRecover = _isRecover
	})
	return

}
func (r *PipesV2) DoHandle(done func(out any, outErr error, isRecover bool)) {

	r.once.Do(func() {
		r.run(done)
	})

}

//
//func (r *PipesV2) DoOut() (outErr error, isRecover bool) {
//
//	r.DoHandle(func(_outErr error, _isRecover bool) {
//		outErr = _outErr
//		isRecover = _isRecover
//	})
//
//	return outErr, isRecover
//}

func (r *PipesV2) run(done func(out any, outErr error, isRecover bool)) {
	var _out any
	defer func() {
		if r.debug {
			return
		}

		err := recover()
		if err != nil {
			r.OutError = err.(error)
			r.isRecover = true
			err = xerr.WithStack(r.OutError, 6)
			fmt.Errorf("[FuncName:%+v]:%+v \n", r.currentFun, err)
			logx.Errorf("[FuncName:%+v]:%+v,\n", r.currentFun, err)
		}
		done(_out, r.OutError, r.isRecover)
	}()

	for _, f := range r.fun {

		if r.stopped {
			break
		}

		r.currentFun = f
		out, err := f(_out)
		if err != nil {
			if r.debug {
				panic(err)
			}
			r.OutError = err
			fName := getFunctionName(f)
			logx.Errorf("[FuncName:%s]:%+v,", fName, xerr.WithStack(r.OutError, 6))
			break
		}
		_out = out
	}

}
