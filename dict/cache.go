// 缓存
// todo 以后增加关闭项目缓存数据，启动时候加载进来
package dict

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex

	ErrNil = errors.New("nil or type error")
)

type CacheTable struct {
	sync.RWMutex
	name    string
	items   map[interface{}]*CacheItem
	addtime time.Time
}

func Cache(table string) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		t = &CacheTable{
			name:    table,
			items:   make(map[interface{}]*CacheItem),
			addtime: time.Now(),
		}
	}
	return t
}

func (p *CacheTable) Count(k interface{}) int {
	return len(p.items)
}

func (p *CacheTable) Set(k, v interface{}, l ...int) *CacheItem {
	lifetime := 0 * time.Second
	if len(l) > 0 {
		lifetime = time.Duration(l[0]) * time.Second
	}
	item := createCacheItem(k, v, lifetime)

	p.Lock()
	p.items[k] = &item
	p.Unlock()
	return &item
}

func (p *CacheTable) Add(k, v interface{}, l ...int) *CacheItem {
	return p.Set(k, v, l...)
}

func (p *CacheTable) Get(k interface{}) (interface{}, error) {
	p.RLock()
	r, ok := p.items[k]
	if !ok {
		p.RUnlock()
		return nil, ErrNil
	}
	p.RUnlock()

	if r.Expired() {
		p.Lock()
		delete(p.items, k)
		p.Unlock()
		return nil, ErrNil
	}
	return r.data, nil
}

func (p *CacheTable) Item(k interface{}) *CacheItem {
	p.RLock()
	r, ok := p.items[k]
	p.RUnlock()

	r.Lock()
	defer r.Unlock()

	if !ok {
		r = nil
	}
	return r
}

func (p *CacheTable) Items() map[interface{}]*CacheItem {
	p.Lock()
	defer p.Unlock()

	return p.items
}

func (p *CacheTable) Exists(k interface{}) bool {
	p.RLock()
	defer p.RUnlock()

	_, ok := p.items[k]
	return ok
}

func (p *CacheTable) Delete(k interface{}) (*CacheItem, error) {
	p.RLock()
	r, ok := p.items[k]
	p.RUnlock()

	if ok {
		p.Lock()
		delete(p.items, k)
		p.Unlock()
		return r, nil
	}
	return nil, ErrNil
}

func (p *CacheTable) Int(reply interface{}, err error) (int, error) {
	return Int(reply, err)
}

func (p *CacheTable) Int64(reply interface{}, err error) (int64, error) {
	return Int64(reply, err)
}

func (p *CacheTable) Float64(reply interface{}, err error) (float64, error) {
	return Float64(reply, err)
}

func (p *CacheTable) String(reply interface{}, err error) (string, error) {
	return String(reply, err)
}

func Int(reply interface{}, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case int64:
		x := int(reply)
		if int64(x) != reply {
			return 0, strconv.ErrRange
		}
		return x, nil
	case int:
		x := int(reply)
		return x, nil
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 0)
		return int(n), err
	case string:
		n, err := strconv.Atoi(reply)
		return n, err
	case nil:
		return 0, ErrNil
	}
	return 0, ErrNil
}

func Int64(reply interface{}, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case int64:
		return reply, nil
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 64)
		return n, err
	case nil:
		return 0, ErrNil
	}
	return 0, ErrNil
}

func Float64(reply interface{}, err error) (float64, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case []byte:
		n, err := strconv.ParseFloat(string(reply), 64)
		return n, err
	case nil:
		return 0, ErrNil
	}
	return 0, ErrNil
}

func String(reply interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}
	switch reply := reply.(type) {
	case []byte:
		return string(reply), nil
	case string:
		return reply, nil
	case int:
		n := strconv.Itoa(reply)
		return n, nil
	case nil:
		return "", ErrNil
	}
	return "", ErrNil
}
