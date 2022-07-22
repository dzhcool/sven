package econf

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dzhcool/sven/etcd"
)

/**
 * 基于etcd实现的配置文件
 */

type Config struct {
	sync.RWMutex
	url        string
	username   string
	password   string
	etcdClient *etcd.Client
	prefix     string
	data       map[string]string
	failedCnt  int // 失败次数
}

var ErrNil = errors.New("nil or type error")

const (
	EconfAliveDuration = 5 // 5秒探活时间
)

// 兼容旧的使用方式，不推荐
func New(prefix, url, username, password string) *Config {
	return newConfig(prefix, url, username, password)
}

func newConfig(prefix, url, username, password string) *Config {
	t := new(Config)
	t.prefix = prefix
	t.url = url
	t.username = username
	t.password = password
	t.data = make(map[string]string)
	if err := t.etcd(); err != nil {
		log.Fatal("connect etcd failed!")
		os.Exit(1)
	}
	t.load()         // 阻塞加载数据,防止后续程序调用失败
	go t.watch()     // 异步监听配置改动，及时刷新数据
	go t.autoAlive() // 启动监听
	return t
}

// 连接etcd
func (p *Config) etcd() error {
	var err error
	var etcdClient *etcd.Client
	if etcdClient, err = etcd.New(strings.Split(p.url, ","), p.username, p.password, 10, 15); err != nil {
		log.Printf("[econf] connect etcd failed: %s \n", err)
		return err
	}
	p.etcdClient = etcdClient
	return nil
}

// 加载数据
func (p *Config) load() {
	data, err := p.etcdClient.GetWithPrefix(p.prefix)
	if err != nil {
		log.Printf("[econf] get etcd data failed: %s \n", err)
		return
	}
	p.Lock()
	defer p.Unlock()

	p.failedCnt = 0

	i := 0
	for k, v := range data {
		p.data[k] = v
		i++
	}
	log.Printf("[econf] load conf num: %d \n", i)
}

// 监听改动
func (p *Config) watch() {
	log.Println("watch econf run ++")
	p.etcdClient.WatchWithPrefix(p.prefix, func(ctx context.Context, tp string, key, val []byte) {
		log.Printf("[econf] watch data key:%s val:%s \n", string(key), string(val))
		p.Lock()
		p.data[string(key)] = string(val)
		p.Unlock()
	})
}

// 遍历全部配置
func (p *Config) Map() (map[string]string, error) {
	buf := make(map[string]string)

	if p.data == nil {
		return nil, ErrNil
	}

	p.Lock()
	defer p.Unlock()

	for k, v := range p.data {
		buf[k] = v
	}

	return buf, nil
}

// 获取配置
func (p *Config) String(name string) string {
	p.RLock()
	defer p.RUnlock()

	if val, ok := p.data[name]; ok {
		return val
	}
	return ""
}

func (p *Config) StringDef(name, def string) string {
	val := p.String(name)
	if len(val) <= 0 {
		return def
	}
	return val
}

func (p *Config) Int(name string) int {
	p.RLock()
	defer p.RUnlock()

	if val, ok := p.data[name]; ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0
		}
		return i
	}
	return 0
}

func (p *Config) IntDef(name string, def int) int {
	val := p.Int(name)
	if val <= 0 {
		return def
	}
	return val
}

func (p *Config) Int64(name string) int64 {
	p.RLock()
	defer p.RUnlock()

	if val, ok := p.data[name]; ok {
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0
		}
		return i
	}
	return 0
}

func (p *Config) Int64Def(name string, def int64) int64 {
	val := p.Int64(name)
	if val <= 0 {
		return def
	}
	return val
}

// etcd探活，并reload conf
func (p *Config) autoAlive() {
	ticker := time.NewTicker(time.Second * time.Duration(EconfAliveDuration))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.alive()
		}
	}
}

func (p *Config) alive() error {
	failedCnt := p.getFailedCnt()

	_, err := p.etcdClient.Get(p.prefix + "/ping")
	if err != nil {
		failedCnt = p.setFailedCnt(1)
		log.Printf("econf alive failed err:%s failedCnt:%d \n", err, failedCnt)
		return err
	}

	if failedCnt > 0 {
		log.Println("econf alive success and reload")
		p.load()
	}
	return nil
}

func (p *Config) getFailedCnt() int {
	p.RLock()
	defer p.RUnlock()

	return p.failedCnt
}

func (p *Config) setFailedCnt(incr int) int {
	p.RLock()
	defer p.RUnlock()

	p.failedCnt = p.failedCnt + incr

	return p.failedCnt
}
