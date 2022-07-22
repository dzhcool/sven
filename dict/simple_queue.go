// 用slice做简单队列
package dict

import (
	"sync"
)

type SimpleQueue struct {
	sync.RWMutex
	initSize int
	data     []interface{}
}

func NewSimpleQueue(initSize int) *SimpleQueue {
	if initSize <= 0 {
		initSize = 50
	}
	t := new(SimpleQueue)
	t.initSize = initSize
	t.data = make([]interface{}, 0, t.initSize)
	return t
}

func (p *SimpleQueue) Count() int {
	p.RLock()
	defer p.RUnlock()
	return len(p.data)
}

func (p *SimpleQueue) Insert(v interface{}) error {
	if p.data == nil {
		p.data = make([]interface{}, 0, p.initSize)
	}
	p.data = append(p.data, v)
	return nil
}

// 从头取数据
func (p *SimpleQueue) Front() (interface{}, error) {
	p.Lock()
	defer p.Unlock()

	l := len(p.data)
	if l <= 0 {
		return nil, ErrNil
	}
	item := p.data[0]
	p.data = p.data[1:]

	return item, nil
}

// 从结尾取数据
func (p *SimpleQueue) Last() (interface{}, error) {
	p.Lock()
	defer p.Unlock()

	l := len(p.data)
	if l <= 0 {
		return nil, ErrNil
	}
	item := p.data[l-1]
	p.data = p.data[:l-2]

	return item, nil
}
