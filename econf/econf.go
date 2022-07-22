package econf

import "sync"

var (
	_configIns  *ConfigSvc
	_configOnce sync.Once
)

type ConfigSvc struct {
	config *Config
}

func InitConfig(prefix, url, username, password string) *ConfigSvc {
	if _configIns == nil {
		_configOnce.Do(func() {
			_configIns = &ConfigSvc{
				config: newConfig(prefix, url, username, password),
			}
		})
	}
	return _configIns
}

// 获取配置
func String(name string) string {
	return _configIns.config.String(name)
}

func StringDef(name, def string) string {
	return _configIns.config.StringDef(name, def)
}

func Int(name string) int {
	return _configIns.config.Int(name)
}

func IntDef(name string, def int) int {
	return _configIns.config.IntDef(name, def)
}

func Int64(name string) int64 {
	return _configIns.config.Int64(name)
}

func Int64Def(name string, def int64) int64 {
	return _configIns.config.Int64Def(name, def)
}
