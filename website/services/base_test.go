package services

import (
	"website/test"
)

// go test -v ./service/ -test.run TestCrond
// go test -v ./service/                   // 测试全部

func init() {
	test.StubInitConfig()
	test.StubInitZapkit()
	test.StubInitRedis()
}
