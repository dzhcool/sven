package dict

import (
	"encoding/json"
	"testing"
)

type Ad struct {
	Ad_id   int    `json:"ad_id"`
	Pos     string `json:"pos"`
	Type    string `json:"type"`
	Jumpurl string `json:"jumpurl"`
	Url     string `json:"url"`
}

type AdRest struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
	Data   []Ad   `json:"data"`
}

func TestCache(t *testing.T) {
	var (
		k = "testkey"
		v = "testval"
	)
	cacheTest := Cache("test")
	//存储
	cacheTest.Set(k, v, 0)

	//查询缓存值
	val, err := cacheTest.String(cacheTest.Get(k))

	if len(val) > 0 {
		t.Log("Cache Ok")
	} else {
		t.Fatal("Cache Error", err)
	}
}

func TestJson(t *testing.T) {

	js := `{"ad_id":1, "pos":"31057001", "ext":"xxx"}`

	var buf Ad
	err := json.Unmarshal([]byte(js), &buf)
	if err != nil {
		t.Fatal("decode json error", err)
	}

	var rest AdRest
	d := []Ad{buf}
	rest.Data = d

	_, err = json.Marshal(rest)
	if err != nil {
		t.Fatal("encode json error", err)
	}
}

func TestDel(t *testing.T) {
	var (
		k = "testdel"
		v = "testval"
	)
	cacheTest := Cache("test")
	//存储
	cacheTest.Set(k, v)

	cacheTest.Delete(k)

	//查询缓存值
	val, err := cacheTest.String(cacheTest.Get(k))
	if val == v {
		t.Fatal("delete error", err)
	}
}
