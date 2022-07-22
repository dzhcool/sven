package services

import (
	"sync"

	cron "github.com/robfig/cron/v3"
)

var (
	_cronIns  *CronSvc
	_cronOnce sync.Once
)

type CronSvc struct {
}

func GetCronSvc() *CronSvc {
	if _cronIns == nil {
		_cronOnce.Do(func() {
			_cronIns = new(CronSvc)
		})
	}
	return _cronIns
}

func (p *CronSvc) Handle() {
	var once sync.Once

	once.Do(func() {
		crond := cron.New(cron.WithSeconds())
		defer crond.Start()

		crond.AddFunc("1 1 1 * * *", func() {
			GetDemoSvc().Say()
		})
	})
}
