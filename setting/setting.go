package setting

// 线上环境
const ONLINE = "online"

var (
	AppName    string
	AppVersion string
	AppDebug   bool
	AppEnv     string
)

func initSetting(appEnv string) {
	AppEnv = appEnv // os.Getenv("EYE_ENV")
	if AppEnv == "" {
		AppEnv = "dev"
	}
	if AppEnv != ONLINE {
		AppDebug = true
	}
	if v, err := Config.GetBool("debug"); err == nil {
		AppDebug = v
	}
	if v, err := Config.GetBool("app.debug"); err == nil {
		AppDebug = v
	}

	AppName = Config.MustString("app.name", "AppName")
	AppVersion = Config.MustString("app.version", "v0.0.0")
}
