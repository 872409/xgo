package xpipes

import (
	"errors"
	"fmt"
	"testing"
)

func tfunc() error {
	return errors.New("errr1")
}

func tfunc2(id string) func() error {
	return func() error {
		if id == "a" {
			return nil
		}
		return errors.New("errr1:" + id)
	}
}
func TestPipes_OnRecover(t *testing.T) {
	var err error
	p := New(
		func() error {

			return nil
		},
		tfunc2("a"), tfunc2("b"),
		tfunc)

	p.Fill(func() error {
		p.Off()
		return nil
	})

	p.DoHandle(func(outErr error, isRecover bool) {
		err = outErr
		fmt.Println("isRecover", isRecover, err)
	})

	fmt.Println(err)
}

func TestPipes_On(t *testing.T) {
	var x []string
	var err error
	err = Do(func() error {
		x[1] = "s"
		return nil
	}, func() error {
		return errors.New("errr1")
	})

	fmt.Println("done:", err, x)
}
