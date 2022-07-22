package dict

import (
	"sync"
	"time"
)

type CacheItem struct {
	sync.RWMutex
	key     interface{}
	data    interface{}
	addTime time.Time
	life    time.Duration
}

func createCacheItem(k interface{}, v interface{}, l time.Duration) CacheItem {
	t := time.Now()
	return CacheItem{
		key:     k,
		data:    v,
		addTime: t,
		life:    l,
	}
}

func (p *CacheItem) Life() time.Duration {
	return p.life
}

func (p *CacheItem) AddTime() time.Time {
	return p.addTime
}

func (p *CacheItem) Key() interface{} {
	return p.key
}

func (p *CacheItem) Data() interface{} {
	return p.data
}

func (p *CacheItem) Expired() bool {
	p.RLock()
	defer p.RUnlock()
	if p.life > 0 && time.Now().Sub(p.addTime) >= p.life {
		return true
	}
	return false
}
