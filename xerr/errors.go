package xerr

import (
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gen"
	"gorm.io/gorm"
	"strings"
	"xgo/xlog"
)

// OK 成功返回
//const OK string = "200"

/**(前3位代表业务,后三位代表具体功能)**/

// GLOBAL 全局错误码
var (
	GLOBAL                  = "100"
	RPCError                = internalNewCodeError(GLOBAL+"000", "internal error")
	ServerCommonError       = internalNewCodeError(GLOBAL+"001", "common error")
	RequestParamError       = internalNewCodeError(GLOBAL+"002", "param error")
	DomainInvalid           = internalNewCodeError(GLOBAL+"003", "domain invalid")
	RequestContentTypeError = internalNewCodeError(GLOBAL+"004", "content-type error")

	DataBaseError          = internalNewCodeError(GLOBAL+"201", "Db OutError")
	DataBaseRecordNotFound = internalNewCodeError(GLOBAL+"202", "record not found")
	DataBaseUpdateError    = internalNewCodeError(GLOBAL+"203", "Db Update OutError")
	SignError              = internalNewCodeError(GLOBAL+"007", "sign error")
	LockError              = internalNewCodeError(GLOBAL+"008", "lock error")
	ChainError             = internalNewCodeError(GLOBAL+"009", "chain error")
	CurrencyDecodeError    = internalNewCodeError(GLOBAL+"010", "currency decode error")

	TokenInvalid       = internalNewCodeError(GLOBAL+"100", "token invalid")
	TokenExpireError   = internalNewCodeError(GLOBAL+"101", "token is expired")
	TokenGenerateError = internalNewCodeError(GLOBAL+"102", "token Generate OutError")
	TOTPCodeInvalid    = internalNewCodeError(GLOBAL+"103", "OTP code invalid")
)

var (
	USER          = "300"
	UserNotExists = internalNewCodeError(USER+"001", "user not exits")
)

var (
	Game              = "400"
	NoFreeGameNumProp = internalNewCodeError(Game+"001", "NoFreeGameNumProp")
)

var (
	System      = "200"
	LevelConfig = internalNewCodeError(System+"001", "LevelConfig not exits")
)

//
//var (
//	Exchange                 = "300"
//	ExchangeAmountMinLimited = internalNewCodeError(Exchange+"001", "min amount limited")
//	ExchangeAmountMaxLimited = internalNewCodeError(Exchange+"002", "max amount limited")
//)

var message map[string]string

// internalNewCodeError 系统错误
// code: 错误码

func internalNewCodeError(code string, msg string) *CodeError {
	if message == nil {
		message = map[string]string{}
	}
	//c, e := strconv.Atoi(code)
	//fmt.Println(code, c, e, msg)
	if _, f := message[code]; f {
		panic(fmt.Sprintf("code:%s exists,msg:%s", code, msg))
	}

	message[code] = msg

	return &CodeError{
		code: code,
		msg:  msg,
	}
}

func IsCodeErr(errCode string) bool {
	if _, ok := message[errCode]; ok {
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
		xlog.Default.Error(err.Error(), err)
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

		logx.Error(err)

		return err
	}
	return nil
}
