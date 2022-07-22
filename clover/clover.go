package clover

/**
 * 守护协程安全退出模块
 * 个别需求需配合defer使用，ctrl+c 时候defer是不执行的，可忽略
 * 比如接收数据存buf异步上报，这时候该模块只负责接收退出，buf提交可由defer执行。因为获取数据可能被阻塞，导致该模块的安全退出不执行
 */

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var ins *clover

type clover struct {
	ctx    context.Context    // 携带done channel根ctx
	cancel context.CancelFunc // 根退出方法
	done   chan bool
	size   int
	sync.Mutex
}

const (
	DoneSize    = 100 // 完成channel buf大小
	DoneTimeout = 10  // 等待协程退出超时时间
)

func newClover() *clover {
	if ins == nil {
		ctx, cancel := context.WithCancel(context.Background())

		ins = &clover{
			ctx:    ctx,
			done:   make(chan bool, DoneSize),
			cancel: cancel,
			size:   0,
		}
	}
	return ins
}

func (p *clover) add() {
	p.Lock()
	defer p.Unlock()

	p.size = p.size + 1
}

func (p *clover) close() {
	p.Lock()
	defer p.Unlock()

	ins.cancel()

	log.Printf("total:%d \n", p.size)

	for i := 0; i < p.size; i++ {
		select {
		case <-p.done:
			log.Printf("ack: %d \n", i)
		case <-time.After(time.Duration(DoneTimeout) * time.Second):
			os.Exit(1)
		}
	}
	os.Exit(0)
}

func Add() (context.Context, chan bool) {
	ins.add()
	return ins.ctx, ins.done
}

func Notify() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigs
		ins.close()
	}()
}

func Ack() {
	ins.done <- true
}

func init() {
	newClover()
}
