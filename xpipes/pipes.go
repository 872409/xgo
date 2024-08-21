package xpipes

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"reflect"
	"runtime"
	"sync"
)

type PipesFun func() error

func Do(fun ...PipesFun) error {
	var err error
	New(fun...).DoHandle(func(outErr error, isRecover bool) {
		err = outErr
	})
	return err
}

//func DoRecover(fun ...PipesFun) (error, bool) {
//	var err error
//	var _isRecover bool
//	New(fun...).DoHandle(func(outErr error, isRecover bool) {
//		err = outErr
//		_isRecover = isRecover
//	})
//	return err, _isRecover
//}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
func New(fun ...PipesFun) *Pipes {
	pipes := &Pipes{}

	if len(fun) > 0 {
		pipes.Fill(fun...)
	}
	return pipes
}

type Pipes struct {
	debug bool
	//recover    bool
	stopped    bool
	isRecover  bool
	currentFun interface{}
	logger     logx.Logger
	//isDone    bool
	OutError error
	once     sync.Once
	fun      []PipesFun
}

func (r *Pipes) Fill(fun ...PipesFun) *Pipes {
	r.fun = append(r.fun, fun...)
	return r
}

func (r *Pipes) Debug(debug bool) *Pipes {
	r.debug = debug
	return r
}

//func (r *Pipes) Recover(recover bool) *Pipes {
//	r.recover = recover
//	return r
//}

func (r *Pipes) Off() {
	r.stopped = true
}

func (r *Pipes) Do() (outErr error) {
	r.DoHandle(func(_outErr error, _ bool) {
		outErr = _outErr
	})
	return outErr

}
func (r *Pipes) DoHandle(done func(outErr error, isRecover bool)) {

	r.once.Do(func() {
		r.run(done)
	})

}

//
//func (r *Pipes) DoOut() (outErr error, isRecover bool) {
//
//	r.DoHandle(func(_outErr error, _isRecover bool) {
//		outErr = _outErr
//		isRecover = _isRecover
//	})
//
//	return outErr, isRecover
//}

func (r *Pipes) run(done func(outErr error, isRecover bool)) {
	defer func() {
		if r.debug {
			return
		}

		err := recover()
		if err != nil {
			r.OutError = err.(error)
			r.isRecover = true
			fmt.Errorf("[FuncName:%+v]:%+v \n", r.currentFun, err)
			logx.Errorf("[FuncName:%+v]:%+v,\n", r.currentFun, err)
		}
		done(r.OutError, r.isRecover)
	}()

	for _, f := range r.fun {

		if r.stopped {
			break
		}

		r.currentFun = f

		if err := f(); err != nil {
			if r.debug {
				panic(err)
			}

			r.OutError = err
			break
		}
	}

}
