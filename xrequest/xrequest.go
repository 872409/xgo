package xrequest

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

type Request struct {
	ctx    context.Context
	Url    string
	Data   map[string]interface{}
	Header map[string]string
}

//// Post 请求
//func Post(req Request) (res string, err error) {
//	maps := make(map[string]interface{})
//	if len(data) > 0 {
//		maps = data[0]
//	}
//	r, err := g.Client().Post(context.Background(), url, maps)
//	if err != nil {
//		xlog.Default.Error(err.Error(), r)
//		return "", err
//	}
//	defer func(r *gclient.Response) {
//		_ = r.Close()
//	}(r)
//	return r.ReadAllString(), nil
//}

// Get 请求
func Get(req Request) (res string, err error) {
	c := g.Client()
	c.SetHeaderMap(req.Header)
	if r, e := c.Get(req.ctx, req.Url, req.Data); e != nil {
		panic(e)
	} else {
		return r.ReadAllString(), nil
	}

}
