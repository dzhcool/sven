package zapkit

import (
	"errors"
)

// 兼容第三方框架初始化方法

func ThirdInit(logPath, logLevel string) error {
	if logkit != nil {
		return errors.New("already initialized")
	}

	zapkitConfig := &ZapkitConfig{
		File:       logPath,
		Level:      logLevel,
		MaxSize:    1024,
		MaxBackups: 10,
		MaxAge:     7,
		Compress:   true,
	}

	return initZapkit(zapkitConfig, "")
}
