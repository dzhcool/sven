package services

import (
	"context"
	"sync"
	"time"

	log "github.com/dzhcool/sven/zapkit"
)

var (
	_baseSvc *baseSvc

	// 控制初始化和关闭方法只执行一次
	initServiceOnce  sync.Once
	closeServiceOnce sync.Once
)

type baseSvc struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// 可注册service收尾工作
func closeService() {
	//TODO
}

// 可注册初始化service工作
func initService() {
	//TODO
}

func InitService() {
	if _baseSvc == nil {
		initServiceOnce.Do(func() {
			_baseSvc = new(baseSvc)

			ctx, cancel := context.WithCancel(context.Background())

			_baseSvc.ctx = ctx
			_baseSvc.cancel = cancel

			initService()
			GetCronSvc().Handle() // 初始化定时任务
		})
	}
}

func CloseService() {
	closeServiceOnce.Do(func() {
		if _baseSvc != nil {
			_baseSvc.cancel()
		}

		log.Info("wait for the service to exit")

		closeService()

		time.Sleep(5 * time.Second)
		log.Info("service exited")
	})
}

func CtxSvc() context.Context {
	return _baseSvc.ctx
}
