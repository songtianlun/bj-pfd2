package main

import (
	"bj-pfd2/com/cache"
	"bj-pfd2/com/cfg"
	"bj-pfd2/com/cli"
	"bj-pfd2/com/db"
	"bj-pfd2/com/log"
	"bj-pfd2/com/v"
	"bj-pfd2/com/web"
	"bj-pfd2/handle"
	"fmt"
)

func runCLI() (isCli bool) {
	cli.RegisterCLI("version", "V", "show version info.", func() {
		fmt.Println(v.GetVersionStr())
	})
	return cli.CheckCLI()
}

func initCfg() {
	// 首先完成配置项的注册
	cfg.RegisterCfg("Port", 6010, "int64")
	cfg.RegisterCfg("ReadTimeout", 10, "int64")
	cfg.RegisterCfg("WriteTimeout", 600, "int64")
	cfg.RegisterCfg("Static", "public", "string")
	cfg.RegisterCfg("SessionTimeoutHour", 6, "int64")
	// log
	cfg.RegisterCfg("log.level", "info", "string")
	cfg.RegisterCfg("log.file_name", "log/minegin.log", "string")
	cfg.RegisterCfg("log.max_size_mb", 1, "int")
	cfg.RegisterCfg("log.max_file_num", 64, "int")
	cfg.RegisterCfg("log.max_file_day", 7, "int")
	cfg.RegisterCfg("log.compress", false, "bool")
	cfg.RegisterCfg("log.stdout", true, "bool")
	cfg.RegisterCfg("log.only_stdout", false, "bool")
	// redis
	cfg.RegisterCfg("redis.addr", "", "string")
	cfg.RegisterCfg("redis.passwd", "", "string")
	cfg.RegisterCfg("redis.db", 0, "int")
	// bjpfd
	cfg.RegisterCfg("bjpfd.notion_token", "", "string")
	cfg.RegisterCfg("bjpfd.account_pid", "", "string")
	cfg.RegisterCfg("bjpfd.bills_pid", "", "string")
	cfg.RegisterCfg("bjpfd.i_account_pid", "", "string")
	cfg.RegisterCfg("bjpfd.investment_pid", "", "string")
	cfg.RegisterCfg("bjpfd.budget_pid", "", "string")

	// 之后再进行初始化
	err := cfg.Init("")
	if err != nil {
		panic(fmt.Sprintf("init cfg failed: %s", err))
	}
}

func initLog() {
	log.InitLogger(
		cfg.GetString("log.file_name"),
		cfg.GetString("log.level"),
		cfg.GetInt("log.max_size_mb"),
		cfg.GetInt("log.max_file_num"),
		cfg.GetInt("log.max_file_day"),
		cfg.GetBool("log.compress"))
}

func initDB() {
	db.InitDB(&db.CfgDb{
		Typ:      cfg.GetString("db.type"),
		Addr:     cfg.GetString("db.addr"),
		Name:     cfg.GetString("db.name"),
		Username: cfg.GetString("db.username"),
		Passwd:   cfg.GetString("db.password"),
	})
}

func initCacheDB() {
	err := cache.InitClient(
		&cache.CfgRedis{
			Addr:   cfg.GetString("redis.addr"),
			Passwd: cfg.GetString("redis.passwd"),
			Db:     cfg.GetInt("redis.db"),
		})
	if err != nil {
		log.ErrorF("Failed to init Cache DB: %s, cache will not be work.", err.Error())
		return
	}
}

func initHandle() {
	// static file
	web.RegisterFile("/static/", cfg.GetString("Static"), true)

	// index
	web.RegisterHandle("/", handle.Index)

	// error
	//web.RegisterHandle("/err", handle.Err)

	// defined in route_auth.go
	//web.RegisterHandle("/login", handle.Login)
	//web.RegisterHandle("/signup", handle.Signup)
	//web.RegisterHandle("/signup_account", handle.SignupAccount)
	//web.RegisterHandle("/authenticate", handle.Authenticate)
	//
	//web.RegisterHandle("/logout", handle.Logout)
}
