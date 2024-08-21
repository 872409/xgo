package xresponse

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// Response TODO: 项目对外的统一返回
func Response(w http.ResponseWriter, resp interface{}, err error) {
	body := struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data,omitempty"`
	}{}
	if err != nil {
		body.Code = -1
		body.Msg = err.Error()
	} else {
		body.Code = 0
		body.Msg = "OK"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}
