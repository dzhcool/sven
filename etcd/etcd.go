package etcd

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	client      *clientv3.Client
	readTimeout time.Duration
}

// 初始化etcd
func New(addr []string, username, password string, readTimeout, dialTimeout int) (*Client, error) {
	conf := clientv3.Config{
		Endpoints:   addr,
		Username:    username,
		Password:    password,
		DialTimeout: time.Duration(dialTimeout) * time.Second,
	}
	client, err := clientv3.New(conf)
	if err != nil {
		return nil, err
	}
	return &Client{
		client:      client,
		readTimeout: time.Duration(readTimeout) * time.Second,
	}, nil
}

// 关闭
func (p *Client) Close() {
	p.client.Close()
}

// 写入数据
func (p *Client) Put(key, val string) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.readTimeout)
	_, err := p.client.Put(ctx, key, val)
	cancel()

	return err
}

// 获取数据
func (p *Client) Get(key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.readTimeout)
	resp, err := p.client.Get(ctx, key)
	cancel()

	if err != nil {
		return nil, err
	}

	buf := make(map[string]string)
	for _, ev := range resp.Kvs {
		buf[string(ev.Key)] = string(ev.Value)
	}
	return buf, nil
}

// 前缀获取数据
func (p *Client) GetWithPrefix(prefix string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.readTimeout)
	resp, err := p.client.Get(ctx, prefix, clientv3.WithPrefix())
	cancel()

	if err != nil {
		return nil, err
	}

	buf := make(map[string]string)
	for _, ev := range resp.Kvs {
		buf[string(ev.Key)] = string(ev.Value)
	}
	return buf, nil
}

// 获取数据，按key排序
func (p *Client) GetKSort(key string) ([]string, []string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.readTimeout)
	resp, err := p.client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	cancel()

	if err != nil {
		return nil, nil, err
	}

	keys := make([]string, 0)
	vals := make([]string, 0)

	for _, ev := range resp.Kvs {
		keys = append(keys, string(ev.Key))
		vals = append(vals, string(ev.Value))
	}
	return keys, vals, nil
}

// 删除数据
func (p *Client) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.readTimeout)
	defer cancel()

	if _, err := p.client.Get(ctx, key, clientv3.WithPrefix()); err != nil {
		return err
	}

	if _, err := p.client.Delete(ctx, key, clientv3.WithPrefix()); err != nil {
		return err
	}
	return nil
}

// 监听数据
func (p *Client) Watch(key string, handle func(ctx context.Context, tp string, key, val []byte)) {
	rch := p.client.Watch(context.Background(), key)

	for wresp := range rch {
		for _, ev := range wresp.Events {
			go handle(context.Background(), ev.Type.String(), ev.Kv.Key, ev.Kv.Value)
		}
	}
}

// 前缀监听数据
func (p *Client) WatchWithPrefix(prefix string, handle func(ctx context.Context, tp string, key, val []byte)) {
	rch := p.client.Watch(context.Background(), prefix, clientv3.WithPrefix())

	for wresp := range rch {
		for _, ev := range wresp.Events {
			go handle(context.Background(), ev.Type.String(), ev.Kv.Key, ev.Kv.Value)
		}
	}
}

// put租约key val
func (p *Client) LeasePutKeep(key, val string, ttl int64) (clientv3.LeaseID, <-chan *clientv3.LeaseKeepAliveResponse, error) {
	ctx := context.Background()
	lease := clientv3.NewLease(p.client)
	resp, err := lease.Grant(ctx, ttl)
	if err != nil {
		return 0, nil, err
	}
	if _, err = p.client.Put(ctx, key, val, clientv3.WithLease(resp.ID)); err != nil {
		return 0, nil, err
	}
	ch, err := lease.KeepAlive(ctx, resp.ID)
	if err != nil {
		return 0, nil, err
	}

	return resp.ID, ch, nil
}
