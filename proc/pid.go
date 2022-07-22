package proc

import (
	"common/utils"
	"log"
	"os"
	"strconv"
	"syscall"
)

const (
	PidFile = "/var/run/gp.pid"
)

var (
	pidFile string
	sigs    chan os.Signal
)

func SetPid(pidfile string) int {
	pidFile = pidfile
	if pidFile == "" {
		pidFile = PidFile
	}

	// 判断pid是否存在，如果存在且有效退出程序
	if Existed() {
		log.Printf("pid file  %s  exist \n", pidFile)
		os.Exit(1)
		return 0
	}

	pid := os.Getpid()
	utils.WriteFile(pidFile, []byte(strconv.Itoa(pid)))

	return pid
}

// 判断pid文件是否存在，以及pid进程是否存在
func Existed() bool {
	if !utils.IsExist(pidFile) {
		return false
	}

	res, err := utils.ReadFile(pidFile)
	if err != nil {
		return false
	}
	pid, _ := strconv.Atoi(string(res))
	if pid <= 0 {
		return false
	}

	if err = syscall.Kill(pid, 0); err == nil {
		return true
	}

	return false
}

// 删除pid文件
func DelPid() {
	if utils.IsExist(pidFile) {
		return
	}
	os.Remove(pidFile)
}
