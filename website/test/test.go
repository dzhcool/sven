package test

import (
	// "agent/middleware/redis"

	"github.com/dzhcool/sven/setting"
	log "github.com/dzhcool/sven/zapkit"
)

// 初始化项目配置
func StubInitConfig() {
	setting.InitSetting("../conf/app.ini", "")
}

// 初始化日志
func StubInitZapkit() {
	zapkitConfig := &log.ZapkitConfig{
		File:       setting.Config.MustString("zapkit.file", "/tmp/zapkit.log"),
		Level:      setting.Config.MustString("zapkit.level", "info"),
		MaxSize:    setting.Config.MustInt("zapkit.maxsize", 512),
		MaxBackups: setting.Config.MustInt("zapkit.maxbackups", 10),
		MaxAge:     setting.Config.MustInt("zapkit.age", 7),
		Compress:   setting.Config.MustBool("zapkit.compress", false),
	}
	log.Init(zapkitConfig)
}

// 初始化redis
// func StubInitRedis() {
// 	redis.InitRedis()
// }
