package utils

import (
	"net/url"
	"testing"
)

func Test_HTTPAuthPostForm(t *testing.T) {
	data := make(url.Values)
	data.Add("name", "davis")
	data["key"] = []string{"this key test"}

	resp, err := HTTPAuthPostForm("http://127.0.0.1:92/post", data, "", "", 10)
	if err != nil {
		t.Fatalf("err:%s", err.Error())
		return
	}
	t.Logf("resp:%s", string(resp))
}
