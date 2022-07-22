package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"website/middlewares"
	"website/routers"
	"website/services"

	"github.com/dzhcool/sven/buildinfo"
	"github.com/dzhcool/sven/setting"
	log "github.com/dzhcool/sven/zapkit"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	confFile string
	confEnv  string
)

func initArgs() {
	flag.StringVar(&confFile, "c", "", "请指定配置文件")
	flag.StringVar(&confEnv, "e", "", "请指定运行环境")
	flag.Parse()
}

func main() {
	initArgs()

	setting.InitSetting(confFile, confEnv)

	// 初始化日志模块
	log.Init(zapkitConf())
	defer log.Sync()

	// 打印系统启动信息
	printStarting()

	//监听程序退出事件
	// subscribeSignal()

	r := gin.New()
	routers.Register(r)

	// 初始化中间件
	middlewares.InitMiddleware(r)

	// 初始化service模块
	services.InitService()

	r.Run(":" + setting.Config.MustString("http.port", "8080"))
}

// 定义close方法
func close() {
	log.Info("start cleaning")
	services.CloseService()

	log.Info("exit")
	os.Exit(0)
}

// 获取zapkit日志配置
func zapkitConf() *log.ZapkitConfig {
	zapkitConfig := log.ZapkitConfig{
		File:       setting.Config.MustString("zapkit.file", "/tmp/zapkit.log"),
		Level:      setting.Config.MustString("zapkit.level", "info"),
		MaxSize:    setting.Config.MustInt("zapkit.maxsize", 512),
		MaxBackups: setting.Config.MustInt("zapkit.maxbackups", 10),
		MaxAge:     setting.Config.MustInt("zapkit.age", 7),
		Compress:   setting.Config.MustBool("zapkit.compress", false),
	}
	return &zapkitConfig
}

// 打印启动日志
func printStarting() {
	log.Info(setting.Config.MustString("app.name", ""), zap.String("env", setting.AppEnv),
		zap.String("version", setting.Config.MustString("app.version", "")),
		zap.String("loglevel", setting.Config.MustString("zapkit.level", "info")),
		zap.String("buildTime", buildinfo.GetBuildTime()),
		zap.String("buildGoVersion", buildinfo.GetBuildGoVersion()))
}

// 监听退出事件
func subscribeSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				close()
			default:
			}
		}
	}()
}
