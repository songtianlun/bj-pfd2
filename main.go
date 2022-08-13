package main

import (
	"bj-pfd2/com/cfg"
	"bj-pfd2/com/log"
	"bj-pfd2/com/utils"
	"bj-pfd2/com/web"
	"strconv"
)

func main() {
	// Step1 - 初始化配置
	initCfg()
	initLog()
	//initDB() // 本应用设计脱离数据库
	initCacheDB()

	// Step2 - 检查是否为命令行允许
	if runCLI() {
		return
	}

	// Step3 - 准备web服务

	initHandle()

	Addr := ":" + strconv.FormatInt(cfg.GetInt64("Port"), 10)
	log.Info("BJ-PFD2 "+utils.Version()+" started at ", Addr)
	web.Run(Addr)
}
