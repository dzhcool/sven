package etcd

import (
	"testing"
)

var etcdHostTest = []string{"127.0.0.1:2379"}

func Test_put(t *testing.T) {
	ins, err := New(etcdHostTest, "","",3, 5)
	if err != nil {
		t.Errorf("connect etcd err:%s", err.Error())
		return
	}
	defer ins.Close()

	if err := ins.Put("/dzhcool/name", "dangzihao"); err != nil {
		t.Errorf("put err:%s", err.Error())
		return
	}

	vals, err := ins.Get("/dzhcool/name")
	if err != nil {
		t.Errorf("get err:%s", err.Error())
		return
	}
	for k, v := range vals {
		t.Logf("get k:%s v:%s", k, v)
	}

	vals, err = ins.GetWithPrefix("/dzhcool")
	if err != nil {
		t.Errorf("get err:%s", err.Error())
		return
	}
	for k, v := range vals {
		t.Logf("getWithPrefix k:%s v:%s", k, v)
	}
}
