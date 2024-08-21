package xutils

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTicker2(t *testing.T) {
	list := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	for i, v := range list {
		go func() {
			fmt.Println(i, v)
		}()
	}
	time.Sleep(10 * time.Second)
}
func TestNewTicker(t *testing.T) {
	i := 0
	tt := NewTicker(1*time.Second, func() {
		i++
		fmt.Println("do", i)
		time.Sleep(2 * time.Second)
		fmt.Println("done", i)
	})
	time.Sleep(10 * time.Second)
	tt <- true
}
