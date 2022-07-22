package clover

import (
	"testing"
	"time"
)

func Test_run(t *testing.T) {
	ctx, done := Notify()
	for {
		select {
		case <-ctx.Done():
			done <- true
			return
		default:
		}
		//TODO
		time.Sleep(1 * time.Second)
	}
}
