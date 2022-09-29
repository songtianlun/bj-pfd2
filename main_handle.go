package main

import (
	"bj-pfd2/handle"
	"bj-pfd2/pkg/cache"
	"bj-pfd2/pkg/cfg"
	"bj-pfd2/pkg/cli"
	"bj-pfd2/pkg/log"
	"bj-pfd2/pkg/v"
	"bj-pfd2/pkg/web"
	"embed"
	"fmt"
)

//go:embed public
var efsStatic embed.FS

//go:embed templates
var tplEFS embed.FS

// Step1 - 初始化配置
func initCfg() {
	// 首先完成配置项的注册
	cfg.RegisterCfg("Port", 6010, "int64")
	cfg.RegisterCfg("ReadTimeout", 10, "int64")
	cfg.RegisterCfg("WriteTimeout", 600, "int64")
	cfg.RegisterCfg("SessionTimeoutHour", 6, "int64")
	// log
	cfg.RegisterCfg("log.level", "info", "string")
	cfg.RegisterCfg("log.file_name", "log/minegin.log", "string")
	cfg.RegisterCfg("log.max_size_mb", 1, "int")
	cfg.RegisterCfg("log.max_file_num", 64, "int")
	cfg.RegisterCfg("log.max_file_day", 7, "int")
	cfg.RegisterCfg("log.compress", false, "bool")
	cfg.RegisterCfg("log.stdout", true, "bool")
	cfg.RegisterCfg("log.only_stdout", true, "bool")
	// cache
	cfg.RegisterCfg("cache.enable", true, "bool")
	cfg.RegisterCfg("cache.type", "memory", "string")
	cfg.RegisterCfg("cache.addr", "127.0.0.1:6379", "string")
	cfg.RegisterCfg("cache.passwd", "", "string")
	cfg.RegisterCfg("cache.db", 0, "int")
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

// Step2 - 初始化日志
func initLog() {
	log.InitGlobal(log.NewLogrus(&log.CfgLog{
		FileName:   cfg.GetString("log.file_name"),
		Level:      cfg.GetString("log.level"),
		MaxSizeMB:  cfg.GetInt("log.max_size_mb"),
		MaxFileNum: cfg.GetInt("log.max_file_num"),
		MaxFileDay: cfg.GetInt("log.max_file_day"),
		Compress:   cfg.GetBool("log.compress"),
		Stdout:     cfg.GetBool("log.stdout"),
		OnlyStdout: cfg.GetBool("log.only_stdout"),
	}))
}

// Step3 - 初始化缓存
func initCacheDB() {
	if cfg.GetBool("cache.enable") {
		cache.Init(
			true,
			cfg.GetString("cache.type"),
			cfg.GetString("cache.addr"),
			cfg.GetString("cache.passwd"),
			cfg.GetInt("cache.db"))
	} else {
		cache.Init(
			false,
			"", "", "", 0)
	}
}

// Step4 - 初始化 CLI 命令
func runCLI() (isCli bool) {
	cli.RegisterBoolCLI("version", "V", "show version info.", func(mapCli cli.MapCli) {
		fmt.Println(v.GetVersionStr())
	})
	cli.RegisterStringCLI("token", "T", "", "Get Report With Notion Token.", func(mapCli cli.MapCli) {
		handle.ReportWithToken(*mapCli["token"].SValue)
	})
	return cli.CheckCLI()
}

// Step5 - 初始化 web 服务
func initHandle() {
	// static file
	//web.RegisterDir("/static/", "public", true)
	web.RegisterEmbedFs("/static/*filepath", &efsStatic)
	web.RegisterTplEmbedFs(&tplEFS)

	web.RegisterDefaultHandles(handle.Log)
	// index
	web.RegisterHandle("get", "/", handle.Index, handle.Auth)
	web.RegisterHandle("get", "/home", handle.Home, handle.Auth)

	// error
	web.RegisterHandle("get", "/err", handle.Err)

	// defined in route_auth.go
	web.RegisterHandle("get", "/login", handle.Login)
	web.RegisterHandle("post", "/authenticate", handle.Authenticate)
	//
	web.RegisterHandle("get", "/logout", handle.Logout)
}
