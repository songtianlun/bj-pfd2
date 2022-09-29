package main

import (
	"bj-pfd2/pkg/cfg"
	"bj-pfd2/pkg/log"
	"bj-pfd2/pkg/v"
	"bj-pfd2/pkg/web"
	"strconv"
)

func main() {
	// Step1 - 初始化配置
	initCfg()
	initLog()
	initCacheDB()

	// Step2 - 检查是否为命令行允许
	if runCLI() {
		return
	}

	// Step3 - 准备web服务
	initHandle()

	Addr := ":" + strconv.FormatInt(cfg.GetInt64("Port"), 10)
	log.Infof("BJ-PFD2[%v] is running on %v", v.GetVersionStr(), Addr)
	web.Run(Addr)
}
