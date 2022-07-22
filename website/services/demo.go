package services

import (
	"sync"
	"time"

	log "github.com/dzhcool/sven/zapkit"
)

var (
	_demoIns  *DemoSvc
	_demoOnce sync.Once
)

type DemoSvc struct{}

func GetDemoSvc() *DemoSvc {
	if _demoIns == nil {
		_demoOnce.Do(func() {
			_demoIns = new(DemoSvc)
		})
	}
	return _demoIns
}

// daemon示例
func (p *DemoSvc) Daemon() {
	for {
		log.Info("daemon demo run")
		time.Sleep(10 * time.Minute)
	}
}

// crontab示例
func (p *DemoSvc) Say() {
	log.Debug("crontab demo run")
}
