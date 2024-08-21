package xpipes

import (
	"fmt"
	"testing"
)

func TestPipes_OnRecoverV2(t *testing.T) {
	var err error
	p := NewV2(
		func(a ...any) (any, error) {
			return 100, nil
		},
		func(a ...any) (any, error) {
			var _a = a[0].(int)
			fmt.Println(_a)
			return "aaaa", nil // ,// errors.New("错了")
		},
		func(a ...any) (any, error) {
			var _a = a[0].(string)
			fmt.Println(_a)
			return 2, nil
		},
	)

	p.Fill(func(a ...any) (any, error) {
		fmt.Println("last 0")
		//p.Off()
		return "last....abc ", nil
	}, func(a ...any) (any, error) {
		fmt.Println("last 1")
		return 1110, nil
	})

	out, err, e := p.Do()
	fmt.Println(out, err, e)
	//
	//p.DoHandle(func(out any, outErr error, isRecover bool) {
	//	err = outErr
	//	fmt.Println("out", out, err, isRecover)
	//})

	//fmt.Println(err)
}
