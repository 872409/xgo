package xpipes

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func ExecTimeout(timeoutMillisecond time.Duration, fun ...PipesFun) error {
	err := make(chan error)

	wg := &sync.WaitGroup{}

	go func() {
		for _, f := range fun {
			wg.Add(1)
			go func(wg *sync.WaitGroup, err chan error, fn PipesFun) {
				defer wg.Done()
				e := fn()
				if e != nil {
					err <- e
				}
			}(wg, err, f)
		}
		wg.Wait()
		close(err)
	}()

	select {
	case e := <-err:
		fmt.Println("e:", e)
		return e
	case <-time.After(timeoutMillisecond * time.Millisecond):
		fmt.Println("timed out..")
		return errors.New("timed out")
	}
}
