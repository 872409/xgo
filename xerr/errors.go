package xerr

import (
	"errors"
	"fmt"
	"gorm.io/gen"
	"gorm.io/gorm"
	"strings"
	"sync"
)

// OK 成功返回
//const OK string = "200"

/**(前3位代表业务,后三位代表具体功能)**/

// GLOBAL 全局错误码
var (
	GLOBAL              = "100"
	DataBaseUpdateError = NewCodeError(GLOBAL+"203", "Db Update OutError")
)
var message = &sync.Map{}

// NewCodeError 系统错误
// code: 错误码

func NewCodeError(code string, msg string) *CodeError {

	//c, e := strconv.Atoi(code)
	//fmt.Println(code, c, e, msg)
	if _, f := message.Load(code); f {
		panic(fmt.Sprintf("code:%s exists,msg:%s", code, msg))
	}

	message.Store(code, msg)

	return &CodeError{
		code: code,
		msg:  msg,
	}
}

func IsCodeErr(errCode string) bool {
	if _, ok := message.Load(errCode); ok {
		return true
	} else {
		return false
	}
}

//func DbSelectError(err error, appendErrorMsg ...string) error {
//	if errors.Is(err, gorm.ErrRecordNotFound) {
//		appendMsg := ""
//		if len(appendErrorMsg) > 0 {
//			appendMsg = strings.Join(appendErrorMsg, " ")
//		}
//		return DataBaseRecordNotFound.WithStack(5).Join(err, fmt.Errorf(appendMsg))
//	}
//	return nil
//}

/*
DbSelectError
单条查询错误信息处理
*/
func DbSelectError(err error) error {
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return nil
}

func DbUpdateError(result gen.ResultInfo, appendErrorMsg ...string) error {
	if result.Error != nil || result.RowsAffected == 0 {
		appendMsg := ""
		if len(appendErrorMsg) > 0 {
			appendMsg = strings.Join(appendErrorMsg, " ")
		}
		err := DataBaseUpdateError.WithStack(5).Join(result.Error, fmt.Errorf("%s OutError:%+v,RowsAffected:%d", appendMsg, result.Error, result.RowsAffected))

		return err
	}
	return nil
}
