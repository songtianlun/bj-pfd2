package main

import (
	"bj-pfd2/com/cfg"
	"bj-pfd2/com/log"
	"bj-pfd2/com/v"
	"bj-pfd2/com/web"
	"bj-pfd2/handle"
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

	// Debug
	fullData := handle.GetAllData("secret_OP7h4rxVR2XE08nz0PPK4qqyGBCwDcrgXpNQgkH1EKI", false)
	fullData.StatisticAll()
	fullData.Report()
	fullData.ShowChartData()

	//fmt.Println(fullData.Accounts.GenerateReport())

	Addr := ":" + strconv.FormatInt(cfg.GetInt64("Port"), 10)
	log.InfoF("BJ-PFD2[%v] is running on %v", v.GetVersionStr(), Addr)
	web.Run(Addr)
}
