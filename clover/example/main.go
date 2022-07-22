package main

import (
	"fmt"
	"time"

	"github.com/dzhcool/sven/clover"
)

func run(i int) {
	defer func() {
		//TODO 如果下方for中逻辑可阻塞，一定要在defer中收尾，否则select可能不会执行
		fmt.Printf("进行收尾工作:%d", i)
	}()

	ctx, done := clover.Add()
	for {
		select {
		case <-ctx.Done():
			//TODO 这里可进行收尾工作
			done <- true
			return
		default:
		}
	}
}

func main() {
	clover.Notify()

	for i := 0; i < 5; i++ {
		go run(i)
	}
	for {
		time.Sleep(10 * time.Second)
	}
}
