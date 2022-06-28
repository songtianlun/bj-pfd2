package main

import (
	"bj-pfd2/com/cfg"
	"bj-pfd2/com/log"
	"bj-pfd2/com/utils"
	"bj-pfd2/com/web"
	"bj-pfd2/handle"
	"strconv"
)

func main() {
	if runCLI() {
		return
	}
	initCfg()
	initLog()
	initDB()
	initCacheDB()
	initHandle()

	handle.TestCode()

	Addr := ":" + strconv.FormatInt(cfg.GetInt64("Port"), 10)
	log.Info("BJ-PFD2 "+utils.Version()+" started at ", Addr)
	web.Run(Addr)
}
