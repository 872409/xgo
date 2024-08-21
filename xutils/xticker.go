package xutils

import (
	"time"
)

func NewTicker(d time.Duration, fun ...func()) chan bool {
	ticker := time.NewTicker(d)
	stopChan := make(chan bool)

	go func(ticker *time.Ticker) {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if len(fun) > 0 {
					for _, f := range fun {
						f()
					}
				}
			case stop := <-stopChan:
				if stop {
					//fmt.Println("Ticker Stop! Channel must be closed")
					return
				}
			}
		}
	}(ticker)

	return stopChan
}
